package html

import (
	"net/http"
)

func (h *Html) Home(w http.ResponseWriter, r *http.Request) {
	var p Page
	p.UpdateWithRequest(r, h)
	p.Title = "hkcam"
	h.HTML(w, r, http.StatusOK, "home", "layout", &p)
}
