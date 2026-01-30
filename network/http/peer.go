package http

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/TheChosenGay/distmem/network"
	"github.com/TheChosenGay/distmem/view"
)

type HttpPeer struct {
	url string
}

func NewHttpPeer(url string) network.Peer {
	return &HttpPeer{
		url: url,
	}
}

func (p *HttpPeer) Get(key string) (any, error) {
	postUrl := p.GetUrl(key)
	resp, err := http.Post(postUrl, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var view view.View
	if err := json.NewDecoder(resp.Body).Decode(&view); err != nil {
		return nil, err
	}
	return view.Value, nil
}

func (p *HttpPeer) Set(key string, value any) error {
	postUrl := p.SetUrl(key, value)
	view := view.View{
		Key:   key,
		Value: value,
	}

	jsonBytes, err := json.Marshal(view)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(jsonBytes)
	resp, err := http.Post(postUrl, "application/json", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (p *HttpPeer) Addr() string {
	return p.url
}

func (p *HttpPeer) Close() error {
	return nil
}

func (p *HttpPeer) Connect(addr string) error {
	url := "http://" + addr + "/cache/connect/" + p.url
	_, err := http.Post(url, "application/json", nil)
	if err != nil {
		return err
	}
	return nil
}

// 约定http请求路由为: http://<peer_url>/cache/get/<key>
func (p *HttpPeer) GetUrl(key string) string {
	return "http://" + p.url + "/cache/get/" + key
}

// 约定http请求路由为: http://<peer_url>/cache/set, 数据通过body传输
func (p *HttpPeer) SetUrl(key string, value any) string {
	return "http://" + p.url + "/cache/set"
}
