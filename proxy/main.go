package proxy

import (
	"log"
	"net/http"
	"os"
	"syscall"
	"os/signal"
	"github.com/urfave/negroni"
)

const (
	DOCKER_SOCKET = "/var/run/docker.sock"
	PROXY_SOCKET = "/tmp/doxy.sock"
	PATTERN_FILE = "/etc/doxy.pattern"
)

var (
	DEF_PAT = []string{
		`^/(v\d\.\d+/)?containers(/\w+)?/(json|stats|top)$`,
		`^/(v\d\.\d+/)?services(/[0-9a-f]+)?$`,
		`^/(v\d\.\d+/)?tasks(/\w+)?$`,
		`^/(v\d\.\d+/)?networks(/\w+)?$`,
		`^/(v\d\.\d+/)?volumes(/\w+)?$`,
		`^/(v\d\.\d+/)?nodes(/\w+)?$`,
		`^/(v\d\.\d+/)?info$`,
		`^/(v\d\.\d+/)?version$`,
		"^/_ping$",
	}
	HPC_PAT = []string{
		`^/(v\d\.\d+/)?containers(/\w+)?/(json|stats|top|create|start|run|kill)$`,
		`^/(v\d\.\d+/)?services(/[0-9a-f]+)?$`,
		`^/(v\d\.\d+/)?tasks(/\w+)?$`,
		`^/(v\d\.\d+/)?networks(/\w+)?$`,
		`^/(v\d\.\d+/)?volumes(/\w+)?$`,
		`^/(v\d\.\d+/)?nodes(/\w+)?$`,
		`^/(v\d\.\d+/)?info$`,
		`^/(v\d\.\d+/)?version$`,
		"^/_ping$",
	}
	PATTERNS = map[string][]string{
		"default": DEF_PAT,
		"hpc": HPC_PAT,
	}
)

type Proxy struct {
	dockerSocket, newSocket string
	debug, gpu 				bool
	patterns 				[]string
	bindMounts,devMappings	[]string
}

func NewProxy(opts ...ProxyOption) Proxy {
	options := defaultProxyOptions
	for _, o := range opts {
		o(&options)
	}
	return Proxy{
		dockerSocket: options.DockerSocket,
		newSocket: options.ProxySocket,
		debug: options.Debug,
		gpu: options.Gpu,
		patterns: options.Patterns,
		bindMounts: options.BindMounts,
		devMappings: options.DevMappings,
	}
}

func (p *Proxy) GetOptions() map[string]interface{} {
	opt := map[string]interface{}{
		"docker-socket": p.dockerSocket,
		"proxy-socket": p.newSocket,
		"debug": p.debug,
		"patterns": p.patterns,
	}
	return opt
}

func (p *Proxy) Run() {
	upstream := NewUpstream(p.dockerSocket, p.patterns, p.bindMounts, p.devMappings, p.gpu)
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

