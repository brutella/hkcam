package app

import (
	"github.com/blang/semver"
	"github.com/brutella/go-github-selfupdate/selfupdate"
	"github.com/brutella/hap"
	"github.com/brutella/hap/log"
	"github.com/brutella/hkcam/ffmpeg"

	"fmt"
	"os"
	"time"
)

type App struct {
	BuildMode string
	BuildDate time.Time
	Version   string
	Launch    time.Time
	Store     hap.Store
	FFMPEG    ffmpeg.FFMPEG
}

func (app App) Restart() {
	log.Info.Println("restart not implemented yet")
}

// SemVersion returns the semantic version of the app.
func (app App) SemVersion() (semver.Version, error) {
	return semver.ParseTolerant(app.Version)
}

func (app App) CheckForUpdate(pre bool) (up *Update, err error) {
	var av, rv semver.Version

	av, err = app.SemVersion()
	if err != nil {
		return
	}

	up, err = app.LatestVersion(pre)
	if err != nil {
		return
	}

	if up == nil {
		log.Debug.Println("check for update: no new version found")
		return
	}

	rv, err = semver.ParseTolerant(up.Version)
	if err != nil {
		log.Debug.Println("check for update:", err)
		return
	}

	if rv.LTE(av) {
		log.Debug.Printf("check for update: %s <= %s\n", rv, av)
		up = nil
		return
	}
	return
}

func (app App) LatestVersion(pre bool) (*Update, error) {
	upt, err := selfupdate.NewUpdater(selfupdate.Config{PreRelease: pre})
	latest, found, err := upt.DetectLatest("brutella/hkcam")
	if err != nil {
		return nil, err
	}

	if !found {
		log.Debug.Println("check for update: no version found")
		return nil, nil
	}

	update := &Update{}
	update.State = UpdateStateDefault
	update.Version = latest.Version.String()
	update.PreRelease = latest.PreRelease
	update.URL = latest.URL

	return update, nil
}

// InstallUpdate performs an update to the latest version.
// If installing fails, the error is stored in up.Err.
func (app *App) InstallUpdate(up *Update) error {
	up.State = UpdateStateInstall
	up.Err = nil

	cmdPath, err := os.Executable()
	if err != nil {
		log.Info.Println(err)
		up.State = UpdateStateFailure
		up.Err = err
		return err
	}

	upt, err := selfupdate.NewUpdater(selfupdate.Config{PreRelease: up.PreRelease})
	if err != nil {
		log.Info.Println(err)
		up.State = UpdateStateFailure
		up.Err = err
		return err
	}

	re, found, err := upt.DetectVersion("brutella/hkcam", up.Version)
	if err != nil {
		log.Info.Println("search version:", err)

		// update failed
		up.State = UpdateStateFailure
		up.Err = err
		return err
	} else if !found {
		err := fmt.Errorf("version %s not found", up.Version)
		log.Info.Println(err)

		// update failed
		up.State = UpdateStateFailure
		up.Err = err
		return err
	}

	err = upt.UpdateTo(re, cmdPath)
	if err != nil {
		log.Info.Println("install update:", err)

		// update failed
		up.State = UpdateStateFailure
		up.Err = err
		return err
	}

	// store version and url of latest version
	up.Version = re.Version.String()
	up.URL = re.URL

	// update succeeded
	up.State = UpdateStateSuccess

	return nil
}
