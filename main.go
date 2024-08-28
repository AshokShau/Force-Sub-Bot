package main

import (
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/modules"
	"time"

	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/config"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func main() {
	b, err := gotgbot.NewBot(config.Token, nil)
	if err != nil {
		config.ErrorLog.Fatal("failed to create new bot:", err)
	}

	updater := ext.NewUpdater(modules.Dispatcher, nil)
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
