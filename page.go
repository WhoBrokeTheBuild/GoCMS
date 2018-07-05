package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/WhoBrokeTheBuild/GoCMS/types"
	"github.com/gorilla/mux"
)

type pageData struct {
	types.Page

	Year int
}

var templates = template.Must(parseTemplates())

func handlePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data := pageData{}
	data.Year = time.Now().Year()

	if path, ok := vars["page"]; ok {
		// Hack for favicon
		res, err := db.Query("SELECT * FROM `Redirects` WHERE `OldPath` = ?", path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if res.Next() {
			var rdr types.Redirect
			rdr.Scan(res)
			http.Redirect(w, r, rdr.NewPath, rdr.Code)
			return
		}

		res, err = db.Query("SELECT * FROM `Pages` WHERE `Path` = ?", path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data.Scan(res)
		res.Close()
	} else {
		// No value for {page}, assume index
		res, err := db.Query("SELECT * FROM `Config` WHERE `Key` = 'index'")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		c := types.Config{}
		c.Scan(res)
		res.Close()

		if c.Value == "" {
			http.Error(w, "`index` has no value", http.StatusInternalServerError)
			return
		}

		id, err := strconv.ParseInt(c.Value, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err = db.Query("SELECT * FROM `Pages` WHERE `ID` = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data.Scan(res)
		res.Close()
	}

	if data.ID == 0 {
		http.NotFound(w, r)
		return
	}

	if data.Template == "" {
		data.Template = "default.gohtml"
	}

	err := templates.ExecuteTemplate(w, data.Template, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func parseTemplates() (*template.Template, error) {
	log.Println("Parsing Templates")
	tmpl := template.New("").Funcs(funcs)
	err := filepath.Walk("./templates", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".gohtml") {
			_, err = tmpl.ParseFiles(path)
			if err != nil {
				log.Println(err)
			}
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
