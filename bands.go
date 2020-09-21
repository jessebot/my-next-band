package main

import (
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
    "net/http"
    "os"
)

// define global Band structure
type Band struct {
    Name      string `json:"Name"`
    Genre     string `json:"Genre"`
    Id        int    `json"Id"`
}

// home page always returns this
func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to my bands page!")
    fmt.Println("Endpoint Hit: homePage")
}

// for just getting ALL bands
func returnAllBands(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: returnAllBands")

    // Here's an array in which you can store the decoded documents
    var results []*Band

    bandsCollection := mongoInit()
    // Passing bson.D{{}} as the filter matches all documents in the collection
    cur, err := bandsCollection.Find(context.TODO(), bson.D{{}})
    if err != nil {
        log.Fatal(err)
    }

    // Iterating through the cursor allows us to decode documents one at a time
    for cur.Next(context.TODO()) {
        // create a value into which the single document can be decoded
        var elem Band
        err := cur.Decode(&elem)
        if err != nil {
            log.Fatal(err)
        }

        results = append(results, &elem)
    }

    json.NewEncoder(w).Encode(results)
}

// getting a single band by ID
func returnSingleBand(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    key := vars["name"]

    var result Band

    bandsCollection := mongoInit()

    filter := bson.D{{"name", key}}
    err := bandsCollection.FindOne(context.TODO(), filter).Decode(&result)
    if err != nil {
        log.Fatal(err)
    }

    json.NewEncoder(w).Encode(result)
}

func createNewBand(w http.ResponseWriter, r *http.Request) {
    // get the body of our POST request
    reqBody, _ := ioutil.ReadAll(r.Body)
    // print this data to logs
    fmt.Println(string(reqBody))
    // create a variable
    var band Band
    // unmarshal this into a new Bands array
    json.Unmarshal(reqBody, &band)

    bandsCollection := mongoInit()
    insertResult, err := bandsCollection.InsertOne(context.TODO(), band)
    if err != nil {
        log.Fatal(err)
    }

    // finally, return the newly added band
    fmt.Println(insertResult)
    fmt.Fprintf(w, "New band added!")
}

func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    // just a default display
    myRouter.HandleFunc("/", homePage)
    // return all bands
    myRouter.HandleFunc("/bands", returnAllBands)
    // create a new band record
    myRouter.HandleFunc("/band", createNewBand).Methods("POST")
    // return band info by name
    myRouter.HandleFunc("/band/{name}", returnSingleBand)
    log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func mongoInit() *mongo.Collection{
    // mongodb
    var mongoURI string
    mongoURI = os.Getenv("MONGO_URI")

    var mongoDB string
    mongoDB = os.Getenv("MONGO_DB")

    clientOptions := options.Client().ApplyURI(mongoURI)
    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), clientOptions)
    // Check the connection
    if err != nil { log.Fatal(err) }

    // get collection
    bandsCollection := client.Database(mongoDB).Collection("bands")

    return bandsCollection
}

func main() {
    handleRequests()
}
