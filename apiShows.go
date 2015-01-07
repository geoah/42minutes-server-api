package main

import (
	"fmt"
	. "github.com/42minutes/api/stores"
	"github.com/codegangsta/martini-contrib/encoder"
	_ "github.com/garfunkel/go-tvdb"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
)

func ApiShowsAll(r *http.Request, enc encoder.Encoder, store Store) (int, []byte) {
	if r.URL.Query().Get("name") == "" {
		Show := store.GetAll()
		return http.StatusOK, encoder.Must(enc.Encode(Show))
	} else {
		name := r.URL.Query().Get("name")
		Show := store.FindAllByName(name, 10)
		if len(Show) == 0 {
			return http.StatusNotFound, encoder.Must(enc.Encode(
				NewError(ErrCodeNotExist, fmt.Sprintf("could not find Show with name '%s'", name))))
		}
		return http.StatusOK, encoder.Must(enc.Encode(Show))
	}
}

func ApiShows(enc encoder.Encoder, store Store, parms martini.Params) (int, []byte) {
	// TODO Check for parms:id
	// Get payload Object from Store
	id, err := strconv.ParseUint(parms["id"], 10, 64)
	if err != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode(err))
	} else {
		Show := store.Get(id)
		// Show.fetchDetails()
		// if Show.Matched == true {
		// 	Show.CheckForExistingEpisodes()
		// 	Show.FetchTorrentLinks()
		// 	// Show.PrintResults()
		// 	// Show.PrintJsonResults()
		// }
		// // TODO Check if payload exists
		return http.StatusOK, encoder.Must(enc.Encode(Show))
	}
}
