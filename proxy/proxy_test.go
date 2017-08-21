package proxy

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewUnixSocket(t *testing.T) {
	us := NewUnixSocket("test")
	exp := UnixSocket{"test"}
	assert.Equal(t, exp, us)
}