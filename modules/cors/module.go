package cors

import (
	"flag"
	"net/http"

	"github.com/melsincostan/menhir/menhir"
)

type Cors struct {
	origin            *string
	originFunc        func(r *http.Request) string
	headers           *string
	handleOptionsFlag *bool
	handleOptions     bool
}

func New() *Cors {
	return &Cors{
		origin:            flag.String("cors.origin", "*", "allowed origin(s). if the value is 'ALL', the value of the origin header on the request will be use"),
		headers:           flag.String("cors.allowedh", "", "a comma-separated lists of headers to be allowed"),
		handleOptionsFlag: flag.Bool("cors.handleo", true, "return early when processing OPTIONS requests"),
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
		if *c.headers != "" {
			rw.Header().Add("Access-Control-Allowed-Headers", *c.headers)
		}
		rw.Header().Add("Access-Control-Allow-Origin", c.originFunc(req.Request))
		rw.WriteHeader(http.StatusNoContent)
	}
}

func (c *Cors) ModifyResponse(res *http.Response) (err error) {
	if *c.headers != "" {
		res.Header.Add("Access-Control-Allowed-Headers", *c.headers)
	}
	res.Header.Add("Access-Control-Allow-Origin", c.originFunc(res.Request))
	return
}
