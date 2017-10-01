## What is an API?
An Application-Programming Interface (API) is a set of programming instructions and standards for accessing a Web-based software application. In simple terms, an API is a software-to-software interface. With APIs, applications talk to each other without any user knowledge or intervention.


## APIs in Cloud 
An API resembles Software as a Service (SaaS), since software developers don't have to start from scratch every time they write a program. Cloud capabilities, such as reusing components, connecting different components or scaling of services on demand, have already begun to shift the focus of APIs from simple RPC-programming models to RESTful web models, and even to what is called "lambda models" of services that can be instantly scaled as needed in the cloud.


## Why REST?
REST (REpresentational State Transfer) is an architectural style, and an approach to communications that is often used in the development of Web services. Since REST does not leverage as much bandwidth as compared to heavyweight SOAP, it is considered to be a better fit for use over the Internet.


## Amazon API Gateway
Amazon API Gateway is an AWS service that enables developers to create, publish, maintain, monitor, and secure their APIs at any scale. API Gateway is an AWS service that supports the following:

* Creating, deploying, and managing a RESTful API to expose backend HTTP endpoints, AWS Lambda functions, or other AWS services.

* Invoking exposed API methods through the frontend HTTP endpoints.

## Usage of API Gateway in Project
API Gateway lets us create, configure, and host a RESTful API to enable applications to access the AWS Cloud and AWS or other web services. Together with AWS Lambda, API Gateway forms the app-facing part of the AWS serverless infrastructure. AWS Lambda runs the code on a highly available computing infrastructure and API Gateway can be used to expose the Lambda functions through API methods.

## Benefits of using Amazon API Gateway
1. Performance at Any Scale
    * Low-latency for API requests and responses.
    * Traffic-control
    * Caching the output of API calls.
2. Easily Monitor API Activity
3. Create RESTful Endpoints for Existing Services
4. Run APIs Without Servers with AWS Lambda