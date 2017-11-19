package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

var mongodb_server1 = "192.168.99.100:27017"
var mongodb_server2 = "192.168.99.100:27018"
var mongodb_server3 = "192.168.99.100:27019"
var mongodb_database = "cmpe281"
var mongodb_collection = "recommend"
var i = 0
var servers = []string{mongodb_server1, mongodb_server2}
var updatedCart = "" 
var current User

type (
	// User represents the structure of our resource
	User struct {
		SerialNumber string `json: "id"`
		Cart         string `json: "cart"`
		Clock        int    `json: "clock`
	}
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	n.UseHandler(mx)
	return n
}

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/order", postHandler(formatter)).Methods("POST")
	mx.HandleFunc("/order", putHandler(formatter)).Methods("PUT")
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

// Helper Functions

func getFromMongo(session *mgo.Session, serialNumber string) User {

	var result User
	//get from mongo
	if session != nil {
		c := session.DB(mongodb_database).C(mongodb_collection)
		err := c.Find(bson.M{"serialnumber": serialNumber}).One(&result)
		if err != nil {
			//could not find in mongo (inserting into mongo for now. TODO: Make proper)
			// c.Insert(bson.M{"SerialNumber": "1", "Name": "Sample"})
			fmt.Println("Some Error in Get, maybe data is not present")
		}
	}
	return result

}

func getNodes(uuid int) []string {
	start := uuid % len(servers)
	end := start + len(servers)/2 + 1
	fmt.Println(start)
	fmt.Println(end)
	if end <= len(servers)-1 {
		return servers[start:end]
	} else {
		end = end - len(servers)
		fmt.Println("end changed")
		return append(servers[start:], servers[:end]...)
	}
	return servers
}

func updateHelper(server_val string, user User) {
	s := getSession(server_val)
	if s != nil {
		defer s.Close()
		fmt.Println("Updating the user")
		var current User
		c := s.DB(mongodb_database).C(mongodb_collection)
		err := c.Find(bson.M{"serialnumber": user.SerialNumber}).One(&current)
		if err != nil {
			fmt.Println("no object to update")
			return
		}
		if current.Cart != "" {
			updatedCart = fmt.Sprint(current.Cart + "," + user.Cart)
		} else {
			updatedCart = fmt.Sprint(user.Cart)
		}
		err2 := c.Update(bson.M{"serialnumber": user.SerialNumber}, bson.M{"$set": bson.M{"cart": updatedCart, "clock": (current.Clock + 1)}})
		fmt.Println("update Helper",updatedCart)
		if err2 != nil {
			fmt.Println("Some Random error")
			return
		}
	}
}

func postHelper(server_val string, user1 User) {
	s := getSession(server_val)
	if s != nil {
		defer s.Close()

		c := s.DB(mongodb_database).C(mongodb_collection)

		err6 := c.Insert(user1)

		if err6 != nil {
			if mgo.IsDup(err6) {
				fmt.Println("exists already")
				//ErrorWithJSON(w, "User already exists", http.StatusBadRequest)
				//return
			}

			//ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed insert user: ", err6)
			return
		}
	}
}

func getSession(mongodb_bal_server string) *mgo.Session {
	// Connect to mongo cluster
	//mongodb_bal_server := Balance()
	fmt.Println("mongo connecting to " + mongodb_bal_server)
	s, err := mgo.Dial(mongodb_bal_server)

	if err == nil {
		s.SetMode(mgo.Monotonic, true)
	}
	return s
}

// Balance returns one of the servers based using round-robin algorithm
func Balance() string {
	server := servers[i]
	i++

	// reset the counter and start from the beginning
	// if we reached the end of servers
	if i >= len(servers) {
		i = 0
	}
	return server
}


func postHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		//get mongodb connection

		var user1 User
		decoder := json.NewDecoder(req.Body)
		err5 := decoder.Decode(&user1)
		if err5 != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		//user1.uid = rand.Int()
		var uuid = rand.Int()
		var resp1 = fmt.Sprintf("{'uuid' : %d }", uuid)
		Jsonvalue, err5 := json.Marshal(resp1)
		// servers := getNodes(uuid)
		user1.SerialNumber = strconv.Itoa(uuid)
		user1.Clock = 1
		for _, value := range servers {
			postHelper(value, user1)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", req.URL.Path+"/"+user1.SerialNumber)
		w.WriteHeader(http.StatusCreated)
		w.WriteHeader(200)

		w.Write(Jsonvalue)

	}
}


func recommendHelper(servers []string, user User) {
	fmt.Println("Starting recommendation module!!!")
	s := getSession(servers[1])
	if s != nil {
		defer s.Close()
		fmt.Println("In recommendation module!!!")
		c := s.DB(mongodb_database).C("rules")
		fmt.Println(updatedCart)
		err := c.Find(bson.M{"cart" : updatedCart}).One(&current)
		if err != nil {
			fmt.Println("no rule found")
			return
	}
	fmt.Println(current)
	}
}


func putHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var user User
		//get mongodb connection
		decoder := json.NewDecoder(req.Body)
		err1 := decoder.Decode(&user)

		if err1 != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			fmt.Println(err1)
			return
		}
		// servers := getNodes(int(serialNumber))
		for _, value := range servers {
			updateHelper(value, user)
		}

		recommendHelper(servers, user)
		Jsonvalue, _ := json.Marshal(current)
		w.Header().Set("Content-Type", "application/json")
		// w.Header().Set("Location", req.URL.Path+"/"+user.Cart)
		w.WriteHeader(http.StatusCreated)
		w.WriteHeader(200)
		w.Write(Jsonvalue)

	}
}