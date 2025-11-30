package database


import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	/* logging */
	"textcat/models"
	"log/slog"
)

var err error
var DB *sql.DB

////////////////////////////////////////////
//
//   DATABASE CHECKS
//
//
///////////////////////////////////////////


func DbInit() {
    /////// --- connect to the database --- ///////

    // Connect to the SQLite database
    DB, err = sql.Open("sqlite3", "./appdata.db")
    if err != nil {
        models.App.Log.Error("Failed to connect to database: ", slog.String("err", err.Error()))
        return
    }

    models.App.Log.Info("[DATABASE] Connected to the SQLite database successfully.")

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
        models.App.Log.Error("[DATABASE CHECKS] Failed to run database checks", slog.String("err", err.Error()))
    }

        channelsDBC := `
    CREATE TABLE IF NOT EXISTS channels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    extraData BLOB NULL
);`
    _, err = DB.Exec(channelsDBC)
    if err != nil {
        models.App.Log.Error("[DATABASE CHECKS] Failed to run database checks", slog.String("err", err.Error()))
    }

}

////////////////////////////////////////////
//
//   CHANNELS
//
//
///////////////////////////////////////////

// used for creating channels that are stored in the database
func GetAllChannels() ([]string, error) {
    rows, err := DB.Query("SELECT name FROM channels")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var result []string
    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            return nil, err
        }
        result = append(result, name)
    }

    return result, rows.Err()
}


func CheckChannel(name string) bool {
    var id int
    err := DB.QueryRow("SELECT id FROM channels WHERE name = ?", name).Scan(&id)
    if err != nil {
        if err == sql.ErrNoRows {
            // channel does not exist
            return false
        }
        models.App.Log.Error("[DATABASE] Failed to check if channel exists", slog.String("err", err.Error()))
    }
    return true // channel exists
}

func AddChannel(name string) bool {
    exists := CheckChannel(name)
    if exists == true {
        // channel exists
        return false
    }
    var extraData string
    query := `
        INSERT INTO channels (name, extraData)
        VALUES (?, ?)
    `
    _, err := DB.Exec(query, name, extraData)
    if err != nil {
        models.App.Log.Error("[DATABASE] Failed to create channel", slog.String("err", err.Error()))
        return false
    }
    return true
}

////////////////////////////////////////////
//
//   USERS
//
//
///////////////////////////////////////////


func CheckUser(username string) bool {
    var id int
    err := DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
    if err != nil {
        if err == sql.ErrNoRows {
            // user does not exist
            return false
        }
        models.App.Log.Error("[DATABASE] Failed to check if user exists", slog.String("err", err.Error()))
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

func CreateUser(username, password string) error {
    var description string = "empty"
    query := `
        INSERT INTO users (username, description, password)
        VALUES (?, ?, ?)
    `
    _, err := DB.Exec(query, username, description, password)
    if err != nil {
        models.App.Log.Error("Failed to insert user", slog.String("err", err.Error()))
        return err
    }
    return nil
}
