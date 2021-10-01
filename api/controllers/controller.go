package controllers

import (
	"api/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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

var jwtKey = []byte("secret_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword, ok := users[credentials.Username]

	if !ok || expectedPassword != credentials.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

}

// Display All the details from the database
func GetEmployees(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	db := createConnection()

	var emp models.Employee
	errr := json.NewDecoder(r.Body).Decode(&emp)
	if errr != nil {
		log.Panic("Unable to decode the request body")
	}
	insertDynStmt := `insert into "employee"( "empname","emppro") values($1, $2) returning id`
	id := 0
	errr = db.QueryRow(insertDynStmt, emp.EmpName, emp.EmpPRO).Scan(&id)
	if errr != nil {
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
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	db := createConnection()

	params := mux.Vars(r)
	var emp models.Employee
	errr := json.NewDecoder(r.Body).Decode(&emp)
	if errr != nil {
		log.Panic("Unable to decode the request body")
	}
	id := 0
	exe := `update employee set empname = $1, emppro = $2 where id = $3 returning id`
	errr = db.QueryRow(exe, emp.EmpName, emp.EmpPRO, params["id"]).Scan(&id)
	if errr != nil {
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
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	db := createConnection()

	params := mux.Vars(r)
	id := 0
	exe := `delete from employee where id = $1 returning id`
	errr := db.QueryRow(exe, params["id"]).Scan(&id)
	if errr != nil {
		log.Panic("Unable to execute the query")
	}
	res := JsonResponse{
		ID:      int64(id),
		Message: "Row Deleted",
	}
	json.NewEncoder(w).Encode(res)

}
func Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	expirationTime := time.Now().Add(time.Minute * 5)

	claims.ExpiresAt = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "refresh_token",
			Value:   tokenString,
			Expires: expirationTime,
		})

}
