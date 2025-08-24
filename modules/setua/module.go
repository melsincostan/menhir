package setua

import (
	"flag"
	"net/http/httputil"
)

type SetUA struct {
	userAgent *string
}

func New() *SetUA {
	return &SetUA{
		userAgent: flag.String("setua.user-agent", "menhir/setua", "user agent to use for the backend"),
	}
}

func (s *SetUA) Name() string {
	return "setua"
}

func (s *SetUA) Init() (err error) {
	return nil
}

func (s *SetUA) Priority() *int {
	return nil
}

func (s *SetUA) Default() bool {
	return false
}

func (s *SetUA) Rewrite(req *httputil.ProxyRequest) {
	req.Out.Header.Set("User-Agent", *s.userAgent)
}
