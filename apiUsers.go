package main

import (
	. "github.com/42minutes/42minutes-server-api/models"
	"github.com/codegangsta/martini-contrib/encoder"
	_ "github.com/garfunkel/go-tvdb"
	"github.com/twinj/uuid"
	"log"
	"net/http"
)

func ApiUsersRegister(r *http.Request, enc encoder.Encoder, store Store) (int, []byte) {
	if r.URL.Query().Get("email") != "" {
		db := GetDbSession()
		user := User{ID: uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), Email: r.URL.Query().Get("email"), Token: uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen)}
		err := db.Insert(&user)
		if err != nil {
			log.Println(err)
			return http.StatusBadRequest, encoder.Must(enc.Encode(err))
		}
		log.Printf("Registering new user with email %s, userid:%s, token:%s", user.Email, user.ID, user.Token)
		return http.StatusOK, encoder.Must(enc.Encode(user))
	}
	return http.StatusBadRequest, encoder.Must(enc.Encode("Missing email param"))
}
