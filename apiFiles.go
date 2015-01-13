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
	"fmt"
	"net/http"
	// "regexp"
	"strings"
)

func ApiProcessFiles(userId string) {
	db := GetDbSession()

	// patterns := []*regexp.Regexp{
	// 	regexp.MustCompile("[Ss]([0-9]+)[][ ._-]*[Ee]([0-9]+)([^\\/]*).(avi|mkv)$"),
	// 	regexp.MustCompile(`[\\/\._ \[\(-]([0-9]+)x([0-9]+)([^\\/]*).(avi|mkv)$`),
	// }
	var userFiles []UserFile

	_, err := db.Select(&userFiles, "select * from users_files where processed=0 and user_id=?", userId)
	if err == nil {
		for index, userFile := range userFiles {
			// fmt.Println(userFile, index)
			var seriesName string
			seps := [2]string{"\\", "/"}
			for _, sep := range seps {
				pathElems := strings.Split(userFile.RelativePath, sep)
				if len(pathElems) > 1 {
					seriesName = pathElems[0]
					break
				} else {
					seriesName = ""
				}
			}
			// fmt.Println(seriesName)
			show, err := ShowFindAllByName(seriesName, 1)
			if err == nil {
				// fmt.Println(show[0].Title)
				userFiles[index].ShowID = show[0].ID
				_, err := db.Update(&userFiles[index])
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	} else {
		fmt.Println(err)
	}
}

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
	//Temp call for testing
	go ApiProcessFiles(user.ID)
	return http.StatusOK, encoder.Must(enc.Encode(userFiles))
}
