package main

import (
	"api/router"
	"fmt"
	"log"
	"net/http"
)

func main() {

	r := router.Router()
	r.Use(CommonMiddleware)
	fmt.Println("Starting server at port 8080")
	//	start the server
	log.Fatal(http.ListenAndServe(":8080", r))

}
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

//alter sequence employee_id_seq restart with n;
