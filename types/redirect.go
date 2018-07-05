package types

import (
	"database/sql"
	"log"
)

// Redirect is a struct representation of the a Redirects table entry
type Redirect struct {
	ID      int64
	OldPath string
	NewPath string
	Code    int
}

// Scan reads the SQL result set into a Redirect struct
func (r *Redirect) Scan(res *sql.Rows) error {
	if res.Next() {
		err := res.Scan(&r.ID, &r.OldPath, &r.OldPath, &r.Code)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
