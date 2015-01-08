package models

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	"github.com/hobeone/gotrakt"
	"github.com/jmcvetta/napping"
	_ "github.com/ziutek/mymysql/autorc"
	_ "github.com/ziutek/mymysql/godrv"
	_ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
	"log"
	"os"
)

var dbmap *gorp.DbMap
var trakt *gotrakt.TraktTV
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
	dbmap.AddTableWithName(UserSeason{}, "users_seasons").SetKeys(false, "user_id", "show_id", "season_id")
	dbmap.AddTableWithName(UserEpisode{}, "users_episodes").SetKeys(false, "user_id", "show_id", "season_id", "episode_id")
	dbmap.AddTableWithName(User{}, "users").SetKeys(false, "id")
	dbmap.AddTableWithName(UserFile{}, "users_files").SetKeys(false, "user_id", "full_path_hash")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func GetTraktSession() *gotrakt.TraktTV {
	if trakt != nil {
		return trakt
	}

	sess := &napping.Session{}
	sess.Params = &napping.Params{
		"testing": "true",
	}
	var err error
	trakt, err = gotrakt.New("testingapi", gotrakt.Session(sess))
	if err != nil {
		return nil
	}
	return trakt
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
