package main

import (
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type AddressInfo struct {
	Street  string `json:"street" bson:"street"`
	Code    string `json:"code" bson:"code"`
	City    string `json:"city" bson:"city"`
	State   string `json:"state" bson:"state"`
	Country string `json:"country" bson:"country"`
}

type Institute struct {
	ID            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	InstituteUID  string        `json:"institute_uid" bson:"institute_uid"`
	InstituteName string        `json:"institute_name" bson:"institute_name"`
	Address       AddressInfo   `json:"address" bson:"address"`
}

const (
	MongoDBHosts = "localhost:27017"
	AuthDatabase = "klog"
	AuthUserName = "klog_user"
	AuthPassword = "klog_pwd"
)

func main() {

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  15 * time.Second,
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("failed to create mongodb session: %s\n", err)
	}
	defer mongoSession.Close()

	mongoSession.SetMode(mgo.Monotonic, true)
	log.Println(mongoSession)

	db := mongoSession.DB(AuthDatabase)
	collections, err := db.CollectionNames()
	if err != nil {
		log.Fatalf("failed to load collections: %s\n", err)
	}

	log.Println(db, collections)

	var institute []Institute
	err = db.C("institutes").Find(bson.M{"institute_uid": bson.M{"$regex": "uid-.*1"}}).All(&institute)
	log.Println(err, institute)
}
