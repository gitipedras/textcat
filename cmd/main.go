package main

import (
	"fmt"

	"log/slog" // logging
	"os" // logging

	"github.com/gorilla/websocket" // websocket
	"net/http"
    "database/sql"
    "github.com/zion8992/textcat/tc"
    "strings"
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

            if err := app.HandleReq(msg); err != nil {
                if strings.HasPrefix(err.Error(), "error") {
                    MakeRequest("status", err.Error(), "error", conn)
                    app.Log.Error("request returned", slog.Any("error", err))
                } else if strings.HasPrefix(err.Error(), "ok") {
                    MakeRequest("status", err.Error(), "ok", conn)
                    app.Log.Error("request ok", slog.Any("error", err))
                } else {
                    MakeRequest("status", err.Error(), "server_error", conn)
                    app.Log.Error("internal server error", slog.Any("error", err))
                }
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

    app.Textcat = &tc.Textcat{
        Function: app, // Application implements Handler
        Sessions: tc.NewSessionManager(),
    }
	return 	app
}

/* 
Function run(*Application) runs a textcat server
*/

func slashHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This server is running textcat :)"))
}


func run(app *Application) {
	defer app.Database.Close() // close the database connection

	var port string = ":8080"

	// Wrap ws(app) with a panic recovery
	handler := http.HandlerFunc(ws(app))
    slash := http.HandlerFunc(slashHandler)
	http.Handle("/textcat", RecoverMiddleware(handler, app))
    http.Handle("/", RecoverMiddleware(slash, app))

	app.Log.Info("started textcat server", slog.String("port", port))
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println(err)
	}
}

// Middleware to catch panics in HTTP handler goroutines
func RecoverMiddleware(next http.Handler, app *Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				app.Log.Error("goroutine panicked: ", slog.Any("error", rec))

				// Check if this is a WebSocket upgrade
				if websocket.IsWebSocketUpgrade(r) {
					// Upgrade the connection so we can write a message
					conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
					if err == nil {
						msg := map[string]string{
							"Req":    "status",
							"Status": "server_error",
						}
						conn.WriteJSON(msg)
						conn.Close()
					}
				} else {
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
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