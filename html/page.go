package html

import (
	"github.com/brutella/hkcam/app"

	"html/template"
	"net/http"
)

// Page is the data shown in an HTML page.
type Page struct {
	// Title of the page
	Title string

	// Referrer of this page without the query part.
	// Use this if you don't want to exclude any messages embed in the url
	Referrer string

	// Referer is `http.Request.Referer()` and includes the query part from the url.
	Referer string

	// Message is the message shown
	Message string

	// Update is a system update.
	Update *app.Update

	// App is the app.
	App *app.App

	Funcs template.FuncMap

	// DebugMode is true when running in debug mode
	DebugMode bool
}

func (p *Page) UpdateWithRequest(r *http.Request, h *Html) {
	p.Referer = r.Referer()
	p.Referrer = delURLMsg(p.Referer)
	p.Message = getURLMsg(r)
	p.Update = h.LatestUpdate()
	p.App = h.App
	p.DebugMode = h.BuildMode == "debug"
}
