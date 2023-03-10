package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetEmployees(t *testing.T) {
	// Créer une fausse requête HTTP
	req, err := http.NewRequest("GET", "/employees", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Créer un ResponseWriter mock
	rw := httptest.NewRecorder()

	// Appeler la fonction getEmployees avec les paramètres requis
	getEmployees(rw, req)

	// Vérifier que le code de réponse est 200 (OK)
	if status := rw.Code; status != http.StatusOK {
		t.Errorf("Le code de réponse est incorrect: got %v, attendu %v", status, http.StatusOK)
	}

	// Vérifier que le Content-Type du header est "application/json"
	expectedContentType := "application/json"
	if contentType := rw.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Le Content-Type du header est incorrect: got %v, attendu %v", contentType, expectedContentType)
	}

	// Vérifier que le contenu de la réponse est une liste d'employés en JSON
	expectedJSON := `[{"id":1,"nom":"Dupont","prenom":"Jean","email":"jean.dupont@example.com"},{"id":2,"nom":"Martin","prenom":"Marie","email":"marie.martin@example.com"}]`
	if rw.Body.String() != expectedJSON {
		t.Errorf("Le contenu de la réponse est incorrect: got %v, attendu %v", rw.Body.String(), expectedJSON)
	}
}

func TestGetEmployee(t *testing.T) {
	// Créer une fausse requête HTTP avec un paramètre d'ID valide
	req, err := http.NewRequest("GET", "/employees/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Créer un ResponseWriter mock
	rw := httptest.NewRecorder()

	// Créer un routeur mux et enregistrer la fonction getEmployee
	router := mux.NewRouter()
	router.HandleFunc("/employees/{id}", getEmployee)

	// Appeler le routeur mux pour gérer la requête
	router.ServeHTTP(rw, req)

	// Vérifier que le code de réponse est 200 (OK)
	if status := rw.Code; status != http.StatusOK {
		t.Errorf("Le code de réponse est incorrect: got %v, attendu %v", status, http.StatusOK)
	}

	// Vérifier que le Content-Type du header est "application/json"
	expectedContentType := "application/json"
	if contentType := rw.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Le Content-Type du header est incorrect: got %v, attendu %v", contentType, expectedContentType)
	}

	// Vérifier que le contenu de la réponse est l'employé avec l'ID correspondant en JSON
	expectedJSON := `{"id":1,"nom":"Dupont","prenom":"Jean","email":"jean.dupont@example.com"}`
	if rw.Body.String() != expectedJSON {
		t.Errorf("Le contenu de la réponse est incorrect: got %v, attendu %v", rw.Body.String(), expectedJSON)
	}

	// Créer une fausse requête HTTP avec un paramètre d'ID invalide
	req, err = http.NewRequest("GET", "/employees/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Réinitialiser le ResponseWriter mock
	rw = httptest.NewRecorder()

	// Appeler le routeur mux pour gérer la requête
	router.ServeHTTP(rw, req)

	// Vérifier que le code de réponse est 400 (Bad Request)
	if status := rw.Code; status != http.StatusBadRequest {
		t.Errorf("Le code de réponse est incorrect: got %v, attendu %v", status, http.StatusBadRequest)
	}

	// Vérifier que le contenu de la réponse est "Invalid employee ID"
	expectedBody := "Invalid employee ID\n"
	if rw.Body.String() != expectedBody {
		t.Errorf("Le contenu de la réponse est incorrect: got %v, attendu %v", rw.Body.String(), expectedBody)
	}

	// Créer une fausse requête HTTP avec un paramètre d'ID non existant
	req, err = http.NewRequest("GET", "/employees/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Réinitialiser le ResponseWriter mock
	rw = httptest.NewRecorder()

	// Appeler le routeur mux pour gérer la requête
	router.ServeHTTP(rw, req)

	// Vérifier que le code de réponse est 404 (Not Found)
	if status := rw.Code; status != http.StatusNotFound {
		t.Errorf("Le code de réponse est incorrect: got %v, attendu %v", status, http.StatusNotFound)
	}

	// Vérifier que le contenu de la réponse est "Employee non trouve"
	expectedBody = "Employee non trouve\n"
	if rw.Body.String() != expectedBody {
		t.Errorf("Le contenu de la réponse est incorrect")
	}
}

func TestPostEmployee(t *testing.T) {
	// Créer un employee à ajouter
	employee := Employee{ID: 1, FirstName: "Doe", LastName: "John", EmailId: "john.doe@example.com"}
	requestBody, err := json.Marshal(employee)
	if err != nil {
		t.Fatal(err)
	}

	// Créer une fausse requête HTTP POST avec le body correspondant
	req, err := http.NewRequest("POST", "/employees", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Créer un ResponseWriter mock
	rw := httptest.NewRecorder()

	// Appeler la fonction postEmployee avec la fausse requête HTTP
	postEmployee(rw, req)

	// Vérifier que le code de réponse est 200 (OK)
	if status := rw.Code; status != http.StatusOK {
		t.Errorf("Le code de réponse est incorrect: got %v, attendu %v", status, http.StatusOK)
	}

	// Vérifier que l'employé a bien été ajouté dans la base de données
	db, err := GormConnectToMySQLDatabase(bddUser, bddPassword, dbname, bddPort)
	if err != nil {
		t.Fatal(err)
	}
	// Fermeture automatique de la connexion BDD
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	var result Employee
	if err := db.First(&result, employee.ID).Error; err != nil {
		t.Fatal(err)
	}
	if result != employee {
		t.Errorf("L'employé ajouté est incorrect: got %v, attendu %v", result, employee)
	}
}
