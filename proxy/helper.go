package proxy

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"io"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"fmt"
)

type Headers map[string][]string



type configWrapper struct {
	*container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
}

func createDevMapping(dev string) (dm container.DeviceMapping, err error) {
	spl := strings.Split(dev, ":")
	switch len(spl) {
	case 2:
		dm = container.DeviceMapping{
			PathOnHost: spl[0],
			PathInContainer: spl[1],
		}
	case 3:
		dm = container.DeviceMapping{
			PathOnHost: spl[0],
			PathInContainer: spl[1],
			CgroupPermissions: spl[2],
		}
	default:
		return dm, fmt.Errorf("string needs to specify <src>:<dst>[:permission]")
		}
 	return
}

func encodeBody(obj interface{}, header http.Header) (io.Reader, http.Header, error) {
	if obj == nil {
		return nil, header, nil
	}

	body, err := encodeData(obj)
	if err != nil {
		return nil, header, err
	}
	return body, header, nil
}

func encodeData(data interface{}) (*bytes.Buffer, error) {
	params := bytes.NewBuffer(nil)
	if data != nil {
		if err := json.NewEncoder(params).Encode(data); err != nil {
			return nil, err
		}
	}
	return params, nil
}