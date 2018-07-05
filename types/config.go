package types

import (
	"database/sql"
	"log"
	"time"
)

// Config is a struct representation of a Config table entry
type Config struct {
	ID      int64
	Key     string
	Value   string
	Created time.Time
	Updated time.Time
}

// Scan reads the SQL result set into a Config struct
func (c *Config) Scan(res *sql.Rows) error {
	if res.Next() {
		var value []byte

		err := res.Scan(&c.ID, &c.Key, &value, &c.Created, &c.Updated)
		if err != nil {
			log.Println(err)
			return err
		}

		c.Value = string(value)
	}
	return nil
}
