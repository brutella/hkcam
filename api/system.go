package api

import (
	"fmt"
	"net/http"
	"syscall"
	"time"
)

// SystemRestartRequest is a request restart the system.
type SystemRestartRequest struct {
}

// SystemRestartResponse is a response to a SystemRestartRequest.
type SystemRestartResponse struct {
	Data SystemRestartResponseData `json:"data"`
}

// SystemRestartResponseData is the response data of a SystemRestartRequest.
type SystemRestartResponseData struct {
	Success bool `json:"success"`
}

// SystemRestart triggers a system restart by terminating the app.
func (a *Api) SystemRestart(w http.ResponseWriter, r *http.Request) {
	var resp = SystemRestartResponse{
		Data: SystemRestartResponseData{
			Success: true,
		},
	}
	if err := WriteJSON(w, r, resp); err != nil {
		fmt.Println("responding failed", err)
	}

	go func() {
		// sleep so we can be sure that the client gets a response
		// before the process is killed
		time.Sleep(1 * time.Second)

		// SIGUSR1 has to be handled in main
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	}()
}

// SystemHeartbeatRequest is a request check the availability of the system.
type SystemHeartbeatRequest struct {
}

// SystemHeartbeatResponse is a response to a SystemHeartbeatRequest.
type SystemHeartbeatResponse struct {
	Data SystemHeartbeatResponseData `json:"data"`
}

// SystemHeartbeatResponseData is the response data of a SystemHeartbeatRequest.
type SystemHeartbeatResponseData struct {
	Success bool `json:"success"`
}

// SystemHeartbeat returns the system heartbeat.
func (a *Api) SystemHeartbeat(w http.ResponseWriter, r *http.Request) {
	var resp = SystemHeartbeatResponse{
		Data: SystemHeartbeatResponseData{
			Success: true,
		},
	}
	if err := WriteJSON(w, r, resp); err != nil {
		fmt.Println("responding failed", err)
	}
}

// SystemInfoRequest is a request check the system info.
type SystemInfoRequest struct{}

// SystemInfoResponse is a response to a SystemInfoRequest.
type SystemInfoResponse struct {
	Data SystemInfoResponseData `json:"data"`
}

// SystemInfoResponseData is the response data of a SystemInfoRequest.
type SystemInfoResponseData struct {
	Version string  `json:"version"`
	Uptime  float64 `json:"uptime"`
}

// SystemInfo returns the system info.
func (a *Api) SystemInfo(w http.ResponseWriter, r *http.Request) {
	var resp = SystemInfoResponse{
		Data: SystemInfoResponseData{
			Version: a.App.Version,
			Uptime:  time.Since(a.App.Launch).Seconds(),
		},
	}

	if err := WriteJSON(w, r, resp); err != nil {
		fmt.Println("responding failed", err)
	}
}
