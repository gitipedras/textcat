package auth

import (
	/* other stuff */
    "log/slog"
    "time"
    "encoding/json"

    /* websockets!!! */
    "github.com/gorilla/websocket"

	/* internal  */
	"textcat/database"
	"textcat/models"
    "textcat/sessions"
    "textcat/validator"
)

var SessionManager = sessions.NewSessionManager()


func UserLogin(conn *websocket.Conn, msg models.WsIncome) {
    goodInput := validator.Username(msg.Username)
    if !goodInput {
        response := models.WsSend {
            Rtype:   "loginStats",
            Status:  "invalidInput",
            Value: msg.Username,
        }
        data, err := json.Marshal(response)
        if err != nil {
            models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
            return
        }
        
        conn.WriteMessage(websocket.TextMessage, data)
        return
    }

	ok := database.CheckUser(msg.Username)
	if ok {

		good := database.CheckPass(msg.Username, msg.SessionToken)
		if good {
			models.App.Log.Info("UserLogin", slog.String("username", msg.Username))
		
            token, err := SessionManager.GenerateToken(16) // 16 bytes = 32 hex chars
            if err != nil {
                models.App.Log.Error("Failed to generate session token", slog.String("err", err.Error()))
            }
            
            // reusable struct, sessionManager wont complain if this is duped
            session := &sessions.Session {
                Username:     msg.Username,
                SessionToken: token,
                Conn:         conn, // the websocket.Conn for this client
                ConnectedAt:  time.Now(),
            }

            SessionManager.Add(session)

            msg := models.WsSend{
                Rtype:   "loginStats",
                Status: "ok",
                Value: token,
                ServerName: models.Config.ServerName,
                ServerDesc: models.Config.ServerDesc,
            }

            data, err := json.Marshal(msg)
            if err != nil {
                models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
            }
            
            err = SessionManager.SendToClient(token, data) 
            if err != nil {
                models.App.Log.Error("Failed to send message", slog.String("err", err.Error()))
            }

        } else {
            response := models.WsSend {
            Rtype:   "loginStats",
            Status:  "invalid",
            }
            data, err := json.Marshal(response)
            if err != nil {
                models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
                return
            }
            conn.WriteMessage(websocket.TextMessage, data)
        }

	} else {
		// user already exists
        response := models.WsSend {
            Rtype:   "loginStats",
            Status:  "invalid",
        }
        data, err := json.Marshal(response)
        if err != nil {
            models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
            return
        }
        
        conn.WriteMessage(websocket.TextMessage, data)
	}
}

func UserRegister(conn *websocket.Conn, msg models.WsIncome) {
    models.App.Log.Info("user register", slog.String("username", msg.Username))
    goodInput := validator.Username(msg.Username)
    if !goodInput {
        response := models.WsSend {
            Rtype:   "loginStats",
            Status:  "invalidInput",
        }
        data, err := json.Marshal(response)
        if err != nil {
            models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
            return
        }
        
        conn.WriteMessage(websocket.TextMessage, data)
        return
    }
	ok := database.CheckUser(msg.Username)
    if ok {
        // user already exists
        response := models.WsSend {
            Rtype:   "registerStats",
            Status:  "alreadyExists",
        }
        data, err := json.Marshal(response)
        if err != nil {
            models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
            return
        }
        
        conn.WriteMessage(websocket.TextMessage, data)
    } else {
        err := database.CreateUser(msg.Username, msg.SessionToken)
        if err != nil {
            response := models.WsSend {
            Rtype:   "registerStats",
            Status:  "isr",
            }
            data, err := json.Marshal(response)
            if err != nil {
                models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
                return
            }
            conn.WriteMessage(websocket.TextMessage, data)

            models.App.Log.Error("Failed to Create user")
        } else {
            response := models.WsSend {
            Rtype:   "registerStats",
            Status:  "ok",
            }
            data, err := json.Marshal(response)
            if err != nil {
                models.App.Log.Error("Failed to marshal JSON", slog.String("err", err.Error()))
                return
            }
            conn.WriteMessage(websocket.TextMessage, data)
        }
    }
}


func SessionTimer() {
    models.App.Log.Info("Started session timer")

    expiration := 5  * time.Hour

    for {
        SessionManager.Mu.Lock()
        for token, session := range SessionManager.Sessions {
            if time.Since(session.ConnectedAt) > expiration {
                delete(SessionManager.Sessions, token)
                session.Conn.Close()
                models.App.Log.Info("Expired session removed", slog.String("user", session.Username))
            }
        }
        SessionManager.Mu.Unlock()

        time.Sleep(10 * time.Second) // wait 10 secs before checking again
    }
}
