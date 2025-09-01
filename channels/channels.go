package channels

import (
	/* data processing */
	"encoding/json"

	/* logging */
	"log/slog"

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
	ok := auth.SessionManager.CheckByUsername(username, token)
	if !ok {
		response := models.WsSend {
	            Rtype:   "invalidToken",
	            Status:  token,
	    }
	    data, err := json.Marshal(response)
	    if err != nil {
	        models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
	        return
	    }
	    conn.WriteMessage(websocket.TextMessage, data)
		return;
	}

	if _, ok := ExistentChannels[channelID]; ok {
		// channel exists

		allValues := []string{}
		sent := make(map[string]bool) // <- move here

		for _, value := range ExistentChannels {
			allValues = append(allValues, value...)
		}

		// loop through every token and send a message
		for _, v := range allValues {
			if sent[v] {
				continue // already sent to this token
			}
			response := models.WsSend{
				Rtype:  "NewMessage",
				Status: "newmsg",
				Value:  message,
				Username: username,
			}
			data, err := json.Marshal(response)
			if err != nil {
				models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
				return
			}

			models.App.Log.Info("broadcasting messages")
			err = auth.SessionManager.SendToClient(v, data)
			if err != nil {
				models.App.Log.Error("Failed to send message", slog.String("err", err.Error()))
			}
			sent[v] = true
			models.App.Log.Info("sent message to token", slog.String("token", v))
		}

		return

	} else {
		models.App.Log.Info("invalid channel", slog.String("chid", channelID))
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

	/*ok := auth.SessionManager.Exists(token)
	if ok {
		if clients, ok := ExistentChannels[channelID]; ok {
    		for _, token := range clients {
    			models.App.Log.Info("broadcasting messages")
			    auth.SessionManager.SendToClient(token, []byte(message))
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
	}*/
}

func ConnectUser(token string, channel string, conn *websocket.Conn) {
    models.App.Log.Info("user connected", slog.String("channel", channel), slog.String("token", token))
    ExistentChannels[channel] = append(ExistentChannels[channel], token)
    /*for channelName, clientTokens := range ExistentChannels {
		fmt.Printf("Channel: %s, Clients: %v\n", channelName, clientTokens)
	}*/

}	