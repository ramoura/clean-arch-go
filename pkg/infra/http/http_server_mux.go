package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type ServerMux struct {
	Server *http.ServeMux
}

func NewHttpServerMux() *ServerMux {
	return &ServerMux{
		Server: http.NewServeMux(),
	}
}

func (httpServer ServerMux) HandlerFunc(method string, path string, handler func(Request) (error, *Response)) {
	httpServer.Server.HandleFunc(
		method+" "+path,
		handlerFunc(handler),
	)
}

type RequestMuxImpl struct {
	r *http.Request
}

func (r RequestMuxImpl) PathValue(param string) string {
	return r.r.PathValue(param)
}
func (r RequestMuxImpl) QueryValue(param string) string {
	return r.r.URL.Query().Get(param)
}

func (r RequestMuxImpl) Body() io.ReadCloser {
	return r.r.Body
}

func (httpServer ServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpServer.Server.ServeHTTP(w, r)
}

func handlerFunc(handler func(Request) (error, *Response)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request = RequestMuxImpl{
			r: r,
		}

		err, resp := handler(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if resp == nil {
			log.Println("Response is nil")
			return
		}

		if resp.Headers != nil {
			for key, value := range resp.Headers {
				w.Header().Set(key, value)
			}
		}

		w.WriteHeader(resp.Status)

		if resp.Body != nil {
			if err := json.NewEncoder(w).Encode(resp.Body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

func (httpServer ServerMux) Start() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Print("Server running on port", port)
	addr := fmt.Sprintf(":%s", port)
	srv := &http.Server{
		Addr:    addr,
		Handler: httpServer.Server,
	}
	log.Print("Starting server at ", addr)
	log.Fatal(srv.ListenAndServe())

}
