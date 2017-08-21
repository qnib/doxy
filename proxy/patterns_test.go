package proxy

import (
	"testing"
	"strings"
	"github.com/stretchr/testify/assert"
)

const (
	p1 = `# List and inspect containers
^/(v\d\.\d+/)?containers(/\w+)?/json$
# List and inspect services
^/(v\d\.\d+/)?services(/[0-9a-f]+)?$
# List and inspect tasks
^/(v\d\.\d+/)?tasks(/\w+)?$
# List and inspect networks
^/(v\d\.\d+/)?networks(/\w+)?$
# List and inspect nodes
^/(v\d\.\d+/)?nodes(/\w+)?$
# Show engine info
^/(v\d\.\d+/)?info$
# Healthcheck
^/_ping$`
)

func TestReadPatterns(t *testing.T) {
	r := strings.NewReader(p1)
	got,err := ReadPatterns(r)
	assert.NoError(t, err, "Should be parsed without problems")
	exp := []string{
		`^/(v\d\.\d+/)?containers(/\w+)?/json$`,
		`^/(v\d\.\d+/)?services(/[0-9a-f]+)?$`,
		`^/(v\d\.\d+/)?tasks(/\w+)?$`,
		`^/(v\d\.\d+/)?networks(/\w+)?$`,
		`^/(v\d\.\d+/)?nodes(/\w+)?$`,
		`^/(v\d\.\d+/)?info$`,
		"^/_ping$",
	}
	assert.Equal(t, exp, got)
}
