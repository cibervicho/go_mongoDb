package main

import (
	"fmt"
	"os"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
)

const (
	url = "localhost:27017"
)

func main() {
	// set the default mongodb database
	my_DB := "local"
	if 1 == len(os.Args) {
		log.Warnf("No db specified, using '%v' for now", my_DB)
	} else {
		// set the specified mongodb database in command line
		my_DB = os.Args[1]
	}

	// connecting to the mongodb server
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal(err)
	}

	defer session.Close()
	log.Infof("Successfully connected to mongodb server at %v", url)

	// list all available databases
	dbNames, err := session.DatabaseNames()
	if err != nil {
		log.Warn(err)
	}
	for i, dbName := range dbNames {
		fmt.Printf("[%2v] - %v\n", i+1, dbName)
	}

	db := session.DB(my_DB)
	if db == nil {
		log.Errorf("db %v could not be found, exiting...", my_DB)
		return
	}

	// iterate the collections
	fmt.Printf("Collections found in db '%v':\n", my_DB)
	cols, err := db.CollectionNames()
	if err != nil {
		log.Warnf("No collections found in db '%v'", my_DB)
	}
	for i, c := range cols {
		fmt.Printf("[%2v] - %v\n", i+1, c)
		//listDocs(db, c)
	}

	colName := "unicorns"
	fmt.Printf("Querying collection %v\n", colName)
	listDocs(db, colName)
}

func listDocs(db *mgo.Database, col string) {
	coll := db.C(col)
	if coll == nil || col == "system.profile" {
		return
	}

	// main.Document{ID:0, Dob:time.Time{wall:0x0, ext:62836091220, loc:(*time.Location)(nil)}, Gender:"m", Loves:[]string{"carrot", "papaya"}, Name:"Horny", Vaccinated:true, Vampires:63, Weight:600}
	type Document struct {
		ID         bson.ObjectId `json:"id,omitempty"         bson:"_id,omitempty"`
		Dob        time.Time     `json:"dob,omitempty"        bson:"dob,omitempty"`
		Gender     string        `json:"gender,omitempty"     bson:"gender,omitempty"`
		Loves      []string      `json:"loves,omitempty"      bson:"loves,omitempty"`
		Name       string        `json:"name,omitempty"       bson:"name,omitempty"`
		Vaccinated bool          `json:"vaccinated,omitempty" bson:"vaccinated,omitempty"`
		Vampires   int           `json:"vampires,omitempty"   bson:"vampires,omitempty"`
		Weight     int           `json:"weight,omitempty"     bson:"weight,omitempty"`
	}

	//var result []interface{}
	//var result []map[string]interface{}     // []bson.M
	var result []Document
	coll.Find(nil).All(&result)
	for i, _ := range result {
		//fmt.Printf("\tDoc%2v - %v\n", i+1, d)
		fmt.Printf("%v - Name: %#v, Gender: %#v\n", result[i].ID, result[i].Name, result[i].Gender)
	}
}
