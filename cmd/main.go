package main

import (
	"fmt"

	"log/slog" // logging
	"os" // logging

	"github.com/gorilla/websocket" // websocket
	"net/http"
    "database/sql"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true }, // Allow all connections
}

/*
Middleware handles checking user sessions, preventing spam, etc...
Middlware is implemented in ws()

ws() returns a http.HandlerFunc function cuz it needs access to our app struct
*/

func ws(app *Application) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            app.Log.Error("upgrade failed", slog.Any("error", err))
            return
        }
        defer conn.Close()

        for {
            _, msg, err := conn.ReadMessage()
            if err != nil {
                break
            }

            MakeRequest("hi", "how are you", "ok", conn)

            if err := app.HandleReq(msg); err != nil {
                MakeRequest("error", "", "server_error", conn)
                app.Log.Error("failed to handle request", slog.Any("error", err))
            }
        }
    }
}

/*
Creates and returns an *Application (declared in models.go)
with all the needed components

Does not return error because any errors here should panic()
*/

func createApp() *Application {
    db, err := sql.Open("sqlite3", "textcat.db")
    if err != nil {
        panic(err)
    }
    //defer db.Close() -> happens in run()
    
    app := &Application {
		Log: slog.New(slog.NewTextHandler(os.Stderr, nil)),
        Database: db,
	}
	return 	app
}

/* 
Function run(*Application) runs a textcat server
*/

func run(app *Application) {
    defer app.Database.Close() // close the database connection

	var port string = ":8080"
	http.HandleFunc("/textcat", ws(app))

    app.Log.Info("started textcat server", slog.String("port", port))
    err := http.ListenAndServe(port, nil)
    if err != nil {
    	fmt.Println(err)
    }
}


/*
'*app := createApp()' won't work:
github.com/golang/go/issues/6842
*/

func main() {
	var app *Application
	app = createApp()
	run(app)
}