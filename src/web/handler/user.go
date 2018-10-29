package handler

import (
	"net/http"
)

func CreateUser(req *http.Request) (data interface{}, errorType int, message string) {
	return "hello world", 0, ""
}
