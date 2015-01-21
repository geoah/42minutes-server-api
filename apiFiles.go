package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"

	. "github.com/42minutes/42minutes-server-api/models"
	"github.com/codegangsta/martini-contrib/encoder"
	_ "github.com/garfunkel/go-tvdb"
	"github.com/go-martini/martini"
	"github.com/twinj/uuid"
	// "log"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func ApiProcessFiles(userId string) {
	store := *GetStoreSession()
	patterns := []*regexp.Regexp{
		regexp.MustCompile("[Ss]([0-9]+)[][ ._-]*[Ee]([0-9]+)([^\\/]*).(avi|mkv)$"),
		regexp.MustCompile(`[\\/\._ \[\(-]([0-9]+)x([0-9]+)([^\\/]*).(avi|mkv)$`),
	}

	var userFiles []UserFile
	seriesIDs := make(map[string]int)

	db := GetDbSession()

	_, err := db.Select(&userFiles, "select * from users_files where processed=0 and user_id=?", userId)
	if err == nil && len(userFiles) > 0 {
		for index, userFile := range userFiles {
			// fmt.Println(userFile, index)
			var seriesName string
			var seasonID int
			var episodeID int
			var seriesID int

			seps := [2]string{"\\", "/"}

			for _, sep := range seps {
				relativePath := strings.TrimLeft(userFile.RelativePath, "/\\")
				pathElems := strings.Split(relativePath, sep)
				if len(pathElems) > 1 {
					seriesName = pathElems[0]
					break
				} else {
					seriesName = ""
				}
			}

			for _, pattern := range patterns {
				matches := pattern.FindAllStringSubmatch(userFile.RelativePath, -1)
				if len(matches) > 0 && len(matches[0]) > 0 {
					seasonID, _ = strconv.Atoi(matches[0][1])
					episodeID, _ = strconv.Atoi(matches[0][2])
					break
				}
			}

			if val, ok := seriesIDs[seriesName]; ok {
				seriesID = val
			} else {
				show, err := store.GetShowOrRetrieveFromTitle(seriesName)
				if err == nil && show != nil {
					seriesIDs[seriesName] = show.ID
					seriesID = show.ID
					userFiles[index].ShowID = seriesID
					userFiles[index].EpisodeID = episodeID
					userFiles[index].SeasonID = seasonID
					userFiles[index].Processed = true

					_, err := db.Update(&userFiles[index])
					if err != nil {
						fmt.Println(err)
					}
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
	go func(userId string) {
		ApiProcessFiles(userId)
	}(user.ID)
	return http.StatusOK, encoder.Must(enc.Encode(userFiles))
}
