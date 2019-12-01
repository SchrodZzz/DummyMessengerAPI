package main

import (
	"DummyMessengerAPI/app"
	"DummyMessengerAPI/controllers"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/user/new", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authorize).Methods("POST")
	router.HandleFunc("/api/user/logout", controllers.Logout).Methods("POST")

	//TODO: fix Friend.Remove func (deletes all the data)
	router.HandleFunc("/api/me/friends/add", controllers.AddFriend).Methods("POST")
	//router.HandleFunc("/api/me/friends/remove", controllers.RemoveFriend).Methods("POST")
	router.HandleFunc("/api/me/friends/get/all", controllers.GetFriendsFor).Methods("GET")

	router.HandleFunc("/api/me/messages/send", controllers.SendMessage).Methods("POST")
	router.HandleFunc("/api/me/messages/get/all", controllers.GetMessages).Methods("POST")
	router.HandleFunc("/api/me/messages/get/last", controllers.GetLastMessage).Methods("POST")

	router.Use(app.JwtAuthorization)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Println(err)
	}
}
