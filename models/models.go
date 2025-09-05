package models


import (
	/* misc */
	"os"
	"log/slog"
	"encoding/json"
	"time"
)

/* App Struct */

type Application struct {
	Log    *slog.Logger // basic logger
}

var App = &Application{
	Log:    slog.New(slog.NewTextHandler(os.Stderr, nil)),
}

type AppConfig struct {
    ServerName string `json:"ServerName"`
    ServerDesc string `json:"ServerDesc"`
    Port       string    `json:"Port"`
    MaxLength  int    `json:"MaxLength"`
    CacheMessages bool `json:"CacheMessages"`
}

/* Ws Sending and Recieving */
// recieved from client
type WsIncome struct {
	Rtype string `"json:Rtype"`
	Username string `"json:Username"`
	SessionToken string `"json:Password"` // change SessionToken to password later pls
	Message string `"json:Message"`
	ChannelID string `"json:ChannelID"`
}

// sent to client
type WsSend struct {
	Rtype string `"json:Rtype"`
	Status string `"json:Status"`
	Value string `"json:Value"`
	Username string `"json:Username"`
	ServerName string `"json:ServerName"`
	ServerDesc string `"json:ServerDesc"`
	MsgCache map[string]string `json:"MsgCache"`
	Time  time.Time `json:"Time"`
}

/* Messages */

type MessageCache struct {
	Channel string
	Cache map[string]string `json:"Cache"`
}

func (mc *MessageCache) AddMessage(username string, msg string) {
    if mc.Cache == nil {
        mc.Cache = make(map[string]string)
    }

    mc.Cache[username] = msg

    // Optional: clear when reaching 10 messages
    if len(mc.Cache) >= 10 {
        mc.Cache = make(map[string]string)
    }
}


/* Configuration */

var Config AppConfig

// LoadConfig reads config.json into Config
func LoadConfig(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    err = decoder.Decode(&Config) // <- decode into Config, not AppConfig
    if err != nil {
        return err
    }

    return nil
}