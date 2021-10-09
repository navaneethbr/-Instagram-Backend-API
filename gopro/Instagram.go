package main

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}
type InstaPost struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption         string             `json:"caption,omitempty" bson:"caption,omitempty"`
	ImgURL          string             `json:"imgurl,omitempty" bson:"imgurl,omitempty"`
	PostedTimestamp string             `json:"postedtimestamp,omitempty" bson:"postedtimestamp,omitempty"`
}

var client *mongo.Client

/****************************
Author Name: Navaneeth B R
Date: 09/10/2021
Description: Method will hash the password and the json structure of user will be passed for
post call out in the main block
******************************/
func createPersonEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var person Person
	s := person.Password
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	fmt.Println(s)
	fmt.Printf("%x\n", bs)
	json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("instagram").Collection("Person")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

/************************
Author Name: Navaneeth B R
Date: 09/10/2021
Description: The json structure(insta post) created will be passed for
post call out in the main block
***********************/
func createPostEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var person InstaPost
	json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("instagram").Collection("InstaPost")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

/**************************
Author Name: Navaneeth B R
Date: 09/10/2021
Description: The json structure of user will be passed for
get call out in the main block
******************************/
func getPersonEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Person
	collection := client.Database("instagram").Collection("Person")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(person)
}

/****************************
Author Name: Navaneeth B R
Date: 09/10/2021
Description: The json structure(insta post) created will be passed for
get call out in the main block
***********************************/
func getPostEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person InstaPost
	collection := client.Database("instagram").Collection("InstaPost")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, InstaPost{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(person)
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()

	router.HandleFunc("/users", createPersonEndPoint).Methods("POST")
	router.HandleFunc("/users/{id}", getPersonEndPoint).Methods("GET")
	router.HandleFunc("/posts", createPostEndPoint).Methods("POST")
	router.HandleFunc("/posts/{id}", getPostEndPoint).Methods("GET")
	http.ListenAndServe(":12345", router)

}
