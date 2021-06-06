package Controller

import (
	"basicServer/GuestController"
	"net/http"
)

func GuestlistHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GuestController.ShowGuest(w,r)
		return
	case "POST":
		GuestController.CreateGuest(w,r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func GuestLogHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GuestController.ShowGuestLog(w,r)
		return
	case "PUT":
		GuestController.CreateGuestLog(w,r)
		return
	case "DELETE":
		GuestController.DeleteGuestLog(w,r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func SeatHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GuestController.GetCapacity(w,r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}


