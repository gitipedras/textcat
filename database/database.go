package database


import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	/* logging */
	"textcat/models"
	"log/slog"
)

var DB *sql.DB

func DbInit() {
	/////// --- connect to the database --- ///////

	// Connect to the SQLite database
    DB, err := sql.Open("sqlite3", "./appdata.db")
    if err != nil {
        models.App.Log.Error("Failed to connect to database: ", slog.String("err", err.Error()))
        return
    }

    defer DB.Close()
    models.App.Log.Info("Connected to the SQLite database successfully.")

    /////// --- database checks --- ///////
    /* DBC means database checks */

    // Create the users table if it doesnt exist
    userDBC := `
    CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    description TEXT NOT NULL,
    password TEXT NOT NULL
);`
    _, err = DB.Exec(userDBC)
    if err != nil {
        models.App.Log.Error("Failed to run database checks", slog.String("err", err.Error()))
    }

}

func CheckUser(username string) bool {
    var id int
    err := DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
    if err != nil {
        if err == sql.ErrNoRows {
            // user does not exist
            return false
        }
        models.App.Log.Error("Failed to check if user exists", slog.String("err", err.Error()))
    }
    return true // user exists
}

func CheckPass(username, password string) bool {
     var storedPassword string

    // Query the user's password
    err := DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
    if err != nil {
        if err == sql.ErrNoRows {
            return false // user does not exist
        }
        models.App.Log.Error("[database.go:68] failed to run CheckPass", slog.String("err", err.Error()))
    }

    // Compare the stored password with the provided one
    return storedPassword == password
}