package main

import (
	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/log"
	"github.com/brutella/hkcam"
	"github.com/brutella/hkcam/ffmpeg"

	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	// Date is build date.
	Date string

	// Version is the app version.
	Version string
)

func main() {

	// Platform dependent flags
	var inputDevice *string
	var inputFilename *string
	var loopbackFilename *string
	var h264Encoder *string
	var h264Decoder *string

	if runtime.GOOS == "linux" {
		inputDevice = flag.String("input_device", "v4l2", "video input device")
		inputFilename = flag.String("input_filename", "/dev/video0", "video input device filename")
		loopbackFilename = flag.String("loopback_filename", "/dev/video99", "video loopback device filename")
		h264Decoder = flag.String("h264_decoder", "", "h264 video decoder")
		h264Encoder = flag.String("h264_encoder", "h264_v4l2m2m", "h264 video encoder")
	} else if runtime.GOOS == "darwin" { // macOS
		inputDevice = flag.String("input_device", "avfoundation", "video input device")
		inputFilename = flag.String("input_filename", "default", "video input device filename")
		// loopback is not needed on macOS because avfoundation provides multi-access to the camera
		loopbackFilename = flag.String("loopback_filename", "", "video loopback device filename")
		h264Decoder = flag.String("h264_decoder", "", "h264 video decoder")
		h264Encoder = flag.String("h264_encoder", "h264_videotoolbox", "h264 video encoder")
	} else {
		log.Info.Fatalf("%s platform is not supported", runtime.GOOS)
	}

	var minVideoBitrate *int = flag.Int("min_video_bitrate", 0, "minimum video bit rate in kbps")
	var multiStream *bool = flag.Bool("multi_stream", false, "Allow mutliple clients to view the stream simultaneously")
	var dataDir *string = flag.String("data_dir", "db", "Path to data directory")
	var verbose *bool = flag.Bool("verbose", false, "Verbose logging")
	var pin *string = flag.String("pin", "00102003", "PIN for HomeKit pairing")
	var port *string = flag.String("port", "", "Port on which transport is reachable")

	flag.Parse()

	if *verbose {
		log.Debug.Enable()
		ffmpeg.EnableVerboseLogging()
	}

	log.Info.Printf("version %s (built at %s)\n", Version, Date)

	switchInfo := accessory.Info{Name: "Camera", Firmware: Version, Manufacturer: "Matthias Hochgatterer"}
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

	ffmpeg := hkcam.SetupFFMPEGStreaming(cam, cfg)

	// Add a custom camera control service to record snapshots
	cc := hkcam.NewCameraControl()
	cam.Control.AddC(cc.Assets.C)
	cam.Control.AddC(cc.GetAsset.C)
	cam.Control.AddC(cc.DeleteAssets.C)
	cam.Control.AddC(cc.TakeSnapshot.C)

	s, err := hap.NewServer(hap.NewFsStore(*dataDir), cam.A)
	if err != nil {
		log.Info.Panic(err)
	}

	s.Pin = *pin
	s.Addr = fmt.Sprintf(":%s", *port)

	s.ServeMux().HandleFunc("/resource", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Info.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		r := struct {
			Type   string `json:"resource-type"`
			Width  uint   `json:"image-width"`
			Height uint   `json:"image-height"`
		}{}

		err = json.Unmarshal(body, &r)
		if err != nil {
			log.Info.Println(err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Debug.Printf("%+v\n", r)

		switch r.Type {
		case "image":
			b, err := snapshot(r.Width, r.Height, ffmpeg)
			if err != nil {
				log.Info.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			res.Header().Set("Content-Type", "image/jpeg")
			wr := hap.NewChunkedWriter(res, 2048)
			wr.Write(b)
		default:
			log.Info.Printf("unsupported resource request \"%s\"\n", r.Type)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	cc.SetupWithDir(*dataDir)
	cc.CameraSnapshotReq = func(width, height uint) (*image.Image, error) {
		return ffmpeg.Snapshot(width, height)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		signal.Stop(c) // stop delivering signals
		cancel()
	}()

	s.ListenAndServe(ctx)
}

func snapshot(width, height uint, ffmpeg ffmpeg.FFMPEG) ([]byte, error) {
	log.Debug.Printf("snasphot %dw x %dh\n", width, height)

	img, err := ffmpeg.Snapshot(width, height)
	if err != nil {
		return nil, fmt.Errorf("snapshot: %v", err)
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, *img, nil); err != nil {
		return nil, fmt.Errorf("encode: %v", err)
	}

	return buf.Bytes(), nil
}
