package models

import (
	u "DummyMessengerAPI/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type Token struct {
	UserId         uint
	ExpirationTime int64
	jwt.StandardClaims
}

type User struct {
	gorm.Model
	Login    string `json:"login"`
	Password string `json:"password"`
	Token    string `sql:"-";json:"token"`
}

func (user *User) Validate() (map[string]interface{}, bool) {
	tmp := &User{}
	err := GetDB().Table("users").Where("login = ?", user.Login).First(tmp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}

	if tmp.Login != "" {
		return u.Message(false, "Login already exists."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (user *User) Create() map[string]interface{} {
	if resp, ok := user.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	GetDB().Create(user)

	if user.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}
	return u.Message(true, "Account has been created")
}

func Login(login, password string) map[string]interface{} {
	user := &User{}
	err := GetDB().Table("users").Where("login = ?", login).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Login not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid login credentials. Please try again")
	}

	tk := &Token{UserId: user.ID, ExpirationTime: time.Now().Add(6 * time.Hour).Unix()}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

	resp := u.Message(true, "Logged In")
	resp["token"] = tokenString
	return resp
}

func Logout(token string) map[string]interface{} {
	GetDB().Exec("INSERT INTO tokens_blacklist(token) VALUES(\"" + token + "\")")
	return u.Message(true, "Logged Out")
}

func getUser(id uint) *User {
	user := &User{}
	err := GetDB().Table("users").Where("id = ?", id).First(user).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}

	user.Password = ""
	return user
}
