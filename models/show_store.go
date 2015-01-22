package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/coopernurse/gorp"
)

type Model interface {
}

// The Store interface defines methods to manipulate items.
type Store interface {
	GetShows() ([]*Show, error)
	GetShowOrRetrieve(id int) (*Show, error)
	GetSeasonsByShowId(id int) ([]*Season, error)
	GetSeasonsOrRetrieveByShowId(id int) ([]*Season, error)
	GetEpisodesByShowIdAndSeason(id int, seasonNumber int) ([]*Episode, error)
	GetEpisodesOrRetrieveByShowIdAndSeason(id int, seasonNumber int) ([]*Episode, error)
	GetEpisodeOrRetrieveByShowIdAndSeasonAndEpisode(id int, seasonNumber int, episodeNumber int) (*Episode, error)
	GetShowOrRetrieveFromTitle(showName string) (*Show, error)
	UserShowUpsert(p *UserShow) error
	Delete(p *Show) (int, error)
}

type ShowStore struct {
	Db *gorp.DbMap
}

func (store *ShowStore) GetShows() ([]*Show, error) {
	var shows []*Show = make([]*Show, 0)
	_, err := store.Db.Select(&shows, "select * from shows where rating>0 order by id desc")
	return shows, err
}

func (store *ShowStore) GetShowOrRetrieve(showId int) (*Show, error) {
	var show Show = Show{}
	err := store.Db.SelectOne(&show, "select * from shows where id=?", showId)
	if err == sql.ErrNoRows {
		log.Printf(" > Show does not exist locally")
		trakt := GetTraktSession()
		traktShow, result := trakt.Shows().One(showId)
		if result.Err != nil {
			return &show, result.Err
		} else {
			show.MapInfo(*traktShow)
			// Cache
			go func(show *Show) {
				db := GetDbSession()
				log.Printf("Trying to insert show:%d", show.ID)
				err := db.Insert(show)
				if err != nil {
					log.Println("ERR:", err)
				}
			}(&show)
		}
		return &show, nil
	} else if err != nil {
		return &show, err
	}
	return &show, nil
}

func (store *ShowStore) GetSeasonsByShowId(showId int) ([]*Season, error) {
	var seasons []*Season = make([]*Season, 0)
	_, err := store.Db.Select(&seasons, "select * from seasons where show_id=?", showId)
	return seasons, err
}

func (store *ShowStore) GetSeasonsOrRetrieveByShowId(showId int) ([]*Season, error) {
	log.Printf("Trying to retrieve seasons for showid:%d\n", showId)
	var seasons []*Season = make([]*Season, 0)
	_, err := store.Db.Select(&seasons, "select * from seasons where show_id=?", showId)
	if err == sql.ErrNoRows || len(seasons) == 0 {
		log.Printf(" > Show's seasons do not exist locally")
		trakt := GetTraktSession()
		traktSeasons, result := trakt.Seasons().All(showId)
		if result.Err != nil {
			return seasons, result.Err
		} else {
			for _, traktSeason := range traktSeasons {
				season := Season{}
				season.MapInfo(traktSeason)
				season.ShowID = showId
				seasons = append(seasons, &season)
				// Cache
				go func(season *Season) {
					db := GetDbSession()
					log.Printf("Trying to insert season:%d:%d", season.ShowID, season.Season)
					err := db.Insert(season)
					if err != nil {
						log.Println("ERR:", err)
					}
				}(&season)
			}
			return seasons, nil
		}
	} else if err != nil {
		return seasons, err
	}
	return seasons, nil
}

func (store *ShowStore) GetEpisodesByShowIdAndSeason(showId int, seasonNumber int) ([]*Episode, error) {
	var episodes []*Episode = make([]*Episode, 0)
	_, err := store.Db.Select(&episodes, "select * from episodes where show_id=? and season=?", showId, seasonNumber)
	return episodes, err
}

func (store *ShowStore) GetEpisodesOrRetrieveByShowIdAndSeason(showId int, seasonNumber int) ([]*Episode, error) {
	log.Printf("Trying to retrieve episodes for showid:%d season:%d\n", showId, seasonNumber)
	var episodes []*Episode = make([]*Episode, 0)
	_, err := store.Db.Select(&episodes, "select * from episodes where show_id=? and season=?", showId, seasonNumber)
	if err == sql.ErrNoRows || len(episodes) == 0 {
		log.Printf(" > Show's episodes do not exist locally")
		trakt := GetTraktSession()
		traktEpisodes, result := trakt.Episodes().AllBySeason(showId, seasonNumber)
		if result.Err != nil {
			return episodes, result.Err
		} else {
			for _, traktEpisode := range traktEpisodes {
				episode := Episode{}
				episode.MapInfo(traktEpisode)
				episode.ShowID = showId
				episode.Season = seasonNumber
				episodes = append(episodes, &episode)
				// Cache
				go func(episode *Episode) {
					db := GetDbSession()
					log.Printf("Trying to insert episode:%d:%d:%d", episode.ShowID, episode.Season, episode.Episode)
					err := db.Insert(episode)
					if err != nil {
						log.Println("ERR:", err)
					}
				}(&episode)
			}
			return episodes, nil
		}
	} else if err != nil {
		return episodes, err
	}
	return episodes, nil
}

func (store *ShowStore) GetEpisodeOrRetrieveByShowIdAndSeasonAndEpisode(showId int, seasonNumber int, episodeNumber int) (*Episode, error) {
	log.Printf("Trying to retrieve episodes for showid:%d season:%d episode:%d\n", showId, seasonNumber, episodeNumber)
	var episode Episode = Episode{}
	err := store.Db.SelectOne(&episode, "select * from episodes where show_id=? and season=? and episode=? limit 1", showId, seasonNumber, episodeNumber)
	if err == sql.ErrNoRows {
		log.Printf(" > Show's episode does not exist locally")
		// TODO Cache the episode
		// trakt := GetTraktSession()
		// traktEpisodes, err := trakt.Episodes().AllBySeason(showId, seasonNumber)
		// if err.HasError() == false {
		// 	for _, traktEpisode := range traktEpisodes {
		// 		episode := Episode{}
		// 		episode.MapInfo(traktEpisode)
		// 		episode.ShowID = showId
		// 		episode.Season = seasonNumber
		// 		episodes = append(episodes, &episode)
		// 		// Cache
		// 		go func(episode *Episode) {
		// 			db := GetDbSession()
		// 			log.Printf("Trying to insert episode:%d:%d:%d", episode.ShowID, episode.Season, episode.Episode)
		// 			db.Insert(episode)
		// 		}(&episode)
		// 	}
		// }
	} else if err != nil {
		log.Println("TODO error", err)
		return &episode, err
	}
	return &episode, nil
}

// Deletes removes Show and returns count of removed records
func (store *ShowStore) Delete(show *Show) (int, error) {
	_, err := store.Db.Delete(show)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (store *ShowStore) GetShowOrRetrieveFromTitle(showTitle string) (*Show, error) {
	var show *Show
	var showMatch ShowMatch = ShowMatch{}
	err := store.Db.SelectOne(&showMatch, "select * from shows_matches where title=? limit 1", showTitle)
	if err == nil && showMatch.ShowID == 0 {
		showMatch.ShowID = 0
		// err = errors.New("Has been cached as unmatched.")
		log.Printf("GetShowOrRetrieveFromTitle: Show '%s' has been found in cache with showid:0 as it could not be matched last time", showTitle)

	} else if err == nil && showMatch.ShowID > 0 {
		show, err = store.GetShowOrRetrieve(showMatch.ShowID)
		log.Printf("GetShowOrRetrieveFromTitle: Show '%s' has been found in cache with showid:%d", showTitle, show.ID)
		if err != nil {
			showMatch.ShowID = 0
			err = nil
		}

	} else if err == sql.ErrNoRows {
		log.Printf("GetShowOrRetrieveFromTitle: Show '%s' could not be found in cache", showTitle)
		shows, err := ShowFindAllByName(showTitle, 1)

		if len(shows) == 0 && err == nil {
			log.Printf("GetShowOrRetrieveFromTitle: Show '%s' could not be found in Trakt", showTitle)
			showMatch.ShowID = 0
		} else if len(shows) > 0 && err == nil {
			show = shows[0]
			for show_i, show_r := range shows {
				if show_r.Rating > show.Rating {
					show = shows[show_i]
				}
			}
			showMatch.ShowID = show.ID
		} else {
			showMatch.ShowID = 0
			err = errors.New("Failed with matching with Trakt")
			fmt.Println("GetShowOrRetrieveFromTitle:ShowFindAllByName>err", err)
		}

		// Cache
		if err == nil {
			showMatch.Title = showTitle
			go func(showMatch *ShowMatch) {
				db := GetDbSession()
				log.Printf("GetShowOrRetrieveFromTitle: Caching '%s' with showid:%d", showMatch.Title, showMatch.ShowID)
				db.Insert(showMatch)
			}(&showMatch)
		}

	} else {
		fmt.Println("GetShowOrRetrieveFromTitle>err", err)
	}

	if err == nil && (show == nil || show.ID == 0) {
		err = errors.New("Could not be matched")
	}
	return show, err
}

func ShowFindAllByName(name string, maxResults int) ([]*Show, error) {
	shows := make([]*Show, 0)
	trakt := GetTraktSession()
	// store := *GetStoreSession()
	showResults, result := trakt.Shows().Search(name)
	if result.HasError() == true {
		return shows, result.Err
	}
	for _, traktShow := range showResults {
		// TODO: Add additional checks
		if traktShow.Show.Title != "" && traktShow.Show.Ids.Imdb != "" {
			// TODO Currently the api doesn't support properly getting extended info on search
			// so season and episodes were missing a lot of data.
			// newShow, err := store.GetOrRetrieveByTraktShow(&traktShow)
			newShow, err := store.GetShowOrRetrieve(traktShow.Show.Ids.Trakt)
			if err == nil && newShow.TraktID > 0 && newShow.TvdbID > 0 && newShow.ImdbID != "" && newShow.Rating > 0 {
				shows = append(shows, newShow)
			}
		}
	}
	return shows, nil
}

func (store *ShowStore) UserShowUpsert(userShow *UserShow) error {
	err := store.Db.Insert(userShow)
	if err != nil {
		_, err = store.Db.Update(userShow)
	}
	return err
}
