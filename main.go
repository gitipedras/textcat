package main

import (
	"fmt"
	"sync"

	/* websocket */
	"github.com/gorilla/websocket"
	"net/http"

	/* logging */
	"log/slog"

	/* internal */
	"textcat/messages"	
	"textcat/models"
	"textcat/auth"
	"textcat/database"

)


var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan []byte)            // Broadcast channel
var mutex = &sync.Mutex{}                    // Protect clients map
//var Sessions = auth.NewSessionManager()

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
       messages.HandleMSG(conn, message)

       //broadcast <- message
    }
}

/*func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		message := <-broadcast

		messages.HandleMSG(message)
	}
}*/


func main() {
	loaderr := models.LoadConfig("config.json")
   if loaderr != nil {
      panic(loaderr)
   }

   models.App.Log.Info("Server Details", slog.String("ServerName", models.Config.ServerName), slog.String("ServerDesc", models.Config.ServerDesc))

	http.HandleFunc("/ws", wsHandler)
	var port string = models.Config.Port

	database.DbInit()

	models.App.Log.Info("starting network server...", slog.String("port", port))

	// if you put this goroutine after it will never be executed
	go auth.SessionTimer()

	err := http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}

	defer database.DB.Close()
}
