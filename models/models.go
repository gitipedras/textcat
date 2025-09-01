package models


import (
	"os"
	"log/slog"
	"encoding/json"
)

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
}

type WsIncome struct {
	Rtype string `"json:Rtype"`
	Username string `"json:Username"`
	SessionToken string `"json:Password"` // change SessionToken to password later pls
	Message string `"json:Message"`
	ChannelID string `"json:ChannelID"`
}

type WsSend struct {
	Rtype string `"json:Rtype"`
	Status string `"json:Status"`
	Value string `"json:Value"`
	Username string `"json:Username"`
	ServerName string `"json:ServerName"`
	ServerDesc string `"json:ServerDesc"`
}

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