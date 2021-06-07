package TableController

import "basicServer/Table"

/*
*	TableController
*	@Functions
*	GetTableCapacity(table int) (int,error)
*		Input params
*		@table : integer value depicting a tableid
*		Output params
*		@param1 : integer value depicting the capacity of the table with id as tableid (provided as input)
*		@param2 : error in case some failure occurs during the detail fetch else nil.
*
*	Flow
*	1.	Makes call to respective function and returns the pair of (tableCapacity, error)
 */

func GetTableCapacity(table int) (int, error) {
	return Table.GetCapacityByTable(table)
}

/*
*	GetCapacity() (int, error)
*		Input params
*		Output params
*		@int : total capacity of the venue in terms of seat count
*		@error : valid error in case of any failure. nil in case of success
*
*	Flow
*	1.	Makes call to respective function and returns the pair of (totalCapacity, error)
 */

func GetCapacity() (int, error) {
	return Table.GetCapacity()
}
