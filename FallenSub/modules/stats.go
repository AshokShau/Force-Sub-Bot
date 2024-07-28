package modules

import (
	"fmt"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/config"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/db"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// getChats is the handler for the /getChats command
func getChats(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	if msg.From.Id != config.OwnerId {
		_, _ = msg.Reply(b, "You are not allowed to use this command.", nil)
		return ext.EndGroups
	}

	// Get all chats
	chats, err := db.AllChats(b.Id)
	if err != nil {
		return logError("Error getting chats", err)
	}

	// Create a message with the list of chats
	var chatList string
	for _, chat := range chats {
		chatList += fmt.Sprintf("» %d\n", chat)
	}

	// Send the list of chats as a message
	_, err = b.SendMessage(msg.Chat.Id, fmt.Sprintf("Total chats: %d\n\n%s", len(chats), chatList), nil)
	if err != nil {
		return logError("Error sending message", err)
	}

	return ext.EndGroups
}

// getUsers is the handler for the /getUsers command
func getUsers(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	if msg.From.Id != config.OwnerId {
		_, _ = msg.Reply(b, "You are not allowed to use this command.", nil)
		return ext.EndGroups
	}

	// Get all users
	users, err := db.AllUsers(b.Id)
	if err != nil {
		return logError("Error getting users", err)
	}

	// Create a message with the list of users
	var userList string
	for _, user := range users {
		userList += fmt.Sprintf("» %d\n", user)
	}

	// Send the list of users as a message
	_, err = b.SendMessage(msg.Chat.Id, fmt.Sprintf("Total users: %d\n\n%s", len(users), userList), nil)
	if err != nil {
		return logError("Error sending message", err)
	}

	return ext.EndGroups
}

// getAllBots is the handler for the /getAllBots command
func getAllBots(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	if msg.From.Id != config.OwnerId {
		_, _ = msg.Reply(b, "You are not allowed to use this command.", nil)
		return ext.EndGroups
	}

	// Get all bots
	bots, err := db.GetAllBots()
	if err != nil {
		return logError("Error getting bots", err)
	}

	// Create a message with the list of bots
	var botList string
	for _, bot := range bots {
		botList += fmt.Sprintf("» %d\n", bot)
	}

	// Send the list of bots as a message
	_, err = b.SendMessage(msg.Chat.Id, fmt.Sprintf("Total bots: %d\n\n%s", len(bots), botList), nil)
	if err != nil {
		return logError("Error sending message", err)
	}

	return ext.EndGroups
}
