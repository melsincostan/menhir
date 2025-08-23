package menhir

import "net/http"

type Request struct {
	Request *http.Request
	done    bool
}

func (r *Request) Done() {
	r.done = true
}

func (r *Request) ShouldContinue() bool {
	return r.done
}
