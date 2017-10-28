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
var mongodb_server = "192.168.99.100"
var mongodb_database = "cmpe281"
var mongodb_collection = "redistest"

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
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("here")
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// API GET Handler
func getHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		var serialNumber string = params["order_id"]
		cacheFlag := false
		//connect to redis
		conn, err := redis.Dial("tcp", redis_connect)
		if err != nil {
			log.Fatal("redis failed to connect")
			log.Fatal(err)
		}
		//get from redis
		name, err := conn.Cmd("HGET", serialNumber, "Name").Str()
		if err != nil {
			//not in redis
			cacheFlag = true
		}

		if cacheFlag {
			fmt.Println("It got from Mongo")
			//connect to mongo
			session, err := mgo.Dial(mongodb_server)
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
			//print from mongo
			formatter.JSON(w, http.StatusOK, result)
			//store in redis
			conn.Cmd("HMSET", result["SerialNumber"], "Name", result["Name"])

		} else {
			fmt.Println("It got from REDIS")
			//print from redis
			formatter.JSON(w, http.StatusOK, name)
		}
	}

}
