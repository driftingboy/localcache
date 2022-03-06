package localcache

import (
	"fmt"
	"hash/crc32"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/driftingboy/localcache/consistenthash"
)

const (
	defaultBasePath = "/cache/"
	defaultReplicas = 50
)

var defaultHashF = func(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

var _ PeerPicker = (*HTTPPool)(nil)

// HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	// ip:port, e.g. "localhost:8000"
	endpoint string
	basePath string

	// protect peers and loaderMap
	mu        sync.RWMutex
	peers     *consistenthash.Map
	replicas  int
	hashFn    consistenthash.Hash
	loaderMap map[string]*httpLoader
}

// HTTPPoolOptions are the configurations of a HTTPPool.
type HTTPPoolOptions struct {
	// BasePath specifies the HTTP path that will serve groupcache requests.
	// If blank, it defaults to "/_groupcache/".
	BasePath string

	// Replicas specifies the number of key replicas on the consistent hash.
	// If blank, it defaults to 50.
	Replicas int

	// HashFn specifies the hash function of the consistent hash.
	// If blank, it defaults to crc32.ChecksumIEEE.
	HashFn consistenthash.Hash
}

// NewHTTPPoolOpts initializes an HTTP pool of peers with the given options.
// The returned *HTTPPool implements http.Handler and must be registered using http.Handle.
func NewHTTPPoolOpts(endpoint string, o *HTTPPoolOptions) *HTTPPool {
	p := &HTTPPool{
		basePath:  defaultBasePath,
		replicas:  defaultReplicas,
		endpoint:  endpoint,
		hashFn:    defaultHashF,
		loaderMap: make(map[string]*httpLoader),
	}
	if o != nil {
		if o.BasePath != "" {
			p.basePath = o.BasePath
		}
		if o.Replicas > 0 {
			p.replicas = o.Replicas
		}
		if o.HashFn != nil {
			p.hashFn = o.HashFn
		}
	}

	p.peers = consistenthash.New(p.replicas, p.hashFn)

	// RegisterPeerPicker(func() PeerPicker { return p })
	return p
}

// Set full-update, covering the previous peer
// Each peer value should be a valid base URL,
// for example "http://example.net:8000".
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.peers.Add(peers...)
	p.loaderMap = make(map[string]*httpLoader, len(peers))
	for _, peer := range peers {
		p.loaderMap[peer] = &httpLoader{name: peer, baseURL: p.basePath}
	}
}

func (p *HTTPPool) Pick(key string) (rl RemoteLoader, ok bool, isSelf bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	peer := p.peers.Get(key)
	p.Log("pick key:%s, peer:%s", key, peer)

	if p.endpoint == peer {
		return nil, true, true
	}

	rl, ok = p.loaderMap[peer]
	return
}

// Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.endpoint, fmt.Sprintf(format, v...))
}

// ServeHTTP handle all http requests
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Log("%s %s", r.Method, r.URL.Path)

	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		http.Error(w, "url format: /<basepath>/<dbName>/<key>", http.StatusBadRequest)
		return
	}

	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "url format: /<basepath>/<dbName>/<key>", http.StatusBadRequest)
		return
	}

	dbName, key := parts[0], parts[1]

	db, ok := GetDB(dbName)
	if !ok {
		http.Error(w, "no such db: "+dbName, http.StatusNotFound)
		return
	}

	view, err := db.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(view.Bytes())
}
