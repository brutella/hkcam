package app

import (
	"time"
)

type UpdateState int

const (
	UpdateStateDefault UpdateState = iota
	UpdateStateInstall
	UpdateStateSuccess
	UpdateStateCancelled
	UpdateStateFailure
)

type Update struct {
	State      UpdateState `json:"state"`
	Version    string      `json:"version"`
	PreRelease bool        `json:"pre"`
	URL        string      `json:"url"`
	CreatedAt  time.Time   `json:"created_at"`
	Err        error       `json:"error"`
}

func (u Update) Installing() bool {
	return u.State == UpdateStateInstall
}

func (u Update) Cancelled() bool {
	return u.State == UpdateStateCancelled
}
func (u Update) Failure() bool {
	return u.State == UpdateStateFailure
}

func (u Update) Success() bool {
	return u.State == UpdateStateSuccess
}
