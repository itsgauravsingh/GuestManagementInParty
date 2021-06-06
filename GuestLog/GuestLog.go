package GuestLog

import (
	"basicServer/Guest"
	"database/sql"
	"time"
)

const (
	DB_NAME  = "partydb"
	DB_TABLE = "GuestLog"
)

type GuestLog struct {
	ID                 int    `json:"id,omitempty"`
	GuestID            string `json:"guest_id,omitempty"`
	AccompanyingGuests int    `json:"accompanying_guests,omitempty"`
	IsPresent          int    `json:"is_present,omitempty"`
	TimeArrived        string `json:"time_arrived"`
	TimeDeparted       string `json:"time_departed"`
}

type AbstractGuestLog struct {
	GuestID            string `json:"guest_id,omitempty"`
	AccompanyingGuests int    `json:"accompanying_guests,omitempty"`
	TimeArrived        string `json:"time_arrived"`
}

func CreateNew(guest Guest.Guest) error {
	//	ToDO
	db, err := sql.Open("mysql", "root:bhaPP@123@tcp(127.0.0.1:3306)/"+DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	//	we need to close the db handle
	defer db.Close()

	// Check if the user is already present
	res, err := db.Query("SELECT  * FROM "+DB_TABLE+" WHERE guestid=?", guest.Name)
	if err != nil {
		panic(err.Error())
	}
	if res.Next() == false {
		_, err = db.Query("INSERT into "+DB_TABLE+"(guestid,ispresent,time_arrived,accompanying_guest) VALUES (?,?,?,?)", guest.Name, 1, time.Now().Format(time.RFC3339), guest.AccompanyingGuests)
		if err != nil {
			panic(err.Error())
		}
		return nil
	}
	//	Update the already present entry
	_, err = db.Query("UPDATE "+DB_TABLE+" SET ispresent = 1, time_arrived = ?, accompanying_guest = ? WHERE guestid = ?", time.Now().Format(time.RFC3339), guest.AccompanyingGuests, guest.Name)
	if err != nil {
		return err
	}
	return nil
}

func ShowGuestLog() ([]AbstractGuestLog, error) {
	db, err := sql.Open("mysql", "root:bhaPP@123@tcp(127.0.0.1:3306)/"+DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	//	we need to close the db handle
	defer db.Close()
	//	Insert the new entry into the database
	results, err := db.Query("SELECT guestid, accompanying_guest, time_arrived FROM " + DB_TABLE + " WHERE ispresent = 1")
	if err != nil {
		return []AbstractGuestLog{}, err
	}
	var guests []AbstractGuestLog
	for results.Next() {
		var abstractGuestLog AbstractGuestLog
		err = results.Scan(&abstractGuestLog.GuestID, &abstractGuestLog.AccompanyingGuests, &abstractGuestLog.TimeArrived)
		if err != nil {
			panic(err.Error())
		}
		guests = append(guests, abstractGuestLog)
	}
	return guests, nil
}

func Delete(name string) error {
	db, err := sql.Open("mysql", "root:bhaPP@123@tcp(127.0.0.1:3306)/"+DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	//	we need to close the db handle
	defer db.Close()
	//	Insert the new entry into the database
	_, err = db.Query("UPDATE "+DB_TABLE+" SET ispresent = 0, time_departed = ? WHERE guestid = ?", time.Now().Format(time.RFC3339), name)
	if err != nil {
		return err
	}
	return nil
}
