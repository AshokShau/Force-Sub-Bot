package modules

import (
	"fmt"
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

func isUserExempt(ctx *ext.Context, fSub db.FSub) bool {
	user := ctx.EffectiveUser
	return fSub.ForceSubChannel == 0 || user.IsBot || user.Id == 777000 || user.Id == 1087968824 || ctx.EffectiveSender.IsAnonymousAdmin()
}

func handleGetChatMemberError(b *gotgbot.Bot, chat *gotgbot.Chat, fSub db.FSub, err error) error {
	_ = db.SetFSub(chat.Id, false)
	text := "Force Sub disabled because I can't get your chat member status. Please add me as an admin."
	_, _ = b.SendMessage(chat.Id, text, nil)
	return logError(fmt.Sprintf("[fSubWatcher]Error getting chat member: %s [chatId: %d]", err, fSub.ForceSubChannel), err)
}

func restrictUser(b *gotgbot.Bot, chat *gotgbot.Chat, user *gotgbot.User) error {
	_, err := b.RestrictChatMember(chat.Id, user.Id, chatMutePermissions, &gotgbot.RestrictChatMemberOpts{UseIndependentChatPermissions: false})
	if err != nil {
		return logError(fmt.Sprintf("[fSubWatcher]Error restricting user: %s [chatId: %d]", err, chat.Id), err)
	}
	return db.UpdateMuted(chat.Id, user.Id)
}

func sendJoinChannelMessage(b *gotgbot.Bot, msg *gotgbot.Message, chat *gotgbot.Chat, user *gotgbot.User, inviteLink string) error {
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
	_, err := msg.Reply(b, text, &gotgbot.SendMessageOpts{ReplyMarkup: button, ReplyParameters: &gotgbot.ReplyParameters{AllowSendingWithoutReply: true}})
	if err != nil {
		return logError(fmt.Sprintf("[fSubWatcher]Error replying to message: %s [chatId: %d]", err, chat.Id), err)
	}
	return nil
}

func fSubWatcher(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	msg := ctx.EffectiveMessage
	user := ctx.EffectiveUser
	fSub, _ := db.GetFSubSetting(chat.Id)

	if !fSub.ForceSub || isUserExempt(ctx, *fSub) || isAdmin(ctx.EffectiveChat, ctx.EffectiveUser, b) {
		return ext.EndGroups
	}

	member, err := b.GetChatMember(fSub.ForceSubChannel, user.Id, nil)
	if err != nil {
		return handleGetChatMemberError(b, chat, *fSub, err)
	}

	if member.GetStatus() == "member" || member.GetStatus() == "administrator" || member.GetStatus() == "creator" {
		return ext.EndGroups
	}

	err = restrictUser(b, chat, user)
	if err != nil {
		return logError(fmt.Sprintf("[fSubWatcher]Error restricting user: %s [chatId: %d]", err, chat.Id), err)
	}

	channel, err := b.GetChat(fSub.ForceSubChannel, nil)
	if err != nil {
		text := fmt.Sprintf("Something went wrong. Looks like I am not admin In Your Fsub Channel: %d\nDo <code>/fsub off</code>\nError: %s", fSub.ForceSubChannel, err)
		_, _ = b.SendMessage(chat.Id, text, nil)
		return logError(fmt.Sprintf("[fSubWatcher]Error getting channel info: %s [chatId: %d]", err, chat.Id), err)
	}

	return sendJoinChannelMessage(b, msg, chat, user, channel.InviteLink)
}
