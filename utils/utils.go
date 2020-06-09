package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"io/ioutil"
)

func ParseJSONReq(r *http.Request, ref interface{}) error {
	return json.NewDecoder(r.Body).Decode(ref)
}

func ParseJSONRes(r *http.Response, ref interface{}) error {
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

func SendEquemeral(text string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Text string `json:"text"`
	}{
		Text: text,
	})
}

func SendPlainText(text string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(text))
}