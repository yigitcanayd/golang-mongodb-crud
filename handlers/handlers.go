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
	w.Header().Set("Content-Type", "application/json")
	name := strings.TrimPrefix(r.URL.Path, "/people/")

	if len(name) > 0 {
		var result primitive.M
		err := collection.FindOne(context.TODO(), bson.D{{"name.firstname", name}}).Decode(&result)
		if err != nil {
			w.WriteHeader(404)
			log.Fatal(err)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
	} else {
		var results []primitive.M
		cur, err := collection.Find(context.TODO(), bson.D{{}})
		if err != nil {
			w.WriteHeader(404)
			log.Fatal(err)
		}
		for cur.Next(context.TODO()) {
			var elem primitive.M
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}
			results = append(results, elem)
		}
		cur.Close(context.TODO())
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(results)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var person Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		w.WriteHeader(404)
		log.Fatal(err)
	}
	insertResult, err := collection.InsertOne(context.TODO(), person)
	if err != nil {
		w.WriteHeader(404)
		log.Fatal(err)
	}

	log.Println("Inserted a single document: ", insertResult)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(insertResult.InsertedID)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/people/")

	if len(name) < 1 {
		w.WriteHeader(400)
	} else {
		res, err := collection.DeleteOne(context.TODO(), bson.D{{"name.firstname", name}})
		if err != nil {
			w.WriteHeader(404)
		}
		log.Println(res)
		w.WriteHeader(204)
	}
}
