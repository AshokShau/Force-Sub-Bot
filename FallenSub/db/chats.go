package db

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type State struct {
	BotID int64   `json:"botId"`
	Chats []int64 `json:"chats"`
	Users []int64 `json:"users"`
}

func stateKey(botId int64) string {
	return "stats:" + strconv.FormatInt(botId, 10)
}

func GetStats(botId int64) (*State, error) {
	data, err := rdb.Get(ctx, stateKey(botId)).Result()
	if errors.Is(err, redis.Nil) {
		return &State{BotID: botId}, nil
	} else if err != nil {
		return nil, err
	}
	var state State
	if err = json.Unmarshal([]byte(data), &state); err != nil {
		return nil, err
	}
	return &state, nil
}

// setStats sets the stats for a bot
func setStats(state *State) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, stateKey(state.BotID), string(data), 0).Err()
}

// AddChat adds a chat to the stats
func AddChat(botId int64, chatId int64) error {
	state, err := GetStats(botId)
	if err != nil {
		return err
	}
	if !findInInt64Slice(state.Chats, chatId) {
		state.Chats = append(state.Chats, chatId)
		return setStats(state)
	}
	return nil
}

// AddUser adds a user to the stats
func AddUser(botId int64, userId int64) error {
	state, err := GetStats(botId)
	if err != nil {
		return err
	}
	if !findInInt64Slice(state.Users, userId) {
		state.Users = append(state.Users, userId)
		return setStats(state)
	}
	return nil
}

// AllChats gets all chats in the stats
func AllChats(botId int64) ([]int64, error) {
	state, err := GetStats(botId)
	if err != nil {
		return nil, err
	}
	return state.Chats, nil
}

// AllUsers gets all users in the stats
func AllUsers(botId int64) ([]int64, error) {
	state, err := GetStats(botId)
	if err != nil {
		return nil, err
	}
	return state.Users, nil
}

// GetAllBots gets all bots in the stats
func GetAllBots() ([]int64, error) {
	keys, _, err := rdb.Scan(ctx, 0, "stats:*", 0).Result()
	if err != nil {
		return nil, err
	}
	var bots []int64
	for _, key := range keys {
		botId, err := strconv.ParseInt(key[6:], 10, 64)
		if err != nil {
			return nil, err
		}
		bots = append(bots, botId)
	}
	return bots, nil
}
