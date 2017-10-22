## What is Golang? How is it evolved?

Go is a programming language that didn’t evolve from the existing programming languages, such as Java and C#. This programming language focuses on real-world practices for building the next-generation applications on the cloud, computing environments, and system programming.  
Go is a statically typed, garbage-collected, natively compiled programming language that mostly belongs to the C family of languages in terms of basic syntax. Go provides an expressive syntax with it’s lightweight type system and comes with concurrency as a built-in feature at the language level. Go provides performance comparable to C and C++ while offering rapid development.

# Benefits of Programming with Golang

Some of the important benefits/advantages of Golang are as follows:

#1 - Highly Concurrent  
#2 - Approachable for the Whole Team  
#3 - Excellent Performance

There are other advantages such as: 
* Go compiles very quickly.
* Go supports concurrency at the language level.
* Functions are first class objects in Go.
* Go has garbage collection.
* Strings and maps are built into the language. 

# Implementing Go Framework for AWS Lambda

Go offers a solid toolchain and set of primitives to write AWS Lambda services. The following are the features that Go offers:

- Single binary deployment
- Excellent concurrency primitives
- An official AWS Go SDK
- Extremely fast compilation
- Well-defined error handling patterns
- Static types
- Minimal startup overhead
- Cross-platform compilation
- Rich standard library

# AWS Serverless Architecture | Lambda + Go+ API Gateway

## Why Serverless?  

- <b>Move fast, innovate:</b> Focus on application logic, not on infrastructure.  
- <b>Cost Savings:</b> Save on devops resources. Pay for exactly the number of requests and invocations needed.  
- <b>Scale without worry:</b> No additional capacity needs provisioning to handle your workload at peak. During lulls in activity costs are proportional to usage.

## Serverless Architecture

The first thing to understand about serverless architecture is that it's not about the absence of server. What it means is that as a developer you are not concerned with server. You provide a code piece to an environment and it will be executed and results will be returned to you. You are just responsible for providing the code piece and generally the code piece has to adhere to some contract, so that the execution environment can understand it. AWS Lambda is an example of serverless architecture.  

Serverless architecture is important from the perspective that it saves the developers from dealing with servers and in dealing with deployments. Write code, load it in execution environment and run it. This helps developers in focusing on the core stuff than on other aspects of development.


## go-lambda

go-lambda is a multi-purpose tool for creating and managing AWS Lambda instances backed by arbitrary Go code. Since there is no official support of Go, this tool automatically generates a wrappig module in Python 2.7 which is able to pass data back and forth to the Go land.

Features at glance:

<10ms startup time, feels like native experience;
Resulting source.zip size is only 1.0M and in most cases has 2 files;
Easy to use: start writing your own lambdas in Go just in few minutes;
Relies on the official AWS SDK for Go while making all the requests, security is guaranteed;
No any boilerplate or "all-in-one" aims: the tool is made to do its job and nothing else. 
Installation

$ go get github.com/xlab/go-lambd

## Connecting MongoDB with Golang

The Golang driver for MongoDB is called mgo. Using Apex to create a function that connects to Compose's MongoDB is almost as straightforward as the simpleGo function which we have been reviewing. 
MongoDB is a document datastore. Rather than storing spreadsheet like tables (columns and rows), it’s more like a set of folders (or buckets) into which JSON files (documents) can be put, then queried.

In this section, we’ll write a simple comments API in Go that:

Connects to MongoDB
Inserts some comments data
Reads that comments data

## Starting MongoDB

Before we get started, be sure to:
- Install MongoDB  
- Get the mgo package — a ‘driver’ that will let us interact with MongoDB  
- Create a new folder called `commentsapp` — this is where our Go code will live. Inside that, create a subfolder called `db` which is where we’ll ask MongoDB to keep the data.  
- Start MongoDB by running the following in a command terminal after navigating to the `commentsapp` folder:
mongod --dbpath=”./db”

You should see some output including something like the line: “waiting for connections on port 27017” — then we know we’re good to go.

## Interacting with data
Now that the structure for our program is in place, we are ready to write the two handlers that will actually interact with MongoDB to do the work for us.

## Inserting data
Now we will write our `handleInsert` function that will take the data from an http.Request, and insert it into the database. We’re going to decode the request body manually here, but you might want to consider some patterns for decoding and validating input. Once we’ve decided the comment data, we’ll set the time and a give it a unique ID before inserting it into the database. Finally, we’ll redirect the user to a path that uniquely describes the new comment.

To write a function that returns an Adapter that will setup (and teardown) the database session for our handlers and store it in a context for our handlers to get it later.
Add the following code to main.go:

func withDB(db *mgo.Session) Adapter {
  // return the Adapter
  return func(h http.Handler) http.Handler {
    // the adapter (when called) should return a new handler
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      // copy the database session      
      dbsession := db.Copy()
      defer dbsession.Close() // clean up 
      // save it in the mux context
      context.Set(r, "database", dbsession)
      // pass execution to the original handler
      h.ServeHTTP(w, r)
    })
  }
}

## Setup the handler
Since this is the only handler our API will use, we can now modify our main function to connect to MongoDB, adapt the handle function we just added, and tell the http package to serve it on port :8080.

Modify the main function:
func main() {
  // connect to the database
  db, err := mgo.Dial("localhost")
  if err != nil {
    log.Fatal("cannot dial mongo", err)
  }
  defer db.Close() // clean up when we’re done
  // Adapt our handle function using withDB
  h := Adapt(http.HandlerFunc(handle), withDB(db))
  // add the handler
  http.Handle("/comments", context.ClearHandler(h))
  // start the server
  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal(err)
  }
}




# References

[https://www.pluralsight.com/blog/software-development/golang-get-started](https://www.pluralsight.com/blog/software-development/golang-get-started)

[https://thenewstack.io/a-closer-look-at-golang-from-an-architects-perspective/](https://thenewstack.io/a-closer-look-at-golang-from-an-architects-perspective/)

