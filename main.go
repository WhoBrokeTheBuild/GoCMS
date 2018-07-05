package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/WhoBrokeTheBuild/GoCMS/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var funcs = template.FuncMap{
	"getMenu": getMenu,
}

var db *sql.DB

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

func handleAPI(w http.ResponseWriter, r *http.Request) {

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
