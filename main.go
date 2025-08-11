package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	//"strings"
	"database/sql"   // Package for SQL database interactions
	"log"            // Package for logging
	_ "github.com/mattn/go-sqlite3"
	// "encoding/json"
	// "log/slog"
	"textcat/messageHandler"
)

// Upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var DB *sql.DB
var err error

func dbInit() {
	DB, err = sql.Open("sqlite3", "./appdata.db")

	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL,
	sessiontoken TEXT NOT NULL
	);`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

}

func isValidLogin(DB *sql.DB, username, sessionToken string) (bool, error) {
	var count int
	err := DB.QueryRow(`
	SELECT COUNT(*) FROM users
	WHERE username = ? AND sessiontoken = ?`, username, sessionToken).Scan(&count)

	if err != nil {
		return false, err
	}
	return count > 0, nil
}


var conn *websocket.Conn

func wsHandler(w http.ResponseWriter, r *http.Request) {
	//var conn *websocket.Conn
	var err error
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()
	// Listen for incoming messages
	for {
		// Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		messageString := string(message[:])
		handler.ProcessMsg(DB, messageString, conn)

		// Echo the message back to the client
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}

func main() {
	dbInit()

	http.HandleFunc("/ws", wsHandler)
	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

	defer DB.Close()
}
