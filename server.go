package main

import (
	"fmt"

	"github.com/google/uuid"
)

type Channel struct {
	id      uuid.UUID
	clients []Client

	logger LoggerInterface
}

func NewChannel(logger LoggerInterface) *Channel {
	hardUuid, _ := uuid.Parse("5e306e54-05b8-498d-8678-cba88822b42d")
	return &Channel{
		//id:      uuid.New(),
		id:      hardUuid,
		clients: nil,
		logger:  logger,
	}
}

func (c *Channel) Connect(client Client) {
	fmt.Println("current clients", c.clients)

	c.clients = append(c.clients, client)

	c.logger.Debug(fmt.Sprintf("New client connected: %s, channel: %s", client.userName, c.id))
}

type ChatServerInterface interface {
	Connect(channel Channel, client Client)
	ChannelById(id uuid.UUID) *Channel
	Logger() LoggerInterface
	NewChannel() *Channel
}

type channels map[uuid.UUID]*Channel

type ChatServer struct {
	channels channels

	broadcast chan []byte

	logger LoggerInterface

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewChatServer(logger LoggerInterface) *ChatServer {
	return &ChatServer{
		channels:   make(channels),
		broadcast:  make(chan []byte),
		logger:     logger,
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (s *ChatServer) Connect(channel Channel, client Client) {
	channel.Connect(client)
}

func (s *ChatServer) ChannelById(id uuid.UUID) *Channel {
	if ch, ok := s.channels[id]; ok {
		return ch
	}

	return nil
}
func (s *ChatServer) Logger() LoggerInterface {
	return s.logger
}

func (s *ChatServer) NewChannel() *Channel {
	ch := NewChannel(s.logger)
	s.channels[ch.id] = ch

	return ch
}
