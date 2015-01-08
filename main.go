package main

import (
	. "github.com/42minutes/api/models"
	"github.com/codegangsta/martini-contrib/encoder"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"log"
	"net/http"
)

// Martini instance
var m *martini.Martini
var store Store

func init() {
	// Initialize store
	store = *GetStoreSession()

	// Initialize martini
	m = martini.New()

	// Setup martini middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())

	// Setup routes
	r := martini.NewRouter()
	r.Get(`/shows`, ApiShowsAll)
	// r.Get(`/shows/:id`, ApiShows)
	r.Post(`/files`, ApiFilesPost)

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
