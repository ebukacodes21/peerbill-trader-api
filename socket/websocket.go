package socket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn         *websocket.Conn
	Subscription string
}

type WebSocketManager struct {
	clients map[*websocket.Conn]*Client
	lock    sync.Mutex
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients: make(map[*websocket.Conn]*Client),
	}
}

func (m *WebSocketManager) AddClient(client *websocket.Conn, subscription string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.clients[client] = &Client{
		Conn:         client,
		Subscription: subscription,
	}
	log.Printf("client added %v with subscription: %s", client.RemoteAddr(), subscription)
}

func (m *WebSocketManager) RemoveClient(client *websocket.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.clients, client)
	log.Printf("client removed %v", client.RemoteAddr())
}

func (m *WebSocketManager) Broadcast(message interface{}, subscription string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, client := range m.clients {
		if client.Subscription == subscription {
			if err := client.Conn.WriteJSON(message); err != nil {
				log.Printf("Error sending message to client: %v", err)
				client.Conn.Close()
				delete(m.clients, client.Conn)
			}
		}
	}

	return nil
}
