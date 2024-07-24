package modules

import (
	"fmt"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/config"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/db"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// SetFSub enables or disables Force Sub in the chat.
func setFSub(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveChat.Type == gotgbot.ChatTypePrivate {
		_, _ = ctx.EffectiveMessage.Reply(b, "This command is only available in groups.", nil)
		return ext.EndGroups
	}

	if !isAdmin(ctx.EffectiveChat, ctx.EffectiveUser, b) {
		_, _ = ctx.EffectiveMessage.Reply(b, "You must be an admin to use this command.", nil)
		return ext.EndGroups
	}

	if repliedMsg := ctx.EffectiveMessage.ReplyToMessage; repliedMsg != nil && repliedMsg.ForwardOrigin != nil {
		return handleForwardedMessage(ctx, b, repliedMsg)
	}

	return handleCommand(ctx, b)
}

// isAdmin checks if the user is an admin in the chat or the owner of the bot.
func isAdmin(chat *gotgbot.Chat, user *gotgbot.User, b *gotgbot.Bot) bool {
	userMember, _ := chat.GetMember(b, user.Id, nil)
	return userMember.GetStatus() == "administrator" || config.OwnerId == user.Id || userMember.GetStatus() == "creator"
}

// handleForwardedMessage enables Force Sub in the chat using a forwarded message.
func handleForwardedMessage(ctx *ext.Context, b *gotgbot.Bot, repliedMsg *gotgbot.Message) error {
	msgOrigen := repliedMsg.ForwardOrigin.MergeMessageOrigin()

	if msgOrigen.Chat.Type != gotgbot.ChatTypeChannel {
		_, _ = ctx.EffectiveMessage.Reply(b, "Reply to a forwarded message from a channel.", nil)
		return ext.EndGroups
	}

	botMember, err := msgOrigen.Chat.GetMember(b, b.Id, nil)
	if err != nil || botMember.GetStatus() != "administrator" {
		_, _ = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("I must be an admin in the %s to use this command.", msgOrigen.Chat.Title), nil)
		return err
	}

	go db.SetFSubChannel(ctx.EffectiveChat.Id, msgOrigen.Chat.Id)
	db.SetFSub(ctx.EffectiveChat.Id, true)

	text := fmt.Sprintf("Force Sub enabled in %s.\nChannelID: %d", ctx.EffectiveChat.Title, msgOrigen.Chat.Id)
	_, err = ctx.EffectiveMessage.Reply(b, text, nil)
	return err
}

// handleCommand enables or disables Force Sub in the chat.
func handleCommand(ctx *ext.Context, b *gotgbot.Bot) error {
	fSub := db.GetFSubSetting(ctx.EffectiveChat.Id)
	var text string

	if len(ctx.Args()) == 1 {
		text = fmt.Sprintf("Force Sub is %s in %s.", onOff(fSub.ForceSub), ctx.EffectiveChat.Title)
	} else {
		switch ctx.Args()[1] {
		case "enable", "on", "true", "y", "yes":
			if !fSub.ForceSub {
				go db.SetFSub(ctx.EffectiveChat.Id, true)
				text = "Force Sub enabled."
			} else {
				text = "Force Sub is already enabled."
			}
		case "disable", "off", "false", "n", "no":
			if fSub.ForceSub {
				go db.SetFSub(ctx.EffectiveChat.Id, false)
				text = "Force Sub disabled."
			} else {
				text = "Force Sub is already disabled."
			}
		default:
			text = "Invalid argument. Use /fsub on or /fsub off."
		}
	}

	_, err := ctx.EffectiveMessage.Reply(b, text, nil)
	return err
}

// onOff returns "enabled" if state is true, "disabled" otherwise.
func onOff(state bool) string {
	if state {
		return "enabled"
	}
	return "disabled"
}
