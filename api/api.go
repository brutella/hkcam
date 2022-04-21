package api

import (
	"github.com/brutella/hkcam/app"
	"github.com/go-chi/chi"

	"net/http"
)

const (
	ErrorInvalidPayload = 1
	ErrorInvalidRequest = 2
	ErrorUnknown        = 3
)

type Api struct {
	App *app.App
}

func (a *Api) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/system/heartbeat", a.SystemHeartbeat)
	r.Get("/system/info", a.SystemInfo)
	r.Post("/system/restart", a.SystemRestart)
    r.Get("/snapshots/recent", a.RecentSnapshot)
	r.Get("/snapshots/new", a.NewSnapshot)

	return r
}

// RestartApp restarts the app.
func (a *Api) RestartApp() {
	a.App.Restart()
}
