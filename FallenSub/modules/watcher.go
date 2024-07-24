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
	fSub := db.GetFSubSetting(chat.Id)

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
		go db.SetFSub(chat.Id, false)
		text := "Force Sub disabled because I can't get your chat member status. Please add me as an admin."
		_, _ = b.SendMessage(chat.Id, text, nil)
		config.ErrorLog.Printf("[fSubWatcher]Error getting chat member: %s [chatId: %d]", err, chat.Id)
		return err
	}

	if member.GetStatus() == "member" || member.GetStatus() == "administrator" || member.GetStatus() == "creator" {
		return ext.EndGroups
	}

	_, err = b.RestrictChatMember(chat.Id, user.Id, chatMutePermissions, &gotgbot.RestrictChatMemberOpts{UseIndependentChatPermissions: false})
	if err != nil {
		config.ErrorLog.Printf("[fSubWatcher]Error restricting user: %s [chatId: %d]", err, chat.Id)
		return err
	}

	db.UpdateMuted(chat.Id, user.Id)

	channel, err := b.GetChat(fSub.ForceSubChannel, nil)
	if err != nil {
		config.ErrorLog.Printf("[fSubWatcher]Error getting channel info: %s [chatId: %d]", err, chat.Id)
		return err
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
		config.ErrorLog.Printf("[fSubWatcher]Error sending message: %s [chatId: %d]", err, chat.Id)
	}

	return ext.EndGroups
}

// unMuteMe unMutes the user if they have joined the channel.
func unMuteMe(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.Update.CallbackQuery
	user := ctx.EffectiveUser
	chat := ctx.EffectiveChat

	if !db.IsMuted(chat.Id, user.Id) {
		_, err := query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "You are not muted by me.", ShowAlert: true})
		if err != nil {
			return err
		}
		return ext.EndGroups
	}
	fSub := db.GetFSubSetting(chat.Id)

	member, err := b.GetChatMember(fSub.ForceSubChannel, user.Id, nil)
	if err != nil {
		return err
	}

	stats := member.MergeChatMember()
	config.InfoLog.Printf("status: %s", stats.Status)

	if stats.Status != "member" && stats.Status != "administrator" && stats.Status != "creator" {
		_, err = query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "You are not a member of the channel.\nTap on Join Channel Button", ShowAlert: true})
		if err != nil {
			return err
		}
		return ext.EndGroups
	}

	c, err := b.GetChat(chat.Id, nil)
	if err != nil {
		config.ErrorLog.Printf("[unMuteMe]Error getting chat info: %s [chatId: %d]", err, chat.Id)
		return err
	}

	_, err = b.RestrictChatMember(chat.Id, user.Id, *c.Permissions, &gotgbot.RestrictChatMemberOpts{UseIndependentChatPermissions: true})
	if err != nil {
		config.ErrorLog.Printf("[unMuteMe]Error unrestricting user: %s [chatId: %d]", err, chat.Id)
		return err
	}

	db.RemoveMuted(chat.Id, user.Id)

	_, err = query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "You are unMuted now.", ShowAlert: true})
	if err != nil {
		return err
	}

	_, _, _ = query.Message.EditText(b, "You are unMuted now.", nil)

	return ext.EndGroups
}
