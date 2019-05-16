package main

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
)

const (
	dbName         = "github"
	collectionName = "awesome"
	mongoURL       = "mongodb://mongo_user:mongo_secret@localhost:27017/github?authSource=admin"
	mongoURLTest   = "mongodb://mongo_user:mongo_secret@localhost:27017/github_test?authSource=admin"
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

	_, err := coll.Upsert(bson.M{"reponame": r.RepoName}, r)
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
	// collation := mgo.Collation{Locale: "en", Strength: 1}

	err := coll.Find(bson.M{"reponame": r.RepoName}).One(&r)
	if err != nil {
		log.Printf("Find %s in MongoDB error: %v", r.FullName, err)
	}
	return err
}

func newMongoSession(url string) (*mgo.Session, error) {
	return mgo.Dial(url)
}

// NewMongoStore return a PRSaver use MongoDB
func NewMongoStore(test bool) RepoInfoStore {
	url := mongoURL
	if test {
		url = mongoURLTest
	}
	session, err := newMongoSession(url)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v\n", err)
	}

	return MongoStore{
		session: session,
	}
}
