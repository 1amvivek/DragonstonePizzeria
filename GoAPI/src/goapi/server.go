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
	"math/rand"
	"net/http"
	"strconv"
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
	mx.HandleFunc("/order/{order_id}", deleteHandler(formatter)).Methods("DELETE")
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

func updateHelper(server_val string, user User) {
	s := getSession(server_val)
	if s != nil {
		defer s.Close()
		fmt.Println("Updating the user")
		var current User
		// It wont update the data from redis and mongo at the same time, as the server_val is not properly set. Hence we have to run update command again.
		c := s.DB(mongodb_database).C(mongodb_collection)
		err := c.Find(bson.M{"serialnumber": user.SerialNumber}).One(&current)
		if err != nil {
			fmt.Println("no object to update")
			return
		}
		err2 := c.Update(bson.M{"serialnumber": user.SerialNumber}, bson.M{"$set": bson.M{"name": user.Name, "clock": (current.Clock + 1)}})

		if err2 != nil {
			fmt.Println("Some Random error")
			return
		}
	}
}

func deleteHelper(server_val string, serialNumber string) {
	s := getSession(server_val)
	if s != nil {
		defer s.Close()
		fmt.Println("Deleting the user")
		user := getFromMongo(s, serialNumber)
		//deleting in mongo
		c := s.DB(mongodb_database).C(mongodb_collection)
		err2 := c.Update(bson.M{"serialnumber": serialNumber}, bson.M{"$set": bson.M{"name": user.Name, "clock": -1}})

		if err2 != nil {
			fmt.Println("Some Random error")
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

// API GET Handler
func getHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		//get request param
		params := mux.Vars(req)
		var serialNumber string = params["order_id"]
		cacheFlag := false
		var connections [10]*mgo.Session
		var objects [10]User
		//connect to redis
		conn, cacheFlag, name := connectToRedis(redis_connect, serialNumber)
		serialNumberInt, err1 := strconv.Atoi(serialNumber)
		if err1 != nil {
			fmt.Println("could not convert to integer")
		}
		if cacheFlag {
			fmt.Println("The values are fetched from Mongo")
			var result User
			max := 0
			servers := getNodes(serialNumberInt)
			for index, value := range servers {
				connections[index] = getSession(value)
				objects[index] = getFromMongo(connections[index], serialNumber)
			}
			for _, object := range objects {
				if int(object.Clock) > max {
					result = object
					max = int(object.Clock)
				}
			}
			formatter.JSON(w, http.StatusOK, result)

			for index, object := range objects {
				if object.Clock != result.Clock {
					if connections[index] != nil {
						c := connections[index].DB(mongodb_database).C(mongodb_collection)
						err2 := c.Update(bson.M{"serialnumber": result.SerialNumber}, bson.M{"$set": bson.M{"name": result.Name, "clock": (result.Clock)}})
						if err2 != nil {
							fmt.Println("Some Random error")
							return
						}
					}
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
		servers := getNodes(uuid)
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

		if err1 != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			fmt.Println(err1)
			return
		}

		serialNumber, err := strconv.Atoi(user.SerialNumber)
		if err != nil {
			fmt.Println("could not convert serialnumber to int")
		}
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
		servers := getNodes(int(serialNumber))
		for _, value := range servers {
			updateHelper(value, user)
		}

	}
}