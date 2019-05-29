package hkcam

import (
	"fmt"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/log"
	"github.com/brutella/hc/rtp"
	"github.com/brutella/hc/service"
	"github.com/brutella/hc/tlv8"
	"net"
	"reflect"

	"github.com/brutella/hkcam/ffmpeg"
)

// SetupFFMPEGStreaming configures a camera to use ffmpeg to stream video.
// The returned handle can be used to interact with the camera (start, stop, take snapshot…).
func SetupFFMPEGStreaming(cam *accessory.Camera, cfg ffmpeg.Config) ffmpeg.FFMPEG {
	ff := ffmpeg.New(cfg)

	setupStreamManagement(cam.StreamManagement1, ff, cfg.MultiStream)
	setupStreamManagement(cam.StreamManagement2, ff, cfg.MultiStream)

	return ff
}

func first(ips []net.IP, filter func(net.IP) bool) net.IP {
	for _, ip := range ips {
		if filter(ip) == true {
			return ip
		}
	}

	return nil
}

func setupStreamManagement(m *service.CameraRTPStreamManagement, ff ffmpeg.FFMPEG, multiStream bool) {
	status := rtp.StreamingStatus{rtp.StreamingStatusAvailable}
	setTLV8Payload(m.StreamingStatus.Bytes, status)
	setTLV8Payload(m.SupportedRTPConfiguration.Bytes, rtp.NewSupportedRTPConfiguration(rtp.CryptoSuite_AES_CM_128_HMAC_SHA1_80))
	setTLV8Payload(m.SupportedVideoStreamConfiguration.Bytes, rtp.DefaultVideoStreamConfiguration())
	setTLV8Payload(m.SupportedAudioStreamConfiguration.Bytes, rtp.DefaultAudioStreamConfiguration())

	m.SelectedRTPStreamConfiguration.OnValueRemoteUpdate(func(buf []byte) {
		var cfg rtp.SelectedRtpStreamConfiguration
		err := tlv8.Unmarshal(buf, &cfg)
		if err != nil {
			log.Debug.Fatalf("SelectedRTPStreamConfiguration: Could not unmarshal tlv8 data: %s\n", err)
		}

		log.Debug.Printf("%+v\n", cfg)

		id := ffmpeg.StreamID(cfg.Command.Identifier)
		switch cfg.Command.Type {
		case rtp.SessionControlCommandTypeEnd:
			ff.Stop(id)

			if ff.ActiveStreams() == 0 {
				// Update stream status when no streams are currently active
				setTLV8Payload(m.StreamingStatus.Bytes, rtp.StreamingStatus{rtp.StreamingStatusAvailable})
			}

		case rtp.SessionControlCommandTypeStart:
			ff.Start(id, cfg.Video, cfg.Audio)

			if multiStream == false {
				// If only one video stream is suppported, set the status to busy.
				// This way HomeKit knows that nobody is allowed to connect anymore.
				// If multiple streams are supported, the status is always availabe.
				setTLV8Payload(m.StreamingStatus.Bytes, rtp.StreamingStatus{rtp.StreamingStatusBusy})
			}
		case rtp.SessionControlCommandTypeSuspend:
			ff.Suspend(id)
		case rtp.SessionControlCommandTypeResume:
			ff.Resume(id)
		case rtp.SessionControlCommandTypeReconfigure:
			ff.Reconfigure(id, cfg.Video, cfg.Audio)
		default:
			log.Debug.Printf("Unknown command type %d", cfg.Command.Type)
		}
	})

	m.SetupEndpoints.OnValueUpdateFromConn(func(conn net.Conn, c *characteristic.Characteristic, new, old interface{}) {
		buf := m.SetupEndpoints.GetValue()
		var req rtp.SetupEndpoints
		err := tlv8.Unmarshal(buf, &req)
		if err != nil {
			log.Debug.Fatalf("SetupEndpoints: Could not unmarshal tlv8 data: %s\n", err)
		}

		log.Debug.Printf("%+v\n", req)

		iface, err := ifaceOfConnection(conn)
		if err != nil {
			log.Debug.Println(err)
			return
		}
		ip, err := ipAtInterface(*iface, req.ControllerAddr.IPVersion)
		if err != nil {
			log.Debug.Println(err)
			return
		}

		// TODO ssrc is different for every stream
		ssrcVideo := int32(1)
		ssrcAudio := int32(2)

		resp := rtp.SetupEndpointsResponse{
			SessionId: req.SessionId,
			Status:    rtp.SessionStatusSuccess,
			AccessoryAddr: rtp.Addr{
				IPVersion:    req.ControllerAddr.IPVersion,
				IPAddr:       ip.String(),
				VideoRtpPort: req.ControllerAddr.VideoRtpPort,
				AudioRtpPort: req.ControllerAddr.AudioRtpPort,
			},
			Video:     req.Video,
			Audio:     req.Audio,
			SsrcVideo: ssrcVideo,
			SsrcAudio: ssrcAudio,
		}

		ff.PrepareNewStream(req, resp)

		log.Debug.Printf("%+v\n", resp)

		// After a write, the characteristic should contain a response
		setTLV8Payload(m.SetupEndpoints.Bytes, resp)
	})
}

func ipAtInterface(iface net.Interface, version uint8) (net.IP, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		log.Debug.Println(err)
		return nil, err
	}

	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			log.Debug.Println(err)
			continue
		}

		switch version {
		case rtp.IPAddrVersionv4:
			if ip.To4() != nil {
				return ip, nil
			}
		case rtp.IPAddrVersionv6:
			if ip.To16() != nil {
				return ip, nil
			}
		default:
			break
		}
	}

	return nil, fmt.Errorf("%s: No ip address found for version %d", iface.Name, version)
}

func ifaceOfConnection(conn net.Conn) (*net.Interface, error) {
	host, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return nil, fmt.Errorf("unable to parse ip %s", host)
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			addrIP, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				return nil, err
			}

			if reflect.DeepEqual(addrIP, ip) {
				return &iface, nil
			}
		}
	}

	return nil, fmt.Errorf("Could not find interface for connection")
}

func setTLV8Payload(c *characteristic.Bytes, v interface{}) {
	if tlv8, err := tlv8.Marshal(v); err == nil {
		c.SetValue(tlv8)
	} else {
		log.Debug.Fatal(err)
	}
}
