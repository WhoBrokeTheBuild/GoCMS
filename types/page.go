package types

import (
	"database/sql"
	"html/template"
	"log"
	"time"
)

// Page is a struct representation of the a Pages table entry
type Page struct {
	ID       int64
	Template string
	Path     string
	Title    string
	Content  template.HTML
	Created  time.Time
	Updated  time.Time
}

// Scan reads the SQL result set into a Page struct
func (p *Page) Scan(res *sql.Rows) error {
	if res.Next() {
		var tmpl []byte
		var path []byte
		var title []byte
		var content []byte

		err := res.Scan(&p.ID, &tmpl, &path, &title, &content, &p.Created, &p.Updated)
		if err != nil {
			log.Println(err)
			return err
		}

		p.Template = string(tmpl[:])
		p.Path = string(path[:])
		p.Title = string(title[:])
		p.Content = template.HTML(string(content[:]))
	}
	return nil
}
