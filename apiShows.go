package main

import (
	"fmt"
	. "github.com/42minutes/42minutes-server-api/models"
	"github.com/codegangsta/martini-contrib/encoder"
	_ "github.com/garfunkel/go-tvdb"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"strconv"
)

func ApiShowsGetAll(r *http.Request, enc encoder.Encoder, store Store) (int, []byte) {
	if r.URL.Query().Get("name") == "" {
		show, err := store.GetShows()
		if err != nil {
			log.Println(err)
			return http.StatusNotFound, encoder.Must(enc.Encode(
				NewError(ErrCodeNotExist, fmt.Sprintf("TODO error"))))
		}
		return http.StatusOK, encoder.Must(enc.Encode(show))
	} else {
		name := r.URL.Query().Get("name")
		fmt.Println("Trying to look for show.", name)
		show, err := ShowFindAllByName(name, 5)
		if err != nil {
			log.Println(err)
			return http.StatusNotFound, encoder.Must(enc.Encode(
				NewError(ErrCodeNotExist, fmt.Sprintf("TODO error"))))
		}
		if len(show) == 0 {
			return http.StatusNotFound, encoder.Must(enc.Encode(
				NewError(ErrCodeNotExist, fmt.Sprintf("could not find Show with name '%s'", name))))
		}
		return http.StatusOK, encoder.Must(enc.Encode(show))
	}
}

func ApiShowsGetOne(r *http.Request, enc encoder.Encoder, store Store, parms martini.Params) (int, []byte) {
	id, err := strconv.Atoi(parms["id"])
	if err != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode(err))
	} else {
		show, err := store.GetShowOrRetrieve(id)
		if err != nil {
			return http.StatusBadRequest, encoder.Must(enc.Encode(err))
		}
		return http.StatusOK, encoder.Must(enc.Encode(show))
	}
}
