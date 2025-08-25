func dbInit() {
	// Connect to the SQLite database
    db, err := sql.Open("sqlite3", "./appdata.db")
    if err != nil {
        fmt.Println(err)
        return
    }

    defer db.Close()
    fmt.Println("Connected to the SQLite database successfully.")

}