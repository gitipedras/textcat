package models


import (
	"os"
	"log/slog"
)

type Application struct {
	Log *slog.Logger // basic logger
}

var App = &Application {
	    Log: slog.New(slog.NewTextHandler(os.Stderr, nil)),
}

type WsIncome struct {
	Rtype string `"json:Rtype"`
	Username string `"json:Username"`
	SessionToken string `"json:sessionToken"`
}

/*

// USE A METHOD
func (a *Application) InitAppFields() {
	a.Log = slog.New(slog.NewTextHandler(os.Stderr, nil))
}

var App *Application

func InitApp() {
	App = &Application{}
	App.InitAppFields()
}

*/