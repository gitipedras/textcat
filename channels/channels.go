package channels

import (
	/* data processing */
	"encoding/json"

	/* logging */
	"log/slog"
	"fmt"

	/* websocket ?? */
	"github.com/gorilla/websocket"

	/* internal */
	"textcat/auth"
	"textcat/models"
)

var ExistentChannels = map[string][]string{
	"main":      {},
	"minecraft": {},
	"minecraft-bedwars": {},
}


func HandleMSG(username string, token string, message string, channelID string, conn *websocket.Conn) {
	ok := auth.SessionManager.Exists(token)
	if ok {
		if clients, ko := ExistentChannels[channelID]; ko {

			for key, value := range ExistentChannels {
			    fmt.Println("Channel:", key)
			    for i, token := range value {
			        fmt.Printf("  [%d] %s\n", i, token)
			    }
			}

			for _, clientToken := range clients {
			    models.App.Log.Info("Sending message to client", slog.String("token", clientToken))
			    send := []byte(message)

			    err := auth.SessionManager.SendToClient(clientToken, send)
			    if err != nil {
			        models.App.Log.Error("SendToClient failed",
			            slog.String("token", clientToken),
			            slog.String("error", err.Error()))
			    }
			}
			
		} else {
			response := models.WsSend {
	            Rtype:   "invalidChannel",
	            Status:  token,
	        }
	        data, err := json.Marshal(response)
	        if err != nil {
	            models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
	            return
	        }
	        conn.WriteMessage(websocket.TextMessage, data)
		}
		
	} else {
		response := models.WsSend {
            Rtype:   "invalidSession",
            Status:  token,
        }
        data, err := json.Marshal(response)
        if err != nil {
            models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
            return
        }
        conn.WriteMessage(websocket.TextMessage, data)
	}
}

func ConnectUser(token string, channel string, conn *websocket.Conn) {
	models.App.Log.Info("user connected", slog.String("channel", channel), slog.String("token", token))
	ExistentChannels[channel] = append(ExistentChannels[channel], token)
}	