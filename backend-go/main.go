package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	FirstName string `json:"firstName" gorm:"column:firstName"`
	LastName  string `json:"lastName" gorm:"column:lastName"`
	EmailId   string `json:"emailId" gorm:"column:emailId"`
}

var employees []Employee

var dbname = "mygoapp"
var bddUser = "root"
var bddPassword = "samir"
var bddPort = 3306

func main() {
	// Connexion à la base de données MySQL
	db, err := GormConnectToMySQLDatabase(bddUser, bddPassword, "nil", bddPort)
	if err != nil {
		panic("Erreur lors de la connexion a la base de donnees")
	}

	// Fermeture automatique de la connexion BDD
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	// Création de la BDD si elle n'existe pas
	if err := createDatabaseIfNotExists(db, "mygoapp"); err != nil {
		fmt.Println(err.Error())
		return
	}

	// Selection de la BDD mygoapp sur la connexion existante
	if err := db.Exec("USE mygoapp;").Error; err != nil {
		panic("BDD non trouve")
	}

	// Création de la table employee si elle n'existe pas
	if err := createTableIfNotExists(db); err != nil {
		fmt.Println(err.Error())
		return
	}

	// Création d'un Routeur http
	rtr := mux.NewRouter()

	// Mise en place des routes
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

	db, err := GormConnectToMySQLDatabase(bddUser, bddPassword, dbname, bddPort)
	if err != nil {
		log.Printf("Erreur lors de la connexion à la base de données: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Fermeture automatique de la connexion BDD
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	// Récupération des employés depuis la base de donnée.
	var employees []Employee
	result := db.Find(&employees)
	if result.Error != nil {
		log.Printf("Erreur lors de la récupération des employés depuis la base de données: %s", result.Error)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal employees slice to JSON
	employeesJSON, err := json.Marshal(employees)
	if err != nil {
		log.Printf("Erreur lors de la conversion des données des employés en JSON: %s", err)
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

	// Connexion BDD
	db, err := GormConnectToMySQLDatabase(bddUser, bddPassword, dbname, bddPort)
	if err != nil {
		log.Printf("Erreur lors de la connexion à la base de données: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Fermeture automatique de la connexion BDD
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	// La requette sql pour trouvver l'employee par ID
	var employee Employee
	if err := db.Where("id = ?", id).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(rw, "Employee non trouve", http.StatusNotFound)
			return
		}
		log.Printf("Erreur de la requette: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	employeeJSON, err := json.Marshal(employee)
	if err != nil {
		log.Printf("Erreur lors du marshaling employee en JSON: %s", err)
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
		fmt.Printf("Erreur lors de la lecture du body: %v", err)
		http.Error(rw, "impossible de lire de body", http.StatusBadRequest)
	}

	// Unmarshal du json envoyé par le front pour le traiter en GO
	var employee Employee
	err = json.Unmarshal(bytesBody, &employee)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Connexion BDD
	db, err := GormConnectToMySQLDatabase(bddUser, bddPassword, dbname, bddPort)
	if err != nil {
		log.Printf("Erreur lors de la connexion à la base de données: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Fermeture automatique de la connexion BDD
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	// Vérification du dernier ID.
	var lastEmployee Employee
	if err := db.Order("id desc").Limit(1).Find(&lastEmployee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			employee.ID = 1
		} else {
			log.Printf("Erreur dans la requette employee: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		employee.ID = lastEmployee.ID + 1
	}

	// Insertion de l'employé dans la base de données
	if err := db.Create(&employee).Error; err != nil {
		log.Printf("Erreur lors de la creation de l'employee: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	// Connexion BDD
	db, err := GormConnectToMySQLDatabase(bddUser, bddPassword, dbname, bddPort)
	if err != nil {
		log.Printf("Erreur lors de la connexion à la base de données: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Fermeture automatique de la connexion BDD
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	// Suppresion de l'employee avec l'ID
	result := db.Delete(&Employee{}, id)
	if result.Error != nil {
		log.Fatalf("Error deleting employee: %v", result.Error)
	}

	// Vérification de la suppression
	rowsAffected := result.RowsAffected
	if rowsAffected == 0 {
		http.Error(rw, "L'employee n'as pas ete trouve", http.StatusNotFound)
		return
	}

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
		http.Error(rw, "ID de l'employee invalide", http.StatusBadRequest)
		return
	}

	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Erreur lors de la lecture du body: %v", err)
		http.Error(rw, "impossible de lire le body", http.StatusBadRequest)
	}

	// Transformation du JSON en struc pour un traitement en GO
	var updatedEmployee Employee
	err = json.Unmarshal(bytesBody, &updatedEmployee)
	if err != nil {
		fmt.Println(err)
	}

	// Connexion BDD
	db, err := GormConnectToMySQLDatabase(bddUser, bddPassword, dbname, bddPort)
	if err != nil {
		log.Printf("Erreur lors de la connexion à la base de données: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Fermeture automatique de la connexion BDD
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	// Recherche de l'employee dans la bdd
	var employee Employee
	result := db.First(&employee, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(rw, "Employee non trouve", http.StatusNotFound)
		return
	} else if result.Error != nil {
		log.Printf("Erreur dans la recherche d'un employe dans la base de donnees: %v", result.Error)
		http.Error(rw, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Mise a jour de l'employee en base
	result = db.Model(&employee).Updates(updatedEmployee)
	if result.Error != nil {
		log.Printf("Erreur lors de la mise a jour d'un employé dans la base de donnees: %v", result.Error)
		http.Error(rw, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Vérification de la mise a jour
	rowsAffected := result.RowsAffected
	if rowsAffected == 0 {
		http.Error(rw, "Employee non trouve", http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func createDatabaseIfNotExists(db *gorm.DB, dbname string) error {
	// Vérifier si la base de données existe déjà
	var count int64
	db.Raw(fmt.Sprintf("SELECT count(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%s'", dbname)).Count(&count)
	if count > 0 {
		return nil
	}

	// Créer la base de données si elle n'existe pas
	if err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname)).Error; err != nil {
		return err
	}

	return nil
}

func createTableIfNotExists(db *gorm.DB) error {
	// Vérifier si la table "employees" existe déjà
	if db.Migrator().HasTable(&Employee{}) {
		return nil
	}

	// Créer la table "employees" si elle n'existe pas encore
	err := db.AutoMigrate(&Employee{})
	if err != nil {
		return err
	}

	return nil
}

func GormConnectToMySQLDatabase(user, password, dbname string, port int) (*gorm.DB, error) {
	var dsn string
	if dbname == "nil" {
		dsn = fmt.Sprintf("%s:%s@tcp(tools-mybddapp-1:%d)/?charset=utf8mb4&parseTime=True&loc=Local", user, password, port)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(tools-mybddapp-1:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, port, dbname)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
