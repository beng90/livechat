package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	server ChatServerInterface

	conn *websocket.Conn

	userName string
}

func NewClient(server ChatServerInterface, conn *websocket.Conn, userName string) *Client {
	return &Client{
		server:   server,
		conn:     conn,
		userName: userName,
	}
}
