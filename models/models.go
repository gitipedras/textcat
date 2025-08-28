package models


import (
	"os"
	"log/slog"
)

type Application struct {
	Log    *slog.Logger // basic logger
}

var App = &Application{
	Log:    slog.New(slog.NewTextHandler(os.Stderr, nil)),
}


type WsIncome struct {
	Rtype string `"json:Rtype"`
	Username string `"json:Username"`
	SessionToken string `"json:sessionToken"`
	Message string `"json:message"`
	ChannelID string `"json:channelID"`
}

type WsSend struct {
	Rtype string `"json:Rtype"`
	Status string `"json:Status"`
	Value string `"json:Value"`
}