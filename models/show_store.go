package models

import (
	"database/sql"
	// "fmt"
	"github.com/coopernurse/gorp"
	"log"
	// "github.com/hobeone/gotrakt"
	// "sync"
)

// The Store interface defines methods to manipulate items.
type Store interface {
	Get(id int) (*Show, error)
	GetOrRetrieve(id int) (*Show, error)
	GetAll() ([]Show, error)
	Insert(p *Show) (int, error)
	Update(p *Show) (int, error)
	Delete(p *Show) (int, error)
	// FindAllByName(name string, maxResults int) []*Show
}

type ShowStore struct {
	// Trakt *gotrakt.TraktTV
	Db *gorp.DbMap
}

// GetAll returns all Shows
func (store *ShowStore) GetAll() ([]Show, error) {
	var shows []Show
	_, err := store.Db.Select(&shows, "select * from shows order by id desc")
	if err != nil {
		return nil, err
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
	return &show, nil
}

// Get returns a single Show identified by its id, if the episode doesn't exist it retrieves it and stores it
func (store *ShowStore) GetOrRetrieve(id int) (*Show, error) {
	show := Show{}
	err := store.Db.SelectOne(&show, "select * from shows where id=?", id)
	if err == sql.ErrNoRows {
		show.UpdateInfoByTvdbID(id)
		store.Insert(&show)
	} else if err != nil {
		log.Println("TODO error", err)
		// show.UpdateInfoByTvdbID(id)
		// store.Insert(&show)
		return nil, err
	}
	return &show, nil
}

// Insert stores a new Show and returns nil
func (store *ShowStore) Insert(show *Show) (int, error) {
	err := store.Db.Insert(show)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return 1, nil
}

// Update updates Show and returns count of updated records
func (store *ShowStore) Update(show *Show) (int, error) {
	_, err := store.Db.Update(show)
	if err != nil {
		return 0, err
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
