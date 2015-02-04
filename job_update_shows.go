package main

import (
	"fmt"
	"time"

	. "github.com/42minutes/42minutes-server-api/models"
	"github.com/42minutes/go-torrentlookup"
)

type ShowsUpdateJob struct {
	Msg string
}

func (j *ShowsUpdateJob) Run() error {
	db := *GetDbSession()
	store := *GetStoreSession()

	var shows []*Show
	_, err := db.Select(&shows, "SELECT shows.* FROM shows LEFT JOIN users_shows ON shows.id = users_shows.show_id WHERE users_shows.library = true ORDER BY shows.title asc")
	if err != nil {
		fmt.Println("ERRRR1>>>>>>>>>", err)
		return nil
	}

	for _, show := range shows {
		fmt.Printf("Updating %s.\n", show.Title)
		seasons, err := store.GetSeasonsOrRetrieveByShowId(show.ID)
		if err != nil {
			fmt.Println("ERRRR1>>>>>>>>>", err)
			continue
		}
		fmt.Printf(" > Retrieved %d seasons\n", len(seasons))
		for _, season := range seasons {
			if season.Season > 0 {
				fmt.Printf(" > > Updating season %d\n", season.Season)
				episodes, err := store.GetEpisodesOrRetrieveByShowIdAndSeason(show.ID, season.Season)
				if err != nil {
					fmt.Println("ERRRR2>>>>>>>>>", err)
					continue
				}
				fmt.Printf(" > > Retrieved %d episodes\n", len(seasons))
				for _, episode := range episodes {
					// go func(show *Show, episode *Episode) {
					// db := *GetDbSession()
					if episode.InfoHashHd720p != "" {
						continue
					}
					// We don't care about specials for now or episodes that have not yet aired
					if episode.Episode > 0 && time.Now().Sub(*episode.FirstAired).Hours() > 0 {
						fmt.Printf(" > > Updating episode %d\n", episode.Episode)
						query := fmt.Sprintf("%s S%02dE%02d 720p HDTV", show.Title, episode.Season, episode.Episode)
						// fmt.Printf(" > > Looking for %s\n", query)
						name720p, infohash720p := torrentlookup.Search(query)

						if infohash720p != "" {
							fmt.Printf(" > > > Found 720p %s as %s\n", name720p, infohash720p)
							episode.InfoHashHd720p = infohash720p
						}

						query = fmt.Sprintf("%s S%02dE%02d HDTV", show.Title, episode.Season, episode.Episode)
						// fmt.Printf(" > > Looking for %s\n", query)
						name480p, infohash480p := torrentlookup.Search(query)

						if infohash480p != "" {
							fmt.Printf(" > > > Found 480p %s as %s\n", name480p, infohash480p)
							episode.InfoHashSd480p = infohash480p
						}

						if infohash480p != "" || infohash720p != "" {
							db.Update(episode)
						}
					}
					// }(show, episode)
				}
			}
		}

		// // TODO Check if we actually need anything from eztv else skip it
		// eztvShow := show.GetEztvShow()
		// _, _ = store.GetSeasonsOrRetrieveByShowId(show.ID) // TODO Check err
		// // episodes, _ := store.GetEpisodesOrRetrieveByShowIdAndSeason(show.ID, season.Season) // TODO Check err
		// for _, ezEpisode := range eztvShow.Episodes {
		// 	fmt.Printf(" > EzTv has S%dE%d.\n", int(ezEpisode.Season), int(ezEpisode.Episode))
		// 	episode, err := store.GetEpisodeOrRetrieveByShowIdAndSeasonAndEpisode(show.ID, int(ezEpisode.Season), int(ezEpisode.Episode)) // TODO Check err
		// 	if err == nil {
		// 		if ezEpisode.Torrents.Hd720p.URL != "" {
		// 			episode.EztvInfoHashHd720p = ezEpisode.Torrents.Hd720p.URL
		// 		} else {
		// 			filename, infohash :=torrentlookup.Search()
		// 		}
		// 		if ezEpisode.Torrents.Sd480p.URL != "" {
		// 			episode.EztvInfoHashSd480p = ezEpisode.Torrents.Sd480p.URL
		// 		}
		// 		if ezEpisode.Torrents.Sd.URL != "" {
		// 			episode.EztvInfoHashSd = ezEpisode.Torrents.Sd.URL
		// 		}
		// 		_, err := db.Update(episode)
		// 		if err != nil {
		// 			fmt.Println(">>>>>>>>>>>> ERRPR ", err)
		// 		}
		// 	} else {
		// 		fmt.Println(">>>>>>>>>>>> ERRPR ", err)
		// 	}
		// }

		// fmt.Printf(" > Status is %s.\n", show.Status)
		switch show.Status {
		case "returning series": // airing right now
		case "in production": // airing soon
		case "canceled": // canceled
		case "ended": // ended
		}
	}

	return nil
}
