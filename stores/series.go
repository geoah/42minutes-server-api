package stores

import (
	. "github.com/42minutes/api/models"
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

// Thread-safe in-memory map.
type SeriesStore struct {
	sync.RWMutex
	M map[uint64]*Series
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
	if store.M[id] == nil {
		newSeries := Series{}
		newSeries.FetchInfoByID(id)
		store.Add(&newSeries)
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
