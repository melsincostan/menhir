package menhir

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	ErrDuplicateModule = errors.New("a module with this name was already registered")
	ErrModuleUnusable  = errors.New("provided module doesn't match any known interface and won't be used")
)

type Wrapper struct {
	destination *url.URL
	modules     map[string]ModuleBase
	handlers    []Handler
	rewriters   []Rewriter
	proxy       *httputil.ReverseProxy
}

func (w *Wrapper) Register(mod ModuleBase) (err error) {
	if _, ok := w.modules[mod.Name()]; ok {
		return fmt.Errorf("module %s: %w", mod.Name(), ErrDuplicateModule)
	}
	w.modules[mod.Name()] = mod
	return
}

func (w *Wrapper) Init() {
	w.proxy = &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(w.destination)
			for _, rewriter := range w.rewriters {
				rewriter.Rewrite(pr)
			}
		},
		ModifyResponse: func(r *http.Response) (err error) {
			for _, rewriter := range w.rewriters {
				if err := rewriter.ModifyResponse(r); err != nil {
					return err
				}
			}
			return
		},
	}
}

func (w *Wrapper) Enable(mod ModuleBase) (err error) {
	if _, ok := w.modules[mod.Name()]; !ok {
		if err := w.Register(mod); err != nil {
			return err
		}
	}

	if err := mod.Init(); err != nil {
		return fmt.Errorf("module %s could not be initialized: %w", mod.Name(), err)
	}

	used := false
	if hm, ok := mod.(Handler); ok {
		used = true
		w.handlers = append(w.handlers, hm)
	}

	if rm, ok := mod.(Rewriter); ok {
		used = true
		w.rewriters = append(w.rewriters, rm)
	}

	if !used {
		return fmt.Errorf("module %s: %w", mod.Name(), ErrModuleUnusable)
	}

	return
}

func (w *Wrapper) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	wrappedRequest := &Request{
		done:    false,
		Request: req,
	}

	for _, handler := range w.handlers {
		handler.ServeHTTP(rw, wrappedRequest)
		if wrappedRequest.done {
			return
		}
	}

	w.proxy.ServeHTTP(rw, wrappedRequest.Request)
}
