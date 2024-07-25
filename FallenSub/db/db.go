package db

import (
	"context"
	"encoding/json"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/config"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func init() {

	opt, err := redis.ParseURL(config.DatabaseURI)
	if err != nil {
		log.Fatalf("[Database][Connect]: %v", err)
	}
	rdb = redis.NewClient(opt)

	if _, err = rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("[Database][Ping]: %v", err)
	}

	log.Println("[Database][Connect]: Connected to Redis")
}

type FSub struct {
	ChatId          int64   `json:"chatId"`
	ForceSub        bool    `json:"forceSub"`
	ForceSubChannel int64   `json:"forceSubChannel"`
	FSubMuted       []int64 `json:"mutedUsers"`
}

func redisKey(chatId int64) string {
	return "forceSub:" + strconv.FormatInt(chatId, 10)
}

// GetFSubSetting gets the FSub setting for a chat
func GetFSubSetting(chatId int64) (*FSub, error) {
	data, err := rdb.Get(ctx, redisKey(chatId)).Result()
	if err == redis.Nil {
		return &FSub{ChatId: chatId}, nil
	} else if err != nil {
		return nil, err
	}
	var fSub FSub
	if err := json.Unmarshal([]byte(data), &fSub); err != nil {
		return nil, err
	}
	return &fSub, nil
}

// SetFSubSetting sets the FSub setting for a chat
func SetFSubSetting(fSub *FSub) error {
	data, err := json.Marshal(fSub)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, redisKey(fSub.ChatId), string(data), 0).Err()
}

// UpdateMuted adds a user to the muted list
func UpdateMuted(chatId int64, userid int64) error {
	fSub, err := GetFSubSetting(chatId)
	if err != nil {
		return err
	}
	if !findInInt64Slice(fSub.FSubMuted, userid) {
		fSub.FSubMuted = append(fSub.FSubMuted, userid)
		return SetFSubSetting(fSub)
	}
	return nil
}

// RemoveMuted removes a user from the muted list
func RemoveMuted(chatId int64, userid int64) error {
	fSub, err := GetFSubSetting(chatId)
	if err != nil {
		return err
	}
	for i, v := range fSub.FSubMuted {
		if v == userid {
			fSub.FSubMuted = append(fSub.FSubMuted[:i], fSub.FSubMuted[i+1:]...)
			return SetFSubSetting(fSub)
		}
	}
	return nil
}

// SetFSub sets the FSub setting for a chat
func SetFSub(chatId int64, fSub bool) error {
	fSubUpdate, err := GetFSubSetting(chatId)
	if err != nil {
		return err
	}
	fSubUpdate.ForceSub = fSub
	return SetFSubSetting(fSubUpdate)
}

// SetFSubChannel sets the FSub setting for a chat
func SetFSubChannel(chatId int64, channel int64) error {
	fSubUpdate, err := GetFSubSetting(chatId)
	if err != nil {
		return err
	}
	fSubUpdate.ForceSubChannel = channel
	return SetFSubSetting(fSubUpdate)
}

// IsMuted checks if a user is muted
func IsMuted(chatId int64, userId int64) (bool, error) {
	fSub, err := GetFSubSetting(chatId)
	if err != nil {
		return false, err
	}
	return findInInt64Slice(fSub.FSubMuted, userId), nil
}

// findInInt64Slice checks if a value is in a slice or not
func findInInt64Slice(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
