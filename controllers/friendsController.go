package controllers

import (
	"DummyMessengerAPI/models"
	u "DummyMessengerAPI/utils"
	"encoding/json"
	"net/http"
)

var AddFriend = func(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(uint)
	friend := &models.Friend{}
	err := json.NewDecoder(r.Body).Decode(friend)
	if err != nil {
		u.Response(w, http.StatusForbidden, u.Message(false, "Incorrect request"))
		return
	}

	friend.OwnerId = userId
	friendRev := &models.Friend{OwnId: friend.OwnerId, OwnerId: friend.OwnId}
	friendRev.Add()
	u.Response(w, http.StatusOK, friend.Add())
}

var RemoveFriend = func(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(uint)
	friend := &models.Friend{}
	err := json.NewDecoder(r.Body).Decode(friend)
	if err != nil {
		u.Response(w, http.StatusForbidden, u.Message(false, "Incorrect request"))
		return
	}

	friend.OwnerId = userId
	friendRev := &models.Friend{OwnId: friend.OwnerId, OwnerId: friend.OwnId}
	friendRev.Remove()
	u.Response(w, http.StatusOK, friend.Remove())
}

var GetFriendsFor = func(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(uint)
	friends := models.GetFriends(userId)
	resp := u.Message(true, "Friend list generated")
	resp["friends"] = friends
	u.Response(w, http.StatusOK, resp)
}
