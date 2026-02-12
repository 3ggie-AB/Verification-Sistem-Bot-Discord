package sse

import (
	"sync"
)

type Client struct {
	UserID *uint
	Stream chan []byte
}

type Hub struct {
	Clients map[*Client]bool
	Mutex   sync.Mutex
}

var NotificationHub = Hub{
	Clients: make(map[*Client]bool),
}

func (h *Hub) AddClient(c *Client) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	h.Clients[c] = true
}

func (h *Hub) RemoveClient(c *Client) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	delete(h.Clients, c)
	close(c.Stream)
}

func (h *Hub) Broadcast(userID *uint, data []byte) {

	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	for client := range h.Clients {

		if userID != nil {
			if client.UserID == nil || *client.UserID != *userID {
				continue
			}
		}

		select {
		case client.Stream <- data:
		default:
		}
	}
}