package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/WhoBrokeTheBuild/GoCMS/types"
	"github.com/gorilla/mux"
)

type adminPageData struct {
	types.User
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error
	var uid int64
	var token string

	for _, c := range r.Cookies() {
		if c.Name == "gocmsuser" {
			uid, err = strconv.ParseInt(c.Value, 10, 64)
			if err != nil {
				log.Println(err)
			}
		} else if c.Name == "gocmsauth" {
			token = c.Value
		}
	}

	res, err := db.Query("SELECT * FROM `Sessions` WHERE `UserID` = ? AND `Token` = ?", uid, token)
	if err != nil {
		log.Println(err)
		return
	}

	s := types.Session{}
	s.Scan(res)

	if !s.IsExpired() {

	}

	_ = vars
	//	if path, ok := vars["page"]; ok {

	//	}
}
