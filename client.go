package main

import "github.com/gorilla/websocket"

type Client struct {
	server ChatServerInterface

	ws *websocket.Conn

	userName string
}

func NewClient(server ChatServerInterface, ws *websocket.Conn, userName string) *Client {
	return &Client{
		server:   server,
		ws:       ws,
		userName: userName,
	}
}
