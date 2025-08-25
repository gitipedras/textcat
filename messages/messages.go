package messages

import (
	//"textcat/auth"
	"textcat/models"
)

func HandleMSG(message []byte) {
	var data models.WsIncome

	err := json.Unmarshal(msg, &data)
	if err != nil {
		models.App.Log.Error("[auth.go:8] Failed to Unmarshal json: %s", err)
	}
	
	switch msg.Rtype {
		case "login":
			fmt.Printf("log the user in")

		case "register":
			fmt.Printf("create account for user")

		case "message":
			fmt.Printf("handle the message")
	}
}