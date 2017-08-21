package proxy

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestProxy_AddPattern(t *testing.T) {
	p := NewProxy("new", "old", false)
	p.AddPattern("mypat1")
	assert.Equal(t, []string{"mypat1"}, p.patterns)
}

func TestProxy_AddPatterns(t *testing.T) {
	p := NewProxy("new", "old", false)
	p.AddPatterns([]string{"mypat1","mypat2"})
	assert.Equal(t, []string{"mypat1","mypat2"}, p.patterns)
}
