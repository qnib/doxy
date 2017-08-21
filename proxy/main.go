package proxy

import (
	"log"
	"net/http"
	"os"
	"syscall"
	"os/signal"
	"github.com/urfave/negroni"

)

type Proxy struct {
	dockerSocket, newSocket string
	debug 					bool
	patterns 				[]string
}

func NewProxy(newsock, oldsock string, debug bool) Proxy {
	return Proxy{
		dockerSocket: oldsock,
		newSocket: newsock,
		debug: debug,
		patterns: []string{},
	}
}

func (p *Proxy) AddPatterns(patterns []string) {
	p.patterns = append(p.patterns,patterns...)
}

func (p *Proxy) AddPattern(pattern string) {
	p.patterns = append(p.patterns,pattern)
}

func (p *Proxy) Run() {
	upstream := NewUpstream(p.dockerSocket, p.patterns)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	l, err := ListenToNewSock(p.newSocket, sigc)
	if err != nil {
		panic(err)
	}
	n := negroni.New()
	if p.debug {
		n.Use(negroni.NewLogger())
	}
	n.UseHandler(upstream)
	log.Printf("Serving proxy on '%s'", p.newSocket)
	if err = http.Serve(l, n); err != nil {
		panic(err)
	}
}

