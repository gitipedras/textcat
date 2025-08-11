package handler

import (
	"fmt"
	"log/slog"
	"encoding/json"
	"database/sql"   // Package for SQL database interactions
	"github.com/gorilla/websocket"
)

type User struct {
	ID int
	username string
	sessionToken string
}

/*
func userExists(DB *sql.DB, username string) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
*/

func checkSession(DB *sql.DB, username string, sessionToken string) (bool, error) {
	var count int
	err := DB.QueryRow(`
	SELECT COUNT(*) FROM users
	WHERE username = ? AND sessiontoken = ?`,
	username, sessionToken).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func insertUser(DB *sql.DB, user User) error {
	stmt, err := DB.Prepare("INSERT INTO users(username, sessiontoken) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.username, user.sessionToken)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)
	slog.Info("Created user with ID: %d\n", user.ID)

	return nil
}

func auth(DB *sql.DB, username string, sessionToken string) string {
	fmt.Printf("[auth] recieved user %s with session token %s \n", username, sessionToken)

	//check SessionToken
	status, err := checkSession(DB, username, sessionToken)
	if err != nil {
		fmt.Println(err)
		return "err"
	}

	if status == true {
		return "ok"
	} else {
		return "invalid"
	}

	return "err"
}


func ProcessMsg(DB *sql.DB, message string, conn *websocket.Conn) {
	type wsRequest struct {
		Rtype string `json:"rtype"`
		Username string `json:"username"`
		SessionToken string `json:"sessionToken"`
		Message string `json:"message"`
		//clientID string `json:"clientID"`
	}

	type wsSend struct {
		Rtype string `json:"rtype"`
		Value string `json:"value"`
		//ClientAlert string `json:"clientalert"`
	}

	var out wsRequest
	err := json.Unmarshal([]byte(message), &out)
	if err != nil {
		fmt.Println("Error reading json: ", err)
		return
	}


	/*
	fmt.Println("Rtype:", out.Rtype)
	fmt.Println("Username:", out.Username)
	fmt.Println("SessionToken:", out.SessionToken)
	*/

	// Example case statement
	switch out.Rtype {
		case "loginRequest":
			slog.Info("Username login request for: ", out.Username)
			status := auth(DB, out.Username, out.SessionToken)
			if status == "ok" {
				un_returnmsg := wsSend{
					Rtype: "goodCredentials",
					Value: "",
				}

				// Convert struct → []byte (JSON)
				msgBytes, err := json.Marshal(un_returnmsg)
				if err != nil {
					return
				}

				// Send over WebSocket
				if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
					fmt.Println("Error writing message:", err)
					return
				}


			} else {
				un_returnmsg := wsSend{
					Rtype: "invalidCredentials",
					Value: "",
				}

				// Convert struct → []byte (JSON)
				msgBytes, err := json.Marshal(un_returnmsg)
				if err != nil {
					return
				}

				// Send over WebSocket
				if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
					fmt.Println("Error writing message:", err)
					return
				}
			}

		case "register":
			newUser := User{
				username:     out.Username,
				sessionToken: out.SessionToken,
			}

			err := insertUser(DB, newUser)
			if err != nil {
				slog.Error("Failed to insert user:", err)
			}

		case "sendMessage":
			fmt.Println("[MessageProcessor] ", out.Message)

			status := auth(DB, out.Username, out.SessionToken)
			if status == "ok" {
				un_returnmsg := wsSend{
					Rtype: "goodCredentials",
					Value: "",
				}

				// Convert struct → []byte (JSON)
				msgBytes, err := json.Marshal(un_returnmsg)
				if err != nil {
					return
				}

				// Send over WebSocket
				if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
					fmt.Println("Error writing message:", err)
					return
				}


			} else {
				un_returnmsg := wsSend{
					Rtype: "invalidCredentials",
					Value: "",
				}

				// Convert struct → []byte (JSON)
				msgBytes, err := json.Marshal(un_returnmsg)
				if err != nil {
					return
				}

				// Send over WebSocket
				if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
					fmt.Println("Error writing message:", err)
					return
				}
			}

		default:
			fmt.Println("Unknown request type")
	}
}
