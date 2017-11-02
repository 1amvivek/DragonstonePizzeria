package main

import (
	"fmt"
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/unrolled/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"github.com/gocql/gocql"
)

var redis_connect = "192.168.99.100:6379"
var mongodb_server1 = "192.168.99.100:27017"
var mongodb_server2 = "192.168.99.100:27018"
var mongodb_server3 = "192.168.99.100:27019"
var cassandra_server = "192.168.99.100:32769"
var mongodb_database = "cmpe281"
var mongodb_collection = "redistest"
var i = 0
var servers = []string{mongodb_server1, mongodb_server2, mongodb_server3, cassandra_server}


type (
	// User represents the structure of our resource
	User struct {
		SerialNumber string	`json: "id"`
		Name  string `json: "name"`
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
	//mx.HandleFunc("/order", putHandler(formatter)).Methods("PUT")
	//mx.HandleFunc("/order", deleteHandler(formatter)).Methods("DELETE")
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {  
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(code)
    fmt.Fprintf(w, "{message: %q}", message)
}


// Helper Functions

func getFromMongo(mongodb string, serialNumber string) bson.M {
	fmt.Println("mongo connecting to " + mongodb)
	session, err := mgo.Dial(mongodb)

	if err != nil {
		log.Fatal("mongo failed to connect")
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongodb_database).C(mongodb_collection)
	var result bson.M
	//get from mongo
	err = c.Find(bson.M{"SerialNumber": serialNumber}).One(&result)
	if err != nil {
		//could not find in mongo (inserting into mongo for now. TODO: Make proper)
		c.Insert(bson.M{"SerialNumber": "1", "Name": "Sample"})
	}
	return result

}

func connectToRedis(redis_connect string, serialNumber string) (*redis.Client, bool, string) {

	conn, err := redis.Dial("tcp", redis_connect)
	if err != nil {
		log.Fatal("redis failed to connect")
		log.Fatal(err)
	}
	cacheFlag := false
	//get from redis
	name, err := conn.Cmd("HGET", serialNumber, "Name").Str()
	if err != nil {
		//not in redis
		fmt.Println("couldn't find values in Redis")

		cacheFlag := true
		return conn, cacheFlag, name

	}

	return conn, cacheFlag, name
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
			//connect to mongo
			result := getFromMongo(mongodb_server, serialNumber)
			result2 := getFromMongo(mongodb_server2, serialNumber)
			if result["serialNumber"] == result2["serialNumber"] {
				//print from mongo
				formatter.JSON(w, http.StatusOK, result)
				//store in redis
				conn.Cmd("HMSET", result["SerialNumber"], "Name", result["Name"])
			} else {
				fmt.Println("things went wrong")
			}
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

func getSession2(cassandra_bal_server string) *gocql.Session {  
    // Connect to cassandra cluster
    
    fmt.Println("cassandra connecting to " + cassandra_bal_server)
    
	cluster := gocql.NewCluster(cassandra_bal_server)
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	
    return session
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
    var errs []string
    var user User
    decoder := json.NewDecoder(req.Body)
    err1 := decoder.Decode(&user)

    server_val := Balance()
    if server_val== cassandra_server {
    s := getSession2(server_val) 
    defer s.Close()

    fmt.Println("creating a new user")

    // generate a unique UUID for this user
    var id = gocql.TimeUUID()
    // write data to Cassandra
    if err1 := s.Query(`
     INSERT INTO user (SerialNumber, id, Name) VALUES (?, ?, ?)`,
	 user.SerialNumber,id, user.Name).Exec(); err1 != nil {
		      errs = append(errs, err1.Error())
    } 
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Location", req.URL.Path+"/"+user.SerialNumber)
    w.WriteHeader(http.StatusCreated)
    w.WriteHeader(200)
	 

    return
    }
    
    s := getSession(server_val)
    defer s.Close()
    
    if err1 != nil {
        ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
        return
    }

    
    c := s.DB(mongodb_database).C(mongodb_collection)

    err2 := c.Insert(user)
    
    if err2 != nil {
        if mgo.IsDup(err2) {
        	    fmt.Println("exists already")
                ErrorWithJSON(w, "User already exists", http.StatusBadRequest)
                return
        }


        ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
        log.Println("Failed insert user: ", err2)
        return
        }
    
       
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Location", req.URL.Path+"/"+user.SerialNumber)
    w.WriteHeader(http.StatusCreated)
    w.WriteHeader(200)
	 

	}
}

