package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/unrolled/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

var redis_connect = "192.168.99.100:6379"
var mongodb_server = "192.168.99.100:27017"
var mongodb_server2 = "192.168.99.100:27018"
var mongodb_database = "cmpe281"
var mongodb_collection = "redistest"

type (
	// User represents the structure of our resource
	User struct {
		SerialNumber string `json: "id"`
		Name         string `json: "name"`
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
	//mx.HandleFunc("/order", postHandler(formatter)).Methods("POST")
	//mx.HandleFunc("/order", putHandler(formatter)).Methods("PUT")
	//mx.HandleFunc("/order", deleteHandler(formatter)).Methods("DELETE")
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
			fmt.Println("It got from Mongo")
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
