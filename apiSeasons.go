package main

import (
	. "github.com/42minutes/42minutes-server-api/models"
	"github.com/codegangsta/martini-contrib/encoder"
	_ "github.com/garfunkel/go-tvdb"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
)

func ApiSeasonsGetAllByShow(r *http.Request, enc encoder.Encoder, store Store, parms martini.Params) (int, []byte) {
	id, err := strconv.Atoi(parms["showId"])
	if err != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode(err))
	} else {
		seasons, err := store.GetSeasonsOrRetrieveByShowId(id)
		if err != nil {
			return http.StatusBadRequest, encoder.Must(enc.Encode(err))
		}
		return http.StatusOK, encoder.Must(enc.Encode(seasons))
	}
}
