package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Employee struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	EmailId   string `json:"emailId"`
}

var employees []Employee

func main() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/api/v1/employees", getEmployees).Methods("GET")
	rtr.HandleFunc("/api/v1/employees/{id}", getEmployee).Methods("GET")
	rtr.HandleFunc("/api/v1/employees", postEmployee).Methods("POST")
	rtr.HandleFunc("/api/v1/employees/{id}", deleteEmployee).Methods("DELETE")
	rtr.HandleFunc("/api/v1/employees/{id}", putEmployee).Methods("PUT")

	// Ajouter un middleware pour la gestion des en-têtes CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
	})

	handler := c.Handler(rtr)

	http.Handle("/", handler)

	log.Println("Listening...")
	http.ListenAndServe(":8081", nil)

	//handler := cors.Default().Handler(mux)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func getEmployees(rw http.ResponseWriter, r *http.Request) {
	// Set Content-Type header to application/json
	rw.Header().Set("Content-Type", "application/json")

	// Marshal employees slice to JSON
	employeesJSON, err := json.Marshal(employees)
	if err != nil {
		log.Printf("Error marshaling employees to JSON: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write employeesJSON to response body
	rw.Write(employeesJSON)
}

func getEmployee(rw http.ResponseWriter, r *http.Request) {
	// Set Content-Type header to application/json
	rw.Header().Set("Content-Type", "application/json")

	// Récupération de l'ID dans l'URL
	vars := mux.Vars(r)
	strId := vars["id"]
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(rw, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	// Find employee by ID
	index := -1
	for i, emp := range employees {
		if emp.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		http.Error(rw, "Employee not found", http.StatusNotFound)
		return
	}

	// Marshal employees slice to JSON
	employeeJSON, err := json.Marshal(employees[index])
	if err != nil {
		log.Printf("Error marshaling employees to JSON: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write employeesJSON to response body
	rw.Write(employeeJSON)
}

func postEmployee(rw http.ResponseWriter, r *http.Request) {
	// Set Content-Type header to application/json
	rw.Header().Set("Content-Type", "application/json")

	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v", err)
		http.Error(rw, "can't read body", http.StatusBadRequest)
	}

	var employee Employee
	err = json.Unmarshal(bytesBody, &employee)
	if err != nil {
		fmt.Println(err)
	}

	// Vérification du dernier ID.
	if len(employees) > 0 {
		lastEmployee := employees[len(employees)-1]
		employee.ID = lastEmployee.ID + 1
	} else {
		employee.ID = 1
	}

	employees = append(employees, employee)
	rw.WriteHeader(http.StatusOK)
}

func deleteEmployee(rw http.ResponseWriter, r *http.Request) {
	// Set Content-Type header to application/json
	rw.Header().Set("Content-Type", "application/json")

	// Récupération de l'ID dans l'URL
	vars := mux.Vars(r)
	strId := vars["id"]
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(rw, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	// Find employee by ID and remove from slice
	index := -1
	for i, emp := range employees {
		if emp.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		http.Error(rw, "Employee not found", http.StatusNotFound)
		return
	}
	employees = append(employees[:index], employees[index+1:]...)

	rw.WriteHeader(http.StatusOK)
}

func putEmployee(rw http.ResponseWriter, r *http.Request) {
	// Set Content-Type header to application/json
	rw.Header().Set("Content-Type", "application/json")

	// Récupération de l'ID dans l'URL
	vars := mux.Vars(r)
	strId := vars["id"]
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(rw, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v", err)
		http.Error(rw, "can't read body", http.StatusBadRequest)
	}

	var updatedEmployee Employee
	err = json.Unmarshal(bytesBody, &updatedEmployee)
	//fmt.Println(updatedEmployee.FirstName)
	if err != nil {
		fmt.Println(err)
	}

	for i, employee := range employees {
		if employee.ID == id {
			updatedEmployee.ID = id
			if updatedEmployee.FirstName == "" {
				updatedEmployee.FirstName = employees[i].FirstName
			}
			if updatedEmployee.LastName == "" {
				updatedEmployee.LastName = employees[i].LastName
			}
			if updatedEmployee.EmailId == "" {
				updatedEmployee.EmailId = employees[i].EmailId
			}
			employees[i] = updatedEmployee
			break
		}
	}

	rw.WriteHeader(http.StatusOK)
}
