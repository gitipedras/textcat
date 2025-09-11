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
			sendOk := channels.Channels.SendMessage(data.ChannelID, data.Message, data.Username, data.SessionToken, conn)
			if sendOk == false {
				// error occurred while trying to send
				response := models.WsSend{
		            Rtype:  "isr",
		            Status: "sendMessage",
		        }
		        data, err := json.Marshal(response)
		        if err != nil {
		            models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
		            return
		        }
		        conn.WriteMessage(websocket.TextMessage, data)
			}
		/* channels */
		case "connect":
			channels.Channels.AddUser(data.ChannelID, data.SessionToken, data.Username, conn)

		case "disconnect":
			channels.Channels.RemoveUser(data.ChannelID, data.SessionToken, data.Username)
	}
}