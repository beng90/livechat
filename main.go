package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"

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

var conns = []*websocket.Conn{}

func serveWs(chatServer ChatServerInterface, w http.ResponseWriter, r *http.Request, logger LoggerInterface) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)

		return
	}

	// create new chat client
	userName := fmt.Sprintf("randomUser%d", rand.Intn(100))
	client := NewClient(chatServer, conn, userName)

	// create new channel
	channel := NewChannel(logger)

	// connect user to channel
	channel.Connect(*client)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		for _, client := range channel.clients {
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

	chatServer := NewChatServer()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(chatServer, w, r, customLogger)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
