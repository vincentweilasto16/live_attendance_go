package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (server *Server) initializeRoutes() {
	server.Router = mux.NewRouter()

	server.Router.HandleFunc("/", server.Home).Methods("GET")
	server.Router.HandleFunc("/clock-in", server.ClockIn).Methods("POST")
	server.Router.HandleFunc("/clock-out", server.ClockOut).Methods("POST")

	server.Router.HandleFunc("/register", server.Register).Methods("GET")

	server.Router.HandleFunc("/register", server.AddEmployee).Methods("POST")

	server.Router.HandleFunc("/login", server.Login).Methods("GET")
	server.Router.HandleFunc("/login", server.DoLogin).Methods("POST")
	server.Router.HandleFunc("/logout", server.Logout).Methods("GET")

	server.Router.HandleFunc("/logout", server.Logout).Methods("GET")

	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/public/", http.FileServer(staticFileDirectory))
	server.Router.PathPrefix("/public/").Handler(staticFileHandler).Methods("GET")
}
