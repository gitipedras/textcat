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
	"textcat/core"
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
          removedToken := auth.SessionManager.RemoveByConn(conn)
		  core.Channels.RemoveTokenFromAllChannels(removedToken)
          mutex.Unlock()
          break
       }
       messages.HandleMSG(conn, message)
    }
}


func main() {
	loaderr := models.LoadConfig("config.json")
   if loaderr != nil {
      fmt.Println(loaderr)
      fmt.Println("Please create a config.json")
   }

   models.App.Log.Info("Server Details", slog.String("ServerName", models.Config.ServerName), slog.String("ServerDesc", models.Config.ServerDesc))

   database.DbInit()
   core.ChannelsInit()

	http.HandleFunc("/ws", wsHandler)
	var port string = models.Config.Port


	models.App.Log.Info("[INIT] Starting network server...", slog.String("port", port))

	// if you put this goroutine after it will never be executed
	go auth.SessionTimer()

	err := http.ListenAndServe(port, nil)

	if err != nil {
		//models.App.Log.Error("[ERROR] Error starting server:", slog.Any("err", err))
		panic(err)
	}

	models.App.Log.Info("Stopping...")
	defer database.DB.Close()
}
