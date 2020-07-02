package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	contactBook = append(contactBook, Contact{ID: 1, Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	contactBook = append(contactBook, Contact{ID: 2, Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	contactBook = append(contactBook, Contact{ID: 3, Firstname: "Francis", Lastname: "Sunday"})

	router := mux.NewRouter()
	router.HandleFunc("/contact/", GetContacts).Methods("GET")
	router.HandleFunc("/contact/{id}", GetContact).Methods("GET")
	router.HandleFunc("/contact/", CreateContact).Methods("POST")
	router.HandleFunc("/contact/{id}", DeleteContact).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// GetContacts privodes an array of contacts
func GetContacts(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(contactBook) }

// GetContact provides a contact by id
func GetContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range contactBook {
		ID, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if item.ID == ID {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Contact{})
}

// CreateContact creates a new contact by id
func CreateContact(w http.ResponseWriter, r *http.Request) {
	var newContact Contact
	_ = json.NewDecoder(r.Body).Decode(&newContact)
	newContact.ID = getNextID()
	contactBook = append(contactBook, newContact)
	json.NewEncoder(w).Encode(newContact)
}

// DeleteContact deletes a contact by id
func DeleteContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var found bool = false
	for index, item := range contactBook {
		if item.ID == ID {
			contactBook = append(contactBook[:index], contactBook[index+1:]...)
			found = true
			break
		}
	}

	if found {
		json.NewEncoder(w).Encode(contactBook)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	return

}

// Contact struct
type Contact struct {
	ID        int      `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

// getNextID method iterates over contactBook and get the next valid ID
func getNextID() int {
	var greater int = 0
	for _, item := range contactBook {
		if item.ID > greater {
			greater = item.ID
		}
	}
	return greater + 1
}

// Address struct
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var contactBook []Contact
