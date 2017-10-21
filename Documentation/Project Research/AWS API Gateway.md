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

## Getting Started with Amazon API Gateway having Lambda integration
### Overview of Lambda Proxy integration
* The Lambda proxy integration allows the client to call a single Lambda function in the backend. The function accesses many resources or features of other AWS services such as EC2 or S3, including calling other Lambda functions.
* When a client submits an API request, API Gateway passes to the integrated Lambda function the raw request as-is along with the request data containing request headers, query string parameters, URL path variables, payload, and API configuration data.
* The backend Lambda function being called parses the incoming request data to determine the response that it has to return to the client.
* The integrated Lambda function initially verifies all of the input sources before processing the request and responding to the client with meaningful error messages if any of the required input is missing.

### Setting up API Gateway Proxy Integration
The following three tasks need to be performed:
* Create a proxy resource with a greedy path variable of {proxy+}. This path parameter represents any of the child resources under its parent resource of an API. In other words, /parent/{proxy+} can stand for any resource matching the path pattern of "/parent/*".
* A special method, named ANY, used to define the same integration set up for all supported methods: DELETE, GET, HEAD, OPTIONS, PATCH, POST, and PUT.
* Integrate the resource and method with a backend using the HTTP or Lambda integration type.
   * The HTTP proxy integration, designated by HTTP_PROXY in the API Gateway REST API, is for integrating a method request with a backend HTTP endpoint.
   * The Lambda proxy integration, designated by AWS_PROXY in the API Gateway REST API, is for integrating a method request with a Lambda function in the backend.
   
## Testing an API
1. Test a Method with the API Gateway Console
    * After the creation of the APIs, click on the "Test" option in the "Method Execution" pane.
        * The following information will be displayed:
           * "Request" is the resource's path that was called for the method.
           * "Status" is the response's HTTP status code.
           * "Latency" is the time between the receipt of the request from the caller and the returned response.
           * "Response Body" is the HTTP response body.
           * "Response Headers" are the HTTP response headers.
           * "Logs" are the simulated Amazon CloudWatch Logs entries that would have been written if this method were called outside of the API Gateway console.
           
2. Use Postman to Call an API
   * Enter the endpoint URL of a request in the address bar and choose the appropriate HTTP method from the drop-down list to the left of the address bar.
   * If required, provide the appropriate Authorization using AWS credentials.
   * Click on Send to test the API.

## Reference:
http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-create-api-as-simple-proxy.html
http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-set-up-simple-proxy.html
