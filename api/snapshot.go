package api

import (
	"github.com/brutella/hkcam/api/apiutil"

	"bytes"
	"fmt"
	"image/jpeg"
	"net/http"
	"time"
)

// SnapshotRequest is a request restart the system.
type SnapshotRequest struct {
	Width  uint `schema:"width"`
	Height uint `schema:"height"`
}

// SnapshotResponse is a response to a SnapshotRequest.
type SnapshotResponse struct {
	Data *SnapshotResponseData `json:"data"`
}

// SnapshotResponseData is the response data of a SnapshotRequest.
type SnapshotResponseData struct {
	Date  *time.Time `json:"date,omitempty"`
	Bytes []byte     `json:"bytes"`
}

// RecentSnapshot responds with the recent snapshot.
func (a *Api) RecentSnapshot(w http.ResponseWriter, r *http.Request) {
	req := SnapshotRequest{
		Width:  1920,
		Height: 1080,
	}
	var resp interface{}

	if err := apiutil.DecodeURLQuery(w, r, &req); err != nil {
		resp = NewErrResponse(fmt.Errorf("invalid payload"), ErrorInvalidPayload)
	} else if req.Width == 0 || req.Height == 0 {
		resp = NewErrResponse(fmt.Errorf("invalid payload"), ErrorInvalidPayload)
	} else if snapshot := a.App.FFMPEG.RecentSnapshot(req.Width, req.Height); snapshot != nil {
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, snapshot.Image, nil); err != nil {
			resp = NewErrResponse(fmt.Errorf("encode: %v", err), ErrorUnknown)
		}
		resp = SnapshotResponse{
			Data: &SnapshotResponseData{
				Bytes: buf.Bytes(),
				Date:  &snapshot.Date,
			},
		}
	} else {
		resp = SnapshotResponse{}
	}

	if err := WriteJSON(w, r, resp); err != nil {
		fmt.Println("responding failed", err)
	}
}

// NewSnapshot create a new snapshot.
func (a *Api) NewSnapshot(w http.ResponseWriter, r *http.Request) {
	req := SnapshotRequest{
		Width:  1920,
		Height: 1080,
	}
	var resp interface{}

	if err := apiutil.DecodeURLQuery(w, r, &req); err != nil {
		resp = NewErrResponse(fmt.Errorf("invalid payload"), ErrorInvalidPayload)
	} else if req.Width == 0 || req.Height == 0 {
		resp = NewErrResponse(fmt.Errorf("invalid payload"), ErrorInvalidPayload)
	} else if snapshot, err := a.App.FFMPEG.Snapshot(req.Width, req.Height); err != nil {
		resp = NewErrResponse(fmt.Errorf("snapshot: %v", err), ErrorUnknown)
	} else {
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, snapshot.Image, nil); err != nil {
			resp = NewErrResponse(fmt.Errorf("encode: %v", err), ErrorUnknown)
		} else {
			resp = SnapshotResponse{
				Data: &SnapshotResponseData{
					Bytes: buf.Bytes(),
				},
			}
		}
	}

	if err := WriteJSON(w, r, resp); err != nil {
		fmt.Println("responding failed", err)
	}
}
