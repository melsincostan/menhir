package menhir

import (
	"net/http"
	"net/http/httputil"
)

type ModuleBase interface {
	Name() string
	Init() error
}

type Handler interface {
	ModuleBase
	ServeHTTP(rw http.ResponseWriter, req *Request)
}

type Rewriter interface {
	ModuleBase
	Rewrite(req *httputil.ProxyRequest)
	ModifyResponse(res *http.Response) error
}
