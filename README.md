# textcat client
A simple chat application made using golang and javascript

## Features
- Connect to any server (using ws://) only
- Login and register on any server


## Documentation

### Requests

**`login`**

Sent to the server when a client wants to login

Send Paramaters:
`Rtype`: `login`
`Username`: username
`SessionToken`: password

Response Paramaters:
`Rtype`: `loginStats`

`Status`: response status
ok = went well
invalid = invalid json data, username or password
already = already logged in / username already has a session
isr = internal server error

`Value`: session token -- **NOTE**: you must save this

**`register`**

Send Paramaters:
`Rtype`: `register`

`Username` // chosen username

`SessionToken` // password

Response Paramaters:
`Rtype`: `registerStats`

`Status`: response status
ok = went well
invalid = invalid json data, username or password
alreadyExists = already logged in / username already has a session
isr = internal server error

`Value`: session token -- **NOTE**: you must save this