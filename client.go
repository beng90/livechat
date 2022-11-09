package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	id     uuid.UUID
	server ChatServerInterface

	conn *websocket.Conn

	userName string
}

func NewClient(server ChatServerInterface, conn *websocket.Conn, userName string) *Client {
	return &Client{
		id:       uuid.New(),
		server:   server,
		conn:     conn,
		userName: userName,
	}
}
