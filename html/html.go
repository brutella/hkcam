package html

import (
	"github.com/brutella/hap"
	"github.com/brutella/hkcam/api"
	"github.com/brutella/hkcam/app"
	"github.com/go-chi/chi"
	"github.com/unrolled/render"

	"fmt"
	"html/template"
	"net/http"
	neturl "net/url"
)

type Html struct {
	Store      hap.Store
	BuildMode  string
	Api        *api.Api
	App        *app.App
	FileSystem render.FileSystem
	Render     *render.Render
	u          *app.Update
}

func (h *Html) HTML(w http.ResponseWriter, r *http.Request, status int, tmpl string, layout string, binding interface{}) {
	opt := render.HTMLOptions{
		Layout: layout,
		Funcs: template.FuncMap{
			"T": func(format string, args ...interface{}) string {
				return fmt.Sprintf(format, args...)
			},
		},
	}

	h.Render.HTML(w, status, tmpl, binding, opt)
}

func (h *Html) Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/", h.Home)
	r.Post("/cleanup-update-and-restart", h.CleanupUpdateAndRestart)
	// r.Get("/system/log", h.Log)

	r.Post("/update/check", h.CheckForUpdate)
	r.Post("/update/install/latest", h.InstallLatestVersion)
	r.Post("/update/install", h.InstallUpdate)

	return r
}

// setURLParam returns a new url which contains val for key in url
func setURLParam(url, key, val string) string {
	u, err := neturl.Parse(url)
	if err == nil {
		vals := u.Query()
		vals.Set(key, val)
		u.RawQuery = vals.Encode()
		return u.String()
	}

	return url
}

// getURLValForKey returns the value for key in r.
func getURLParamForKey(r *http.Request, key string) string {
	return r.FormValue(key)
}

// setURLMsg set a message using the "msg" url value.
func setURLMsg(url, msg string) string {
	return setURLParam(url, "msg", msg)
}

// delURLMsg deletes a message from an url.
func delURLMsg(s string) string {
	u, err := neturl.Parse(s)
	if err != nil {
		return s
	}
	qu := u.Query()
	qu.Del("msg")
	u.RawQuery = qu.Encode()
	return u.String()
}

// getURLMsg returns the msg encoded into the request's url,
// or an empty string if no message is present.
func getURLMsg(r *http.Request) string {
	return getURLParamForKey(r, "msg")
}
