package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	. "github.com/42minutes/42minutes-server-api/models"
	"github.com/codegangsta/martini-contrib/encoder"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
)

// Martini instance
var m *martini.Martini
var store Store

func Authenticate(w http.ResponseWriter, r *http.Request, c martini.Context, enc encoder.Encoder) {
	db := GetDbSession()
	token := r.Header.Get("X-API-TOKEN")
	user := User{}
	err := db.SelectOne(&user, "select * from users where token=?", token)
	if err != nil {
		http.Error(w, "Auth Error", http.StatusUnauthorized)
	}
	c.Map(user)
}

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
	r.Get(`/shows`, Authenticate, ApiShowsGetAll)
	r.Get(`/shows/:id`, Authenticate, ApiShowsGetOne)
	r.Put(`/shows/:id`, Authenticate, ApiShowsPutOne)
	r.Get(`/shows/:showId/seasons`, ApiSeasonsGetAllByShow)
	r.Get(`/shows/:showId/seasons/:seasonNumber/episodes`, ApiEpisodesGetAllByShowAndSeason)
	r.Get(`/shows/:showId/seasons/:seasonNumber/episodes/:episodeNumber`, ApiEpisodesGetOneByShowAndSeasonAndEpisode)
	r.Post(`/files`, ApiFilesPost)
	r.Patch(`/files`, ApiProcessFiles)
	r.Get(`/register`, ApiUsersRegister)
	r.Post(`/login`, ApiUsersLogin)
	r.Get(`/rss`, ApiRss)

	// Allow CORS
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-API-TOKEN"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
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
	// Run show update job
	go func() {
		s := ShowsUpdateJob{}
		s.Run()
	}()

	// Startup HTTP server
	fmt.Println("Starting HTTP server")
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	if err := http.ListenAndServe("0.0.0.0:"+port, m); err != nil {
		log.Fatal(err)
	}
}
