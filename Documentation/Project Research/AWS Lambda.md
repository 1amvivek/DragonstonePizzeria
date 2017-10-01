# What is AWS Lambda?
Lambdas are the epitome of Ephemeral Computing which, as we learned in the CMPE 281 class, is a key cloud selection strategy. Lambdas are run on containers that are created and destroyed instantaneously on demand. Charges are only applied on the time the Lambda spends running which is very limited. The serverless architecture of AWS Lambda makes is highly scalable and constantly available.

## Lifecycle of a Lambda
Lambdas are associated with AWS triggers which can be AWS services, mobile apps or an HTTP endpoint and can perform activities on most available AWS services including backends like S3, DynamoDB, and Mongo DB instances.
  
Lambdas have a handler method that accepts an event object which is passed to it by the trigger. It performs computations based on this event object and returns results back to the trigger point.
  
When first triggered, AWS fetches a zipped version of the required lambda code, launches a container in a cluster, runs the code inside this container and the event trigger is then passed to the handler function in the lambda. The lambda performs the tasks it is required to do and returns the required result to the endpoint. At this point, AWS can either immediately turn off the container or keep it temporarily in a "freeze" or "sync" state preventing it from performing tasks in the background and ready to process any subsequent triggers to the lambda. If kept inactive, the container is disposed of.
  
    
## Usage of Lambda in Project
Given the nature of the project and the means of the functioning of AWS Lambda, it is ideal to create a DB cluster with MongoDB and have these accessed by lambda functions made in Go. These lambda functions will have REST endpoints with various URLs and REST verbs compliant with Level 3 of Richardson's Maturity Model. At each endpoint, a standalone lambda will be called by AWS API gateway which will, in turn, perform required computations and actions on the MongoDB database following which it will return a JSON object back to the REST endpoint containing data compliant to Richardson's Maturity Model -i.e. Hypermedia.
  
The front-end system will be designed only to hit the exposed REST endpoints with required data thereby making it a completely changeable component.

## Advantages of using Lambda
The ephemeral nature of Lambda ensures that it will not store any state within the function itself. Given the strategies (possibly hashing to decide which node in the cluster to store data and replication to ensure that data will be available in case a node is lost) that will be used to store data efficiently in a cluster, AWS Lambda will be able to function concurrently while maintaining consistency within itself. The lack of a server directly implies there will be no network partitions in the lambda layer in a truly REST-ful fashion.
  
The Lambdas can also be designed to perform additional data validation tasks to ensure the quality of data remains compliant with the requirement.
  
It's serverless architecture and function-oriented programming nature (only one handler method which can call other functions within the Lambda code) makes it ideal to use with the Go programming language as well.

## References
* [AWS Lambda - Serverless Compute](https://aws.amazon.com/lambda/?sc_channel=PS&sc_campaign=pac_ps_q3&sc_publisher=google&sc_medium=lambda_b_pac_q32017&sc_content=lambda_e&sc_detail=aws%20lambda&sc_category=lambda&sc_segment=webp&sc_matchtype=e&sc_county=US&sc_geo=namer&sc_outcome=pac&s_kwcid=AL!4422!3!217987870463!e!!g!!aws%20lambda&ef_id=WajifAAAAHiQf3xO:20171001010813:s)

