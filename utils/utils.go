package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"io/ioutil"
)

func ParseJSONReq(r *http.Request, ref interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(ref)
}

func ParseJSONRes(r *http.Response, ref interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(ref)
}

func MvReqBodyToStr(r *http.Request) (string, error) {
	body, err := ioutil.ReadAll(r.Body)
	return string(body), err
}

func CpReqBodyToStr(r *http.Request) (string, error) {
	body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return string(body), err
    }
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return string(body), err
}
