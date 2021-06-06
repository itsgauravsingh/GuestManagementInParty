package main

func main() {
	//	To setup routes and start listening at localhost:8080
	initilize()
}

/*
	User {
		type : admin / guest / cordinator
	}
	Guest {
		name, accompanying_guest, tableId
 	}
	Table {
		tableID, capacity, state, venueId
	}
	GuestHistory {
		guestID, isPresent, time_arrived, time_departed, accompanying_guests
	}
	Venue {
		name, id, address
	}

createGuest()
createTable()

### Add a guest to the guestlist
	POST /guest_list/name
	create guest and add to table

### Get the guest list
	GET /guest_list
	read the guestTable and return the list
	need to check for pagination

### Guest Arrives
	PUT /guests/name
	Validate(name)	//	Guest exists, Table capacity constraint, alreadyPresent
	insert into GuestHistory //	update the arrived time, isPresent, accompanying_list

### Guest Leaves
	DELETE /guests/name

### Get arrived guests
	return isPresent=true
	// Add for pagination

### Count number of empty seats
	GET /seats_empty
	check GuestHistory for occupiedSeat
	check Table for totalCapacity
	return the difference

concurrency
goroutine
channel
*/