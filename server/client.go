package main

import (
	"net/http"
)

type Client struct {
	messages chan []byte
}

func NewClient() *Client {
	return &Client{
		messages: make(chan []byte, 5),
	}
}

func (c *Client) Run(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rc := http.NewResponseController(w)

	for {
		select {
		case <-ctx.Done():
			return
		case event := <-c.messages:
			if event == nil {
				return
			}
			w.Write(event)
			if rc.Flush() != nil {
				return
			}
		}
	}
}
