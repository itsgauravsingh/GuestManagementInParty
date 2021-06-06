package TableController

import "basicServer/Table"

func GetTableCapacity(table int) (int,error) {
	return Table.GetCapacityByTable(table)
}

func GetCapacity() (int, error) {
	return Table.GetCapacity()
}


