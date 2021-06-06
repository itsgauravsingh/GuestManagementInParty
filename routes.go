package main

import (
	"basicServer/Controller"
	"net/http"
)

func initilize()  {
	http.HandleFunc("/guests", Controller.GuestLogHandle)
	http.HandleFunc("/guests/", Controller.GuestLogHandle)
	http.HandleFunc("/seats_empty", Controller.SeatHandle)
	http.HandleFunc("/guest_list", Controller.GuestlistHandle)
	http.ListenAndServe(":8080",nil)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}