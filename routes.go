package main

import (
	"basicServer/Controller"
	"basicServer/constants"
	"net/http"
)

func initilize() {
	/*
	*	Routing the requests for different endpoints to Different handlers via Controller
	*	GuestlistHandle	: To handle requests for '/guest_list' requests
	*	GuestLogHandle	: To handle requests for '/guests' and '/guests/'
	*	SeatHandle		: To handle requests for '/seats_empty'.
	*	And started to listen on 'localhost:8080'
	 */
	http.HandleFunc("/guests", Controller.GuestLogHandle)
	http.HandleFunc("/guests/", Controller.GuestLogHandle)
	http.HandleFunc("/seats_empty", Controller.SeatHandle)
	http.HandleFunc("/guest_list", Controller.GuestlistHandle)
	http.HandleFunc("/guest_list/", Controller.GuestlistHandle)
	err := http.ListenAndServe(constants.APPLICATION_HOSTIP+":"+constants.APPLICATION_PORT, nil)
	if err != nil {
		panic(err.Error())
	}
}
