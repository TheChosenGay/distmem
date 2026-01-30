package http

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHttpTransport(t *testing.T) {
	tr := NewHttpTransport()
	data := make(map[string]any)
	tr.OnKey(func(key string) (any, error) {
		value, ok := data[key]
		if !ok {
			return nil, errors.New("key not found")
		}
		return value, nil
	})

	tr.OnSetKey(func(key string, value any) error {
		data[key] = value
		return nil
	})
	go func() {
		tr.Listen("127.0.0.1:8088")
	}()

	time.Sleep(100 * time.Millisecond)

	peer := NewHttpPeer("127.0.0.1:8088")
	peer2 := NewHttpPeer("127.0.0.1:8088")
	if err := peer.Set("key", "value"); err != nil {
		t.Errorf("expected nil, got error: %v", err)
	}
	if err := peer2.Set("key2", "value2"); err != nil {
		t.Errorf("expected nil, got error: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	value, err := peer2.Get("key")
	value2, err := peer.Get("key2")
	if err != nil {
		t.Errorf("expected nil, got error: %v", err)
	}
	if value != "value" {
		t.Errorf("expected value, got %v", value)
	}
	if value2 != "value2" {
		t.Errorf("expected value2, got %v", value2)
	}

	fmt.Printf("key: %s, value: %v\n", "key", value)
	fmt.Printf("key: %s, value: %v\n", "key2", value2)

}

func TestHttpTransportConnect(t *testing.T) {
	tr := NewHttpTransport()
	go tr.Listen("127.0.0.1:8089")

	tr1 := NewHttpTransport()
	go tr1.Listen("127.0.0.1:8090")

	tr2 := NewHttpTransport()
	go tr2.Listen("127.0.0.1:8091")

	time.Sleep(100 * time.Millisecond)

	tr.Connect("127.0.0.1:8090")
	tr.Connect("127.0.0.1:8091")
	time.Sleep(100 * time.Millisecond)
	peer, err := tr.GetPeer("127.0.0.1:8090")
	assert.NotNil(t, peer)
	assert.Nil(t, err)
	peer, err = tr.GetPeer("127.0.0.1:8089")
	assert.NotNil(t, peer)
	assert.Nil(t, err)

	peer, err = tr1.GetPeer("127.0.0.1:8091")
	assert.NotNil(t, peer)
	assert.Nil(t, err)

}
