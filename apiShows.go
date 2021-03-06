package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	. "github.com/42minutes/42minutes-server-api/models"
	"github.com/codegangsta/martini-contrib/encoder"
	_ "github.com/garfunkel/go-tvdb"
	"github.com/go-martini/martini"
)

func ApiShowsGetAll(r *http.Request, enc encoder.Encoder, store Store, user User) (int, []byte) {
	db := GetDbSession()

	var shows []*Show
	if r.URL.Query().Get("name") == "" {
		_, err := db.Select(&shows, "SELECT shows.* FROM shows LEFT JOIN users_shows ON shows.id = users_shows.show_id WHERE users_shows.user_id = ? AND (users_shows.favorite = true OR users_shows.library = true) ORDER BY shows.title asc", user.ID)
		if err != nil {
			return http.StatusNotFound, encoder.Must(enc.Encode(NewError(ErrCodeNotExist, fmt.Sprintf("TODO error"))))
		}
	} else {
		name := r.URL.Query().Get("name")
		fmt.Println("Looking for show...", name)
		shows, _ = ShowFindAllByName(name, 5)
		if len(shows) == 0 {
			return http.StatusNotFound, encoder.Must(enc.Encode(NewError(ErrCodeNotExist, fmt.Sprintf("could not find Show with name '%s'", name))))
		}
	}

	for show_i, _ := range shows {
		shows[show_i].Personalize(user.ID)
	}

	return http.StatusOK, encoder.Must(enc.Encode(shows))
}

func ApiShowsGetOne(r *http.Request, enc encoder.Encoder, store Store, parms martini.Params, user User) (int, []byte) {

	id, err := strconv.Atoi(parms["id"])
	if err != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode(err))
	}

	show, err := store.GetShowOrRetrieve(id)
	if err != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode(err))
	}
	show.Personalize(user.ID)
	return http.StatusOK, encoder.Must(enc.Encode(show))

}

func ApiShowsPutOne(r *http.Request, enc encoder.Encoder, store Store, parms martini.Params, user User) (int, []byte) {
	id, err := strconv.Atoi(parms["id"])
	if err != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode(err))
	}
	show, err := store.GetShowOrRetrieve(id)
	if err != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode(err))
	}

	var showPost Show
	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &showPost)
	if err != nil {
		return http.StatusNotFound, encoder.Must(enc.Encode(NewError(ErrCodeNotExist, "Could not decode body")))
	}

	userShow := UserShow{UserID: user.ID, ShowID: show.ID, Favorite: showPost.Favorite, Library: showPost.Library}
	err = store.UserShowUpsert(&userShow)

	show.Personalize(user.ID)

	return http.StatusOK, encoder.Must(enc.Encode(show))
}
