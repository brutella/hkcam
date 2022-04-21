package html

import (
	"fmt"
	"net/http"
)

// ErrorPage is an error page.
type ErrorPage struct {
	Page
	CallbackUrl string
	Error       string
	Reason      string
}

func (p *ErrorPage) UpdateWithRequest(r *http.Request, h *Html) {
	p.Page.UpdateWithRequest(r, h)
	p.Error = "Error"
	p.Title = p.Error
}

func (h *Html) Error(w http.ResponseWriter, r *http.Request, err error) {
	p := ErrorPage{
		Reason: fmt.Sprintf("%s", err),
	}
	p.UpdateWithRequest(r, h)
	h.HTML(w, r, http.StatusOK, "error", "layout", &p)
}
