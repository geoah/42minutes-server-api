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

func ApiSeriesAll(r *http.Request, enc encoder.Encoder, store Store) (int, []byte) {
	if r.URL.Query().Get("name") == "" {
		series := store.GetAll()
		return http.StatusOK, encoder.Must(enc.Encode(series))
	} else {
		name := r.URL.Query().Get("name")
		series := store.FindAllByName(name, 10)
		if len(series) == 0 {
			return http.StatusNotFound, encoder.Must(enc.Encode(
				NewError(ErrCodeNotExist, fmt.Sprintf("could not find series with name '%s'", name))))
		}
		return http.StatusOK, encoder.Must(enc.Encode(series))
	}
}

func ApiSeries(enc encoder.Encoder, store Store, parms martini.Params) (int, []byte) {
	// TODO Check for parms:id
	// Get payload Object from Store
	id, err := strconv.ParseUint(parms["id"], 10, 64)
	if err != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode(err))
	} else {
		series := store.Get(id)
		// series.fetchDetails()
		// if series.Matched == true {
		// 	series.CheckForExistingEpisodes()
		// 	series.FetchTorrentLinks()
		// 	// series.PrintResults()
		// 	// series.PrintJsonResults()
		// }
		// // TODO Check if payload exists
		return http.StatusOK, encoder.Must(enc.Encode(series))
	}
}
