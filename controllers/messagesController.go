package controllers

import (
	"DummyMessengerAPI/models"
	u "DummyMessengerAPI/utils"
	"encoding/json"
	"net/http"
)

var SendMessage = func(w http.ResponseWriter, r *http.Request) {
	message := &models.Message{}
	userId := r.Context().Value("user_id").(uint)
	err := json.NewDecoder(r.Body).Decode(message)
	if err != nil {
		u.Response(w, http.StatusForbidden, u.Message(false, "Incorrect request"))
		return
	}

	message.SenderId = userId

	if !models.AreFriends(message.SenderId, message.ReceiverId) {
		friend := &models.Friend{OwnId: message.ReceiverId, OwnerId: message.SenderId}
		friendRev := &models.Friend{OwnId: friend.OwnerId, OwnerId: friend.OwnId}
		friend.Add()
		friendRev.Add()
	}

	u.Response(w, http.StatusOK, message.Create())
}

var GetLastMessage = func(w http.ResponseWriter, r *http.Request) {
	message := &models.Message{}
	userId := r.Context().Value("user_id").(uint)
	err := json.NewDecoder(r.Body).Decode(message)
	if err != nil {
		u.Response(w, http.StatusForbidden, u.Message(false, "Incorrect request"))
		return
	}

	message = models.GetLastMessage(userId, message.ReceiverId)

	resp := u.Message(true, "Last message generated")
	resp["message"] = message
	u.Response(w, http.StatusOK, resp)
}

var GetMessages = func(w http.ResponseWriter, r *http.Request) {
	message := &models.Message{}
	userId := r.Context().Value("user_id").(uint)
	err := json.NewDecoder(r.Body).Decode(message)
	if err != nil {
		u.Response(w, http.StatusForbidden, u.Message(false, "Incorrect request"))
		return
	}

	messages := models.GetMessages(userId, message.ReceiverId)

	resp := u.Message(true, "Messages list generated")
	resp["messages"] = messages
	u.Response(w, http.StatusOK, resp)
}
