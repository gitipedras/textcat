package channels

import (
    "time"
	"sync"
	"github.com/gorilla/websocket"
	"encoding/json"
    "log/slog"

	/* internal packages */
	"textcat/auth"
	"textcat/models"
    "textcat/database"
)   

type ChannelHandler struct {
	Mu sync.RWMutex
	StartedAt time.Time
	Channels map[string]*Channel
    MessageCache map[string][]CachedMessage

    // do this here since global variables are defined at compile time
    // meaning var <xyz> would be nil if we put it up there
    MaxCachedMessages int
    MessageCacheEnabled bool
}

type Channel struct {
	// name is in the map
	Description string
	Connected map[string]string
	Permissions map[string][]string
	// permissions are not pre-specified
	// meaning there are no default permissions
	// the default perm is nothing
	// use an addon to stop people from chatting
	// on specific channels
}

type CachedMessage struct {
	Username string
	Message  string
	Time     time.Time
}




// backend
func (ch *ChannelHandler) NewChannel(channelName string) {
    ch.Mu.Lock()
    defer ch.Mu.Unlock() // wait until method brackets end
    channel := Channel{
    	Description: "",
    	Connected: make(map[string]string),
    	Permissions: make(map[string][]string),
    }
    ch.Channels[channelName] = &channel

    ok := database.AddChannel(channelName)
    if ok == false {
        //models.App.Log.Error("Failed to create channel!")
        //panic("Failed to create a channel: Possible database or channels error")
        // db can return false if channel exists, no need to panic
    }
}

func (h *ChannelHandler) BuildChannelList() map[string]int {
    h.Mu.RLock()
    defer h.Mu.RUnlock()

    result := make(map[string]int)
    for name, ch := range h.Channels {
        result[name] = len(ch.Connected)
    }
    return result
}


func (h *ChannelHandler) AddUser(channelName, token, username string, conn *websocket.Conn) {
    h.Mu.Lock()

    ch, ok := h.Channels[channelName] // ch is *Channel
    if !ok {
        models.App.Log.Info("invalid channel", slog.String("chid", channelName))
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
        return
    }

    if auth.SessionManager.Exists(token) {
    	// Add user to the channel's Connected map
    	models.App.Log.Info("[CONNECT] user connected", slog.String("token", token))
    	ch.Connected[username] = token
    	models.App.Log.Info("valid channel", slog.String("chid", channelName))
		response := models.WsSend {
	            Rtype:   "connectStats",
	            Status:  "ok",
	    }
	    data, err := json.Marshal(response)
	    if err != nil {
	        models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
	        return
	    }
	    conn.WriteMessage(websocket.TextMessage, data)

        if h.MessageCacheEnabled {
            var cached []CachedMessage

            if h.MessageCacheEnabled {
                cached = append([]CachedMessage(nil), h.MessageCache[channelName]...)
            }

            for _, msg := range cached {
                response := models.WsSend{
                    Rtype: "NewMessage",
                    Status: "newmsg",
                    Value: msg.Message,
                    Username: msg.Username,
                    Time: msg.Time,
                }
                data, _ := json.Marshal(response)
                conn.WriteMessage(websocket.TextMessage, data)
            }
        }

        h.Mu.Unlock()


    } else {
    	models.App.Log.Info("invalid token to connect to channel", slog.String("chid", channelName))
		response := models.WsSend {
	            Rtype:   "connectStats",
	            Status:  "invalidToken",
	    }
	    data, err := json.Marshal(response)
	    if err != nil {
	        models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
	        return
	    }
	    conn.WriteMessage(websocket.TextMessage, data)
    }
}

func (h *ChannelHandler) RemoveTokenFromAllChannels(token string) {
    h.Mu.Lock()
    defer h.Mu.Unlock()

    for _, ch := range h.Channels {
        for username, currentToken := range ch.Connected {
            if currentToken == token {
                delete(ch.Connected, username)
            }
        }
    }
}



func (h *ChannelHandler) RemoveUser(channelName, token, username string) {
    h.Mu.Lock()
    defer h.Mu.Unlock()

    ch, ok := h.Channels[channelName]
    if !ok {
        // channel doesn't exist
        return
    }

    // Only remove if the token matches the one stored for this username
    if currentToken, exists := ch.Connected[username]; exists && currentToken == token {
        delete(ch.Connected, username)
    }
}



func (h *ChannelHandler) ChannelExists(channelName string) bool {
    h.Mu.RLock()
    defer h.Mu.RUnlock()

    _, ok := h.Channels[channelName]
    if !ok {
        return false
    }
    return true
}


func (h *ChannelHandler) CheckPerm(channelName, username, permission string) bool {
    h.Mu.RLock()
    defer h.Mu.RUnlock()

    ch, ok := h.Channels[channelName]
    if !ok {
        return false
    }

    for _, perm := range ch.Permissions[username] {
        if perm == permission {
            return true
        }
    }
    return false
}

// used for users sending stuff
func (h *ChannelHandler) SendMessage(channelName, message, username, token string, conn *websocket.Conn) bool {
    h.Mu.Lock()

    models.App.Log.Info("received message!", slog.String("message", message), slog.String("channel", channelName))

    ch, ok := h.Channels[channelName]
    if !ok {
        models.App.Log.Info("invalid channel", slog.String("chid", channelName))
        response := models.WsSend{
            Rtype:  "invalidChannel",
            Status: "non-existent",
        }
        data, err := json.Marshal(response)
        if err != nil {
            models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
            return false
        }
        conn.WriteMessage(websocket.TextMessage, data)
        return false
    }

    h.Mu.Unlock()

        if message == "/hi" {

        message4Client := models.WsSend{
				Rtype:  "NewMessage",
				Status: "newmsg",
				Value:  "Hello!",
				Username: "Server",
				Time: time.Now(),

	    }			
		data, failed := json.Marshal(message4Client)
		if failed != nil {
			models.App.Log.Error("Failed to parse json! ", slog.Any("error", failed))
			return false
		}
		
        err := auth.SessionManager.SendToClient(token, []byte(data))
        if err != nil {
            models.App.Log.Error("failed to send", slog.String("err", err.Error()), slog.String("token", token))
        }
        return true

    } else if message == "/about" {
        message4Client := models.WsSend{
				Rtype:  "NewMessage",
				Status: "newmsg",
				Value:  "This server is running textcat server official!",
				Username: "Server",
				Time: time.Now(),

	    }			
		data, failed := json.Marshal(message4Client)
		if failed != nil {
			models.App.Log.Error("Failed to parse json! ", slog.Any("error", failed))
			return false
		}
		
        err := auth.SessionManager.SendToClient(token, []byte(data))
        if err != nil {
            models.App.Log.Error("failed to send", slog.String("err", err.Error()), slog.String("token", token))
        }
        return true

    } else if message == "/source" {

        message4Client := models.WsSend{
				Rtype:  "NewMessage",
				Status: "newmsg",
				Value:  "Source -> github.com/gitipedras/textcat",
				Username: "Server",
				Time: time.Now(),

	    }			
		data, failed := json.Marshal(message4Client)
		if failed != nil {
			models.App.Log.Error("Failed to parse json! ", slog.Any("error", failed))
			return false
		}
		
        err := auth.SessionManager.SendToClient(token, []byte(data))
        if err != nil {
            models.App.Log.Error("failed to send", slog.String("err", err.Error()), slog.String("token", token))
        }
        return true
    }

    if h.MessageCacheEnabled == true {

        h.MessageCache[channelName] = append(
            h.MessageCache[channelName],
            CachedMessage{
                Username: username,
                Message:  message,
                Time:     time.Now(),
            },
        )

        msgs := h.MessageCache[channelName]
        if len(msgs) > h.MaxCachedMessages {
            h.MessageCache[channelName] = msgs[len(msgs)-h.MaxCachedMessages:]
        }
    }


    //debug stuff that isn't really interesting
    //models.App.Log.Info("[MESSAGES] message OK, sending...")


    sent := make(map[string]struct{})
    for _, userToken := range ch.Connected {
        if _, seen := sent[userToken]; seen {
            continue
        }
        models.App.Log.Info("channel state", slog.Any("connected", ch.Connected))
        
        message4Client := models.WsSend{
				Rtype:  "NewMessage",
				Status: "newmsg",
				Value:  message,
				Username: username,
				Time: time.Now(),

	    }			
		data, failed := json.Marshal(message4Client)
		if failed != nil {
			models.App.Log.Error("Failed to parse json! ", slog.Any("error", failed))
			return false
		}
		
        err := auth.SessionManager.SendToClient(userToken, []byte(data))
        if err != nil {
            models.App.Log.Error("failed to send", slog.String("err", err.Error()), slog.String("token", userToken))
        }
        sent[userToken] = struct{}{}
    }

    return true
}


