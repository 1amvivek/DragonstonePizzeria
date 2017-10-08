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


# References

[https://www.pluralsight.com/blog/software-development/golang-get-started](https://www.pluralsight.com/blog/software-development/golang-get-started)

[https://thenewstack.io/a-closer-look-at-golang-from-an-architects-perspective/](https://thenewstack.io/a-closer-look-at-golang-from-an-architects-perspective/)





