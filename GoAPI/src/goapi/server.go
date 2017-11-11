package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/unrolled/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"strings"
)

var redis_connect = "192.168.99.100:6379"
var mongodb_server1 = "192.168.99.100:27017"
var mongodb_server2 = "192.168.99.100:27018"
var mongodb_server3 = "192.168.99.100:27019"
var mongodb_database = "cmpe281"
var mongodb_collection = "redistest"
var i = 0
var servers = []string{mongodb_server1, mongodb_server2, mongodb_server3}

type (
	// User represents the structure of our resource
	User struct {
		SerialNumber string `json: "id"`
		Name         string `json: "name"`
		Clock        int    `json: "clock`
	}
	Id struct {
		SerialNumber string `json: "id"`
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
	mx.HandleFunc("/orders/{order_id}", getHandler(formatter)).Methods("GET")
	mx.HandleFunc("/order", postHandler(formatter)).Methods("POST")
	mx.HandleFunc("/order", putHandler(formatter)).Methods("PUT")
	mx.HandleFunc("/order", deleteHandler(formatter)).Methods("DELETE")
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

// Helper Functions

func getFromMongo(mongodb string, serialNumber string) User {
	fmt.Println("mongo connecting to " + mongodb)
	session, err := mgo.Dial(mongodb)

	if err != nil {
		log.Fatal("mongo failed to connect")
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongodb_database).C(mongodb_collection)
	var result User
	//get from mongo
	err = c.Find(bson.M{"serialnumber": serialNumber}).One(&result)
	if err != nil {
		//could not find in mongo (inserting into mongo for now. TODO: Make proper)
		// c.Insert(bson.M{"SerialNumber": "1", "Name": "Sample"})
		fmt.Println("Some Error in Get, maybe data is not present")
	}
	return result

}

func connectToRedis(redis_connect string, serialNumber string) (*redis.Client, bool, User) {
	var result User
	conn, err := redis.Dial("tcp", redis_connect)
	if err != nil {
		log.Fatal("redis failed to connect")
		log.Fatal(err)
	}
	cacheFlag := false
	//get from redis
	val, err := conn.Cmd("HGET", serialNumber, "object").Str()
	if err != nil {
		//not in redis
		fmt.Println("couldn't find values in Redis")

		cacheFlag = true

	}
	json.Unmarshal([]byte(val), &result)
	fmt.Println("cacheFlag")
	fmt.Println(cacheFlag)

	return conn, cacheFlag, result
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
}

func deleteHelper(server_val string, serialNumber string) {
	s := getSession(server_val)
	defer s.Close()
	fmt.Println("Deleting the user")
	user := getFromMongo(server_val, serialNumber)
	//deleting in mongo
	c := s.DB(mongodb_database).C(mongodb_collection)
	err2 := c.Update(bson.M{"serialnumber": serialNumber}, bson.M{"$set": bson.M{"name": user.Name, "clock": -1}})

	if err2 != nil {
		fmt.Println("Some Random error")
	}

}

// API GET Handler
func getHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		//get request param
		params := mux.Vars(req)
		var serialNumber string = params["order_id"]
		cacheFlag := false
		//connect to redis
		conn, cacheFlag, name := connectToRedis(redis_connect, serialNumber)

		if cacheFlag {
			fmt.Println("The values are fetched from Mongo")
			var result User
			max := 0
			var id big.Int
			id.SetString(strings.Replace(serialNumber, "-", "", 4), 16)
			servers := getNodes(id)
			for _, value := range servers {
				temp := getFromMongo(value, serialNumber)
				if int(temp.Clock) > max {
					result = temp
					max = int(temp.Clock)
				}
			}
			formatter.JSON(w, http.StatusOK, result)

			for _, value := range servers {
				temp := getFromMongo(value, serialNumber)
				if temp.Clock != result.Clock {
					//WRITE VALUE BACK TO MONGO
				}
			}

			redstore, err := json.Marshal(result)
			if err != nil {
				fmt.Println("Couldn't marshall the result")
			}
			conn.Cmd("HMSET", result.SerialNumber, "object", redstore)
		} else {
			fmt.Println("The values are fetched from REDIS")
			//print from redis
			formatter.JSON(w, http.StatusOK, name)
		}
	}

}

func getSession(mongodb_bal_server string) *mgo.Session {
	// Connect to mongo cluster
	//mongodb_bal_server := Balance()
	fmt.Println("mongo connecting to " + mongodb_bal_server)
	s, err := mgo.Dial(mongodb_bal_server)

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	s.SetMode(mgo.Monotonic, true)
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

		server_val := Balance()

		s := getSession(server_val)
		defer s.Close()

		if err5 != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		//user1.uid = rand.Int()
		var uuid = rand.Int()
		fmt.Println(uuid)

		var resp1 = fmt.Sprintf("{'uuid' : %d }", uuid)

		Jsonvalue, err5 := json.Marshal(resp1)

		c := s.DB(mongodb_database).C(mongodb_collection)

		err6 := c.Insert(
			struct{ SerialNumber, Name interface{} }{
				SerialNumber: rand.Int(),
				Name:         user1.Name})

		if err6 != nil {
			if mgo.IsDup(err6) {
				fmt.Println("exists already")
				ErrorWithJSON(w, "User already exists", http.StatusBadRequest)
				return
			}

			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed insert user: ", err6)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", req.URL.Path+"/"+user1.SerialNumber)
		w.WriteHeader(http.StatusCreated)
		w.WriteHeader(200)

		w.Write(Jsonvalue)

	}
}

func deleteHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		//get variables
		params := mux.Vars(req)
		var serialNumber string = params["order_id"]

		//connect to redis
		conn, cacheFlag, _ := connectToRedis(redis_connect, serialNumber)

		//deleting in redis
		if cacheFlag {
			fmt.Println("There aren't any values in Redis")
		} else {
			fmt.Println("Deleting values at Redis End")
			//delete in redis
			conn.Cmd("DEL", serialNumber)

		}
		serialNumberInt, err1 := strconv.Atoi(serialNumber)
		if err1 != nil {
			fmt.Println("could not convert to integer")
		}
		// It wont delete the data from redis and mongo, as the server_val is not properly set. Hence we have to run del again.
		servers := getNodes(serialNumberInt)
		for _, value := range servers {
			deleteHelper(value, serialNumber)
		}
	}
}

func putHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var user User
		//get mongodb connection
		decoder := json.NewDecoder(req.Body)
		err1 := decoder.Decode(&user)
		server_val := Balance()

		//connect to redis
		conn, _, name := connectToRedis(redis_connect, user.SerialNumber)

		//deleting in redis
		if name.SerialNumber != "" {
			fmt.Println("Deleting values at Redis End")
			//delete in redis
			conn.Cmd("DEL", name.SerialNumber)
		} else {
			fmt.Println("There aren't any values in Redis")
		}

		s := getSession(server_val)
		defer s.Close()
		fmt.Println("Updating the user")

		if err1 != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			fmt.Println(err1)
			return
		}
		fmt.Println(user.SerialNumber)

		// It wont update the data from redis and mongo at the same time, as the server_val is not properly set. Hence we have to run update command again.
		c := s.DB(mongodb_database).C(mongodb_collection)
		err2 := c.Update(bson.M{"serialnumber": user.SerialNumber}, bson.M{"$set": bson.M{"name": user.Name}})

		if err2 != nil {
			fmt.Println("Some Random error")
			return
		}
	}
}
