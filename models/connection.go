package models

import (
	"database/sql"
	"log"
	"os"

	"github.com/42minutes/go-trakt"
	"github.com/coopernurse/gorp"
	_ "github.com/ziutek/mymysql/autorc"
	_ "github.com/ziutek/mymysql/godrv"
	_ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
)

var dbmap *gorp.DbMap
var traktClient *trakt.Client
var store Store

func GetDbSession() *gorp.DbMap {
	if dbmap != nil {
		return dbmap
	}

	db_uri := "tcp:"

	if os.Getenv("DB_HOST") != "" {
		db_uri += os.Getenv("DB_HOST")
	} else {
		db_uri += "localhost"
	}

	if os.Getenv("DB_PORT") != "" {
		db_uri += ":" + os.Getenv("DB_PORT")
	} else {
		db_uri += ":3306"
	}

	if os.Getenv("DB_NAME") != "" {
		db_uri += "*" + os.Getenv("DB_NAME")
	} else {
		db_uri += "*42minutes"
	}

	if os.Getenv("DB_USER") != "" && os.Getenv("DB_PASS") != "" {
		db_uri += "/" + os.Getenv("DB_USER") + "/" + os.Getenv("DB_PASS")
	}

	log.Println("Connecting to", db_uri)

	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("mymysql", db_uri)
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	// defer dbmap.Db.Close()

	// add tables
	dbmap.AddTableWithName(Show{}, "shows").SetKeys(false, "id")
	dbmap.AddTableWithName(Season{}, "seasons").SetKeys(false, "show_id", "season")
	dbmap.AddTableWithName(Episode{}, "episodes").SetKeys(false, "show_id", "season", "episode")
	dbmap.AddTableWithName(UserShow{}, "users_shows").SetKeys(false, "user_id", "show_id")
	dbmap.AddTableWithName(UserSeason{}, "users_seasons").SetKeys(false, "user_id", "show_id", "season")
	dbmap.AddTableWithName(UserEpisode{}, "users_episodes").SetKeys(false, "user_id", "show_id", "season", "episode")
	dbmap.AddTableWithName(User{}, "users").SetKeys(false, "id")
	dbmap.AddTableWithName(UserFile{}, "users_files").SetKeys(false, "user_id", "relative_path")
	dbmap.AddTableWithName(ShowMatch{}, "shows_matches").SetKeys(false, "title")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func GetTraktSession() *trakt.Client {
	if traktClient != nil {
		return traktClient
	}

	var trakt_api_key, trakt_access_token string
	if os.Getenv("TRAKT_API_KEY") != "" {
		trakt_api_key = os.Getenv("TRAKT_API_KEY")
	} else {
		log.Fatal("Missing TRAKT_API_KEY")
	}
	if os.Getenv("TRAKT_ACCESS_TOKEN") != "" {
		trakt_access_token = os.Getenv("TRAKT_ACCESS_TOKEN")
	} else {
		log.Fatal("Missing TRAKT_ACCESS_TOKEN")
	}

	traktClient = trakt.NewClient(
		trakt_api_key,
		trakt.TokenAuth{AccessToken: trakt_access_token},
	)

	return traktClient
}

func GetStoreSession() *Store {
	if store != nil {
		return &store
	}

	var db *gorp.DbMap = GetDbSession()
	store = &ShowStore{
		Db: db,
	}
	return &store
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
