package types

import (
	"database/sql"
	"errors"
	"log"
)

// Menu is a struct representation of the a Menus table entry
type Menu struct {
	ID   int64
	Name string
}

// Scan reads the SQL result set into a Menu struct
func (m *Menu) Scan(res *sql.Rows) error {
	if res.Next() {
		err := res.Scan(&m.ID, &m.Name)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

// GetItems gets the list of Menu Items associated with this Menu
func (m *Menu) GetItems(db *sql.DB) []MenuItem {
	items := []MenuItem{}

	res, err := db.Query("SELECT * FROM `MenuItems` WHERE `MenuID` = ? ORDER BY `Order`", m.ID)
	if err != nil {
		log.Println(err)
		return nil
	}

	for {
		mi := MenuItem{}
		err = mi.Scan(res)
		if err != nil {
			break
		}
		items = append(items, mi)
	}

	return items
}

// MenuItem is a struct representation of the a MenuItems table entry
type MenuItem struct {
	ID     int64
	MenuID int64
	Order  int
	Path   string
	Text   string
}

// Scan reads the SQL result set into a MenuItem struct
func (mi *MenuItem) Scan(res *sql.Rows) error {
	if res.Next() {
		var text []byte

		err := res.Scan(&mi.ID, &mi.MenuID, &mi.Order, &mi.Path, &text)
		if err != nil {
			log.Println(err)
			return err
		}

		mi.Text = string(text)
		return nil
	}
	return errors.New("Out of rows")
}
