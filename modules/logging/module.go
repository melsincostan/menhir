package logging

import (
	"log"
	"net/http"

	"github.com/melsincostan/menhir/menhir"
)

type Logging struct{}

func New() *Logging {
	return &Logging{}
}

func (l *Logging) Name() string {
	return "logging"
}

func (l *Logging) Priority() *int {
	prio := 0
	return &prio
}

func (l *Logging) Init() (err error) {
	return
}

func (l *Logging) ServeHTTP(rw http.ResponseWriter, req *menhir.Request) {
	log.Printf("IN %s %s %s", req.Request.Method, req.Request.Host, req.Request.RequestURI)
}

func (l *Logging) ModifyResponse(res *http.Response) (err error) {
	log.Printf("OUT %s %s %s %d", res.Request.Method, res.Request.URL.Host, res.Request.RequestURI, res.StatusCode)
	return
}
