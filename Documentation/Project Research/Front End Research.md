# PROJECT RESEARCH – FRONT END
Communication on the web is usually unsynchronized. The Internet has always been this way: the client requests and the server responds. The server can’t decide for itself to send something to the client. It must be the client who reloads the page or takes action to call the server because server doesn’t have the right to talk to the client on its own. Even with Ajax, it is the client who requests first. 
If the application has to be real time, we need a mechanism to push messages to a client(s) from the server without a request from client(s). For e.g. The Shopping Cart / Online Ordering system must support multiple users accessing the Cart concurrently from different Web Browsers. If one user updates the cart, other users’ cart should get updated at real time without any requests.  


## Challenge: 

When data changes on the server, let the clients know without asking.   
## Possible Solutions:
Polling using Ajax requests is not a good option as the server will be flooded with requests and it is not a scalable solution. For the project, we are researching about the following two solutions to tackle the challenge.
* HTML5 Server Sent Events.
* Web sockets.

## SSE
SSE allows you to stream events continuously from your web server to the visitor’s browser. But this is also a one-way communication – From server to client. Online stock quotes, or twitters updating timeline or feed are good examples of an application that could benefit from SSE.
For the client to server communication we still have to send Ajax requests. When the server receives a request from a client, it updates the database and notifies all the subscribed clients. By this way, when a user adds a product to a cart, an Ajax request is sent to the server and server send responses to all the subscribed clients to update the cart using SSE. So, all the Users’ carts are updated concurrently.  

**PROS**  
*	Concurrent updates to cart is achieved
*	SSEs are sent over traditional HTTP. That means they do not require a special protocol or server implementation to get working.
*	Built in support for re-connection

**CONS**  
*	One-way communication 
*	Ajax request and SSE messaging will take a good few seconds, so the cart update will not be real time.
*	SSE suffers from a limitation to the maximum number of open connections, which can be especially painful when opening various tabs as the limit is per browser and set to a very low number (6).

## WEB SOCKETS
Websockets connections can both send data to the browser and receive data from the browser. A good example of an application that could use websockets is a real-time chat application. In practice since everything that can be done with SSE can also be done with Websockets. When a user adds a product to the cart, a request is sent to the server through the socket and the server updates the database. Server sends messages to all the clients subscribed through the same open socket connection.   

**PROS**   
*	Two-way communication.
*	Fast and real time.
*	Concurrent Updates to cart is achieved.
*	Allows any number of connections.
*	Supported by more browsers than SSE.

**CONS**
*	require full-duplex connections and new Web Socket servers to handle the protocol.
*	No built-in support for re-connection

Based on this research and analysing the pros and cons of each, we have decided to go with the websockets for this project to achieve real time multi user shopping cart.
