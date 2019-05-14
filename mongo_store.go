package main

import (
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	dbName         = "github"
	collectionName = "awesome"
	mongoURL       = "mongodb://mongo_user:mongo_secret@localhost:27017/github?authSource=admin"
)

// GetCollectionName return collectionName
func GetCollectionName() string {
	return collectionName
}

// MongoStore is an implement of interface RepoInfoStore,
// store the RepoInfo in MongoDB
type MongoStore struct {
	session *mgo.Session
}

// Save the RepoInfo to MongoDB
func (s MongoStore) Save(r RepoInfo) error {
	session := s.session.Copy()
	defer session.Close()
	coll := session.DB("").C(collectionName)

	_, err := coll.Upsert(bson.M{"full_name": r.FullName}, r)
	if err != nil {
		log.Printf("Upsert to MongoDB error: %v", err)
		return err
	}
	return nil
}

// Load RepoInfo from MongoDB
func (s MongoStore) Load(r *RepoInfo) error {
	session := s.session.Copy()
	defer session.Close()
	coll := session.DB("").C(collectionName)

	err := coll.Find(bson.M{"full_name": r.FullName}).One(&r)
	if err != nil {
		log.Printf("Find in MongoDB error: %v", err)
	}
	return err
}

func newMongoSession() (*mgo.Session, error) {
	return mgo.Dial(mongoURL)
}

// NewMongoStore return a PRSaver use MongoDB
func NewMongoStore() RepoInfoStore {
	session, err := newMongoSession()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v\n", err)
	}

	return MongoStore{
		session: session,
	}
}
