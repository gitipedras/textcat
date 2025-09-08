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
	"log/slog"
)

func HandleMSG(conn *websocket.Conn, msg []byte) {
	var data models.WsIncome

	err := json.Unmarshal(msg, &data)
	if err != nil {
		models.App.Log.Error("[messages.go:21] Invalid json recieved from client:", slog.Any("error", err))
	}
	
	switch data.Rtype {
		/* authentication */
		case "login":
			auth.UserLogin(conn, data)

		case "register":
			auth.UserRegister(conn, data)

		/* messaging */
		case "message":
			channels.Channels.SendMessage(data.ChannelID, data.Message, data.Username, data.SessionToken, conn)
		
		/* channels */
		case "connect":
			channels.Channels.AddUser(data.ChannelID, data.SessionToken, data.Username, conn)

		case "disconnect":
			channels.Channels.RemoveUser(data.ChannelID, data.SessionToken, data.Username)
	}
}