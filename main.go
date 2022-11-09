package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(chatServer *ChatServer, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)

		return
	}

	// create new chat client
	userName := fmt.Sprintf("randomUser%d", rand.Intn(100))
	client := NewClient(chatServer, conn, userName)

	channelId, _ := uuid.Parse("5e306e54-05b8-498d-8678-cba88822b42d")
	channel := chatServer.ChannelById(channelId)

	if channel == nil {
		channel = chatServer.NewChannel()
		chatServer.Logger().Debug("connected to new channel")
	}

	// connect user to channel
	channel.Connect(*client)

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		chatServer.Logger().Debug(fmt.Sprintf("msg: %s", p))

		var message Message
		if err := json.Unmarshal(p, &message); err != nil {
			log.Println("error", err)
		}

		if message.ChannelId == nil {
			chatServer.logger.Error("missing channel id")
		}

		chatServer.broadcast <- message
	}
}

func main() {
	flag.Parse()

	customLogger := NewCustomLogger()
	chatServer := NewChatServer(customLogger)

	go chatServer.Run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(chatServer, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
