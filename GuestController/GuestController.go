/*
*	Guest Controller will handle all the requests for Guest and GuestLog tables.
*	This is the only way to communicate with Guest and GuestLog entities.
*	To keep things simple, added the seats_empty handling also in this file. This can be moved to Venue controller
 */

package GuestController

import (
	"basicServer/Guest"
	"basicServer/GuestLog"
	"basicServer/TableController"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

/*
*	API : POST /guest_list/name
*
*	CreateGuest(w http.ResponseWriter, r *http.Request)
*	Input params
*	@http.ResponseWriter	:	standard response writer object
*	@http.Request			:	standard http request object
*
*	Flow
*	1. Fetch the 'name' param from the URL
*	2. Check for content-type of the request body
*	3. Check for the request body validity
*	4. Call the Guest.CreateGuest() with the params to create a record in the table
*
 */

func CreateGuest(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		//	ToDo
	}
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}
	if len(parts) < 3 || len(parts[2]) < 1 {
		//	ToDo
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//	Replacing the %20 with a space
	name := strings.Replace(parts[2], "%20", " ", -1)
	inputJsonModel := make(map[string]int)
	err = json.Unmarshal(bodyBytes, &inputJsonModel)
	if err != nil {
		errResponse := make(map[string]string)
		errResponse[name] = err.Error()
		jsonbytes, _ := json.Marshal(errResponse)
		w.Write(jsonbytes)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var guest Guest.Guest
	guest.Name = name
	for key, val := range inputJsonModel {
		if key == "accompanying_guests" {
			guest.AccompanyingGuests = val
		} else if key == "table" {
			guest.Table = val
		}
	}
	//	Calling the Create method for creation the actual record
	_, err = Guest.Create(guest)

	if err != nil {
		errResponse := make(map[string]string)
		errResponse[guest.Name] = err.Error()
		jsonbytes, _ := json.Marshal(errResponse)
		w.Write(jsonbytes)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

/*
*	API : GET /guest_list
*
*	ShowGuest(w http.ResponseWriter, r *http.Request)
*	Input params
*	@http.ResponseWriter	:	standard response writer object
*	@http.Request			:	standard http request object
*
*	Flow
*	1. Fetch all guest records using Guest.Guestlist()
*	2. create a map to create response in desired format
*	3. set the content-type as application/json
*	4. return with statusOK in success
*
 */

func ShowGuest(w http.ResponseWriter, r *http.Request) {
	guests, err := Guest.Guestlist()
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("Failed to pull the records due to %s", err.Error())))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	guestlist := make(map[string][]Guest.AbstractGuest)
	guestlist["guests"] = guests
	jsonBytes, err := json.Marshal(guestlist)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("Failed to pull the records due to %s", err.Error())))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("content-type", "application/json")
	_, _ = w.Write(jsonBytes)
	w.WriteHeader(http.StatusOK)
	return
}

/*
*	API : PUT /guests/name
*
*	CreateGuestLog(w http.ResponseWriter, r *http.Request)
*	Input params
*	@http.ResponseWriter	:	standard response writer object
*	@http.Request			:	standard http request object
*
*	Flow
*	1. Fetch the 'name' param from the URL
*	2. Check for content-type of the request body
*	3. Check for the request body validity
*	4. Check if the guest is valid guest using isValidGuest()
*	5. If Valid guest, fetch the guest details.
*	6. Check if table can accommodate any change in the guest along with the accompanying guest. If yes, allow and create/update record in GuestLog, else return error with response.
*	7. return statusOK if success else return error with message.
*
 */

func CreateGuestLog(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.EscapedPath())
	parts := strings.Split(r.URL.String(), "/")
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		//	ToDo
	}
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}
	if len(parts) < 3 || len(parts[2]) < 1 {
		//	ToDo
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	accompanyingGuests := make(map[string]int)
	aguestKey := "accompanying_guests"
	err = json.Unmarshal(bodyBytes, &accompanyingGuests)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Could read the input body due to %s", err.Error())))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	name := strings.Replace(parts[2], "%20", " ", -1)
	if isValidGuest(name) == false {
		w.Write([]byte("Invalid guest"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	guest, err := getGuestDetails(name)
	checkForErrorAndPanic(err)
	allowedCount, err := isGuestCountAcceptable(accompanyingGuests[aguestKey], guest)
	checkForErrorAndPanic(err)
	if allowedCount == 0 {
		//	ToDo
		w.Write([]byte("Guest count exceeds the allowed limit"))
		w.WriteHeader(http.StatusOK)
		return
	}
	err = GuestLog.CreateNew(guest)
	checkForErrorAndPanic(err)
	w.WriteHeader(http.StatusOK)
	return
}

/*
*	API : Get /guests
*
*	CreateGuestLog(w http.ResponseWriter, r *http.Request)
*	Input params
*	@http.ResponseWriter	:	standard response writer object
*	@http.Request			:	standard http request object
*
*	Flow
*	1. Get all the records from the GuestLog table, where the ispresent flag is set to 1. That means the guest is present at the party.
*	2. create a response in the desired format
*	3. return the response with StatusOK
*
 */

func ShowGuestLog(w http.ResponseWriter, r *http.Request) {
	guests, err := GuestLog.ShowGuestLog()
	checkForErrorAndPanic(err)
	guestlist := make(map[string][]GuestLog.AbstractGuestLog)
	guestlist["guests"] = guests
	jsonBytes, err := json.Marshal(guestlist)
	checkForErrorAndPanic(err)
	w.Write(jsonBytes)
	w.WriteHeader(http.StatusOK)
}

/*
*	API : DELETE /guests/name
*
*	DeleteGuestLog(w http.ResponseWriter, r *http.Request)
*	Input params
*	@http.ResponseWriter	:	standard response writer object
*	@http.Request			:	standard http request object
*
*	Flow
*	1. Fetch the 'name' param from the URL
*	2. Send a request to mark the guest as departed
*	3. return the response with StatusOK
*
 */

func DeleteGuestLog(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")

	//	Check if the URL contains name OR not.
	if len(parts) < 3 || len(parts[2]) < 1 {
		w.Write([]byte("Invalid URL"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//	For now handing the space by simply replacing the %20 with whitespace. But we can create a method to sanitise the input.
	name := strings.Replace(parts[2], "%20", " ", -1)

	//	Sending request to delete the Guest entry from the GuestLog table
	err := GuestLog.Delete(name)
	if err != nil {
		errResponse := make(map[string]string)
		errResponse[name] = "Error occurred while updating the database. Error =" + err.Error()
		jsonbytes, err := json.Marshal(errResponse)
		if err == nil {
			w.Write(jsonbytes)
			w.Header().Add("content-type", "application/json")
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

/*
*	API : GET /seats_empty
*
*	GetCapacity(w http.ResponseWriter, r *http.Request)
*	Input params
*	@http.ResponseWriter	:	standard response writer object
*	@http.Request			:	standard http request object
*
*	Flow
*	1. Fetch the list of guests currently present in the party
*	2. Calculate the count of guests present at party i.e. presentGuestCount.
*	3. Get the total capacity of the venue from the TableInfo entity i.e. venueCapacity
*	4. Calculate the difference and return with statusOK
*
 */

func GetCapacity(w http.ResponseWriter, r *http.Request) {
	// Run through guestlog and calculate the presentGuestCount
	// Ask TableController for the tatalCapacity
	// return the difference
	guests, err := GuestLog.ShowGuestLog()
	if err != nil {
		w.Write([]byte("Failed to get seats count"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	presentGuestCount := 0
	for idx := range guests {
		presentGuestCount += guests[idx].AccompanyingGuests + 1
	}
	venueCapacity, err := TableController.GetCapacity()
	if err != nil {
		w.Write([]byte("Failed to get seats count"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	emptySeatCount := venueCapacity - presentGuestCount
	emptyseat := make(map[string]int)
	emptyseat["seats_empty"] = emptySeatCount
	jsonBytes, err := json.Marshal(emptyseat)
	if err != nil {
		w.Write([]byte("Failed to get seats count"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(jsonBytes)
	w.WriteHeader(http.StatusOK)
	return
}

/*
*	Internal utility method to fetch the guest details from Guest list i.e. Guest table.
*
*	getGuestDetails(name string) (Guest.Guest, error)
*	Input params
*	@name	:	string value depicting name of a guest
*	Output params
*	@Guest	:	guest details wrapped in a Guest struct
*	@error	:	error in case of any failure else nil.
*
*	Flow
*	1. request for a specific record from the Guest table.
*
 */

func getGuestDetails(name string) (Guest.Guest, error) {
	return Guest.GetGuestDetails(name)
}

/*
*	Internal utility method to fetch the guest details from Guest list i.e. Guest table.
*
*	isGuestCountAcceptable(guestCount int, guest Guest.Guest) (int, error)
*	Input params
*	@guestCount	:	int value depicting the count of people arrived at the venue i.e. accompanying guests+1
*	@guest	:	guest details wrapped in a Guest struct
*	Output params
*	@int	:	int depicting the number of seats blocked for the guest and his company.
*	@error	:	error in case of any failure else nil.
*
*	Flow
*	1. Check if arrived guest count is less than what is expected. If yes, return (guestCount, nil)
*	2. If no, fetch the table capacity from TableInfo and see if guests can be accommodated. If yes, return (guestCount, nil)
*	3. If no, return 0, nil
*
 */

func isGuestCountAcceptable(guestCount int, guest Guest.Guest) (int, error) {
	if guestCount <= guest.AccompanyingGuests {
		return guestCount, nil
	}
	//	Check against the table capacity
	tableCapacity, err := TableController.GetTableCapacity(guest.Table)
	if err != nil {
		return 0, err
	}
	if tableCapacity >= guestCount {
		return guestCount, nil
	}
	return 0, nil
}

/*
*	Internal utility method to fetch the guest details from Guest list i.e. Guest table.
*
*	isValidGuest(name string) bool
*	Input params
*	@name	:	string value depicting name of a guest
*	Output params
*	@bool	:	true if guest is a valid guest else false
*
*	Flow
*	1. Check the guest identity against the guest list.
*	2. If guest is present in guest list, return true
*	3. If no, return false
*
 */

func isValidGuest(name string) bool {
	_, err := getGuestDetails(name)
	checkForErrorAndPanic(err)
	return true
}

//	This method is created to create panic if there is an error.
//	Called at different places in the code while implementing the feature. Expected to be replaced with a more graceful handling at the end
func checkForErrorAndPanic(err error) {
	if err != nil {
		panic(err.Error())
	}
}
