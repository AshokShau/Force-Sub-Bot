package modules

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"time"
)

// start is the handler for the /start command
func start(b *gotgbot.Bot, ctx *ext.Context) error {
	text := fmt.Sprintf("Hello, %s!\n\nI am a bot that can help you manage your group by forcing users to join a channel before they can send messages in the group.\n\nTo get started, add me to your group and make me an admin with ban users permission. Then, set the channel that you want users to join using /fsub command.\n\nFor more information, click the button below.", ctx.EffectiveMessage.From.FirstName)
	button := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{{
				Text: "Add Me To Your Group",
				Url:  fmt.Sprintf("https://t.me/%s?startgroup=true", b.Username),
			}},
			{
				{
					Text: "Update Channel",
					Url:  "https://t.me/FallenProjects",
				},
				{
					Text: "Help",
					Url:  "https://abishnoi69.github.io/Force-Sub-Bot/#commands",
				},
			},
		},
	}

	if _, err := ctx.EffectiveMessage.Reply(b, text, &gotgbot.SendMessageOpts{ReplyMarkup: button}); err != nil {
		return logError(fmt.Sprintf("[Start] Error sending message - %v", err), err)
	}

	return ext.EndGroups
}

// ping is the handler for the /ping command
func ping(b *gotgbot.Bot, ctx *ext.Context) error {
	startTime := time.Now()

	if msg, err := ctx.EffectiveMessage.Reply(b, "Pong!", nil); err != nil {
		return logError(fmt.Sprintf("[Ping] Error sending message - %v", err), err)
	} else if _, _, err := msg.EditText(b, "Pong! "+time.Since(startTime).String(), nil); err != nil {
		return logError(fmt.Sprintf("[Ping] Error editing message - %v", err), err)
	}

	return ext.EndGroups
}
