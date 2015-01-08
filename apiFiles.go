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
	var userFiles []UserFile

	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	err := json.Unmarshal(body, &userFiles)

	if err != nil {
		return http.StatusNotFound, encoder.Must(enc.Encode(
			NewError(ErrCodeNotExist, "Could not decode body")))
	}

	for userFile_i, userFile := range userFiles {
		err = db.SelectOne(&userFiles[userFile_i], "select * from users_files where user_id=? and full_path_hash=?", userFile.UserID, userFile.FullPathHash)
		if err == sql.ErrNoRows {
			userFiles[userFile_i].ID = uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen)
			db.Insert(&userFile)
		} else if err != nil {
			return http.StatusNotFound, encoder.Must(enc.Encode(
				NewError(ErrCodeNotExist, "Could not insert file record")))
		}
	}
	return http.StatusOK, encoder.Must(enc.Encode(userFiles))
}
