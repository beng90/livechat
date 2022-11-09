package main

import (
	"fmt"

	"github.com/google/uuid"
)

type ChatServerInterface interface {
	Connect(channel Channel, client Client)
}

type Channel struct {
	id      uuid.UUID
	clients []Client

	logger LoggerInterface
}

func NewChannel(logger LoggerInterface) Channel {
	return Channel{
		id:      uuid.New(),
		clients: nil,
		logger:  logger,
	}
}

func (c *Channel) Connect(client Client) {
	c.clients = append(c.clients, client)

	c.logger.Debug(fmt.Sprintf("New client connected: %s", client.userName))
}

type ChatServer struct {
	channels []Channel

	broadcast chan []byte

	logger LoggerInterface
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		channels:  []Channel{},
		broadcast: make(chan []byte),
	}
}

func (s *ChatServer) Connect(channel Channel, client Client) {
	channel.Connect(client)
}
