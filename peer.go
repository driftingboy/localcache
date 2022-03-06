package localcache

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// var globalPicker = func()

// pick a node Client
type PeerPicker interface {
	Pick(key string) (rl RemoteLoader, ok bool, isSelf bool)
}

// func RegisterPeerPicker(func() PeerPicker)

// for load remote data with group and key
type RemoteLoader interface {
	Name() string
	Load(ctx context.Context, group, key string) ([]byte, error)
}

var _ RemoteLoader = (*httpLoader)(nil)

type httpLoader struct {
	name    string
	baseURL string
}

func (h *httpLoader) Name() string {
	return h.name
}

func (h *httpLoader) Load(ctx context.Context, group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"http://%v%v%v/%v",
		h.name,
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

// grpc loader
