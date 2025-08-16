package main

import (
	"encoding/json"
	"fmt"
)

type Hub struct {
	broadcast  chan []byte
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

func NewHub(cards *Cards) *Hub {
	return &Hub{
		broadcast:  make(chan []byte, 10),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client, 10),
		unregister: make(chan *Client, 10),
	}
}

func (h *Hub) SendEvent(event string, data any) {
	switch data.(type) {
	case string:
		h.broadcast <- fmt.Appendf(nil, "event: %s\ndata: %s\n\n", event, data)
	default:
		dataJson, _ := json.Marshal(data)
		h.broadcast <- fmt.Appendf(nil, "event: %s\ndata: %s\n\n", event, dataJson)
	}
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
		}
	}
}
