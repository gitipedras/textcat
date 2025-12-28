package main

import (
	"log/slog"
	"github.com/zion8992/textcat/tc"
	"encoding/json"
	"fmt"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Application struct {
	Log *slog.Logger
	Database *sql.DB
	//Middleware appMiddleware
	//Sessions tc.SessionManager
	//Auth tc.Auth
}

func (app *Application) Store(table string, record any) error {
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	fmt.Println(data)
	return nil
}

func (app *Application) LogMsg(level string, message string, args ...any) {
	switch level {
		case "info":
			app.Log.Info(message, args...)
		case "warn":
			app.Log.Warn(message, args...)
		case "error":
			app.Log.Error(message, args...)
		default:
			app.Log.Info(message, args...)
	}
}


func (app *Application) HandleReq(msg []byte) error {
    var data tc.Recieve

    if err := json.Unmarshal(msg, &data); err != nil {
        return err
    }

    fmt.Println(data)
    return nil
}
