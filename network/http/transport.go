package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/TheChosenGay/distmem/hash"
	"github.com/TheChosenGay/distmem/network"
	"github.com/TheChosenGay/distmem/view"
	"github.com/gorilla/mux"
)

type HttpTransport struct {
	addr string

	mx    sync.Mutex
	peers map[string]network.Peer

	onKey    network.OnKeyFunc
	onSetKey network.OnSetKeyFunc

	hashCircle *hash.HashCircle
}

func NewHttpTransport() *HttpTransport {
	tr := &HttpTransport{
		mx:    sync.Mutex{},
		peers: make(map[string]network.Peer),
		hashCircle: hash.NewHashCircle(hash.HashCircleOpts{
			Replicas: 10,
		}),
		onKey:    nil,
		onSetKey: nil,
	}

	return tr
}

func (t *HttpTransport) Listen(addr string) error {
	t.addr = addr
	t.hashCircle.Add(addr)

	t.mx.Lock()
	t.peers[addr] = NewHttpPeer(addr)
	t.mx.Unlock()

	r := mux.NewRouter()

	r.HandleFunc("/cache/get/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]
		if key == "" {
			http.Error(w, "key is required", http.StatusBadRequest)
			return
		}

		value, err := t.onKey(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		view := view.View{
			Key:   key,
			Value: value,
		}

		if err := json.NewEncoder(w).Encode(view); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	r.HandleFunc("/cache/set", func(w http.ResponseWriter, r *http.Request) {
		var view view.View
		if err := json.NewDecoder(r.Body).Decode(&view); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if view.Key == "" || view.Value == "" {
			http.Error(w, "key and value are required", http.StatusBadRequest)
			return
		}

		if err := t.onSetKey(view.Key, view.Value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// handshake protocol for discovery between peers.
	r.HandleFunc("/cache/connect/{addr}", func(w http.ResponseWriter, r *http.Request) {
		t.mx.Lock()
		defer t.mx.Unlock()
		addr := mux.Vars(r)["addr"]
		if addr == "" {
			http.Error(w, "addr is required", http.StatusBadRequest)
			return
		}

		// reture if exist
		if _, ok := t.peers[addr]; ok {
			fmt.Fprintf(w, "hello from %s, already connected to you!", t.addr)
			return
		}

		// broadcast to all peers
		go t.broadcast(addr)

		// add to peer
		peer := NewHttpPeer(addr)
		t.hashCircle.Add(addr)
		t.peers[addr] = peer
	})

	return http.ListenAndServe(addr, r)
}

func (t *HttpTransport) Connect(addr string) error {
	return t.broadcast(addr)
}

func (t *HttpTransport) GetPeer(addr string) (network.Peer, error) {
	t.mx.Lock()
	defer t.mx.Unlock()
	if peer, ok := t.peers[addr]; ok {
		return peer, nil
	}

	return nil, network.ErrPeerNotFound
}

func (t *HttpTransport) OnKey(keyFunc network.OnKeyFunc) {
	t.onKey = keyFunc
}

func (t *HttpTransport) OnSetKey(setKeyFunc network.OnSetKeyFunc) {
	t.onSetKey = setKeyFunc
}

func (t *HttpTransport) Close() error {
	return nil
}

func (t *HttpTransport) broadcast(addr string) error {
	t.mx.Lock()
	defer t.mx.Unlock()
	for _, peer := range t.peers {
		if peer.Addr() == addr {
			continue
		}
		go peer.Connect(addr)
	}
	return nil
}
