package messages

import (
	/* internal */
	"textcat/auth"
	"textcat/models"
	"textcat/channels"

	/* websocket*/
	"github.com/gorilla/websocket"

	/* data processing */
	"encoding/json"
)

func HandleMSG(conn *websocket.Conn, msg []byte) {
	var data models.WsIncome

	err := json.Unmarshal(msg, &data)
	if err != nil {
		models.App.Log.Error("[messages.go:21] Invalid json recieved from client: %s", err)
	}
	
	switch data.Rtype {
		case "login":
			auth.UserLogin(conn, data)

		case "register":
			auth.UserRegister(conn, data)

		case "message":
			channels.HandleMSG(data.Username, data.SessionToken, data.Message, data.ChannelID, conn)
		
		case "connect":
			channels.ConnectUser(data.SessionToken, data.ChannelID, conn)
	}
}