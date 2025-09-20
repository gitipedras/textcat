### TODO

**Server**

- prevent users from logging in twice (session token spam)
- prevent message spam
- add validation EVERYWHERE

**Client**

- fix some theming issues whatever
- make a desktop  version with electron or some other shit

**Server code to eventually be used/added (ordered by importance)**

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