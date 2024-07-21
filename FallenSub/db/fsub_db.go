package db

import (
	"errors"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FSub struct {
	ChatId          int64   `bson:"_id,omitempty" json:"_id,omitempty"`
	ForceSub        bool    `bson:"forceSub" json:"forceSub"`
	ForceSubChannel int64   `bson:"forceSubChannel" json:"forceSubChannel"`
	FSubMuted       []int64 `bson:"muted_users" json:"muted_users" default:"nil"`
}

func checkFSubSetting(chatId int64) (fSubScr *FSub) {
	err := findOne(fSubCell, bson.M{"_id": chatId}).Decode(&fSubScr)
	if errors.Is(err, mongo.ErrNoDocuments) {
		fSubScr = &FSub{ChatId: chatId, ForceSub: false, ForceSubChannel: 0, FSubMuted: []int64{}}
	} else if err != nil {
		config.ErrorLog.Printf("[Database][checkFSubSetting]: %v", err)
		return
	}
	return
}

func GetFSubSetting(chatId int64) *FSub {
	return checkFSubSetting(chatId)
}

func UpdateMuted(chatId int64, userid int64) {
	fSub := checkFSubSetting(chatId)
	foundUser := config.FindInInt64Slice(fSub.FSubMuted, userid)
	if foundUser {
		return // Already muted
	} else {
		fSub.FSubMuted = append(fSub.FSubMuted, userid)
		err := updateOne(fSubCell, bson.M{"_id": chatId}, fSub)
		if err != nil {
			config.ErrorLog.Printf("[Database][UpdateMuted]: %v", err)
			return
		}
	}
}

func RemoveMuted(chatId int64, userid int64) {
	fSub := checkFSubSetting(chatId)
	foundUser := config.FindInInt64Slice(fSub.FSubMuted, userid)
	if !foundUser {
		return // User not muted
	}

	fSub.FSubMuted = config.RemoveFromInt64Slice(fSub.FSubMuted, userid)

	err := updateOne(fSubCell, bson.M{"_id": chatId}, fSub)
	if err != nil {
		config.ErrorLog.Printf("[Database][RemoveMuted]: %v", err)
		return
	}
}

func IsMuted(chatId int64, userId int64) bool {
	fSub := checkFSubSetting(chatId)
	return config.FindInInt64Slice(fSub.FSubMuted, userId)
}

func SetFSub(chatId int64, fSub bool) {
	fSubUpdate := checkFSubSetting(chatId)
	fSubUpdate.ForceSub = fSub

	err := updateOne(fSubCell, bson.M{"_id": chatId}, fSubUpdate)
	if err != nil {
		config.ErrorLog.Printf("[Database][SetFSub]: %v", err)
	}
}

func SetFSubChannel(chatId int64, channel int64) {
	fSubUpdate := checkFSubSetting(chatId)
	fSubUpdate.ForceSubChannel = channel

	err := updateOne(fSubCell, bson.M{"_id": chatId}, fSubUpdate)
	if err != nil {
		config.ErrorLog.Printf("[Database][SetFSubChannel]: %v", err)
	}
}
