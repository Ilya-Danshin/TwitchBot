package channel_interaction

import (
	"io"
	"net/http"
)

type ChatClient struct {
	cli *http.Client
}

type options struct {
	endpoint          string
	additionalHeaders map[string]string
}

var defaultHeaders = map[string]string{}

func NewClient() *ChatClient {
	return &ChatClient{
		cli: &http.Client{},
	}
}

func (cc *ChatClient) post(opt options, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", opt.endpoint, body)
	if err != nil {
		return nil, err
	}

	for key, value := range defaultHeaders {
		req.Header.Add(key, value)
	}

	for key, value := range opt.additionalHeaders {
		req.Header.Add(key, value)
	}

	resp, err := cc.cli.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (cc *ChatClient) get(opt options) (*http.Response, error) {
	req, err := http.NewRequest("GET", opt.endpoint, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range defaultHeaders {
		req.Header.Add(key, value)
	}

	for key, value := range opt.additionalHeaders {
		req.Header.Add(key, value)
	}

	resp, err := cc.cli.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
