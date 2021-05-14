package ffmpeg

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/brutella/hc/log"
)

// loopback copies data from the inpute filename to the loopback filename.
// This is needed to provide simultaneous access to a v4l2 device.
// On a Raspberry Pi you can create a video loopback device at /dev/video1 via [v4l2loopback](https://github.com/umlaeute/v4l2loopback).
// The data from /dev/video0 is then copied to /dev/video1 via `ffmpeg -f v4l2 -i /dev/video0 -codec:v copy -f v4l2 /dev/video1`.
// The loopback device can then be access simultanously from multiple ffmpeg processes.
type loopback struct {
	inputDevice      string
	inputFilename    string
	h264Decoder      string
	loopbackFilename string

	mutex *sync.Mutex
	cmd   *exec.Cmd
	out   io.ReadWriter
}

// NewLoopback returns a new video loopback.
func NewLoopback(inputDevice, inputFilename, loopbackFilename string) *loopback {
	return &loopback{
		inputDevice:      inputDevice,
		inputFilename:    inputFilename,
		loopbackFilename: loopbackFilename,
		mutex:            &sync.Mutex{},
	}
}

// Start starts the loopback.
// This method waits until the ffmpeg process is running.
func (l *loopback) Start() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.cmd == nil {
		log.Debug.Println("Starting loopback")
		cmd := l.execCmd()
		pr, pw := io.Pipe()
		// cmd.Stdout = pw
		cmd.Stderr = pw

		if err := cmd.Start(); err != nil {
			return err
		}

		done := make(chan struct{}, 0)
		go func() {
			r := bufio.NewReader(pr)
			for {
				line, _, err := r.ReadLine()
				if err != nil {
					if err == io.EOF {
						log.Info.Println("ffmpeg: process stopped")
					} else {
						log.Info.Println("ffmpeg:", err)
					}
					return
				}
				log.Debug.Println(string(line))
				if strings.Contains(string(line), "Press [q] to stop, [?] for help") {
					log.Debug.Println("ffmpeg is now running")
					done <- struct{}{}
				}
			}
		}()

		select {
		case <-done:
			log.Debug.Println("Loopback started")
			l.cmd = cmd
			return nil
		case <-time.After(20 * time.Second):
			err := errors.New("Loopback failed to start")
			log.Debug.Println(err)
			cmd.Process.Signal(syscall.SIGINT)
			cmd.Wait()
			return err
		}
	}

	return nil
}

// Stop stops the loopback.
func (l *loopback) Stop() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.cmd != nil {
		log.Debug.Println("Stopping loopback")
		l.cmd.Process.Signal(syscall.SIGINT)
		l.cmd.Wait()
		l.cmd = nil
	}
}

// cmd returns a new command to stream video from the input file to the loopback file.
func (l *loopback) execCmd() *exec.Cmd {
	var cmd *exec.Cmd
	if l.inputDevice != "rtsp" {
		cmd = exec.Command("ffmpeg", "-f", l.inputDevice, "-i", l.inputFilename, "-codec:v", "copy", "-f", l.inputDevice, l.loopbackFilename)
	}

	log.Debug.Println(cmd)

	return cmd
}
