package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
)

type ServerGin struct {
	Server *gin.Engine
}

func NewHttpServerGin() *ServerGin {
	gin.SetMode(gin.ReleaseMode)
	return &ServerGin{
		Server: gin.Default(),
	}
}

type RequestGinImpl struct {
	c *gin.Context
}

func (r RequestGinImpl) PathValue(param string) string {
	return r.c.Param(param)
}

func (r RequestGinImpl) QueryValue(param string) string {
	return r.c.Query(param)
}

func (r RequestGinImpl) Body() io.ReadCloser {
	return r.c.Request.Body
}

var re = regexp.MustCompile(`\{(\w+)}`)

func (httpServer ServerGin) HandlerFunc(method string, path string, handler func(Request) (error, *Response)) {
	path = re.ReplaceAllString(path, `:$1`)
	httpServer.Server.Handle(method,
		path,
		handlerFuncGin(handler),
	)
}

func (httpServer ServerGin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpServer.Server.ServeHTTP(w, r)
}

func handlerFuncGin(handler func(param Request) (error, *Response)) gin.HandlerFunc {
	return func(context *gin.Context) {
		var req Request = RequestGinImpl{
			c: context,
		}

		err, resp := handler(req)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(context.Writer, err.Error(), http.StatusBadRequest)
			return
		}
		for key, value := range resp.Headers {
			context.Writer.Header().Set(key, value)
		}
		context.JSON(resp.Status, resp.Body)
	}
}

func (httpServer ServerGin) Start() {

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
