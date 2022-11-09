package main

import (
	"log"

	"github.com/google/uuid"
)

type Channel struct {
	id      uuid.UUID
	clients map[uuid.UUID]Client

	logger LoggerInterface
}

func NewChannel(logger LoggerInterface) *Channel {
	return &Channel{
		id:      uuid.New(),
		clients: make(map[uuid.UUID]Client),
		logger:  logger,
	}
}

func (c *Channel) HasClient(userId uuid.UUID) bool {
	if _, ok := c.clients[userId]; ok {
		return true
	}

	return false
}
func (c *Channel) Connect(client Client) {
	c.clients[client.id] = client

	//c.logger.Debug(fmt.Sprintf("New client connected: %s, channel: %s", client.userName, c.id))
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

			ch := s.ChannelById(*message.ChannelId)
			if ch == nil {
				resp := "channel does not exist"
				s.logger.Error(resp)

				if err := message.Client.conn.WriteMessage(1, []byte(resp)); err != nil {
					log.Println(err)

					return
				}

				break
			}

			// check if user is member of this chat
			if !ch.HasClient(message.Client.id) {
				notification := "client is not a member of channel"
				s.logger.Error(notification)

				if err := message.Client.conn.WriteMessage(1, []byte(notification)); err != nil {
					log.Println(err)

					return
				}
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

func (s *ChatServer) Connect(channel Channel, client Client) {
	channel.Connect(client)
}

func (s *ChatServer) ChannelById(id uuid.UUID) *Channel {
	s.logger.Debug("current channels", s.channels)

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
