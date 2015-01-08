package main

import (
	"database/sql"
	"encoding/json"
	. "github.com/42minutes/42minutes-server-api/models"
	"github.com/codegangsta/martini-contrib/encoder"
	_ "github.com/garfunkel/go-tvdb"
	"github.com/go-martini/martini"
	"github.com/twinj/uuid"
	"io/ioutil"
	// "log"
	"net/http"
)

func ApiFilesPost(r *http.Request, enc encoder.Encoder, store Store, parms martini.Params) (int, []byte) {
	db := GetDbSession()

	token := r.Header.Get("X-API-TOKEN")
	user := User{}
	err := db.SelectOne(&user, "select * from users where token=?", token)
	if err != nil {
		return http.StatusUnauthorized, encoder.Must(enc.Encode(
			NewError(ErrCodeNotExist, "Error")))
	}

	var userFiles []UserFile

	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &userFiles)

	if err != nil {
		return http.StatusNotFound, encoder.Must(enc.Encode(
			NewError(ErrCodeNotExist, "Could not decode body")))
	}

	for userFile_i, userFile := range userFiles {
		err = db.SelectOne(&userFiles[userFile_i], "select * from users_files where user_id=? and relative_path=?", userFile.UserID, userFile.RelativePath)
		if err == sql.ErrNoRows {
			userFiles[userFile_i].ID = uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen)
			userFiles[userFile_i].UserID = user.ID
			db.Insert(&userFiles[userFile_i])
			// TODO Error
		} else if err != nil {
			// TODO Error
		}
	}
	return http.StatusOK, encoder.Must(enc.Encode(userFiles))
}
