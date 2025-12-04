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

## Building a zip (for github releases)
`make`

This builds for windows, darwin (macos) and linux for x86 and arm64, then compresses them into a zip file.

## Server Configuration

Configuration is stored in `config.json`

```json
{
    "ServerName": "Server Name", // Short and readable name
    "ServerDesc": "Short Description", // Description, recommended up to 3 lines (for official client)
    "Port": ":8080", // port to host http server HOST:PORT. Use :PORT for localhost
    "MaxLength": 2 // maximum message length
}

```

<!--
random option that does not work yet
"CacheMessages": false
-->

## Version Control

The official server and client versions **are not matched**!

Version system:
`MAJOR`.`MINOR`.`PATCH`
