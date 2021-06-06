package Guest

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	DB_TABLE = "Guest"
	DB_NAME = "partydb"
)

type Guest struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Table              int    `json:"table_id"`
	AccompanyingGuests int    `json:"accompanying_guests"`
}

type AbstractGuest struct {
	Name               string `json:"name"`
	Table              int    `json:"table_id"`
	AccompanyingGuests int    `json:"accompanying_guests"`
}

func Create(guest *Guest) (*Guest, error) {
	//	ToDO
	db, err := sql.Open("mysql","root:bhaPP@123@tcp(127.0.0.1:3306)/"+DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	//	we need to close the db handle
	defer db.Close()

	//	Insert the new entry into the database
	_ , err = db.Query("INSERT into "+DB_TABLE + "(guestname, accompanying_guests, tableId) VALUES (:name, :accompanying_guests, :table_id", *guest)
	if err != nil {
		return &Guest{},err
	}
	return guest, err
}

func Guestlist() ([]AbstractGuest, error) {
	db, err := sql.Open("mysql","root:bhaPP@123@tcp(127.0.0.1:3306)/"+DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	//	we need to close the db handle
	defer db.Close()
	var guestlist []AbstractGuest
	results, err := db.Query("SELECT guestname,accompanying_guests,tableId FROM "+DB_TABLE)
	if err != nil {
		return []AbstractGuest{}, err
	}
	for results.Next() {
		var guest AbstractGuest
		err = results.Scan(&guest.Name, &guest.Table, &guest.AccompanyingGuests)
		if err != nil {
			panic(err.Error())
		}
		guestlist = append(guestlist, guest)
	}
	return guestlist, err
}

func GetGuestDetails(name string) (Guest, error) {
	db, err := sql.Open("mysql","root:bhaPP@123@tcp(127.0.0.1:3306)/"+DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	//	we need to close the db handle
	defer db.Close()
	var guest Guest
	err = db.QueryRow("SELECT id, guestname, tableId, accompanying_guests FROM "+DB_TABLE + " WHERE guestname = ?",name).
		Scan(&guest.ID, &guest.Name, &guest.Table, &guest.AccompanyingGuests)
	if err != nil {
		return Guest{}, err
	}
	return guest, err
}