package modules

import (
	"fmt"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/config"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/db"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func setFSub(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveChat.Type == gotgbot.ChatTypePrivate || !isAdmin(ctx.EffectiveChat, ctx.EffectiveUser, b) {
		_, _ = ctx.EffectiveMessage.Reply(b, "This command is only available in groups and requires admin privileges.", nil)
		return ext.EndGroups
	}

	if repliedMsg := ctx.EffectiveMessage.ReplyToMessage; repliedMsg != nil && repliedMsg.ForwardOrigin != nil {
		return handleForwardedMessage(ctx, b, repliedMsg)
	}

	return handleCommand(ctx, b)
}

func isAdmin(chat *gotgbot.Chat, user *gotgbot.User, b *gotgbot.Bot) bool {
	userMember, _ := chat.GetMember(b, user.Id, nil)
	return userMember.GetStatus() == "administrator" || config.OwnerId == user.Id || userMember.GetStatus() == "creator"
}

func handleForwardedMessage(ctx *ext.Context, b *gotgbot.Bot, repliedMsg *gotgbot.Message) error {
	msgOrigin := repliedMsg.ForwardOrigin.MergeMessageOrigin()

	if msgOrigin.Chat.Type != gotgbot.ChatTypeChannel || !isBotAdminInChannel(b, msgOrigin.Chat) {
		_, _ = ctx.EffectiveMessage.Reply(b, "Reply to a forwarded message from a channel where the bot is an admin.", nil)
		return ext.EndGroups
	}

	return enableForceSub(b, ctx, msgOrigin.Chat.Id)
}

func isBotAdminInChannel(b *gotgbot.Bot, chat *gotgbot.Chat) bool {
	botMember, err := chat.GetMember(b, b.Id, nil)
	return err == nil && botMember.GetStatus() == "administrator"
}

func enableForceSub(b *gotgbot.Bot, ctx *ext.Context, channelId int64) error {
	if err := db.SetFSubChannel(ctx.EffectiveChat.Id, channelId); err != nil {
		return logError("Error setting Force Sub channel", err)
	}

	if err := db.SetFSub(ctx.EffectiveChat.Id, true); err != nil {
		return logError("Error enabling Force Sub", err)
	}

	_, _ = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Force Sub enabled in %s.\nChannelID: %d", ctx.EffectiveChat.Title, channelId), nil)
	return nil
}

func handleCommand(ctx *ext.Context, b *gotgbot.Bot) error {
	fSub, err := db.GetFSubSetting(ctx.EffectiveChat.Id)
	if err != nil {
		return logError("Error getting Force Sub setting", err)
	}
	var text string
	if len(ctx.Args()) == 1 {
		text = fmt.Sprintf("Force Sub is %s in %s. \nChannelID: %d", onOff(fSub.ForceSub), ctx.EffectiveChat.Title, fSub.ForceSubChannel)
	} else {
		switch ctx.Args()[1] {
		case "enable", "on", "true", "y", "yes":
			if !fSub.ForceSub {
				return enableForceSub(b, ctx, fSub.ForceSubChannel)
			}
			text = "Force Sub is already enabled."
		case "disable", "off", "false", "n", "no":
			if fSub.ForceSub {
				if err := db.SetFSub(ctx.EffectiveChat.Id, false); err != nil {
					return logError("Error disabling Force Sub", err)
				}
				text = "Force Sub disabled."
			} else {
				text = "Force Sub is already disabled."
			}
		default:
			text = "Invalid argument. Use /fsub on or /fsub off."
		}
	}
	_, _ = ctx.EffectiveMessage.Reply(b, text, nil)
	return nil
}

func onOff(state bool) string {
	if state {
		return "enabled"
	}
	return "disabled"
}

func logError(message string, err error) error {
	config.ErrorLog.Printf("[ForceSub] %s: %v", message, err)
	return err
}
