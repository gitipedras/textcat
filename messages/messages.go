package messages

import (
	"textcat/auth"
	"textcat/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

func HandleMSG(conn *websocket.Conn, msg []byte) {
	var data models.WsIncome

	err := json.Unmarshal(msg, &data)
	if err != nil {
		models.App.Log.Error("[auth.go:8] Failed to Unmarshal json: %s", err)
	}
	
	switch data.Rtype {
		case "login":
			auth.UserLogin(conn, data)

		case "register":
			auth.UserRegister(conn, data)

		case "message":
			fmt.Printf("handle the message")
	}
}