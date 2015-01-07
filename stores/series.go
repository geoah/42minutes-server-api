package stores

import (
	"fmt"
	. "github.com/42minutes/api/models"
	"github.com/coopernurse/gorp"
	"github.com/garfunkel/go-tvdb"
	"sync"
)

// The Store interface defines methods to manipulate items.
type Store interface {
	Get(id uint64) *Show
	GetAll() []*Show
	Add(p *Show) (uint64, error)
	Update(p *Show) error
	FindAllByName(name string, maxResults int) []*Show
}

type ShowStore struct {
	sync.RWMutex
	M  map[uint64]*Show
	Db *gorp.DbMap
}

// GetAll returns all Shows from memory
func (store *ShowStore) GetAll() []*Show {
	store.RLock()
	defer store.RUnlock()
	if len(store.M) == 0 {
		return nil
	}
	ar := make([]*Show, len(store.M))
	i := 0
	for _, v := range store.M {
		ar[i] = v
		i++
	}
	return ar
}

// Get returns a single Show identified by its id, or nil.
func (store *ShowStore) Get(id uint64) *Show {
	// Check if Show is in memory else try to find in storage
	if store.M[id] == nil {
		err := store.Pull(id)
		if err != nil {
			fmt.Println(err)
		}
	}
	// If neither is available retrieve it from tvdb
	if store.M[id] == nil {
		newShow := Show{}
		newShow.FetchInfoByID(id)
		// Add it in memory and storage
		store.Add(&newShow)
		err := store.Push(&newShow)
		if err != nil {
			fmt.Println(err)
		}
	}
	store.RLock()
	defer store.RUnlock()
	return store.M[id]
}

// Add stores a new Show and returns its newly generated id, or an error.
func (store *ShowStore) Add(p *Show) (uint64, error) {
	store.Lock()
	defer store.Unlock()
	// Store it
	store.M[p.ID] = p
	return p.ID, nil
}

// Update updates Show and returns nil
func (store *ShowStore) Update(p *Show) error {
	store.Lock()
	defer store.Unlock()
	store.M[p.ID] = p
	return nil
}

func (store *ShowStore) Delete(id uint64) {
	store.Lock()
	defer store.Unlock()
	delete(store.M, id)
}

// TODO: Is this really needed?
// Pull retrieves all Show from tiedot and saves them in memory
func (store *ShowStore) PullAll() error {
	var Show []Show
	_, err := store.Db.Select(&Show, "select * from Show order by id desc")
	if err != nil {
		return err
	}
	for _, ser := range Show {
		store.Add(&ser)
	}
	return nil
}

// Pull retrieves a single Show identified by its id from tiedot and saves it in memory
func (store *ShowStore) Pull(id uint64) error {
	show := Show{}
	err := store.Db.SelectOne(show, "select * from Show where id=?", id)
	if err != nil {
		return err
	}
	store.Add(&show)
	return nil
}

// Push stores a single Show in storage
func (store *ShowStore) Push(show *Show) error {
	err := store.Db.Insert(show)
	if err != nil {
		return err
	}
	store.Add(show)
	return nil
}

func (store *ShowStore) FindAllByName(name string, maxResults int) []*Show {
	ShowResults := make([]*Show, 0)
	ShowList, err := tvdb.SearchSeries(name, maxResults)
	if err == nil {
		for _, series := range ShowList.Series {
			// TODO: Add additional checks
			if series.SeriesName != "" && series.ImdbID != "" {
				newShow := store.Get(series.ID)
				ShowResults = append(ShowResults, newShow)
			}
		}
	}
	return ShowResults
}
