package models

import (
	u "DummyMessengerAPI/utils"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Friend struct {
	gorm.Model
	OwnerId uint `json:"owner_id"`
	OwnId   uint `json:"own_id"`
}

func (friend *Friend) Validate() (map[string]interface{}, bool) {
	if friend.OwnId <= 0 {
		return u.Message(false, "Non-existing user"), false
	}
	if friend.OwnerId <= 0 {
		return u.Message(false, "User is not recognized"), false
	}
	if friend.OwnerId == friend.OwnId {
		return u.Message(false, "Attempt to add yourself as user"), false
	}
	if AreFriends(friend.OwnId, friend.OwnerId) {
		return u.Message(false, "Users are already friends"), false
	}

	cnt := 0
	err := GetDB().Table("users").Where("id = ?", friend.OwnId).Count(&cnt).Error
	if err != nil {
		fmt.Println(err)
	}
	if cnt == 0 {
		return u.Message(false, "Attempt to add nonexistent user"), false
	}

	return u.Message(true, "Requirement passed"), true
}

func (friend *Friend) Add() map[string]interface{} {
	if resp, ok := friend.Validate(); !ok {
		return resp
	}
	GetDB().Create(friend)
	return u.Message(true, "Friend has been added")
}

func (friend *Friend) Remove() map[string]interface{} {
	GetDB().Delete(friend)
	return u.Message(true, "Friend has been removed")
}

func AreFriends(ownId, ownerId uint) bool {
	cnt := 0
	GetDB().Table("friends").Where("owner_id = ? and own_id = ? and deleted_at is null", ownerId, ownId).Count(&cnt)
	return cnt > 0
}

func GetFriends(userId uint) []*User {
	friendsId := make([]*Friend, 0)
	err := GetDB().Table("friends").Where("owner_id = ?", userId).Find(&friendsId).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	friends := make([]*User, len(friendsId))
	for i, curFriend := range friendsId {
		tmp := getUser(curFriend.OwnId)
		friends[i] = tmp
	}
	return friends
}
