package ffmpeg

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

// snapshot returns an image by grapping a frame of the video stream.
func snapshot(width, height uint, inputDevice, inputFilename string) (*image.Image, error) {
	fileName := fmt.Sprintf("snapshot_%s.jpeg", time.Now().Format(time.RFC3339))
	filePath := path.Join(os.TempDir(), fileName)

	// height "-2" keeps the aspect ratio
	var arg string
	if inputDevice == "rtsp" {
		arg = fmt.Sprintf("-i %s -vf scale=%d:-2 -frames:v 1 %s", inputFilename, width, filePath)
	} else {
		// local
		arg = fmt.Sprintf("-f %s -framerate 30 -i %s -vf scale=%d:-2 -frames:v 1 %s", inputDevice, inputFilename, width, filePath)
	}
	args := strings.Split(arg, " ")

	cmd := exec.Command("ffmpeg", args[:]...)
	cmd.Stdout = Stdout
	cmd.Stderr = Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return loadImage(filePath)
}

func loadImage(path string) (*image.Image, error) {
	reader, _ := os.Open(path)
	defer reader.Close()
	img, _, err := image.Decode(reader)
	return &img, err
}
