package handler

import (
	"net/http"
)

func RegisterRoutes() {
	http.HandleFunc("/events", HandleEvents)
	http.HandleFunc("/button", AddToSlack)
	http.HandleFunc("/redirect", Redirect)
	http.HandleFunc("/", Ping)
}
