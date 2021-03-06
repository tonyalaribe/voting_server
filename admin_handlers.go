package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/tonyalaribe/voting_server/config"
	"gopkg.in/mgo.v2/bson"
)

type Election struct {
	ID              string
	ElectionName    string
	ElectionType    string
	VotersCount     int
	RunningElection bool
	Positions       []struct {
		Key   int
		Title string
	}
}

type Contestant struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Key      int
	Position string
	Name     string
	Votes    int
}

type Message struct {
	Message string
	Code    int
}

func NewSessionHandler(w http.ResponseWriter, r *http.Request) {
	conf := config.Get()

	election := Election{}
	err := json.NewDecoder(r.Body).Decode(&election)
	if err != nil {
		log.Println(err)
	}

	mgoSession := conf.Database.Session.Copy()
	defer mgoSession.Close()

	collection := conf.Database.C(config.ELECTION_COLLECTION).With(mgoSession)

	election.ID = "election"

	_, err = collection.Upsert(
		bson.M{
			"id": "election",
		}, election)
	if err != nil {
		log.Println(err)
		return
	}

	message := Message{
		Message: "Session created successfully",
	}

	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		log.Println(err)
	}
}

func EndSessionHandler(w http.ResponseWriter, r *http.Request) {
	conf := config.Get()

	mgoSession := conf.Database.Session.Copy()
	defer mgoSession.Close()

	// err := conf.Database.With(mgoSession).DropDatabase()
	// if err != nil {
	// 	log.Println(err)
	// }

	conf.Database.C(config.ELECTION_COLLECTION).With(mgoSession).DropCollection()

	conf.Database.C(config.CONTESTANTS_COLLECTION).With(mgoSession).DropCollection()

	conf.Database.C(config.VOTERS_COLLECTION).With(mgoSession).DropCollection()

	message := struct {
		Message string
	}{
		Message: "Session ended successfully",
	}
	err := json.NewEncoder(w).Encode(message)
	if err != nil {
		log.Println(err)
	}

}

func GetSessionHandler(w http.ResponseWriter, r *http.Request) {
	conf := config.Get()

	election := Election{}

	mgoSession := conf.Database.Session.Copy()
	defer mgoSession.Close()

	collection := conf.Database.C(config.ELECTION_COLLECTION).With(mgoSession)

	err := collection.Find(bson.M{
		"id": "election",
	}).One(&election)
	if err != nil {
		log.Println(err)
	}

	votersCollection := conf.Database.C(config.VOTERS_COLLECTION).With(mgoSession)

	election.VotersCount, err = votersCollection.Count()
	if err != nil {
		log.Println(err)
		//election.VotersCount = 0
	}

	json.NewEncoder(w).Encode(&election)
}

func StartVotingHandler(w http.ResponseWriter, r *http.Request) {
	conf := config.Get()

	mgoSession := conf.Database.Session.Copy()
	defer mgoSession.Close()

	collection := conf.Database.C(config.ELECTION_COLLECTION).With(mgoSession)

	_, err := collection.Upsert(
		bson.M{
			"id": "election",
		}, bson.M{
			"$set": bson.M{
				"runningelection": true,
			},
		})
	if err != nil {
		log.Println(err)
		return
	}

	message := struct {
		Message string
	}{
		Message: "Voting started successfully",
	}
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		log.Println(err)
	}

}

func EndVotingHandler(w http.ResponseWriter, r *http.Request) {
	conf := config.Get()

	mgoSession := conf.Database.Session.Copy()
	defer mgoSession.Close()

	collection := conf.Database.C(config.ELECTION_COLLECTION).With(mgoSession)

	_, err := collection.Upsert(
		bson.M{
			"id": "election",
		}, bson.M{
			"$set": bson.M{
				"runningelection": false,
			},
		})
	if err != nil {
		log.Println(err)
		return
	}

	message := struct {
		Message string
	}{
		Message: "Voting ended successfully",
	}
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		log.Println(err)
	}
}

func GetPositionHandler(w http.ResponseWriter, r *http.Request) {
	conf := config.Get()
	election := Election{}

	key, err := strconv.Atoi(r.URL.Query().Get("key"))
	if err != nil {
		log.Println(err)
	}
	mgoSession := conf.Database.Session.Copy()
	defer mgoSession.Close()

	collection := conf.Database.C(config.ELECTION_COLLECTION).With(mgoSession)
	err = collection.Find(bson.M{
		"id": "election",
	}).One(&election)
	if err != nil {
		log.Println(err)
	}

	contestants := []Contestant{}

	contestantsCollection := conf.Database.C(config.CONTESTANTS_COLLECTION).With(mgoSession)

	err = contestantsCollection.Find(bson.M{
		"key": key,
	}).All(&contestants)
	if err != nil {
		log.Println(err)
	}

	data := struct {
		Election    Election
		Contestants []Contestant
	}{
		Election:    election,
		Contestants: contestants,
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
	}
}

func AddContestantHandler(w http.ResponseWriter, r *http.Request) {
	key, err := strconv.Atoi(r.URL.Query().Get("key"))
	if err != nil {
		log.Println(err)
	}
	name := r.URL.Query().Get("name")
	position := r.URL.Query().Get("position")

	contestant := Contestant{
		Name:     name,
		Position: position,
		Key:      key,
	}

	conf := config.Get()
	mgoSession := conf.Database.Session.Copy()
	defer mgoSession.Close()

	contestantsCollection := conf.Database.C(config.CONTESTANTS_COLLECTION).With(mgoSession)
	err = contestantsCollection.Insert(contestant)
	if err != nil {
		log.Println(err)
	}

	message := struct {
		Message string
	}{
		Message: "Contestant added successfully",
	}
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		log.Println(err)
	}
}
