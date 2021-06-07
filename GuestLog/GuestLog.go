package GuestLog

import (
	"basicServer/Guest"
	"basicServer/constants"
	"database/sql"
	"time"
)

const (
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

/*
*	CreateNew(guest Guest.Guest) error
*	Input params
*	@guest	:	guest details wrapped in a Guest struct
*	Output params
*	@error	:	error in case of any failure else nil.
*
*	Flow
*	1.	Create a DB connection and defer the call to close the connection
*	2.	Check if entry for this guest is already present. If no, create and entry for this and return
*	3.	If the entry is already present, update the record
*	4.	In case of any error, return error
*	5.	In case of success, return nil
*
 */

func CreateNew(guest Guest.Guest) error {
	//	ToDO
	db, err := sql.Open("mysql", constants.DB_USER+":"+constants.DB_PASSWORD+"@tcp("+constants.DB_HOSTIP+":"+constants.DB_PORT+")/"+constants.DB_NAME)
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

/*
*	ShowGuestLog() ([]AbstractGuestLog, error)
*	Input params
*	Output params
*	@[]AbstractGuest	:	guest list wrapped in AbstractGuest containing all the relevant details only
*	@error	:	error in case of any failure else nil.
*
*	Flow
*	1.	Create a DB connection and defer the call to close the connection
*	2.	Using the db handle run select query to get all the records where ispresent is set to 1 (indicating the present of guest) but only desired fields from GuestLog table
*	3.	In case of any error, return (empty slice of AbstractGuest,error)
*	4.	In case of success, return (valid slice of AbstractGuest, nil)
*
 */

func ShowGuestLog() ([]AbstractGuestLog, error) {
	db, err := sql.Open("mysql", constants.DB_USER+":"+constants.DB_PASSWORD+"@tcp("+constants.DB_HOSTIP+":"+constants.DB_PORT+")/"+constants.DB_NAME)
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

/*
*	Delete(name string) error
*	Input params
*	@name	:	string value depicting name of a guest to be deleted
*	Output params
*	@error	:	error in case of any failure else nil.
*
*	Flow
*	1.	Create a DB connection and defer the call to close the connection
*	2.	Update the ispresent flag, and time_departed for the guest in the GuestLog table. Marking the guest as departed
*	3.	In case of any error, return error
*	4.	In case of success, return nil
*
 */

func Delete(name string) error {
	db, err := sql.Open("mysql", constants.DB_USER+":"+constants.DB_PASSWORD+"@tcp("+constants.DB_HOSTIP+":"+constants.DB_PORT+")/"+constants.DB_NAME)
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
