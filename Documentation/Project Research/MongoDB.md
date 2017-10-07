# MongoDB Overview: 

MongoDB is a cross-platform, document oriented NoSQL database that provides, high performance, high availability, and easy scalability. MongoDB uniquely allows users to mix and match multiple storage engines within a single deployment. This flexibility provides a more simple and reliable approach to meeting diverse application needs for data. Traditionally, multiple database technologies would need to be managed to meet these needs, with complex, custom integration code to move data between the technologies, and to ensure consistent, secure access.

MongoDB stores data as documents in a binary representation called BSON (Binary JSON). Documents that share a similar structure are typically organized as collections. You can think of collections as being analogous to a table in a relational database: documents are similar to rows, and fields are similar to columns.

## Advantages

* MongoDB stores data in flexible, JSON-like documents, meaning fields can vary from document to document and data structure can be changed over time.

* The document model maps to the objects in your application code, making data easy to work with

* Ad hoc queries, indexing, and real time aggregation provide powerful ways to access and analyze your data

* MongoDB is a distributed database at its core, so high availability, horizontal scaling, and geographic distribution are built in and easy to use

* Schema less − MongoDB is a document database in which one collection holds different documents. Number of fields, content and size of the document can differ from one document to another.


## Concurrency Control in MongoDB [source: MongoDB docs]

Concurrency control allows multiple applications to run concurrently without causing data inconsistency or conflicts. This can be achieved in our project with either of the two options: 

* Create a unique field that can have unique values during a write operation to track the changes. This prevents insertions or updates from creating duplicate data. A unique index is created on multiple fields to force uniqueness on that combination of field values.

* Specify the expected current value of a field in the query predicate for the write operations. The two-phase commit pattern provides a variation where the query predicate includes the application identifier as well as the expected state of the data in the write operation.


## Consistency in MongoDB [source: MongoDB docs]

#### Monotonic Writes

MongoDB provides monotonic write guarantees for standalone mongod instances, replica sets, and sharded clusters. Suppose an application performs a sequence of operations that consists of a write operation W1 followed later in the sequence by a write operation W2. MongoDB guarantees that W1 operation precedes W2.

#### Real Time Order

For read and write operations on the primary, issuing read operations with "linearizable" read concern and write operations with "majority" write concern enables multiple threads to perform reads and writes on a single document as if a single thread performed these operations in real time; that is, the corresponding schedule for these reads and writes is considered linearizable.


# MongoDB on the AWS Cloud: 

The MongoDB cluster (version 2.6 or 3.0) makes use of Amazon Elastic Compute Cloud (EC2) and Amazon Virtual Private Cloud, and is launched via a AWS CloudFormation template. You can use the template directly or you can copy and then customize it as needed.  The template creates the following resources:

* VPC with private and public subnets (you can also launch the cluster into an existing VPC).

* A NAT instance in the public subnet to support SSH access to the cluster and outbound Internet connectivity.

* An IAM instance role with fine-grained permissions.

* Security groups

* A fully customized MongoDB cluster with replica sets, shards, and config servers, along with customized EBS storage, all running in the private subnet.

* The document examines scaling, replication, and performance tradeoffs in depth, and provides guidance to help you to choose appropriate types of EC2 instances and EBS volumes.





# Reference: 

https://www.mongodb.com/mongodb-architecture

https://aws.amazon.com/blogs/aws/mongodb-on-the-aws-cloud-new-quick-start-reference-deployment/

https://docs.mongodb.com/manual/core/write-operations-atomicity/

https://docs.mongodb.com/manual/core/read-isolation-consistency-recency/
