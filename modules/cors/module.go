package cors

import (
	"flag"
	"net/http"

	"github.com/melsincostan/menhir/menhir"
)

type Cors struct {
	origin            *string
	originFunc        func(r *http.Request) string
	allowedHeaders    *string
	exposedHeaders    *string
	handleOptionsFlag *bool
	handleOptions     bool
}

func New() *Cors {
	return &Cors{
		origin:            flag.String("cors.origin", "*", "allowed origin(s). if the value is 'ALL', the value of the origin header on the request will be use"),
		allowedHeaders:    flag.String("cors.allowed-headers", "", "a comma-separated list of headers to be allowed from the browser (Access-Control-Allow-Headers)"),
		exposedHeaders:    flag.String("cors.exposed-headers", "", "a comma-separated list of headers to be exposed to the browser (Access-Control-Expose-Headers)"),
		handleOptionsFlag: flag.Bool("cors.intercept-options", true, "return early when processing OPTIONS requests"),
	}
}

func (c *Cors) Name() string {
	return "cors"
}

func (c *Cors) Priority() *int {
	return nil
}

func (c *Cors) Default() bool {
	return false
}

func (c *Cors) Init() (err error) {
	if *c.origin == "ALL" {
		c.originFunc = func(r *http.Request) string {
			return r.Header.Get("Origin")
		}
	} else {
		c.originFunc = func(r *http.Request) string {
			return *c.origin
		}
	}
	c.handleOptions = *c.handleOptionsFlag
	return
}

func (c *Cors) ServeHTTP(rw http.ResponseWriter, req *menhir.Request) {
	if c.handleOptions && req.Request.Method == http.MethodOptions {
		req.Done()
		if *c.allowedHeaders != "" {
			rw.Header().Add("Access-Control-Allowed-Headers", *c.allowedHeaders)
		}

		if *c.exposedHeaders != "" {
			rw.Header().Add("Access-Control-Expose-Headers", *c.exposedHeaders)
		}
		rw.Header().Add("Access-Control-Allow-Origin", c.originFunc(req.Request))
		rw.WriteHeader(http.StatusNoContent)
	}
}

func (c *Cors) ModifyResponse(res *http.Response) (err error) {
	if *c.allowedHeaders != "" {
		res.Header.Add("Access-Control-Allowed-Headers", *c.allowedHeaders)
	}
	if *c.exposedHeaders != "" {
		res.Header.Add("Access-Control-Expose-Headers", *c.exposedHeaders)
	}
	res.Header.Add("Access-Control-Allow-Origin", c.originFunc(res.Request))
	return
}
