package main

import (
	"fmt"
	"log"
	"testing"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestMongoDB(t *testing.T) {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("people")
	err = c.Insert(&Person{"superWang", "13478808311"},
		&Person{"David", "15040268074"})
	if err != nil {
		log.Fatal(err)
	}
	result := Person{}
	err = c.Find(bson.M{"name": "superWang"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Name:", result.Name)
	fmt.Println("Phone:", result.Phone)
}
