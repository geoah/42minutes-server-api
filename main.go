package main

import (
	"database/sql"
	. "github.com/42minutes/api/models"
	. "github.com/42minutes/api/stores"
	"github.com/codegangsta/martini-contrib/encoder"
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	_ "github.com/ziutek/mymysql/autorc"
	_ "github.com/ziutek/mymysql/godrv"
	_ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
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
	// initialize the DbMap
	dbmap := initDb()
	// defer dbmap.Db.Close()

	// Initialize store
	store = &ShowStore{
		M:  make(map[uint64]*Show),
		Db: dbmap,
	}

	// Initialize martini
	m = martini.New()

	// Setup martini middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())

	// Setup routes
	r := martini.NewRouter()
	r.Get(`/shows`, ApiShowsAll)
	r.Get(`/shows/:id`, ApiShows)

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

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("mymysql", "tcp:localhost:3306*42minutes/root/root")
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	// add tables
	dbmap.AddTableWithName(Show{}, "shows")
	dbmap.AddTableWithName(Episode{}, "episodes")
	dbmap.AddTableWithName(Season{}, "seasons")
	dbmap.AddTableWithName(UserShow{}, "users_shows")
	dbmap.AddTableWithName(UserSeason{}, "users_seasons")
	dbmap.AddTableWithName(UserEpisode{}, "users_episodes")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
