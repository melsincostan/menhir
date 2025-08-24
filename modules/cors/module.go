package cors

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/melsincostan/menhir/menhir"
)

type Cors struct {
	origin            *string
	originFunc        func(r *http.Request) string
	allowedHeaders    *string
	exposedHeaders    *string
	handleOptionsFlag *bool
	handleOptions     bool
	preflightMaxAge   *time.Duration
	allowCredentials  *bool
	allowedMethods    *string
}

func New() *Cors {
	return &Cors{
		origin:            flag.String("cors.origin", "*", "allowed origin(s). if the value is 'ALL', the value of the origin header on the request will be use"),
		allowedHeaders:    flag.String("cors.allowed-headers", "", "a comma-separated list of headers to be allowed from the browser (Access-Control-Allow-Headers)"),
		exposedHeaders:    flag.String("cors.exposed-headers", "", "a comma-separated list of headers to be exposed to the browser (Access-Control-Expose-Headers)"),
		handleOptionsFlag: flag.Bool("cors.intercept-options", true, "return early when processing OPTIONS requests"),
		preflightMaxAge:   flag.Duration("cors.preflight-max-age", 0*time.Second, "duration a preflight should be cached for. Set to 0s to omit. (Access-Control-Max-Age)"),
		allowCredentials:  flag.Bool("cors.allow-credentials", false, "allow credentials (Access-Control-Allow-Credentials)"),
		allowedMethods:    flag.String("cors.allowed-methods", "", "allowed method(s) as a comma-separated list (Access-Control-Allow-Methods)"),
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
		c.headers(rw.Header(), req.Request)
		rw.WriteHeader(http.StatusNoContent)
	}
}

func (c *Cors) ModifyResponse(res *http.Response) (err error) {
	c.headers(res.Header, res.Request)
	return
}

func (c *Cors) headers(h http.Header, r *http.Request) {
	if *c.allowedHeaders != "" {
		h.Set("Access-Control-Allowed-Headers", *c.allowedHeaders)
	}

	if *c.exposedHeaders != "" {
		h.Set("Access-Control-Expose-Headers", *c.exposedHeaders)
	}

	if *c.preflightMaxAge != 0*time.Second {
		h.Set("Access-Control-Max-Age", fmt.Sprintf("%.0f", c.preflightMaxAge.Seconds()))
	}

	if *c.allowCredentials {
		h.Set("Access-Control-Allow-Credentials", "true")
	}

	if *c.allowedMethods != "" {
		h.Set("Access-Control-Allow-Methods", *c.allowedMethods)
	}

	h.Set("Access-Control-Allow-Origin", c.originFunc(r))
}
