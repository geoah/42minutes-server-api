package main

import (
	. "github.com/42minutes/api/models"
	. "github.com/42minutes/api/stores"
	"github.com/codegangsta/martini-contrib/encoder"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"log"
	"net/http"
)

// Martini instance
var m *martini.Martini
var store Store

// Create config struct to hold random things
var config struct {
}

func init() {
	// Initialize store
	store = &SeriesStore{
		M: make(map[uint64]*Series),
	}

	// Initialize martini
	m = martini.New()

	// Setup martini middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())

	// Setup routes
	r := martini.NewRouter()
	r.Get(`/series`, ApiSeriesAll)
	r.Get(`/series/:id`, ApiSeries)

	// Allow CORS
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	// Other stuff
	m.Use(func(c martini.Context, w http.ResponseWriter) {
		// Inject JSON Encoder
		c.MapTo(encoder.JsonEncoder{}, (*encoder.Encoder)(nil))
		// Force Content-Type
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	})
	// Inject database
	m.MapTo(store, (*Store)(nil))
	// Add the router action
	m.Action(r.Handle)
}

func main() {
	// Startup HTTP server
	if err := http.ListenAndServe(":8000", m); err != nil {
		log.Fatal(err)
	}
}
