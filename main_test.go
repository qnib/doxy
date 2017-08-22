package main

import (
	"testing"
	"github.com/zpatrick/go-config"
	"github.com/qnib/doxy/proxy"
	"github.com/stretchr/testify/assert"
)

func TestEvalOptions(t *testing.T) {
	myCfg := map[string]string{
		"proxy-socket": "proxy.sock",
		"docker-socket": "docker.sock",
		"debug": "true",
	}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(myCfg)})
	got := EvalOptions(cfg)
	p := proxy.NewProxy(got...)
	exp := map[string]interface{}{
		"docker-socket": "docker.sock",
		"proxy-socket": "proxy.sock",
		"debug": true,
		"patterns": []string{},
	}
	assert.Equal(t, exp, p.GetOptions())
}

func TestEvalPatternOpts(t *testing.T) {
	cfg := config.NewConfig([]config.Provider{config.NewStatic(map[string]string{
		"docker-socket": "docker.sock",
		"proxy-socket": "proxy.sock",
		"pattern-file": "some.file",
	})})
	got := EvalPatternOpts(cfg)
	p := proxy.NewProxy(got)
	exp := map[string]interface{}{
		"docker-socket": "/var/run/docker.sock",
		"proxy-socket": "/tmp/doxy.sock",
		"debug": false,
		"patterns": proxy.DEFAULT_PATTERNS,
	}
	opts := p.GetOptions()
	assert.Equal(t, exp, opts)
}