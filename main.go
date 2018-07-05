package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"./types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var funcs = template.FuncMap{
	"getMenu": getMenu,
}

var templates = template.Must(parseTemplates())
var db *sql.DB

type pageData struct {
	types.Page

	Year int
}

type debugHandler struct {
	rtr *mux.Router
}

func (d debugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	templates = template.Must(parseTemplates())
	d.rtr.ServeHTTP(w, r)
}

func main() {
	var err error

	port := flag.String("port", "8080", "port")
	debug := flag.Bool("debug", false, "debug")
	flag.Parse()

	str := fmt.Sprintf("%s:%s@%s/%s?parseTime=true", dbUser, dbPass, dbHost, dbName)

	db, err = sql.Open("mysql", str)
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	rtr := mux.NewRouter()
	rtr.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	rtr.HandleFunc("/admin/{page}", handleAdmin)
	rtr.HandleFunc("/api/{api}", handleAPI)
	rtr.HandleFunc("/{page}", handlePage)
	rtr.HandleFunc("/", handlePage)

	log.Printf("Listening on :%s\n", *port)

	if *debug {
		http.Handle("/", debugHandler{rtr})
	} else {
		http.Handle("/", rtr)
	}
	http.ListenAndServe(":"+*port, nil)
}

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

func handleAdmin(w http.ResponseWriter, r *http.Request) {

}

func handleAPI(w http.ResponseWriter, r *http.Request) {

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

func getMenu(name string) []types.MenuItem {
	res, err := db.Query("SELECT * FROM `Menus` WHERE `Name` = ?", name)
	if err != nil {
		log.Println(err)
		return nil
	}

	m := types.Menu{}
	err = m.Scan(res)
	if err != nil {
		log.Println(err)
		return nil
	}

	return m.GetItems(db)
}
