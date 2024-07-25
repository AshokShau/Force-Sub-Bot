package modules

import (
	"fmt"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/config"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/db"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var chatMutePermissions = gotgbot.ChatPermissions{
	CanSendMessages:      false,
	CanSendAudios:        false,
	CanSendDocuments:     false,
	CanSendOtherMessages: false,
	CanSendPolls:         false,
	CanSendPhotos:        false,
	CanSendVideos:        false,
	CanSendVideoNotes:    false,
	CanSendVoiceNotes:    false,
}

// fSubWatcher checks if the user is a member of the channel and restricts them if not.
func fSubWatcher(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	msg := ctx.EffectiveMessage

	user := ctx.EffectiveUser
	fSub, _ := db.GetFSubSetting(chat.Id)

	if !fSub.ForceSub {
		return ext.EndGroups
	}

	if fSub.ForceSubChannel == 0 || user.IsBot || user.Id == 777000 || user.Id == 1087968824 || ctx.EffectiveSender.IsAnonymousAdmin() {
		return ext.EndGroups
	}

	if isAdmin(ctx.EffectiveChat, ctx.EffectiveUser, b) {
		return ext.EndGroups
	}

	member, err := b.GetChatMember(fSub.ForceSubChannel, user.Id, nil)
	if err != nil {
		_ = db.SetFSub(chat.Id, false)
		text := "Force Sub disabled because I can't get your chat member status. Please add me as an admin."
		_, _ = b.SendMessage(chat.Id, text, nil)
		return logError(fmt.Sprintf("[fSubWatcher]Error getting chat member: %s [chatId: %d]", err, fSub.ForceSubChannel), err)
	}

	if member.GetStatus() == "member" || member.GetStatus() == "administrator" || member.GetStatus() == "creator" {
		return ext.EndGroups
	}

	_, err = b.RestrictChatMember(chat.Id, user.Id, chatMutePermissions, &gotgbot.RestrictChatMemberOpts{UseIndependentChatPermissions: false})
	if err != nil {
		return logError(fmt.Sprintf("[fSubWatcher]Error restricting user: %s [chatId: %d]", err, chat.Id), err)
	}

	err = db.UpdateMuted(chat.Id, user.Id)
	if err != nil {
		return logError(fmt.Sprintf("[fSubWatcher]Error updating muted user: %s [chatId: %d]", err, chat.Id), err)
	}

	channel, err := b.GetChat(fSub.ForceSubChannel, nil)
	if err != nil {
		return logError(fmt.Sprintf("[fSubWatcher]Error getting channel info: %s [chatId: %d]", err, chat.Id), err)
	}

	inviteLink := channel.InviteLink
	text := "You must join the channel to continue using this group."
	button := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
			{
				Text: "Join Channel",
				Url:  inviteLink,
			},
			{
				Text:         "Unmute Me",
				CallbackData: fmt.Sprintf("unmuteMe_%d", user.Id),
			},
		}},
	}
	_, err = msg.Reply(b, text, &gotgbot.SendMessageOpts{ReplyMarkup: button, ReplyParameters: &gotgbot.ReplyParameters{AllowSendingWithoutReply: true}})
	if err != nil {
		return logError(fmt.Sprintf("[fSubWatcher]Error replying to message: %s [chatId: %d]", err, chat.Id), err)
	}

	return ext.EndGroups
}

// unMuteMe unMutes the user if they have joined the channel.
func unMuteMe(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.Update.CallbackQuery
	user := ctx.EffectiveUser
	chat := ctx.EffectiveChat

	isMuted, err := db.IsMuted(chat.Id, user.Id)
	if !isMuted {
		_, err = query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "You are not muted by me.", ShowAlert: true})
		return err
	}

	fSub, err := db.GetFSubSetting(chat.Id)
	if err != nil {
		return logError(fmt.Sprintf("[unMuteMe]Error getting fSub setting: %s [chatId: %d]", err, chat.Id), err)
	}

	member, err := b.GetChatMember(fSub.ForceSubChannel, user.Id, nil)
	if err != nil {
		return err
	}

	stats := member.MergeChatMember()
	config.InfoLog.Printf("status: %s", stats.Status)

	if stats.Status != "member" && stats.Status != "administrator" && stats.Status != "creator" {
		_, err = query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "You are not a member of the channel.\nTap on Join Channel Button", ShowAlert: true})
		return err
	}

	c, err := b.GetChat(chat.Id, nil)
	if err != nil {
		return logError(fmt.Sprintf("[unMuteMe]Error getting chat info: %s [chatId: %d]", err, chat.Id), err)
	}

	_, err = b.RestrictChatMember(chat.Id, user.Id, *c.Permissions, &gotgbot.RestrictChatMemberOpts{UseIndependentChatPermissions: true})
	if err != nil {
		return logError(fmt.Sprintf("[unMuteMe]Error restricting user: %s [chatId: %d]", err, chat.Id), err)
	}

	err = db.RemoveMuted(chat.Id, user.Id)
	if err != nil {
		return logError(fmt.Sprintf("[unMuteMe]Error removing muted user: %s [chatId: %d]", err, chat.Id), err)
	}

	_, err = query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "You are unMuted now.", ShowAlert: true})
	if err != nil {
		return logError(fmt.Sprintf("[unMuteMe]Error answering callback: %s [chatId: %d]", err, chat.Id), err)
	}

	_, _, _ = query.Message.EditText(b, "You are now unMuted and can participate in the chat again.", nil)

	return ext.EndGroups
}
