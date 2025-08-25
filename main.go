package main

import (
	"fmt"
	"sync"

	/* websocket */
	"github.com/gorilla/websocket"
	"net/http"

	/* logging */
	"log/slog"

	/* database */
	_ "github.com/mattn/go-sqlite3"
	
	//"database/sql"

	/* internal */
	"textcat/messages"	
	"textcat/models"
)


var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan []byte)            // Broadcast channel
var mutex = &sync.Mutex{}                    // Protect clients map

func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
       fmt.Println("Error upgrading:", err)
       return
    }
    defer conn.Close()

    mutex.Lock()
    clients[conn] = true
    mutex.Unlock()

    for {
       _, message, err := conn.ReadMessage()
       if err != nil {
          mutex.Lock()
          delete(clients, conn)
          mutex.Unlock()
          break
       }
       broadcast <- message
    }
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		message := <-broadcast

		messages.HandleMSG(message)
	}
}


func main() {
	http.HandleFunc("/ws", wsHandler)

	var port string = ":8080"
	models.App.Log.Info("starting network server...", slog.String("port", port))

	// if you put this goroutine after it will never be executed
	go handleMessages()

	err := http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
