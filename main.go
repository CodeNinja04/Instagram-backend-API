package main

//Import necessary standard libraries
import (
	"backend/helper"
	"backend/models"
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"hash"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// function to test server
type server struct{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "get called"}`))
	case "POST":
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "post called"}`))
	case "PUT":
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message": "put called"}`))
	case "DELETE":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "delete called"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}

// function which converts string to hash using sha1
func getHash(byteStr []byte) string {
	var hashVal hash.Hash
	hashVal = sha1.New()
	hashVal.Write(byteStr)

	var bytes []byte

	bytes = hashVal.Sum(nil)
	return string(bytes)
}

// return all users info
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var users []models.User

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection1.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var user models.User
		// & character returns the memory address of the following variable.
		err := cur.Decode(&user) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(users) // encode similar to serialize process.
}

// create user
func createUser(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	w.Header().Set("Content-Type", "application/json")

	var user models.User

	// we decode our body request params

	json.NewDecoder(r.Body).Decode(&user)

	user.Password = getHash([]byte(user.Password))

	// insert our book model.
	result, err := collection1.InsertOne(context.TODO(), user)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)

	w.Write([]byte(`{"message": "user added"}`))

}

// get a particular user via id

func getUser(w http.ResponseWriter, r *http.Request) {

	// set header.
	w.Header().Set("Content-Type", "application/json")
	name := strings.Replace(r.URL.Path, "/user/", "", 1)
	fmt.Println(name)

	var user models.User

	//string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(name)

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": id}
	err := collection1.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// create post
func createPost(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	w.Header().Set("Content-Type", "application/json")

	var post models.Post

	// we decode our body request params

	json.NewDecoder(r.Body).Decode(&post)

	result, err := collection2.InsertOne(context.TODO(), post)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)

	w.Write([]byte(`{"message": "post added"}`))

}

// get particular post

func getPost(w http.ResponseWriter, r *http.Request) {

	// set header.
	w.Header().Set("Content-Type", "application/json")
	name := strings.Replace(r.URL.Path, "/post/", "", 1)
	fmt.Println(name)

	var post models.Post
	// we get params with mux.
	//var params = mux.Vars(r)

	//string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(name)

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": id}
	err := collection2.FindOne(context.TODO(), filter).Decode(&post)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(post)
}

// get particular post via post id
func getUserPost(w http.ResponseWriter, r *http.Request) {

	// set header.
	w.Header().Set("Content-Type", "application/json")
	name := strings.Replace(r.URL.Path, "/post/users/", "", 1)
	fmt.Println(name)

	var post models.Post
	// we get params with mux.
	//var params = mux.Vars(r)

	//string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(name)

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": id}
	err := collection1.FindOne(context.TODO(), filter).Decode(&post)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(post)
}

var collection1, collection2 = helper.ConnectDB()

func main() {

	s := &server{}
	http.Handle("/", s)
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/createuser", createUser)
	http.HandleFunc("/user/", getUser)
	http.HandleFunc("/createpost", createPost)
	http.HandleFunc("/post", getPost)
	http.HandleFunc("/post/users/", getUserPost)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
