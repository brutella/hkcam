package apiutil

import (
	"github.com/gorilla/schema"

	"net/http"
	"strconv"
)

// ParseInt64 converts a string to an 8-byte integer
func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

var decoder = schema.NewDecoder()

func DecodeForm(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	return decoder.Decode(v, r.Form)
}

func DecodeURLQuery(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	return decoder.Decode(v, r.URL.Query())
}

func ToBool(s string) bool {
	switch s {
	case "on":
		return true
	case "off":
		return false
	default:
		v, _ := strconv.ParseBool(s)
		return v
	}
}
