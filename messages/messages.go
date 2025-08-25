package messages

import (
	"textcat/auth"
	"textcat/models"
	"encoding/json"
	"fmt"
)

func HandleMSG(msg []byte) {
	var data models.WsIncome

	err := json.Unmarshal(msg, &data)
	if err != nil {
		models.App.Log.Error("[auth.go:8] Failed to Unmarshal json: %s", err)
	}
	
	switch data.Rtype {
		case "login":
			auth.UserLogin(data)

		case "register":
			fmt.Printf("create account for user")

		case "message":
			fmt.Printf("handle the message")
	}
}