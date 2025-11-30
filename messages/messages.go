package messages

import (
	/* internal */
	"textcat/auth"
	"textcat/models"
	"textcat/validator"
	"textcat/core"

	/* websocket*/
	"github.com/gorilla/websocket"

	/* data processing */
	"encoding/json"
	"log/slog"
	"time"
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
			wentOk := validator.Username(data.Username)
			if wentOk == false {
				response := models.WsSend{
		            Rtype:  "invalidInput",
		            Status: "username",
		        }
		        data, err := json.Marshal(response)
		        if err != nil {
		            models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
		            return
		        }
		        conn.WriteMessage(websocket.TextMessage, data)
		        return
			}
			auth.UserLogin(conn, data)

		case "register":
			auth.UserRegister(conn, data)
			wentOk := validator.Username(data.Username)
			if wentOk == false {
				response := models.WsSend{
		            Rtype:  "invalidInput",
		            Status: "username",
		        }
		        data, err := json.Marshal(response)
		        if err != nil {
		            models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
		            return
		        }
		        conn.WriteMessage(websocket.TextMessage, data)
		        return
			}

		/* messaging */
		case "message":
			wentOk := validator.Message(data.Message)
			if wentOk == false {
				response := models.WsSend{
		            Rtype:  "invalidInput",
		            Status: "message",
		        }
		        data, err := json.Marshal(response)
		        if err != nil {
		            models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
		            return
		        }
		        conn.WriteMessage(websocket.TextMessage, data)
		        return
			}

			sendOk := core.Channels.SendMessage(data.ChannelID, data.Message, data.Username, data.SessionToken, conn)
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
			wentOk := validator.Message(data.Username)
			if wentOk == false {
				response := models.WsSend{
		            Rtype:  "invalidInput",
		            Status: "username",
		        }
		        data, err := json.Marshal(response)
		        if err != nil {
		            models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
		            return
		        }
		        conn.WriteMessage(websocket.TextMessage, data)
		        return
			}
			core.Channels.AddUser(data.ChannelID, data.SessionToken, data.Username, conn)

		case "disconnect":
			wentOk := validator.Message(data.Username)
			if wentOk == false {
				response := models.WsSend{
		            Rtype:  "invalidInput",
		            Status: "username",
		        }
		        data, err := json.Marshal(response)
		        if err != nil {
		            models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
		            return
		        }
		        conn.WriteMessage(websocket.TextMessage, data)
		        return
			}
			core.Channels.RemoveUser(data.ChannelID, data.SessionToken, data.Username)

		case "channelsList":
			response := models.WsSend{
			    Rtype:      "channelList",
			    Status:     "ok",
			    ChannelList: core.Channels.BuildChannelList(),
			    Time:       time.Now(),
			}
			data, err := json.Marshal(response)
			if err != nil {
			    models.App.Log.Error("Failed to marshal channel list", slog.String("err", err.Error()))
			    return
			}
			conn.WriteMessage(websocket.TextMessage, data)
	}
}