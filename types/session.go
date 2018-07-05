package types

import (
	"database/sql"
	"log"
	"time"
)

// Session is a struct representation of the a Sessions table entry
type Session struct {
	ID      int64
	UserID  int64
	Token   []byte
	Expires time.Time
}

// IsExpired checks the current time against the Expires timestamp
func (s *Session) IsExpired() bool {
	return time.Now().After(s.Expires)
}

// Scan reads the SQL result set into a Session struct
func (s *Session) Scan(res *sql.Rows) error {
	if res.Next() {
		err := res.Scan(&s.ID, &s.UserID, &s.Token, &s.Expires)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
