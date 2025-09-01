# textcat (still in beta!)
A simple chat application made using golang and javascript

## Features



## Client
**Textcat does not require any app! You can run it in your browser**

Find on the client on the `client` branch.
You can just download the client and open it in your browser.

## TODO

**Server**

- prevent users from logging in twice (session token spam)
- prevent users from reconnecting to a channel
- prevent users from connecting to two channels at once
- prevent message spam

**Client**

- check if fields are empty

**Server code to eventually be used/added (ordered by importance)**
```
import ("flag")

addr := flag.String("addr", ":4000", "HTTP network address")
^--> variable with our address  ^---> default port
```

```golang
s := &http.Server{
	Addr:           ":8080",
	Handler:        myHandler,
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
}
log.Fatal(s.ListenAndServe())
```

```golang

// Send the message to all connected clients
mutex.Lock()
for client := range clients {
	err := client.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		client.Close()
		delete(clients, client)
	}
}
mutex.Unlock()

```