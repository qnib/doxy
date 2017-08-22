package proxy

type ProxyOptions struct {
	DockerSocket 	string
	ProxySocket 	string
	Debug 			bool
	Patterns 		[]string
}

var defaultProxyOptions = ProxyOptions{
	DockerSocket: DOCKER_SOCKET,
	ProxySocket: PROXY_SOCKET,
	Debug: false,
	Patterns: []string{},
}

type ProxyOption func(*ProxyOptions)

func WithDockerSocket(s string) ProxyOption {
	return func(o *ProxyOptions) {
		o.DockerSocket = s
	}
}

func WithProxySocket(s string) ProxyOption {
	return func(o *ProxyOptions) {
		o.ProxySocket = s
	}
}

func WithDebugValue(d bool) ProxyOption {
	return func(o *ProxyOptions) {
		o.Debug = d
	}
}

func WithDebugEnabled() ProxyOption {
	return func(o *ProxyOptions) {
		o.Debug = true
	}
}

func WithPattern(p string) ProxyOption {
	return func(o *ProxyOptions) {
		o.Patterns = append(o.Patterns, p)
	}
}

func WithPatterns(p []string) ProxyOption {
	return func(o *ProxyOptions) {
		o.Patterns = p
	}
}

