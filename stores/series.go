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
	Get(id uint64) *Series
	GetAll() []*Series
	Add(p *Series) (uint64, error)
	Update(p *Series) error
	FindAllByName(name string, maxResults int) []*Series
}

type SeriesStore struct {
	sync.RWMutex
	M  map[uint64]*Series
	Db *gorp.DbMap
}

// GetAll returns all Seriess from memory
func (store *SeriesStore) GetAll() []*Series {
	store.RLock()
	defer store.RUnlock()
	if len(store.M) == 0 {
		return nil
	}
	ar := make([]*Series, len(store.M))
	i := 0
	for _, v := range store.M {
		ar[i] = v
		i++
	}
	return ar
}

// Get returns a single Series identified by its id, or nil.
func (store *SeriesStore) Get(id uint64) *Series {
	// Check if series is in memory else try to find in storage
	if store.M[id] == nil {
		err := store.Pull(id)
		if err != nil {
			fmt.Println(err)
		}
	}
	// If neither is available retrieve it from tvdb
	if store.M[id] == nil {
		newSeries := Series{}
		newSeries.FetchInfoByID(id)
		// Add it in memory and storage
		store.Add(&newSeries)
		err := store.Push(&newSeries)
		if err != nil {
			fmt.Println(err)
		}
	}
	store.RLock()
	defer store.RUnlock()
	return store.M[id]
}

// Add stores a new Series and returns its newly generated id, or an error.
func (store *SeriesStore) Add(p *Series) (uint64, error) {
	store.Lock()
	defer store.Unlock()
	// Store it
	store.M[p.ID] = p
	return p.ID, nil
}

// Update updates Series and returns nil
func (store *SeriesStore) Update(p *Series) error {
	store.Lock()
	defer store.Unlock()
	store.M[p.ID] = p
	return nil
}

func (store *SeriesStore) Delete(id uint64) {
	store.Lock()
	defer store.Unlock()
	delete(store.M, id)
}

// TODO: Is this really needed?
// Pull retrieves all Series from tiedot and saves them in memory
func (store *SeriesStore) PullAll() error {
	var series []Series
	_, err := store.Db.Select(&series, "select * from series order by id desc")
	if err != nil {
		return err
	}
	for _, ser := range series {
		store.Add(&ser)
	}
	return nil
}

// Pull retrieves a single Series identified by its id from tiedot and saves it in memory
func (store *SeriesStore) Pull(id uint64) error {
	series := Series{}
	err := store.Db.SelectOne(series, "select * from series where id=?", id)
	if err != nil {
		return err
	}
	store.Add(&series)
	return nil
}

// Push stores a single series in storage
func (store *SeriesStore) Push(series *Series) error {
	err := store.Db.Insert(series)
	if err != nil {
		return err
	}
	store.Add(series)
	return nil
}

func (store *SeriesStore) FindAllByName(name string, maxResults int) []*Series {
	seriesResults := make([]*Series, 0)
	seriesList, err := tvdb.SearchSeries(name, maxResults)
	if err == nil {
		for _, series := range seriesList.Series {
			// TODO: Add additional checks
			if series.SeriesName != "" && series.ImdbID != "" {
				newSeries := store.Get(series.ID)
				seriesResults = append(seriesResults, newSeries)
			}
		}
	}
	return seriesResults
}
