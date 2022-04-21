package html

import (
	"github.com/brutella/go-github-selfupdate/selfupdate"
	"github.com/brutella/hkcam/app"

	"encoding/json"
	"net/http"
)

type UdpdatePage struct {
	Page
}

// CheckForUpdate downloads the latest release and stores it as an db.Update in the database.
// The user is then redirected to the previous page.
// If no update was found, a message is shown to the user that no updates are available.
func (h *Html) CheckForUpdate(w http.ResponseWriter, r *http.Request) {
	p := UdpdatePage{}
	p.UpdateWithRequest(r, h)

	selfupdate.EnableLog()

	if u, err := h.App.CheckForUpdate(false); err != nil {
		h.Error(w, r, err)
		return
	} else {
		url := p.Referrer
		if u == nil {
			url = setURLMsg(url, "No Update Available")
		}

		h.SaveUpdate(u)
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

// InstallLatestVersion installs the latest release no matter of the current build version
func (h *Html) InstallLatestVersion(w http.ResponseWriter, r *http.Request) {
	p := UdpdatePage{}
	p.UpdateWithRequest(r, h)

	url := p.Referrer
	if u := p.Update; u == nil {
		// fetch latest version (including pre-releases)
		u, err := h.App.LatestVersion(true)
		if err != nil {
			url = setURLMsg(url, err.Error())
		} else if u == nil {
			url = setURLMsg(url, "No Update Available")
		} else {
			h.SaveUpdate(u)
			go func() {
				h.App.InstallUpdate(u)
				h.SaveUpdate(u)
			}()
		}
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}

// InstallUpdate installs the latest release, if an update is already store din the database.
// If no update is stored in the database or an install process currently running, this method does nothing.
func (h *Html) InstallUpdate(w http.ResponseWriter, r *http.Request) {
	p := UdpdatePage{}
	p.UpdateWithRequest(r, h)
	if u := p.Update; u != nil {
		switch u.State {
		case app.UpdateStateInstall:
			// don't try to install twice
			break
		default:
			go func() {
				h.App.InstallUpdate(u)
				h.SaveUpdate(u)
			}()
		}
	}
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

// CleanupUpdateAndRestart deletes the latest update from the database and render the restart page.
// The page then restarts the system by terminating it via an Api call.
func (h *Html) CleanupUpdateAndRestart(w http.ResponseWriter, r *http.Request) {
	// delete latest update
	u := h.LatestUpdate()
	h.DeleteUpdate(u)

	p := Page{}
	p.UpdateWithRequest(r, h)
	p.Title = "Restart"
	h.HTML(w, r, http.StatusOK, "restart", "layout", &p)
}

func (h *Html) LatestUpdate() *app.Update {
	if h.u != nil {
		return h.u
	}

	b, err := h.Store.Get("update")
	if err != nil {
		return nil
	}

	var u app.Update
	if err := json.Unmarshal(b, &u); err != nil {
		return nil
	}

	h.u = &u

	return &u
}

func (h *Html) SaveUpdate(update *app.Update) error {
	b, err := json.Marshal(&update)
	if err != nil {
		return err
	}
	h.u = update
	return h.Store.Set("update", b)
}

func (h *Html) DeleteUpdate(update *app.Update) {
	h.u = nil
	h.Store.Delete("update")
}
