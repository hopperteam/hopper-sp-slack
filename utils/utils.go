package utils

import (
	"encoding/json"
	"net/http"
)

func ParseJSONReq(r *http.Request, ref interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(ref)
}

func ParseJSONRes(r *http.Response, ref interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(ref)
}
