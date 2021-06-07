package Table

import (
	"basicServer/constants"
	"database/sql"
)

const (
	DB_TABLE = "TableInfo"
)

/*
*	Struct : Table
*	To model the TableInfo entity
 */

type Table struct {
	ID       int    `json:"id,omitempty"`
	State    string `json:"state,omitempty"`
	Capacity int    `json:"capacity,omitempty"`
	VenueID  int    `json:"venue_id,omitempty"`
}

/*
**	Creates and inserts a table record to TableInfo table
**	Input Params
**	@id : integer type value to be used as table identifier
**	@state : string type value to depict the current state of the table. Total states can be 1. available, 2. occupied, 3. broken
**	@capacity : integer type value to depict the capacity i.e. how many members can this table accommodate
**	@venueid : integer type value to associate each table to a specific venue
**	Return Param
**	@error : returns error if there is any else returns nil
**
**	This Method is based on the following understanding
**
 */

func Create(id int, state string, capacity int, venueid int) error {
	db, err := sql.Open("mysql", constants.DB_USER+":"+constants.DB_PASSWORD+"@tcp("+constants.DB_HOSTIP+":"+constants.DB_PORT+")/"+constants.DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	//	we need to close the db handle
	defer db.Close()

	//	Insert the new entry into the database
	table := Table{
		ID: id, State: state, Capacity: capacity, VenueID: venueid,
	}

	_, err = db.Query("INSERT into "+DB_TABLE+"(id, state, capacity, venueid) VALUES (:id, :state, :capacity, :venueid)", table)
	if err != nil {
		return err
	}
	return nil
}

func Update() {
	// ToDo
}

func GetCapacity() (int, error) {
	db, err := sql.Open("mysql", constants.DB_USER+":"+constants.DB_PASSWORD+"@tcp("+constants.DB_HOSTIP+":"+constants.DB_PORT+")/"+constants.DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	//	we need to close the db handle
	defer db.Close()

	results, err := db.Query("SELECT SUM(capacity) FROM " + DB_TABLE)
	if err != nil {
		panic(err.Error())
	}
	capacity := 0
	for results.Next() {
		err = results.Scan(&capacity)
		if err != nil {
			panic(err.Error())
		}
	}
	return capacity, nil
}

func GetCapacityByTable(table int) (int, error) {
	db, err := sql.Open("mysql", constants.DB_USER+":"+constants.DB_PASSWORD+"@tcp("+constants.DB_HOSTIP+":"+constants.DB_PORT+")/"+constants.DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	//	we need to close the db handle
	defer db.Close()
	results, err := db.Query("SELECT capactiy FROM "+DB_TABLE+" WHERE id = ?", table)
	if err != nil {
		panic(err.Error())
	}
	capacity := 0
	for results.Next() {
		err = results.Scan(&capacity)
		if err != nil {
			panic(err.Error())
		}
	}
	return capacity, nil
}
