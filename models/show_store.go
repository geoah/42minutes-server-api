package models

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	"log"
)

type Model interface {
}

// The Store interface defines methods to manipulate items.
type Store interface {
	Get(id int) (*Show, error)
	GetOrRetrieve(id int) (*Show, error)
	GetAll() ([]Show, error)
	Upsert(p *Show) (int, error)
	Delete(p *Show) (int, error)
}

type ShowStore struct {
	Db *gorp.DbMap
}

// GetAll returns all Shows
func (store *ShowStore) GetAll() ([]Show, error) {
	var shows []Show = make([]Show, 0)
	var showsTemp []Show
	_, err := store.Db.Select(&showsTemp, "select id from shows order by id desc")
	if err != nil {
		return nil, err
	}
	for _, show := range showsTemp {
		ashow, _ := store.Get(show.ID)
		shows = append(shows, *ashow)
	}
	return shows, nil
}

// Get returns a single Show identified by its id, or nil
func (store *ShowStore) Get(id int) (*Show, error) {
	show := Show{}
	err := store.Db.SelectOne(&show, "select * from shows where id=?", id)
	if err != nil {
		return nil, err
	}
	show.Seasons = make([]Season, 0)
	_, err = store.Db.Select(&show.Seasons, "select * from seasons where show_id=?", show.ID)
	if err != nil {
		log.Println("Error while filling in seasons", err)
	}
	for season_i, season := range show.Seasons {
		// this is required as rance copies the records
		show.Seasons[season_i].Episodes = make([]Episode, 0)
		_, err = store.Db.Select(&show.Seasons[season_i].Episodes, "select * from episodes where show_id=? and season=?", show.ID, season.Season)
		if err != nil {
			log.Println("Error while filling in episodes", err)
		}
	}
	return &show, nil
}

// Get returns a single Show identified by its id, if the episode doesn't exist it retrieves it and stores it
func (store *ShowStore) GetOrRetrieve(id int) (*Show, error) {
	show := Show{}
	err := store.Db.SelectOne(&show, "select * from shows where id=?", id)
	if err == sql.ErrNoRows {
		show.UpdateInfoByTvdbID(id)
		store.Upsert(&show)
	} else if err != nil {
		log.Println("TODO error", err)
		// show.UpdateInfoByTvdbID(id)
		// store.Insert(&show)
		return nil, err
	}
	return &show, nil
}

// Upsert inserts or updates a Show and returns count of inserted/updated records
func (store *ShowStore) Upsert(show *Show) (int, error) {
	log.Printf("Trying to upsert show:%d", show.ID)
	err := store.Db.SelectOne(&show, "select * from shows where id=?", show.ID)
	if err == sql.ErrNoRows {
		log.Printf("Trying to insert show:%d", show.ID)
		store.Db.Insert(show)
		// TODO Check errors
	} else if err != nil {
		log.Println("TODO error 3", err)
	}
	for _, season := range show.Seasons {
		err := store.Db.SelectOne(&show, "select * from seasons where show_id=? and season=?", show.ID, season.Season)
		if err == sql.ErrNoRows {
			log.Printf("Trying to insert season:%d:%d", show.ID, season.Season)
			store.Db.Insert(&season)
			// TODO Check errors
			for _, episode := range season.Episodes {
				err := store.Db.SelectOne(&show, "select * from episodes where show_id=? and season=? and episode=?", show.ID, season.Season, episode.Episode)
				if err == sql.ErrNoRows {
					log.Printf("Trying to insert episode:%d:%d:%d", show.ID, season.Season, episode.Episode)
					store.Db.Insert(&episode)
					// TODO Check errors
				} else if err != nil {
					log.Println("TODO error 4", err)
				}
			}
		} else if err != nil {
			log.Println("TODO error 3", err)
		}
	}
	return 1, nil
}

// Deletes removes Show and returns count of removed records
func (store *ShowStore) Delete(show *Show) (int, error) {
	_, err := store.Db.Delete(show)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func ShowFindAllByName(name string, maxResults int) ([]*Show, error) {
	shows := make([]*Show, 0)
	trakt := GetTraktSession()
	// store := *GetStoreSession()
	showResults, err := trakt.ShowSearch(name)
	if err != nil {
		return shows, err
	}
	for _, show := range showResults {
		// TODO: Add additional checks
		if show.Title != "" && show.TvdbID > 0 {
			newShow, err := store.GetOrRetrieve(show.TvdbID)
			if err == nil && show.ImdbID != "" {
				shows = append(shows, newShow)
			}
		}
	}
	return shows, nil
}
