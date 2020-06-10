package handler

import (
	"net/http"
)

func RegisterRoutes() {
	http.HandleFunc("/events", HandleEvents)
	http.HandleFunc("/reply", HandleReply)
	http.HandleFunc("/button", AddToSlack)
	http.HandleFunc("/redirect", Redirect)
	http.HandleFunc("/subscription", HandleCommand)
	http.HandleFunc("/callback", Callback)
	http.HandleFunc("/", Ping)
}
