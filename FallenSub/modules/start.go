package modules

import (
	"fmt"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/config"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"time"
)

// start sends a welcome message to the user.
func start(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	text := fmt.Sprintf("Hello, %s!\n\nI am a bot that can help you manage your group by forcing users to join a channel before they can send messages in the group.\n\nTo get started, add me to your group and make me an admin with ban users permission. Then, set the channel that you want users to join using /fsub command.\n\nFor more information, click the button below.", msg.From.FirstName)
	button := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text: "Help",
					Url:  "https://abishnoi69.github.io/Force-Sub-Bot/",
				},
			},
		},
	}

	_, err := msg.Reply(b, text, &gotgbot.SendMessageOpts{ReplyMarkup: button})
	if err != nil {
		config.ErrorLog.Printf("[Start] Error sending message - %v", err)
		return err
	}

	return ext.EndGroups
}

// ping responds to a ping command with "Pong!" and the latency.
func ping(b *gotgbot.Bot, ctx *ext.Context) error {
	startTime := time.Now()
	msg, err := ctx.EffectiveMessage.Reply(b, "Pong!", nil)
	if err != nil {
		config.ErrorLog.Printf("[Ping] Error sending message - %v", err)
		return err
	}

	// Calculate the latency
	latency := time.Since(startTime)

	_, _, err = msg.EditText(b, "Pong! "+latency.String(), nil)
	if err != nil {
		config.ErrorLog.Printf("[Ping] Error editing message - %v", err)
		return err
	}

	return ext.EndGroups
}
