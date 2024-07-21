package modules

import (
	"fmt"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/config"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/db"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// SetFSub enables or disables the force sub setting in a group.
func setFSub(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	msg := ctx.EffectiveMessage
	user := ctx.EffectiveUser

	args := ctx.Args()
	var text string

	if chat.Type == gotgbot.ChatTypePrivate {
		_, _ = msg.Reply(b, "This command is only available in groups.", &gotgbot.SendMessageOpts{})
		return ext.EndGroups
	}

	userMember, _ := chat.GetMember(b, user.Id, nil)
	mem := userMember.MergeChatMember()

	if mem.Status == "member" && config.OwnerId != user.Id {
		_, _ = msg.Reply(b, "You must be an admin to use this command.", &gotgbot.SendMessageOpts{})
		return ext.EndGroups
	}

	fSub := db.GetFSubSetting(chat.Id)

	repliedMsg := msg.ReplyToMessage

	if repliedMsg != nil && repliedMsg.ForwardOrigin != nil {
		msgOrigen := msg.ReplyToMessage.ForwardOrigin.MergeMessageOrigin()
		if msgOrigen.Chat.Type == gotgbot.ChatTypeChannel {
			_, err := b.GetChatMember(msgOrigen.Chat.Id, b.Id, nil)
			if err != nil {
				_, _ = msg.Reply(b, fmt.Sprintf("Looks like I am not admin in %s [%d]", msgOrigen.Chat.Title, msgOrigen.Chat.Id), &gotgbot.SendMessageOpts{})
				config.ErrorLog.Printf("[setFSub]Error getting chat member: %s", err)
				return err
			}

			go db.SetFSubChannel(chat.Id, msgOrigen.Chat.Id)
			text = fmt.Sprintf("Force Sub enabled in %s.\nChannelID: %d", chat.Title, msgOrigen.Chat.Id)
		} else {
			text = "Reply to a forwarded message from a channel."
		}

		_, err := msg.Reply(b, text, &gotgbot.SendMessageOpts{})
		if err != nil {
			config.ErrorLog.Printf("[setFSub]Error sending message: %s", err)
			return err
		}

		return ext.EndGroups
	}

	if len(args) == 1 {
		if fSub.ForceSub {
			text = fmt.Sprintf("Force Sub is enabled in %s.\nChannelID: %d", chat.Title, fSub.ForceSubChannel)
		} else {
			text = fmt.Sprintf("Force Sub is disabled in %s.", chat.Title)
		}
	} else {
		switch args[1] {
		case "enable", "on", "true", "y", "yes":
			if fSub.ForceSub {
				text = "Force Sub is already enabled."
			} else {
				go db.SetFSub(chat.Id, true)
				text = "Force Sub enabled."
			}
		case "disable", "off", "false", "n", "no":
			if !fSub.ForceSub {
				text = "Force Sub is already disabled."
			} else {
				go db.SetFSub(chat.Id, false)
				text = "Force Sub disabled."
			}
		case "unmuteall", "clear":
			return unMuteAll(b, ctx)
		default:
			text = fmt.Sprintf("Invalid argument. Use /fsub on or /fsub off; got %s\nelse reply to a forwarded msg from a channel /fsub.", args[1])
		}
	}

	_, err := msg.Reply(b, text, &gotgbot.SendMessageOpts{})
	if err != nil {
		config.ErrorLog.Printf("[setFSub]Error sending message: %s", err)
		return err
	}
	return ext.EndGroups
}

func unMuteAll(_ *gotgbot.Bot, _ *ext.Context) error {
	// Todo: implement unmute all

	return ext.EndGroups
}
