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
	if _, ok := ExistentChannels[channelID]; ok {
		// channel exists

		// define where we will store our values
		allValues := []string{}

		for _, value := range ExistentChannels {
		    allValues = append(allValues, value...) // use ...
		    models.App.Log.Info("processing values", slog.Any("values", allValues))

		    // loop through every token and send a message
		    for _, v := range allValues {
		    	response := models.WsSend {
		            Rtype:  "newMessage",
		            Value:  message,
			    }
			    data, err := json.Marshal(response)
			    if err != nil {
			        models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
			        return
			    }
			    r := []byte(data)
			    models.App.Log.Info("sending message to client", slog.Any("allValues", v), slog.Any("json response", r))

		    	auth.SessionManager.SendToClient(v, r)
		    }

		}
		return

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

   /* for _, t := range ExistentChannels[channel] {
        if t == token {
            response := models.WsSend {
            Rtype:   "alreadyConnected",
            Status:  token,
	        }
	        data, err := json.Marshal(response)
	        if err != nil {
	            models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
	            return
	        }
	        conn.WriteMessage(websocket.TextMessage, data)
	        return
        }
    }*/
    ExistentChannels[channel] = append(ExistentChannels[channel], token)
    /*for channelName, clientTokens := range ExistentChannels {
		fmt.Printf("Channel: %s, Clients: %v\n", channelName, clientTokens)
	}*/

}	