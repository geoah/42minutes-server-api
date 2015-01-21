package main

import (
	"net/http"
	"strconv"

	. "github.com/42minutes/42minutes-server-api/models"
	"github.com/codegangsta/martini-contrib/encoder"
	_ "github.com/garfunkel/go-tvdb"
	"github.com/go-martini/martini"
)

func ApiEpisodesGetAllByShowAndSeason(r *http.Request, enc encoder.Encoder, store Store, parms martini.Params) (int, []byte) {
	showId, errShow := strconv.Atoi(parms["showId"])
	seasonNumber, errSeason := strconv.Atoi(parms["seasonNumber"])
	if errShow != nil || errSeason != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode("Missing show_id or season"))
	} else {
		seasons, err := store.GetEpisodesOrRetrieveByShowIdAndSeason(showId, seasonNumber)
		if err != nil {
			return http.StatusBadRequest, encoder.Must(enc.Encode(err))
		}
		return http.StatusOK, encoder.Must(enc.Encode(seasons))
	}
}

func ApiEpisodesGetOneByShowAndSeasonAndEpisode(r *http.Request, enc encoder.Encoder, store Store, parms martini.Params) (int, []byte) {
	showId, errShow := strconv.Atoi(parms["showId"])
	seasonNumber, errSeason := strconv.Atoi(parms["seasonNumber"])
	episodeNumber, errEpisode := strconv.Atoi(parms["episodeNumber"])
	if errShow != nil || errSeason != nil || errEpisode != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode("Missing show_id or season or episode"))
	} else {
		episode, err := store.GetEpisodeOrRetrieveByShowIdAndSeasonAndEpisode(showId, seasonNumber, episodeNumber)
		if err != nil || episode == nil {
			return http.StatusBadRequest, encoder.Must(enc.Encode(err))
		}
		return http.StatusOK, encoder.Must(enc.Encode(episode))
	}
}
