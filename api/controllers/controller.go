package controllers

import (
	"api/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "234ajanshul"
	dbname   = "employeeDB"
)

func createConnection() *sql.DB {

	psqlconn := fmt.Sprintf("host=%s port =%d user= %s password =%s dbname =%s sslmode = disable", host, port, user, password, dbname)
	db, err := sql.Open(user, psqlconn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()

	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")

	return db
}

type JsonResponse struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// Display All the details from the database
func GetEmployees(w http.ResponseWriter, r *http.Request) {
	db := createConnection()

	var emps []models.Employee
	emps = nil
	row, _ := db.Query("Select * from employee")
	for row.Next() {
		var emp models.Employee
		row.Scan(&emp.ID, &emp.EmpName, &emp.EmpPRO)
		emps = append(emps, emp)
	}
	defer row.Close()
	json.NewEncoder(w).Encode(emps)
	emps = nil
}

// Get an Employee using a particular id
func GetEmployee(w http.ResponseWriter, r *http.Request) {
	db := createConnection()

	params := mux.Vars(r)
	st := `select * from employee where id = $1`

	row, _ := db.Query(st, params["id"])
	var emp models.Employee
	for row.Next() {
		row.Scan(&emp.ID, &emp.EmpName, &emp.EmpPRO)
	}
	defer row.Close()
	json.NewEncoder(w).Encode(emp)
}

// Create A particular Employee
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	db := createConnection()

	var emp models.Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		log.Panic("Unable to decode the request body")
	}
	insertDynStmt := `insert into "employee"( "empname","emppro") values($1, $2) returning id`
	id := 0
	err = db.QueryRow(insertDynStmt, emp.EmpName, emp.EmpPRO).Scan(&id)
	if err != nil {
		log.Panic("Unable to execute the query")
	}
	res := JsonResponse{
		ID:      int64(id),
		Message: "Row created",
	}
	json.NewEncoder(w).Encode(res)

}

// Update the Employee
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	db := createConnection()

	params := mux.Vars(r)
	var emp models.Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		log.Panic("Unable to decode the request body")
	}
	id := 0
	exe := `update employee set empname = $1, emppro = $2 where id = $3 returning id`
	err = db.QueryRow(exe, emp.EmpName, emp.EmpPRO, params["id"]).Scan(&id)
	if err != nil {
		log.Panic("Unable to execute the query")
	}
	res := JsonResponse{
		ID:      int64(id),
		Message: "Row Updated",
	}
	json.NewEncoder(w).Encode(res)

}

// Delete an Employee
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	db := createConnection()

	params := mux.Vars(r)
	id := 0
	exe := `delete from employee where id = $1 returning id`
	err := db.QueryRow(exe, params["id"]).Scan(&id)
	if err != nil {
		log.Panic("Unable to execute the query")
	}
	res := JsonResponse{
		ID:      int64(id),
		Message: "Row Deleted",
	}
	json.NewEncoder(w).Encode(res)

}
