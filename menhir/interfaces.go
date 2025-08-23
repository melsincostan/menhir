package menhir

import (
	"net/http"
	"net/http/httputil"
)

// ModuleBase represents the minimum interface of a module, before any functionality.
type ModuleBase interface {
	// Name should return the name of the module.
	// The name should not contain any space, or, ideally, any special characters, since it will be used to create command line flags.
	Name() string
	// Init is a function that will be called once command line arguments have been parsed for enabled modules, so they have a chance to act on configuration.
	Init() error
	// Priority is used to determine when a module should be run.
	// If it is nil, the module is considered as not having any special priority and will be lower priority than any module setting a priority at all.
	Priority() *int
	// Default is used to determine whether a module should be enabled by default.
	// If this returns yes, the module is enabled by default.
	Default() bool
}

// Handler is for modules receiving incoming requests and getting a chance to respond to them or read their values them before they hit the reverse proxy.
// They can also choose to bypass the reverse proxy by calling [menhir.Request.Done].
// Modules are called in descending orders of priority (highest priority Handler runs first).
type Handler interface {
	ModuleBase
	ServeHTTP(rw http.ResponseWriter, req *Request)
}

// Rewriter is for modules that modify a request in the reverse proxy before it is sent to the backend.
// Modules are called in ascending order of priority (highest priority Rewriter runs last).
type Rewriter interface {
	ModuleBase
	Rewrite(req *httputil.ProxyRequest)
}

// Responder is for modules that modify the response of the backend before it is passed to the client.
// Modules are called in ascending order of priority (highest priority Responder runs last).
type Responder interface {
	ModuleBase
	ModifyResponse(res *http.Response) error
}
