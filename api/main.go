package main

import (
	"api/controllers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// New Router
	r := mux.NewRouter()
	r.Use(commonMiddleware)
	r.HandleFunc("/employee", controllers.GetEmployees).Methods("GET")
	r.HandleFunc("/employee/{id}", controllers.GetEmployee).Methods("GET")
	r.HandleFunc("/employee", controllers.CreateEmployee).Methods("POST")
	r.HandleFunc("/employee/{id}", controllers.UpdateEmployee).Methods("PUT")
	r.HandleFunc("/employee/{id}", controllers.DeleteEmployee).Methods("DELETE")

	fmt.Println("Starting server at port 8080")
	//	start the server
	log.Fatal(http.ListenAndServe(":8080", r))

}
func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
