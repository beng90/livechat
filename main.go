package main

import (
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

func serveWs(chatServer ChatServerInterface, w http.ResponseWriter, r *http.Request) {
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
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		for _, client := range channel.clients {
			chatServer.Logger().Debug(fmt.Sprintf("msg: %s", p))
			if err := client.ws.WriteMessage(messageType, p); err != nil {
				log.Println(err)

				return
			}
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	customLogger := NewCustomLogger()
	chatServer := NewChatServer(customLogger)

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(chatServer, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
