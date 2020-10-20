package main

import (
	"flag"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/log"
	"github.com/brutella/hc/service"

	"github.com/stianeikeland/go-rpio/v4"

	"image"
	"runtime"

	"time"

	"github.com/brutella/hkcam"
	"github.com/brutella/hkcam/ffmpeg"
)

func main() {

	// Platform dependent flags
	var inputDevice *string
	var inputFilename *string
	var loopbackFilename *string
	var h264Encoder *string
	var h264Decoder *string
	var motionPin int = -1

	if runtime.GOOS == "linux" {
		inputDevice = flag.String("input_device", "v4l2", "video input device")
		inputFilename = flag.String("input_filename", "/dev/video0", "video input device filename")
		loopbackFilename = flag.String("loopback_filename", "/dev/video1", "video loopback device filename")
		h264Decoder = flag.String("h264_decoder", "", "h264 video decoder")
		h264Encoder = flag.String("h264_encoder", "h264_omx", "h264 video encoder")
		flag.IntVar(&motionPin, "motion_pin", -1, "GPIO pin of any attached motion sensor")
	} else if runtime.GOOS == "darwin" { // macOS
		inputDevice = flag.String("input_device", "avfoundation", "video input device")
		inputFilename = flag.String("input_filename", "default", "video input device filename")
		// loopback is not needed on macOS because avfoundation provides multi-access to the camera
		loopbackFilename = flag.String("loopback_filename", "", "video loopback device filename")
		h264Decoder = flag.String("h264_decoder", "", "h264 video decoder")
		h264Encoder = flag.String("h264_encoder", "libx264", "h264 video encoder")
	} else {
		log.Info.Fatalf("%s platform is not supported", runtime.GOOS)
	}

	var minVideoBitrate *int = flag.Int("min_video_bitrate", 0, "minimum video bit rate in kbps")
	var multiStream *bool = flag.Bool("multi_stream", false, "Allow mutliple clients to view the stream simultaneously")
	var dataDir *string = flag.String("data_dir", "Camera", "Path to data directory")
	var verbose *bool = flag.Bool("verbose", true, "Verbose logging")
	var pin *string = flag.String("pin", "00102003", "PIN for HomeKit pairing")
	var port *string = flag.String("port", "", "Port on which transport is reachable")

	flag.Parse()

	if *verbose {
		log.Debug.Enable()
		ffmpeg.EnableVerboseLogging()
	}

	switchInfo := accessory.Info{Name: "Camera", FirmwareRevision: "0.0.9", Manufacturer: "Matthias Hochgatterer"}
	cam := accessory.NewCamera(switchInfo)

	cfg := ffmpeg.Config{
		InputDevice:      *inputDevice,
		InputFilename:    *inputFilename,
		LoopbackFilename: *loopbackFilename,
		H264Decoder:      *h264Decoder,
		H264Encoder:      *h264Encoder,
		MinVideoBitrate:  *minVideoBitrate,
		MultiStream:      *multiStream,
	}

	// Add Motion Sensor
	if motionPin != -1 {
		err := rpio.Open()
		if err != nil {
			log.Info.Panic(err)
		}
		pin := rpio.Pin(10)
		pin.Input()

		motionSensor := service.NewMotionSensor()
		cam.AddService(motionSensor.Service)

		go func() {
			for range time.Tick(time.Second) {
				motionSensor.MotionDetected.SetValue(pin.Read() == rpio.High)
			}
		}()
	}

	ffmpeg := hkcam.SetupFFMPEGStreaming(cam, cfg)

	// Add a custom camera control service to record snapshots
	cc := hkcam.NewCameraControl()
	cam.Control.AddCharacteristic(cc.Assets.Characteristic)
	cam.Control.AddCharacteristic(cc.GetAsset.Characteristic)
	cam.Control.AddCharacteristic(cc.DeleteAssets.Characteristic)
	cam.Control.AddCharacteristic(cc.TakeSnapshot.Characteristic)

	t, err := hc.NewIPTransport(hc.Config{StoragePath: *dataDir, Pin: *pin, Port: *port}, cam.Accessory)
	if err != nil {
		log.Info.Panic(err)
	}

	t.CameraSnapshotReq = func(width, height uint) (*image.Image, error) {
		return ffmpeg.Snapshot(width, height)
	}

	cc.SetupWithDir(*dataDir)
	cc.CameraSnapshotReq = func(width, height uint) (*image.Image, error) {
		return ffmpeg.Snapshot(width, height)
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()
}
