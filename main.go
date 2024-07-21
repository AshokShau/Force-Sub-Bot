package main

import (
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/config"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/modules"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"time"
)

func main() {
	b, err := gotgbot.NewBot(config.Token, nil)
	if err != nil {
		config.ErrorLog.Fatal("failed to create new bot:", err)
	}

	dispatcher := ext.NewDispatcher(nil)
	modules.LoadModules(dispatcher)

	updater := ext.NewUpdater(dispatcher, nil)
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})

	if err != nil {
		config.ErrorLog.Fatal("failed to start polling:", err)
	}

	config.InfoLog.Println("Bot started as @" + b.Username)
	_, _ = b.SendMessage(config.LoggerId, "Bot started;", nil)
	updater.Idle()
}
