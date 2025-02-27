package http_server

import (
	"errors"
	"io"
	"net/http"
	"songlibtest/config"

	"github.com/gin-gonic/gin"
)

var ErrBadMethod = errors.New("method type unrecognized")

const HTTP_GET = 1
const HTTP_POST = 2
const HTTP_PATCH = 3
const HTTP_DELETE = 4

type HttpServer struct {
	Server *gin.Engine
	Port   string
}

func NewServer(Config *config.Config) *HttpServer {
	r := gin.Default()
	return &HttpServer{Server: r, Port: Config.ServicePort}
}

func (Server *HttpServer) Map(RequestType int, Path string, Method func(c *gin.Context)) error {
	switch RequestType {
	case HTTP_GET:
		Server.Server.GET(Path, Method)
	case HTTP_POST:
		Server.Server.POST(Path, Method)
	case HTTP_DELETE:
		Server.Server.DELETE(Path, Method)
	case HTTP_PATCH:
		Server.Server.PATCH(Path, Method)
	default:
		return ErrBadMethod
	}
	return nil
}

func (Server *HttpServer) Run() {
	Server.Server.Run(Server.Port)
}

func SimpleRequest(request string) (int, []byte, error) {
	resp, err := http.Get(request)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	if resp.StatusCode != 200 {
		return resp.StatusCode, nil, resp.Request.Context().Err()
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, body, nil
}
