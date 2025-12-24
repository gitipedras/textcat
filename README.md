# textcat
A simple chat application made using the ultimate golang.

Webiste: [gitipedras.github.io](https://gitipedras.github.io)


## Features

### Client

You can use any client, as long as its compatible with the textcat protocol. <br>
Official client here: [Source Code](https://github.com/gitipedras/textcat-telesto-client)

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
    "ServerName": "Server Name",
    "ServerDesc": "Short Description",
    "Port": ":8080",
    "MaxLength": 150
}

```

**ServerName**: set a name for the server
**ServerDesc**: description of your server
**Port**: port to host the chat/http server
**MaxLength**: maximum length for a message

<!--
random option that does not work yet
"CacheMessages": false
-->

## Version Control

The official server and client versions **are not matched**!

Version system:
`MAJOR`.`MINOR`.`PATCH`
