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
}

type Responder interface {
	ModuleBase
	ModifyResponse(res *http.Response) error
}
