package modules

import (
	"fmt"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/db"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"time"
)

var StartTime = time.Now()

// start is the handler for the /start command
func start(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveChat.Type == gotgbot.ChatTypePrivate {
		_ = db.AddUser(b.Id, ctx.EffectiveMessage.From.Id)
	} else {
		_ = db.AddChat(b.Id, ctx.EffectiveChat.Id)
	}
	
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
	msg := ctx.EffectiveMessage
	startTime := time.Now()

	rest, err := msg.Reply(b, "<code>Pinging</code>", &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	if err != nil {
		return logError(fmt.Sprintf("[Ping] Error sending message - %v", err), err)
	}

	// Calculate latency
	elapsedTime := time.Since(startTime)

	// Calculate uptime
	uptime := time.Since(StartTime)
	formattedUptime := getFormattedDuration(uptime)

	location, _ := time.LoadLocation("Asia/Kolkata")
	responseText := fmt.Sprintf("Pinged in %vms (Latency: %.2fs) at %s\n\nUptime: %s", elapsedTime.Milliseconds(), elapsedTime.Seconds(), time.Now().In(location).Format(time.RFC1123), formattedUptime)

	_, _, err = rest.EditText(b, responseText, nil)
	if err != nil {
		return logError(fmt.Sprintf("[Ping] Error editing message - %v", err), err)
	}

	return ext.EndGroups
}
