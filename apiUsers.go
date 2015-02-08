package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	. "github.com/42minutes/42minutes-server-api/models"
	"github.com/codegangsta/martini-contrib/encoder"
	_ "github.com/garfunkel/go-tvdb"
	"github.com/twinj/uuid"
	"log"
	"net/http"
	"os"
)

func sha1Password(pass string) string {
	salt := os.Getenv("SALT")
	passComb := pass + salt
	hash := sha1.New()
	hash.Write([]byte(passComb))
	shaPassword := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	return shaPassword
}

func ApiUsersRegister(r *http.Request, enc encoder.Encoder, store Store) (int, []byte) {
	if r.URL.Query().Get("email") != "" && r.URL.Query().Get("password") != "" {
		db := GetDbSession()
		user := User{ID: uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), Email: r.URL.Query().Get("email"), Password: sha1Password(r.URL.Query().Get("password")), Token: uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen)}
		err := db.Insert(&user)
		if err != nil {
			log.Println(err)
			return http.StatusBadRequest, encoder.Must(enc.Encode(err))
		}
		log.Printf("Registering new user with email %s, password: %s , userid:%s, token:%s", user.Email, user.Password, user.ID, user.Token)
		return http.StatusOK, encoder.Must(enc.Encode(user))
	}
	return http.StatusBadRequest, encoder.Must(enc.Encode("Missing email param"))
}

func ApiUsersLogin(r *http.Request, enc encoder.Encoder, store Store) (int, []byte) {
	email := r.FormValue("email")
	pass := r.FormValue("password")
	if email != "" && pass != "" {
		db := GetDbSession()
		user := User{}
		passHash := sha1Password(pass)
		err := db.SelectOne(&user, "select * from users where email=? and password=? ", email, passHash)
		if err == nil {
			user.Token = uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen)
			_, err := db.Update(&user)
			if err == nil {
				return http.StatusOK, encoder.Must(enc.Encode(user.Token))
			}
		} else {
			return http.StatusBadRequest, encoder.Must(enc.Encode("Wrong email or password"))
		}
	}
	return http.StatusBadRequest, encoder.Must(enc.Encode("Missing email or pass param"))
}
