package http

import (
	"encoding/json"
	"net/http"
	"sync"

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
}

func NewHttpTransport() *HttpTransport {
	return &HttpTransport{
		mx:       sync.Mutex{},
		peers:    make(map[string]network.Peer),
		onKey:    nil,
		onSetKey: nil,
	}
}

func (t *HttpTransport) Listen(addr string) error {
	t.addr = addr

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
	return http.ListenAndServe(addr, r)
}

func (t *HttpTransport) GetPeer(addr string) (network.Peer, error) {
	t.mx.Lock()
	defer t.mx.Unlock()
	if peer, ok := t.peers[addr]; ok {
		return peer, nil
	}

	peer := NewHttpPeer(addr)
	t.peers[addr] = peer
	return peer, nil
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
