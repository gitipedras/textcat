package channels

import (
	/* data processing */
	"encoding/json"

	/* logging */
	"log/slog"
	"time"

	/* websocket ?? */
	"github.com/gorilla/websocket"

	/* internal */
	"textcat/auth"
	"textcat/models"
	"textcat/validator"
)

var ExistentChannels = map[string][]string{
	"main":      {},
	"minecraft": {},
	"minecraft-bedwars": {},
}

var mCache models.MessageCache

func HandleMSG(username string, token string, message string, channelID string, conn *websocket.Conn) {
	validInput := validator.Message(message)
	if !validInput {
		response := models.WsSend {
	            Rtype:   "invalidInput",
	            Status:  "message",
	    }
	    data, err := json.Marshal(response)
	    if err != nil {
	        models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
	        return
	    }
	    conn.WriteMessage(websocket.TextMessage, data)
		return
	}

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
		if models.Config.CacheMessages {
			mCache.AddMessage(username, message)
		}
		allValues := ExistentChannels[channelID]
		sent := make(map[string]bool)

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
				Time: time.Now(),
			}
			data, err := json.Marshal(response)
			if err != nil {
				models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
				return
			}

			err = auth.SessionManager.SendToClient(v, data)
			if err != nil {
				models.App.Log.Error("Failed to send message", slog.String("err", err.Error()))
			}
			sent[v] = true
			models.App.Log.Info("[MESSAGES] sent message to token", slog.String("token", v))
		}

		return

	} else {
		models.App.Log.Info("invalid channel", slog.String("chid", channelID))
		response := models.WsSend {
	            Rtype:   "invalidChannel",
	            Status:  "non-existent",
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
    // step 1: channel doesn’t exist
        response := models.WsSend {
	           Rtype:   "messageCache",
	           MsgCache: mCache.Cache,
	    }
	    data, err := json.Marshal(response)
	    if err != nil {
	        models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
	        return
	    }
	    conn.WriteMessage(websocket.TextMessage, data)
    models.App.Log.Info("[CONNECT] user connect", slog.String("channel", channel), slog.String("token", token))
    ExistentChannels[channel] = append(ExistentChannels[channel], token)

}

func DisconnectUser(token string, channel string, conn *websocket.Conn) {
    models.App.Log.Info("[DISCONNECT] user disconnect", slog.String("channel", channel), slog.String("token", token))
    users, ok := ExistentChannels[channel]
    if !ok {
        // step 1: channel doesn’t exist
        response := models.WsSend {
	           Rtype:   "disconnectStats",
	           Status:  "NoChannelFound",
	           Value: 	channel,
	    }
	    data, err := json.Marshal(response)
	    if err != nil {
	        models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
	        return
	    }
	    conn.WriteMessage(websocket.TextMessage, data)

        return
    }

    // step 2: check if token is present
    found := false
    newUsers := make([]string, 0, len(users))
    for _, t := range users {
        if t == token {
            found = true
            continue // skip this token → removes it
        }
        newUsers = append(newUsers, t)
    }

    if !found {
        // channel exists but token not found
        response := models.WsSend {
	           Rtype:   "disconnectStats",
	           Status:  "notConnected",
	           Value: 	channel,
	    }
	    data, err := json.Marshal(response)
	    if err != nil {
	        models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
	        return
	    }
	    conn.WriteMessage(websocket.TextMessage, data)
        return
    }

    // step 3: replace with updated user list
    ExistentChannels[channel] = newUsers

    // success
    response := models.WsSend {
	           Rtype:   "disconnectStats",
	           Status:  "ok",
	           Value: 	channel,
	}
	data, err := json.Marshal(response)
	if err != nil {
	    models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
	    return
	}
	conn.WriteMessage(websocket.TextMessage, data)
}

		