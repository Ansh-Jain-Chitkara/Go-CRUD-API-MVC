package models

type Employee struct {
	ID      int64  `json:"id"`
	EmpName string `json:"empname"`
	EmpPRO  string `json:"emppro"`
}
