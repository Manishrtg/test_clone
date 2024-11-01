package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

// Patient represents a patient record
type Patient struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Age         int     `json:"age"`
	Doctor      string  `json:"doctor"`
	Prescription *string `json:"prescription,omitempty"` // Pointer to handle NULL values
}

// DB connection string
const connStr = "postgres://postgres:Lodjadcjn4@localhost:5432/count?sslmode=disable"

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open a DB connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	fmt.Println("Connected to the database!")

	// Set up the HTTP router
	r := mux.NewRouter()

	// Define the API endpoints for receptionists
	r.HandleFunc("/patients", addPatientHandler).Methods("POST")        // Receptionist adds a patient
	r.HandleFunc("/patients/{id}", getPatientHandler).Methods("GET")   // Get patient details
	r.HandleFunc("/patients/{id}", updatePatientHandler).Methods("PUT") // Update patient details (Receptionist)
	r.HandleFunc("/patients/{id}", deletePatientHandler).Methods("DELETE") // Delete patient record

	// Define the API endpoints for doctors
	r.HandleFunc("/patients/{id}/prescription", updatePrescriptionHandler).Methods("PUT") // Update prescription

	// Start the HTTP server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// addPatientHandler adds a new patient
func addPatientHandler(w http.ResponseWriter, r *http.Request) {
	var patient Patient
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var id int
	err := db.QueryRow("INSERT INTO patients(name, age, doctor) VALUES($1, $2, $3) RETURNING id",
		patient.Name, patient.Age, patient.Doctor).Scan(&id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert patient: %v", err), http.StatusInternalServerError)
		return
	}
	patient.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(patient)
}

// getPatientHandler retrieves a patient by ID
func getPatientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	var patient Patient
	var prescription sql.NullString // Use sql.NullString to handle NULL values
	err = db.QueryRow("SELECT id, name, age, doctor, prescription FROM patients WHERE id=$1", id).
		Scan(&patient.ID, &patient.Name, &patient.Age, &patient.Doctor, &prescription)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Patient not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to query patient: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Set the prescription field based on the sql.NullString
	if prescription.Valid {
		patient.Prescription = &prescription.String
	} else {
		patient.Prescription = nil // Set to nil if the prescription is NULL
	}

	json.NewEncoder(w).Encode(patient)
}

// updatePatientHandler updates a patient record (Receptionist)
func updatePatientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	var patient Patient
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	patient.ID = id

	_, err = db.Exec("UPDATE patients SET name=$1, age=$2, doctor=$3 WHERE id=$4",
		patient.Name, patient.Age, patient.Doctor, patient.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update patient: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(patient)
}

// updatePrescriptionHandler updates the prescription for a patient (Doctor)
func updatePrescriptionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	var patient Patient
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	patient.ID = id

	_, err = db.Exec("UPDATE patients SET prescription=$1 WHERE id=$2",
		patient.Prescription, patient.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update prescription: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(patient)
}

// deletePatientHandler removes a patient record
func deletePatientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM patients WHERE id=$1", id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete patient: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}





// package main

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strconv"

// 	_ "github.com/lib/pq"
// 	"github.com/gorilla/mux"
// )

// // Patient represents a patient record
// type Patient struct {
// 	ID     int    `json:"id"`
// 	Name   string `json:"name"`
// 	Age    int    `json:"age"`
// 	Doctor string `json:"doctor"`
// }

// // DB connection string
// const connStr = "postgres://postgres:Lodjadcjn4@localhost:5432/hospital?sslmode=disable"

// var db *sql.DB

// func main() {
// 	var err error
// 	db, err = sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatalf("Failed to open a DB connection: %v", err)
// 	}
// 	defer db.Close()

// 	if err := db.Ping(); err != nil {
// 		log.Fatalf("Failed to connect to the database: %v", err)
// 	}
// 	fmt.Println("Connected to the database!")

// 	// Set up the HTTP router
// 	r := mux.NewRouter()

// 	// Define the API endpoints
// 	r.HandleFunc("/patients", addPatientHandler).Methods("POST")
// 	r.HandleFunc("/patients/{id}", getPatientHandler).Methods("GET")
// 	r.HandleFunc("/patients/{id}", updatePatientHandler).Methods("PUT")
// 	r.HandleFunc("/patients/{id}", deletePatientHandler).Methods("DELETE")

// 	// Start the HTTP server
// 	log.Println("Starting server on :8080")
// 	if err := http.ListenAndServe(":8080", r); err != nil {
// 		log.Fatalf("Failed to start server: %v", err)
// 	}
// }

// // addPatientHandler adds a new patient
// func addPatientHandler(w http.ResponseWriter, r *http.Request) {
// 	var patient Patient
// 	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	var id int
// 	err := db.QueryRow("INSERT INTO patients(name, age, doctor) VALUES($1, $2, $3) RETURNING id",
// 		patient.Name, patient.Age, patient.Doctor).Scan(&id)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to insert patient: %v", err), http.StatusInternalServerError)
// 		return
// 	}
// 	patient.ID = id
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(patient)
// }

// // getPatientHandler retrieves a patient by ID
// func getPatientHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
// 		return
// 	}

// 	var patient Patient
// 	err = db.QueryRow("SELECT id, name, age, doctor FROM patients WHERE id=$1", id).
// 		Scan(&patient.ID, &patient.Name, &patient.Age, &patient.Doctor)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			http.Error(w, "Patient not found", http.StatusNotFound)
// 		} else {
// 			http.Error(w, fmt.Sprintf("Failed to query patient: %v", err), http.StatusInternalServerError)
// 		}
// 		return
// 	}
// 	json.NewEncoder(w).Encode(patient)
// }

// // updatePatientHandler updates a patient record
// func updatePatientHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
// 		return
// 	}

// 	var patient Patient
// 	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	patient.ID = id

// 	_, err = db.Exec("UPDATE patients SET name=$1, age=$2, doctor=$3 WHERE id=$4",
// 		patient.Name, patient.Age, patient.Doctor, patient.ID)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to update patient: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(patient)
// }

// // deletePatientHandler removes a patient record
// func deletePatientHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
// 		return
// 	}

// 	_, err = db.Exec("DELETE FROM patients WHERE id=$1", id)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to delete patient: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }


