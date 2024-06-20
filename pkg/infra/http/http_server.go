package http

import (
	"io"
	"net/http"
)

type Request interface {
	PathValue(param string) string
	QueryValue(param string) string
	Body() io.ReadCloser
}

type Response struct {
	Headers map[string]string
	Body    interface{}
	Status  int
}

type Server interface {
	HandlerFunc(method string, path string, handler func(Request) (error, *Response))
	Start()

	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
