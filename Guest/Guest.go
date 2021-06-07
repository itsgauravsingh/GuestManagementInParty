package Guest

import (
	"basicServer/constants"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	DB_TABLE = "Guest"
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

/*
*	Create(guest *Guest) (*Guest, error)
*	Input params
*	@guest	:	guest details wrapped in a Guest struct
*	Output params
*	@*Guest	:	pointer to guest details wrapped in a Guest struct
*	@error	:	error in case of any failure else nil.
*
*	Flow
*	1.	Create a DB connection and defer the call to close the connection
*	2.	Using the db handle run insert query to create a guest record int the Guest table
*	3.	In case of any error, return (guest,error)
*	4.	In case of success, return (guest, nil)
*
 */

func Create(guest Guest) (Guest, error) {
	//	ToDO
	db, err := sql.Open("mysql", constants.DB_USER+":"+constants.DB_PASSWORD+"@tcp("+constants.DB_HOSTIP+":"+constants.DB_PORT+")/"+constants.DB_NAME)
	if nil != err {
		panic(err.Error())
	}
	defer db.Close()

	//	Insert the new entry into the database
	_, err = db.Query("INSERT into "+DB_TABLE+"(guestname, accompanying_guests, tableId) VALUES (?, ?, ?)", guest.Name, guest.AccompanyingGuests, guest.Table)
	if err != nil {
		return Guest{}, err
	}
	return guest, err
}

/*
*	Guestlist() ([]AbstractGuest, error)
*	Input params
*	Output params
*	@[]AbstractGuest	:	guest list wrapped in AbstractGuest containing all the relevant details only
*	@error	:	error in case of any failure else nil.
*
*	Flow
*	1.	Create a DB connection and defer the call to close the connection
*	2.	Using the db handle run select query to get all the records but only desired fields from Guest table
*	3.	In case of any error, return (empty slice of AbstractGuest,error)
*	4.	In case of success, return (valid slice of AbstractGuest, nil)
*
 */

func Guestlist() ([]AbstractGuest, error) {
	db, err := sql.Open("mysql", constants.DB_USER+":"+constants.DB_PASSWORD+"@tcp("+constants.DB_HOSTIP+":"+constants.DB_PORT+")/"+constants.DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	//	we need to close the db handle
	defer db.Close()
	var guestlist []AbstractGuest
	results, err := db.Query("SELECT guestname,accompanying_guests,tableId FROM " + DB_TABLE)
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

/*
*	GetGuestDetails(name string) (Guest, error)
*	Input params
*	@name	:	name of the guest for which details are asked
*	Output params
*	@Guest	:	guest details wrapped in Guest struct
*	@error	:	error in case of any failure else nil.
*
*	Flow
*	1.	Create a DB connection and defer the call to close the connection
*	2.	Using the db handle run select query to get the record corresponding to the guest we are looking for.
*	3.	In case of any error, return (empty Guest struct,error)
*	4.	In case of success, return (valid Guest struct, nil)
*
 */

func GetGuestDetails(name string) (Guest, error) {
	db, err := sql.Open("mysql", constants.DB_USER+":"+constants.DB_PASSWORD+"@tcp("+constants.DB_HOSTIP+":"+constants.DB_PORT+")/"+constants.DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	//	we need to close the db handle
	defer db.Close()
	var guest Guest
	err = db.QueryRow("SELECT id, guestname, tableId, accompanying_guests FROM "+DB_TABLE+" WHERE guestname = ?", name).
		Scan(&guest.ID, &guest.Name, &guest.Table, &guest.AccompanyingGuests)
	if err != nil {
		return Guest{}, err
	}
	return guest, nil
}
