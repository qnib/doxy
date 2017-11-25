package proxy


type ProxyOptions struct {
	DockerSocket,ProxySocket,PinUser 	string
	Debug,Gpu,PinUserEnabled 			bool
	Patterns,BindMounts,DevMappings 	[]string
}

var defaultProxyOptions = ProxyOptions{
	DockerSocket: DOCKER_SOCKET,
	ProxySocket: PROXY_SOCKET,
	PinUser: "",
	PinUserEnabled: false,
	Debug: false,
	Gpu: false,
	Patterns: []string{},
	BindMounts: []string{},
	DevMappings: []string{},
}

type ProxyOption func(*ProxyOptions)

func WithPinUser(pub bool, pu string) ProxyOption {
	return func(o *ProxyOptions) {
		o.PinUser = pu
		o.PinUserEnabled = pub
	}
}
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

func WithGpuValue(b bool) ProxyOption {
	return func(o *ProxyOptions) {
		o.Gpu = b
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

func WithBindMounts(bm []string) ProxyOption {
	return func(o *ProxyOptions) {
		o.BindMounts = bm
	}
}

func WithDevMappings(dm []string) ProxyOption {
	return func(o *ProxyOptions) {
		o.DevMappings = dm
	}
}
