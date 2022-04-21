package apiutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func JSONEncode(v interface{}) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	err := enc.Encode(v)

	return buf, err
}

func JSONDecode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func WriteJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	buf, err := JSONEncode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(buf.Bytes())
	return err
}

func ReadJSON(rc io.Reader, v interface{}) error {
	if err := JSONDecode(rc, v); err != nil {
		return err
	}

	return nil
}
