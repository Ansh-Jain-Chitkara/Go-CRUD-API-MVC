package main

import (
	"api/router"
	"fmt"
	"log"
	"net/http"
)

func main() {

	r := router.Router()
	r.Use(setHeaders)
	fmt.Println("Starting server at port 8080")
	//	Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}
func setHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//anyone can make a CORS request
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
		//Since I was building a REST API that returned JSON, I set the content type to JSON here.
		w.Header().Set("Content-Type", "application/json")
		//Allow requests to have the following headers
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, cache-control")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

//alter sequence employee_id_seq restart with n;
