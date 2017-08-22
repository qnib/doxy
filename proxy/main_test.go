package proxy

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestProxy_NewProxyWithPattern(t *testing.T) {
	p := NewProxy(
		WithProxySocket("new"),
		WithDockerSocket("old"),
		WithPattern("mypat1"),
	)
	assert.Equal(t, []string{"mypat1"}, p.patterns)
	assert.False(t, p.debug, "Debug should be false by default")
}

func TestProxy_NewProxyWithPatterns(t *testing.T) {
	p := NewProxy(
		WithProxySocket("new"),
		WithDockerSocket("old"),
		WithPatterns([]string{"mypat1","mypat2"}),
	)
	assert.Equal(t, []string{"mypat1","mypat2"}, p.patterns)
	assert.False(t, p.debug, "Debug should be false by default")

}

func TestProxy_NewProxyWithDebugEnabled(t *testing.T) {
	p := NewProxy(
		WithProxySocket("new"),
		WithDockerSocket("old"),
		WithDebugEnabled(),
	)
	assert.True(t, p.debug, "Debug should be set")
}
func TestProxy_NewProxyWithDebugValue(t *testing.T) {
	p := NewProxy(
		WithProxySocket("new"),
		WithDockerSocket("old"),
		WithDebugValue(true),
	)
	assert.True(t, p.debug, "Debug should be set")
}

func TestProxy_GetOptions(t *testing.T) {
	p := NewProxy(
		WithProxySocket("proxy.sock"),
		WithDockerSocket("docker.sock"),
		WithPattern("mypat1"),
	)
	exp := map[string]interface{}{
		"docker-socket": "docker.sock",
		"proxy-socket": "proxy.sock",
		"debug": false,
		"patterns": []string{"mypat1"},
	}
	assert.Equal(t, exp, p.GetOptions())
}