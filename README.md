# textcat
A simple chat application made using the ultimate golang.

Webiste: Coming Soon!


## Features

### Client

You can use any client, as long as its compatible with the textcat protocol.
Official client here: [Source Code](redirect.textcat.net/)

### Server

- Login and create an account
- Talk with other people
- Addons (coming soon)
- Self-host textcat

## Server Configuration

Configuration is stored in `config.json`

```json
{
    "ServerName": "Server Name", // Only used for clients
    "ServerDesc": "Short Description", // Only used for clients
    "Port": ":8080", // port to host http server HOST:PORT use :PORT for localhost
    "MaxLength": 2, // maximum message length
    "CacheMessages": false
}

```

## Version Control

The official server and client versions **are not matched**!

Version system:
`MAJOR`.`MINOR`.`PATCH`
