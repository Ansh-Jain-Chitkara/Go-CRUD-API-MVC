package router

import (
	"api/controllers"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/employee", controllers.GetEmployees).Methods("GET")
	router.HandleFunc("/employee/{id}", controllers.GetEmployee).Methods("GET")
	router.HandleFunc("/employee", controllers.CreateEmployee).Methods("POST")
	router.HandleFunc("/employee/{id}", controllers.UpdateEmployee).Methods("PUT")
	router.HandleFunc("/employee/{id}", controllers.DeleteEmployee).Methods("DELETE")
	router.HandleFunc("/refresh", controllers.Refresh).Methods("POST")

	return router

}
