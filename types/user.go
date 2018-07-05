package types

import (
	"database/sql"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User is a struct representation of the a Users table entry
type User struct {
	ID       int64
	Email    string
	HashPass []byte
	Salt     []byte
	Created  time.Time
	Updated  time.Time
}

// CheckPassword checks a given password against the hashed password of the user
func (u *User) CheckPassword(pass []byte) bool {
	for i := range u.Salt {
		pass = append(pass, u.Salt[i])
	}
	err := bcrypt.CompareHashAndPassword(u.HashPass, pass)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Scan reads the SQL result set into a User struct
func (u *User) Scan(res *sql.Rows) error {
	if res.Next() {
		err := res.Scan(&u.ID, &u.Email, &u.HashPass, &u.Salt, &u.Created, &u.Updated)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
