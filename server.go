package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

type Channel struct {
	id      uuid.UUID
	clients []Client

	logger LoggerInterface
}

func NewChannel(logger LoggerInterface) *Channel {
	return &Channel{
		id:      uuid.New(),
		clients: nil,
		logger:  logger,
	}
}

func (c *Channel) Connect(client Client) {
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

	broadcast chan Message

	logger LoggerInterface

	register chan *Client
}

func NewChatServer(logger LoggerInterface) *ChatServer {
	return &ChatServer{
		channels:  make(channels),
		broadcast: make(chan Message),
		logger:    logger,
		register:  make(chan *Client),
	}
}

func (s *ChatServer) Run() {
	for {
		select {
		case message := <-s.broadcast:
			if message.ChannelId == nil {
				s.logger.Error("missing channel id")
			}

			switch message.Action {
			case "join":
				fmt.Println("chUuid", message.ChannelId)
			case "send":
				ch := s.ChannelById(*message.ChannelId)
				if ch == nil {
					s.logger.Error("channel does not exist")
					break
				}

				for _, client := range ch.clients {
					if err := client.conn.WriteMessage(1, []byte(message.Content)); err != nil {
						log.Println(err)

						return
					}
				}
			}
		}
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
