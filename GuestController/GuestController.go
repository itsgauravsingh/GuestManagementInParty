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

//	POST /guest_list/name
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
	var guest Guest.Guest
	err = json.Unmarshal(bodyBytes, &guest)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Could read the input body due to %s", err.Error())))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = Guest.Create(&guest)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Failed to create new record. Error says %s", err.Error())))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

//	GET /guest_list
func ShowGuest(w http.ResponseWriter, r *http.Request) {
	guests, err := Guest.Guestlist()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Failed to pull the records due to %s", err.Error())))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	guestlist := make(map[string][]Guest.AbstractGuest)
	guestlist["guests"] = guests
	jsonBytes, err := json.Marshal(guestlist)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Failed to pull the records due to %s", err.Error())))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
	w.WriteHeader(http.StatusOK)
	return
}

// PUT /guests/name
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
	name := strings.Replace(parts[2], "%20", " ", -1)
	if isValidGuest(name) == false {
		w.Write([]byte("Invalid guest"))
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

// Get /guests
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

//	DELETE /guests/name
func DeleteGuestLog(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) < 3 || len(parts[2]) < 1 {
		//	ToDo
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//	For now handing the space by simply replacing the %20 with whitespace. But we can create a method to sanitise the input.
	name := strings.Replace(parts[2], "%20", " ", -1)
	err := GuestLog.Delete(name)
	checkForErrorAndPanic(err)
	w.WriteHeader(http.StatusOK)
}

// GET /seats_empty
func GetCapacity(w http.ResponseWriter, r *http.Request) {
	// Run through guestlog and calculate the presentGuestCount
	// Ask TableController for the tatalCapacity
	// return the difference
	guests, err := GuestLog.ShowGuestLog()
	checkForErrorAndPanic(err)
	presentGuestCount := 0
	for idx := range guests {
		presentGuestCount += guests[idx].AccompanyingGuests + 1
	}
	venueCapacity, err := TableController.GetCapacity()
	checkForErrorAndPanic(err)
	emptySeatCount := venueCapacity - presentGuestCount
	emptyseat := make(map[string]int)
	emptyseat["seats_empty"] = emptySeatCount
	jsonBytes, err := json.Marshal(emptyseat)
	checkForErrorAndPanic(err)
	w.Write(jsonBytes)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	return
}

func GetDetailsForGuest(name string) (Guest.Guest, error) {
	return getGuestDetails(name)
}

func getGuestDetails(name string) (Guest.Guest, error) {
	return Guest.GetGuestDetails(name)
}

func isGuestCountAcceptable(guestCount int, guest Guest.Guest) (int, error) {
	if guestCount <= guest.AccompanyingGuests {
		return guestCount, nil
	}
	//	Check against the table capacity
	tableCapacity, err := TableController.GetTableCapacity(guest.Table)
	if err != nil {
		panic(err.Error())
	}
	if tableCapacity >= guestCount {
		return guestCount, nil
	}
	return 0, nil
}

func isValidGuest(name string) bool {
	_, err := getGuestDetails(name)
	checkForErrorAndPanic(err)
	return true
}

func checkForErrorAndPanic(err error) {
	if err != nil {
		panic(err.Error())
	}
}
