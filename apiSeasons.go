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
	var seasons []Season
	db := GetDbSession()
	id, err := strconv.Atoi(parms["showId"])
	if err != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode(err))
	} else {
		seasons, err := db.Select(&seasons, "select * from seasons where show_id=?", id)
		if err != nil {
			return http.StatusBadRequest, encoder.Must(enc.Encode(err))
		}
		return http.StatusOK, encoder.Must(enc.Encode(seasons))
	}
}
