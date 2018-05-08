package customers

import "gopkg.in/mgo.v2/bson"

// Represents a customer, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document
type Customer struct {
    ID          bson.ObjectId `bson:"_id" json:"id"`
    Name        string        `bson:"name" json:"name"`
    Address     string        `bson:"address" json:"address"`
    Phone       string        `bson:"phone" json:"phone"`
    Email       string        `bson:"email" json:"email"`
}
