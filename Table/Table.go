package Table

import "database/sql"

const (
	DB_TABLE = "TableInfo"
	DB_NAME = "partydb"
)

type Table struct {
	ID       int    `json:"id,omitempty"`
	State    string `json:"state,omitempty"`
	Capacity int    `json:"capacity,omitempty"`
	VenueID  int    `json:"venue_id,omitempty"`
}

func Create(id int, state string, capacity int, venueid int) error  {
	db, err := sql.Open("mysql","root:bhaPP@123@tcp(127.0.0.1:3306)/"+DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	//	we need to close the db handle
	defer db.Close()
	//	Check if table exists

	//	Insert the new entry into the database
	table := Table{
		ID:id,State: state,Capacity: capacity,VenueID: venueid,
	}

	_ , err = db.Query("INSERT into "+DB_TABLE + "(id, state, capacity, venueid) VALUES (:id, :state, :capacity, :venueid)", table)
	if err != nil {
		return err
	}
	return nil
}

func Update() {

}

func GetCapacity() (int, error) {
	db, err := sql.Open("mysql","root:bhaPP@123@tcp(127.0.0.1:3306)/"+DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	//	we need to close the db handle
	defer db.Close()

	results, err := db.Query("SELECT SUM(capacity) FROM "+DB_TABLE)
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

func GetCapacityByTable(table int)  (int, error) {
	db, err := sql.Open("mysql","root:bhaPP@123@tcp(127.0.0.1:3306)/"+DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	//	we need to close the db handle
	defer db.Close()
	results, err := db.Query("SELECT capactiy FROM "+DB_TABLE+" WHERE id = ?",table)
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