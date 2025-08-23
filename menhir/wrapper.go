package menhir

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
)

var (
	ErrDuplicateModule    = errors.New("a module with this name was already registered")
	ErrModuleUnusable     = errors.New("provided module doesn't match any known interface and won't be used")
	ErrUnregisteredModule = errors.New("no module was registered with this name")
)

type Wrapper struct {
	destination    *url.URL
	modules        map[string]ModuleBase
	handlers       []Handler
	rewriters      []Rewriter
	responders     []Responder
	proxy          *httputil.ReverseProxy
	requestCounter int64
}

func New() (w *Wrapper) {
	return &Wrapper{
		modules:   map[string]ModuleBase{},
		handlers:  []Handler{},
		rewriters: []Rewriter{},
	}
}

func (w *Wrapper) Modules() (mlist []ModuleBase) {
	mlist = []ModuleBase{}
	for _, mod := range w.modules {
		mlist = append(mlist, mod)
	}
	return
}

func (w *Wrapper) Register(mods ...ModuleBase) (err error) {
	for _, mod := range mods {
		if _, ok := w.modules[mod.Name()]; ok {
			return fmt.Errorf("module %s: %w", mod.Name(), ErrDuplicateModule)
		}
		w.modules[mod.Name()] = mod
	}
	return
}

func (w *Wrapper) Enable(mnames ...string) (err error) {
	for _, mname := range mnames {
		mod, ok := w.modules[mname]
		if !ok {
			return fmt.Errorf("module %s: %w", mname, ErrUnregisteredModule)
		}

		if err := mod.Init(); err != nil {
			return fmt.Errorf("module %s could not be initialized: %w", mod.Name(), err)
		}

		used := false
		if hm, ok := mod.(Handler); ok {
			used = true
			log.Printf("registering module %s as handler", mname)
			w.handlers = append(w.handlers, hm)
		}

		if rm, ok := mod.(Rewriter); ok {
			used = true
			log.Printf("registering module %s as rewriter", mname)
			w.rewriters = append(w.rewriters, rm)
		}

		if rm, ok := mod.(Responder); ok {
			used = true
			log.Printf("registering module %s as responder", mname)
			w.responders = append(w.responders, rm)
		}

		if !used {
			return fmt.Errorf("module %s: %w", mod.Name(), ErrModuleUnusable)
		}
	}
	return
}

func (w *Wrapper) Init(destination string) (err error) {
	// sort modules by priority
	// higher priority handlers should handle first, do modifications / rewrites last
	slices.SortFunc(w.handlers, sortFunc[Handler](true))      // higher priority first
	slices.SortFunc(w.rewriters, sortFunc[Rewriter](false))   // higher priority last
	slices.SortFunc(w.responders, sortFunc[Responder](false)) // higher priority last

	parsedDestination, err := url.Parse(destination)
	if err != nil {
		return fmt.Errorf("destination '%s' could not be parsed: %w", destination, err)
	}
	w.destination = parsedDestination
	w.proxy = &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(w.destination)
			for _, rewriter := range w.rewriters {
				rewriter.Rewrite(pr)
			}
		},
		ModifyResponse: func(r *http.Response) (err error) {
			for _, responder := range w.responders {
				if err := responder.ModifyResponse(r); err != nil {
					return err
				}
			}
			return
		},
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
			log.Printf("Request done early, not passing through the reverse proxy!")
			return
		}
	}

	w.proxy.ServeHTTP(rw, wrappedRequest.Request)
}
