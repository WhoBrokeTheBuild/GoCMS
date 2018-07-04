package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type config struct {
	ID      int64
	KeyName string
	Value   string
	Created time.Time
	Updated time.Time
}

func (c *config) Scan(res *sql.Rows) {
	if res.Next() {
		var value []byte

		err := res.Scan(&c.ID, &c.KeyName, &value, &c.Created, &c.Updated)
		if err != nil {
			log.Println(err)
		}

		c.Value = string(value)
	}
}

type page struct {
	ID       int64
	Template string
	Path     string
	Title    string
	Content  template.HTML
	Created  time.Time
	Updated  time.Time
}

func (p *page) Scan(res *sql.Rows) {
	if res.Next() {
		var tmpl []byte
		var path []byte
		var title []byte
		var content []byte

		err := res.Scan(&p.ID, &tmpl, &path, &title, &content, &p.Created, &p.Updated)
		if err != nil {
			log.Println(err)
		}

		p.Template = string(tmpl[:])
		p.Path = string(path[:])
		p.Title = string(title[:])
		p.Content = template.HTML(string(content[:]))
	}
}

var templates = template.Must(parseTemplates())
var db *sql.DB

func main() {
	var err error

	str := fmt.Sprintf("%s:%s@%s/%s", dbUser, dbPass, dbHost, dbName)

	db, err = sql.Open("mysql", str)
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	rtr := mux.NewRouter()
	rtr.HandleFunc("/", handlePage)
	rtr.HandleFunc("/{page}", handlePage)
	rtr.HandleFunc("/admin/{page}", handleAdmin)
	rtr.HandleFunc("/api/{api}", handleAPI)
	rtr.Handle("/static", http.FileServer(http.Dir("static")))

	log.Println("Listening on :8080")

	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
}

func handlePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data := page{}

	log.Println(vars)

	if path, ok := vars["page"]; ok {
		if path == "favicon.ico" {
			http.Redirect(w, r, "/static/favicon.ico", 301)
			return
		}

		log.Println("Looking for", path)
		res, err := db.Query("SELECT * FROM `Pages` WHERE `Path` = ?", path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data.Scan(res)
		res.Close()
	} else {
		log.Println("No {page}, assume index")
		// No value for {page}, assume index
		res, err := db.Query("SELECT * FROM `Config` WHERE `KeyName` = 'Index'")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		c := config{}
		c.Scan(res)
		res.Close()

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
		log.Println(data)
		res.Close()
	}

	log.Println("Page: ", data.Title, data.Content)

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

func handleAdmin(w http.ResponseWriter, r *http.Request) {

}

func handleAPI(w http.ResponseWriter, r *http.Request) {

}

func parseTemplates() (*template.Template, error) {
	tmpl := template.New("")
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
