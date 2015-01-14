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
	GetShow(id int) (*Show, error)
	GetShowOrRetrieve(id int) (*Show, error)
	GetSeasonsByShowId(id int) ([]*Season, error)
	GetSeasonsOrRetrieveByShowId(id int) ([]*Season, error)
	GetEpisodesByShowIdAndSeason(id int, seasonNumber int) ([]*Episode, error)
	GetEpisodesOrRetrieveByShowIdAndSeason(id int, seasonNumber int) ([]*Episode, error)
	GetShowOrRetrieveFromTitle(showName string) (*Show, error)
	UserShowUpsert(p *UserShow) error
	Delete(p *Show) (int, error)
}

type ShowStore struct {
	Db *gorp.DbMap
}

func (store *ShowStore) GetShows() ([]*Show, error) {
	var shows []*Show = make([]*Show, 0)
	_, err := store.Db.Select(&shows, "select * from shows order by id desc")
	return shows, err
}

func (store *ShowStore) GetShow(showId int) (*Show, error) {
	var show Show = Show{}
	err := store.Db.SelectOne(&show, "select * from shows where id=?", showId)
	return &show, err
}

func (store *ShowStore) GetShowOrRetrieve(showId int) (*Show, error) {
	show, err := store.GetShow(showId)
	if err == sql.ErrNoRows {
		err = nil
		log.Printf(" > Show does not exist locally")
		trakt := GetTraktSession()
		traktShow, _ := trakt.Shows().One(showId) // TODO Error
		show.MapInfo(*traktShow)
		// Cache
		go func(show *Show) {
			db := GetDbSession()
			log.Printf("Trying to insert show:%d", show.ID)
			db.Insert(show)
		}(show)
	}
	return show, err
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
		traktSeasons, err := trakt.Seasons().All(showId)
		if err.HasError() == false {
			for _, traktSeason := range traktSeasons {
				season := Season{}
				season.MapInfo(traktSeason)
				season.ShowID = showId
				seasons = append(seasons, &season)
				// Cache
				go func(season *Season) {
					db := GetDbSession()
					log.Printf("Trying to insert season:%d:%d", season.ShowID, season.Season)
					db.Insert(season)
				}(&season)
			}
		}
	} else if err != nil {
		log.Println("TODO error", err)
		return seasons, err
	}
	return seasons, err
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
		traktEpisodes, err := trakt.Episodes().AllBySeason(showId, seasonNumber)
		if err.HasError() == false {
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
					db.Insert(episode)
				}(&episode)
			}
		}
	} else if err != nil {
		log.Println("TODO error", err)
		return episodes, err
	}
	return episodes, err
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
	err := store.Db.SelectOne(&showMatch, "select * from shows_matches where title=?", showTitle)
	fmt.Println("GetShowOrRetrieveFromTitle:SelectOne(showMatch)>err", err)
	if err == nil {
		show, err = store.GetShowOrRetrieve(showMatch.ShowID)
		fmt.Println("GetShowOrRetrieveFromTitle:GetShowOrRetrieve(show)>err", err)
		// show = showTmp
	} else if err == sql.ErrNoRows {
		log.Printf(" > Show could not be matched locally")
		shows, err := ShowFindAllByName(showTitle, 1)
		fmt.Println("GetShowOrRetrieveFromTitle:ShowFindAllByName()>err", err)
		if len(shows) == 0 {
			err = errors.New("Could not find any matches.")
		}
		if err == nil {
			show = shows[0]
			// Cache
			if err == nil {
				showMatch.ShowID = show.ID
				showMatch.Title = showTitle
				go func(showMatch *ShowMatch) {
					db := GetDbSession()
					log.Printf("Trying to insert showmatch:%d:%s", showMatch.ShowID, showMatch.Title)
					db.Insert(showMatch)
				}(&showMatch)
			}
		}
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
			if err == nil && newShow.TraktID > 0 {
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
