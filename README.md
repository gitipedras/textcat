# textcat (still in beta!)
A simple chat application made using golang and javascript

## Features

### Client
**Textcat does not require any app! You can run it in your browser**

Find on the client on the `client` branch.
You can just download the client and open it in your browser.

### Server

- Login and create an account
- Talk with other people
- Addons (beta)

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

## Addons (beta)
Addons allow you to use lua code to create custom commands!

## Version Control

Client and server versions **are not matched**!

Version system:
`MAJOR`.`MINOR`.`PATCH`