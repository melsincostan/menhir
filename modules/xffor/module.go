package xffor

import (
	"flag"
	"net/http/httputil"
)

type XFFor struct {
	spoofAddr *string
}

func New() *XFFor {
	return &XFFor{
		spoofAddr: flag.String("xffor.for", "127.0.0.1", "Value to set in the X-Forwarded-For Header"),
	}
}

func (x *XFFor) Name() string {
	return "xffor"
}

func (x *XFFor) Init() (err error) {
	return
}

func (x *XFFor) Priority() (prio *int) {
	return
}

func (x *XFFor) Default() bool {
	return false
}

func (x *XFFor) Rewrite(req *httputil.ProxyRequest) {
	req.Out.Header.Set("X-Forwarded-For", *x.spoofAddr)
}
