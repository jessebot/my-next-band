package main

import (
    "net/http"
    "io/ioutil"
    "log"
    "fmt"
    "encoding/json"
    "github.com/gorilla/mux"
)

// define global Bands array
type Band struct {
    Name      string `json:"Name"`
    Genre     string `json:"Genre"`
    Id        int    `json"Id"`
}

// declare global Bands array
var Bands []Band

// home page always returns this
func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to my bands page!")
    fmt.Println("Endpoint Hit: homePage")
}

// for just getting ALL bands
func returnAllBands(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: returnAllBands")
    //for key, value := range Bands {
    //  fmt.Fprintf.(w, element)
    //}
    json.NewEncoder(w).Encode(Bands)
}

// getting a single band by ID
func returnSingleBand(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    key := vars["name"]

    // iterate through the Bands list
    for _, band := range Bands {
        // if band.Name in Bands is the passed in key, print the whole band record
        if band.Name == key {
            json.NewEncoder(w).Encode(band)
        }
    }
}

func createNewBand(w http.ResponseWriter, r *http.Request) {
    // get the body of our POST request
    reqBody, _ := ioutil.ReadAll(r.Body)
    // print this data to logs
    fmt.Println("%+v", string(reqBody))
    // create a variable
    var band Band
    // unmarshal this into a new Bands array
    json.Unmarshal(reqBody, &band)
    // update our global Bands array to include our new Band
    Bands = append(Bands, band)

    // finally, return the newly added band
    json.NewEncoder(w).Encode(band)
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

func main() {
    // test bands 
    Bands = []Band {
        Band{Name: "Love Sledge", Genre: "Glam Rock", Id: 01},
        Band{Name: "Power Pizza", Genre: "Food ballads", Id: 02},
    }
    // parse requests
    handleRequests()
}
