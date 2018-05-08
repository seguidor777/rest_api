package dao

import (
    "log"

    . "../models"
    mgo "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type CustomersDAO struct {
    Server   string
    Database string
}

var db *mgo.Database

const (
    COLLECTION = "customers"
)

// Establish a connection to database
func (m *CustomersDAO) Connect() {
    session, err := mgo.Dial(m.Server)
    if err != nil {
        log.Fatal(err)
    }
    db = session.DB(m.Database)
}

// Find list of customers
func (m *CustomersDAO) FindAll() ([]Customer, error) {
    var customers []Customer
    err := db.C(COLLECTION).Find(bson.M{}).All(&customers)
    return customers, err
}

// Find a customer by its id
func (m *CustomersDAO) FindById(id string) (Customer, error) {
    var customer Customer
    err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&customer)
    return customer, err
}

// Insert a customer into database
func (m *CustomersDAO) Insert(customer Customer) error {
    err := db.C(COLLECTION).Insert(&customer)
    return err
}

// Delete an existing customer
func (m *CustomersDAO) Delete(customer Customer) error {
    err := db.C(COLLECTION).Remove(&customer)
    return err
}

// Update an existing customer
func (m *CustomersDAO) Update(customer Customer) error {
    err := db.C(COLLECTION).UpdateId(customer.ID, &customer)
    return err
}
