package proxy

import (
	"net"
	"fmt"
	"net/url"
	"net/http"
	"net/http/httputil"
	"regexp"
)

type UpStream struct {
	Name  string
	proxy http.Handler
	allowed []*regexp.Regexp
}

// UnixSocket just provides the path, so that I can test it
type UnixSocket struct {
	path string
}

// NewUnixSocket return a socket using the path
func NewUnixSocket(path string) UnixSocket {
	return UnixSocket{
		path: path,
	}
}

func (us *UnixSocket) connectSocket(proto, addr string) (net.Conn, error) {
	conn, err := net.Dial("unix", us.path)
	return conn, err
}

func newReverseProxy(dial func(network, addr string) (net.Conn, error)) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			param := ""
			if len(req.URL.RawQuery) > 0 {
				param = "?" + req.URL.RawQuery
			}
			u, _ := url.Parse("http://docker" + req.URL.Path + param)
			*req.URL = *u
		},
		Transport: &http.Transport{
			Dial: dial,
		},
	}
}

// NewUpstream returns a new socket (magic)
func NewUpstream(socket string, regs []string) *UpStream {
	us := NewUnixSocket(socket)
	a := []*regexp.Regexp{}
	for _, r := range regs {
		p, _ := regexp.Compile(r)
		a = append(a, p)
	}
	return &UpStream{
		Name:  socket,
		proxy: newReverseProxy(us.connectSocket),
		allowed: a,
	}
}

func (u *UpStream) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, fmt.Sprintf("Only GET requests are allowed, req.Method: %s", req.Method), 400)
		return
	}
	for _, a := range u.allowed {
		if a.MatchString(req.URL.Path) {
			u.proxy.ServeHTTP(w, req)
			return
		}
	}
	http.Error(w, fmt.Sprintf("'%s' is not allowed.", req.URL.Path), 403)
}
