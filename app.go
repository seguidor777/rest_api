package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
    . "./config"
    . "./dao"
    . "./models"
)

var config = Config{}
var dao = CustomersDAO{}

// GET list of customers
func CustomersIndex(w http.ResponseWriter, r *http.Request) {
    customers, err := dao.FindAll()
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJson(w, http.StatusOK, customers)
}

// GET a customer by its ID
func ShowCustomer(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    customer, err := dao.FindById(params["id"])
    if err != nil {
	    respondWithError(w, http.StatusBadRequest, "Invalid Customer ID")
        return
    }
    respondWithJson(w, http.StatusOK, customer)
}

// POST a new customer
func CreateCustomer(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    var customer Customer
    if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    customer.ID = bson.NewObjectId()
    if err := dao.Insert(customer); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJson(w, http.StatusCreated, customer)
}

// PUT update an existing customer
func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    defer r.Body.Close()
    var customer Customer
	customer.ID = bson.ObjectIdHex(params["id"])
    if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    if err := dao.Update(customer); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing customer
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    customer, err := dao.FindById(params["id"])
    if err != nil {
	respondWithError(w, http.StatusBadRequest, "Invalid Customer ID")
        return
    }
    if err := dao.Delete(customer); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
    respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
    config.Read()

    dao.Server = config.Server
    dao.Database = config.Database
    dao.Connect()
}

// Define HTTP request routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/customers", CustomersIndex).Methods("GET")
	r.HandleFunc("/customers", CreateCustomer).Methods("POST")
	r.HandleFunc("/customers/{id}", UpdateCustomer).Methods("PUT")
	r.HandleFunc("/customers/{id}", DeleteCustomer).Methods("DELETE")
	r.HandleFunc("/customers/{id}", ShowCustomer).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
