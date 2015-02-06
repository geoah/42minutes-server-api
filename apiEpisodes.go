package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	. "github.com/42minutes/42minutes-server-api/models"
	"github.com/42minutes/go-torrentlookup"
	"github.com/codegangsta/martini-contrib/encoder"
	_ "github.com/garfunkel/go-tvdb"
	"github.com/go-martini/martini"
	"github.com/gorilla/feeds"
)

func ApiEpisodesGetAllByShowAndSeason(r *http.Request, enc encoder.Encoder, store Store, parms martini.Params) (int, []byte) {
	showId, errShow := strconv.Atoi(parms["showId"])
	seasonNumber, errSeason := strconv.Atoi(parms["seasonNumber"])
	if errShow != nil || errSeason != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode("Missing show_id or season"))
	} else {
		episodes, err := store.GetEpisodesOrRetrieveByShowIdAndSeason(showId, seasonNumber)
		if err != nil {
			return http.StatusBadRequest, encoder.Must(enc.Encode(err))
		}
		for _, episode := range episodes {
			episode.Normalize()
		}
		return http.StatusOK, encoder.Must(enc.Encode(episodes))
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
		episode.Normalize()
		return http.StatusOK, encoder.Must(enc.Encode(episode))
	}
}

func ApiRss(res http.ResponseWriter, w http.ResponseWriter, r *http.Request, enc encoder.Encoder, store Store) (int, []byte) {
	db := GetDbSession()

	if r.URL.Query().Get("token") == "" {
		return http.StatusUnauthorized, encoder.Must(enc.Encode(NewError(ErrCodeNotExist, "Error")))
	}

	token := r.URL.Query().Get("token")
	user := User{}
	err := db.SelectOne(&user, "select * from users where token=?", token)
	if err != nil {
		return http.StatusUnauthorized, encoder.Must(enc.Encode(NewError(ErrCodeNotExist, "Error")))
	}

	now := time.Now()
	feed := &feeds.Feed{
		Title:       "42minutes",
		Link:        &feeds.Link{Href: "http://42minutes.tv"},
		Description: "Tv Shows Etc.",
		Created:     now,
	}
	feed.Items = []*feeds.Item{}

	episodesRss := []*EpisodeRss{}
	db.Select(&episodesRss, "SELECT shows.title AS show_title, episodes.title, episodes.season, episodes.episode, episodes.first_aired, episodes.infohash_hd720p, episodes.infohash_sd480p FROM episodes LEFT JOIN shows ON episodes.show_id = shows.id LEFT JOIN users_shows ON shows.id = users_shows.show_id WHERE users_shows.library = true") // AND users_shows.user_id = ""

	for _, episodeRss := range episodesRss {
		magnet := ""
		if episodeRss.InfohashHd != "" {
			magnet = torrentlookup.FakeMagnet(episodeRss.InfohashHd)
		} else if episodeRss.InfohashSd != "" {
			magnet = torrentlookup.FakeMagnet(episodeRss.InfohashSd)
		} else {
			continue
		}
		item := feeds.Item{
			Title:   fmt.Sprintf("%s S%02dE%02d", episodeRss.ShowTitle, episodeRss.Season, episodeRss.Episode),
			Link:    &feeds.Link{Href: magnet},
			Created: *episodeRss.FirstAired,
		}
		feed.Items = append(feed.Items, &item)
	}
	rss, _ := feed.ToRss()
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")

	return http.StatusOK, []byte(rss)
}
