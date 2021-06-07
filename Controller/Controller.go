/*
*	Controller will handle the routing of valid request to specific controller and discard the invalid requests
 */

package Controller

import (
	"basicServer/GuestController"
	"net/http"
)

func GuestlistHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GuestController.ShowGuest(w, r)
		return
	case "POST":
		GuestController.CreateGuest(w, r)
		return
	default:
		_, _ = w.Write([]byte("method not allowed"))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func GuestLogHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GuestController.ShowGuestLog(w, r)
		return
	case "PUT":
		GuestController.CreateGuestLog(w, r)
		return
	case "DELETE":
		GuestController.DeleteGuestLog(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("method not allowed"))
		return
	}
}

func SeatHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GuestController.GetCapacity(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("method not allowed"))
		return
	}
}
