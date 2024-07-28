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

	go func() {
		_ = db.AddChat(b.Id, ctx.EffectiveChat.Id)
	}()

	// If the message is a reply to a forwarded message, handle it
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
		if fSub.ForceSubChannel == 0 {
			text = fmt.Sprintf("Force Sub is %s in %s. \nNo channel set.", onOff(fSub.ForceSub), ctx.EffectiveChat.Title)
		}

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
				if err = db.SetFSub(ctx.EffectiveChat.Id, false); err != nil {
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

// unMuteMe unMutes the user if they have joined the channel.
func unMuteMe(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.Update.CallbackQuery
	user := ctx.EffectiveUser
	chat := ctx.EffectiveChat

	isMuted, err := db.IsMuted(chat.Id, user.Id)
	if !isMuted {
		_, _ = query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "You are not muted by me.", ShowAlert: true})
		return err
	}

	fSub, err := db.GetFSubSetting(chat.Id)
	if err != nil {
		return logError(fmt.Sprintf("[unMuteMe]Error getting fSub setting: %s [chatId: %d]", err, chat.Id), err)
	}

	member, err := b.GetChatMember(fSub.ForceSubChannel, user.Id, nil)
	if err != nil {
		return logError(fmt.Sprintf("[unMuteMe]Error getting chat member: %s [chatId: %d]", err, chat.Id), err)
	}

	stats := member.MergeChatMember()

	if stats.Status != "member" && stats.Status != "administrator" && stats.Status != "creator" {
		_, _ = query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "You are not a member of the channel.\nTap on Join Channel Button", ShowAlert: true})
		return err
	}

	c, err := b.GetChat(chat.Id, nil)
	if err != nil {
		_, _ = query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "Error getting chat info.", ShowAlert: true})
		return logError(fmt.Sprintf("[unMuteMe]Error getting chat info: %s [chatId: %d]", err, chat.Id), err)
	}

	_, err = b.RestrictChatMember(chat.Id, user.Id, *c.Permissions, &gotgbot.RestrictChatMemberOpts{UseIndependentChatPermissions: true})
	if err != nil {
		_, _ = query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "Error unMuting you.\nmaybe i am not admin with enough rights.", ShowAlert: true})
		return logError(fmt.Sprintf("[unMuteMe]Error restricting user: %s [chatId: %d]", err, chat.Id), err)
	}

	_ = db.RemoveMuted(chat.Id, user.Id)

	_, _ = query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "You are unMuted now.", ShowAlert: true})
	_, _, _ = query.Message.EditText(b, "You are now unMuted and can participate in the chat again.", nil)

	return ext.EndGroups
}
