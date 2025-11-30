package core

import (
	"textcat/channels"
    "textcat/database"
    "textcat/models"
	"time"
)

var Channels channels.ChannelHandler

func ChannelsInit() {
    Channels = channels.ChannelHandler{
        StartedAt: time.Now(),
        Channels:  make(map[string]*channels.Channel), // must initialize
    }

    Channels.NewChannel("main")

    // Load channels from DB
    dbChannels, err := database.GetAllChannels()
    if err != nil {
        models.App.Log.Error("Failed to load channels from DB:", err)
        return
    }

    for _, chName := range dbChannels {
        Channels.Channels[chName] = &channels.Channel{
            Description: "",
            Connected:   make(map[string]string),
            Permissions: make(map[string][]string),
        }
    }
}
