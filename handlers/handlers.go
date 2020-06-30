package handlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang-mongodb-crud/mongodb"
	"log"
	"net/http"
	"strings"
)

type Person struct {
	Name        Name   `json:"name"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
}

type Name struct {
	Firstname  string `json:"firstname"`
	Middlename string `json:"middlename"`
	Lastname   string `json:"lastname"`
}

var collection = mongodb.GetMongoClient().Database("golang").Collection("people")

func GetHandler(w http.ResponseWriter, r *http.Request) {
	name := setContentTypeAndGetPathVariable(w, r)
	if len(name) > 0 {
		getRecord(w, name)
	} else {
		getAllRecords(w)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	name := setContentTypeAndGetPathVariable(w, r)

	var person Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		w.WriteHeader(400)
		log.Println(err)
	}
	if len(name) > 0 {
		updateRecord(w, name, person)
	} else {
		insertRecord(w, person)
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	name := setContentTypeAndGetPathVariable(w, r)
	if len(name) > 0 {
		deleteRecord(w, name)
	} else {
		w.WriteHeader(400)
	}
}

func setContentTypeAndGetPathVariable(w http.ResponseWriter, r *http.Request) string {
	w.Header().Set("Content-Type", "application/json")
	name := strings.TrimPrefix(r.URL.Path, "/people/")
	return name
}

func getRecord(w http.ResponseWriter, name string) {
	var result primitive.M
	filter := bson.D{{"name.firstname", name}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		w.WriteHeader(404)
		log.Println(err)
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(result)
}

func getAllRecords(w http.ResponseWriter) {
	var results []primitive.M
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		w.WriteHeader(404)
		log.Println(err)
	}
	for cur.Next(context.TODO()) {
		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
		}
		results = append(results, elem)
	}
	cur.Close(context.TODO())
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(results)
}

func insertRecord(w http.ResponseWriter, person Person) {
	insertResult, err := collection.InsertOne(context.TODO(), person)
	if err != nil {
		w.WriteHeader(400)
		log.Println(err)
	}

	log.Println("Inserted a single document: ", insertResult)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(insertResult.InsertedID)
}

func updateRecord(w http.ResponseWriter, name string, person Person) {
	filter := bson.D{{"name.firstname", name}}
	update := bson.D{
		{"$set", person},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		w.WriteHeader(404)
		log.Println(err)
	}
	log.Println("Updated a document: ", name)
	w.WriteHeader(200)
}

func deleteRecord(w http.ResponseWriter, name string) {
	filter := bson.D{{"name.firstname", name}}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		w.WriteHeader(404)
	}
	log.Println("Deleted a document: ", name)
	w.WriteHeader(204)
}
