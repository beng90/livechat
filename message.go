package main

import "github.com/google/uuid"

type Message struct {
	Client    *Client
	Action    string     `json:"action"`
	ChannelId *uuid.UUID `json:"channelId"`
	Content   string     `json:"content"`
}
