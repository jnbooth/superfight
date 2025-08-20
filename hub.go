package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Hub struct {
	broadcast  chan []byte
	clients    map[*Client]bool
	register   chan *Client
	shutdown   chan struct{}
	unregister chan *Client
}

func NewHub(cards *Cards) *Hub {
	return &Hub{
		broadcast:  make(chan []byte, 10),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client, 10),
		shutdown:   make(chan struct{}),
		unregister: make(chan *Client, 10),
	}
}

func (h *Hub) SendEvent(event string, data any) {
	h.broadcast <- formatEvent(event, data)
}

func (h *Hub) Shutdown() {
	h.shutdown <- struct{}{}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.messages)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				client.messages <- message
				select {
				case client.messages <- message:
				default:
					close(client.messages)
					delete(h.clients, client)
				}
			}
		case <-h.shutdown:
			message := formatEvent("shutdown", time.Now())
			for client := range h.clients {
				client.messages <- message
				close(client.messages)
			}
			clear(h.clients)
			return
		}
	}
}

func formatEvent(event string, data any) []byte {
	switch data.(type) {
	case string:
		return fmt.Appendf(nil, "event: %s\ndata: %s\n\n", event, data)
	default:
		dataJson, _ := json.Marshal(data)
		return fmt.Appendf(nil, "event: %s\ndata: %s\n\n", event, dataJson)
	}
}
