package main

import (
	"encoding/json"
	"fmt"
)

type Hub struct {
	messagePrefix []byte
	messageSuffix []byte
	broadcast     chan []byte
	clients       map[chan []byte]bool
	register      chan chan []byte
	unregister    chan chan []byte
}

func NewHub(cards *Cards) *Hub {
	return &Hub{
		messagePrefix: []byte("event: gameupdate\ndata:"),
		messageSuffix: []byte("\n\n"),
		broadcast:     make(chan []byte),
		clients:       make(map[chan []byte]bool),
		register:      make(chan chan []byte),
		unregister:    make(chan chan []byte),
	}
}

func (h *Hub) SendEvent(event string, data any) {
	switch data.(type) {
	case string:
		h.broadcast <- fmt.Appendf(nil, "event: %s\ndata:%s\n\n", event, data)
	default:
		dataJson, _ := json.Marshal(data)
		h.broadcast <- fmt.Appendf(nil, "event: %s\ndata:%s\n\n", event, dataJson)
	}
}

func (h *Hub) Register(client chan []byte) {
	h.register <- client
}

func (h *Hub) Unregister(client chan []byte) {
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
				close(client)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client <- message:
				default:
					close(client)
					delete(h.clients, client)
				}
			}
		}
	}
}
