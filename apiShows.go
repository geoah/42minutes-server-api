package main

import (
	"fmt"
	. "github.com/42minutes/api/models"
	"github.com/codegangsta/martini-contrib/encoder"
	_ "github.com/garfunkel/go-tvdb"
	// "github.com/go-martini/martini"
	"log"
	"net/http"
	// "strconv"
)

func ApiShowsAll(r *http.Request, enc encoder.Encoder, store Store) (int, []byte) {
	if r.URL.Query().Get("name") == "" {
		show, err := store.GetAll()
		if err != nil {
			log.Println(err)
			return http.StatusNotFound, encoder.Must(enc.Encode(
				NewError(ErrCodeNotExist, fmt.Sprintf("TODO error"))))
		}
		return http.StatusOK, encoder.Must(enc.Encode(show))
	} else {
		name := r.URL.Query().Get("name")
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

// func ApiShows(enc encoder.Encoder, store Store, parms martini.Params) (int, []byte) {
// 	// TODO Check for parms:id
// 	// Get payload Object from Store
// 	id, err := strconv.ParseUint(parms["id"], 10, 64)
// 	if err != nil {
// 		return http.StatusBadRequest, encoder.Must(enc.Encode(err))
// 	} else {
// 		Show := store.Get(id)
// 		// Show.fetchDetails()
// 		// if Show.Matched == true {
// 		// 	Show.CheckForExistingEpisodes()
// 		// 	Show.FetchTorrentLinks()
// 		// 	// Show.PrintResults()
// 		// 	// Show.PrintJsonResults()
// 		// }
// 		// // TODO Check if payload exists
// 		return http.StatusOK, encoder.Must(enc.Encode(Show))
// 	}
// }
