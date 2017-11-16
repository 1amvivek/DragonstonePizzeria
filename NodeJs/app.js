var app = require('express')();
var server = require('http').Server(app);
var io = require('socket.io')(server);
var totalPrice = "$40";
var Client = require('node-rest-client').Client;
var client = new Client();
var request = require("request");

//sample REST API URL and arguments
var usersApiGetUrl = "http://localhost:3000/orders/{SerialNumber}";
var usersApiPostUrl = "http://localhost:3000/order";
var usersApiPutUrl = "http://localhost:3000/order";
//change the data to respective post data
var usersPostArgs = {
          data: {
                 "name":"test group cart",
                 "owner":"Arun Ram",
                 "users":"Arun Ram",
                 "cartserialnumber":"12345678",
                 "logsserialnumber":"23456179"
                },
          headers: { "Content-Type": "application/json" }
  };

var usersPutArgs = {
          data: {
                 "SerialNumber":"123456789",
                 "users":"Arun Ram"
                 },
          headers: { "Content-Type": "application/json" }
  };


//Handle from post data
var bodyParser = require('body-parser');

var path = require('path'); 
var catalog;
var connections = [];
var cartApiPostData;
var logsApiPostData;
var groupCartName;
var owner;
var Users;

app.use(bodyParser.urlencoded({ extended: true })); 

//app.use(express.bodyParser());
app.use(function(req, res, next) {
    res.header("Access-Control-Allow-Origin", "*");
    res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
    next();
});

 
app.post('/login', function (req, res, next) {
  var email = req.body.email;
        var password = req.body.password;

  
      if (req.body.email && req.body.email=== email && req.body.password && req.body.password === password) {
        req.session.authenticated = true;
        res.redirect('/secure');
      } else {
        req.flash('error', 'Username and password are incorrect');
        
      }
   res.send('username sent to Node Server: "' + email + '".'+ '<br/> password sent to Node Server: "' + password + '".');
   response = {
      username:email,
      password:password
      
   };

   console.log(response);

    });

// Register User
app.post('/register', function(req, res){
  var name = req.body.name;
  var email = req.body.email;
  var password = req.body.password;
  var password2 = req.body.cpassword;

  // Validation
  req.checkBody('name', 'Name is required').notEmpty();
  req.checkBody('email', 'Email is required').notEmpty();
  req.checkBody('email', 'Email is not valid').isEmail();
  req.checkBody('password', 'Password is required').notEmpty();
  req.checkBody('cpassword', 'Passwords do not match').equals(req.body.password);

  var errors = req.validationErrors();

  if(errors){
    res.render('register',{
      errors:errors
    });
  } else {
    var newUser = new User({
      name: name,
      email:email,
      password: password
    });

    User.createUser(newUser, function(err, user){
      if(err) throw err;
      console.log(user);
    });

    req.flash('success_msg', 'You are registered and can now login');
  }

res.send('You are registered successfully');
   
console.log('New User successfully registered');

});


server.listen(8080, function() {
  console.log('Server running at http://127.0.0.1:8080/');
});

//handle socket connections
io.sockets.on('connection', function(socket) {
  
  console.log('new client:' + socket.id);
  connections.push(socket.id); 
  //sendRestGetRequest(cartApiUrl,123);
    
  socket.on('createGroupCart', function (socketData) {
      groupCartName = socketData.groupName;
      owner = socketData.yourName;
      //sendRestPostRequest(cartApiPostUrl,cartPostArgs,postCartApiCallBack);
      postCartApiCallBack();
     });

  postCartApiCallBack = function(postData){
      //cartApiPostData = postData;
      //sendRestPostRequest(logsApiPostUrl,logsPostArgs,postLogsApiCallBack);
      postLogsApiCallBack();
  }

  postLogsApiCallBack = function(postData){
      //logsApiPostData = postData;
      var args = usersPostArgs;
      args.name = groupCartName;
      args.owner = owner;
      args.users = owner;
      //args.cartserialnumber = cartApiPostData.CartSerialNumber;
      //args.logsserialnumber = cartApiPostData.LogsSerialNumber;
      sendRestPostRequest(usersApiPostUrl,args,postUsersApiCallBack);
     }

  postUsersApiCallBack = function(postData){
    var logs = (owner +' created the cart');
    socket.emit('createGroupCart',{response:postData, logs : logs});
  }

  socket.on('joinGroupCart', function (socketData) {
      sendRestGetRequest(usersApiGetUrl,getUsersApiCallBack,socketData);
    });

  getUsersApiCallBack = function(postData,socketData){
    var logs = (socketData.yourName +' joined the cart');
    io.sockets.emit('join', {logs : logs });
    socket.emit('joinGroupCart',{response:postData});
    
    //get json and send to update $scope.products in cart.js
    //sendRestGetRequest(cartApiGetUrl,getCartApiCallBack,socketData);
    //get json and send to update $scope.logs in cart.js
    //sendRestGetRequest(logsApiGetUrl,getLogsApiCallBack,socketData);
      
      var args = usersPutArgs;
      args.data.users = socketData.yourName;
      args.data.SerialNumber = socketData.SerialNumber;  
      sendRestPutRequest(usersApiPutUrl,args);
  }

  socket.on('addPizza', function (socketData) {
      //console.log(data.pizzaId);
      //edit args to respective args from selected pizza details
      //selected pizza details is available in socketData
      var logs = (socket.id +' added ' + socketData.pizzaName + ' to the cart');
      io.sockets.emit('addPizza',{pizzaId:socketData.pizzaId,pizzaName : socketData.pizzaName, user : socket.id,logs : logs});
   
     });
  
  socket.on('addQuantity', function (socketData) {
      //console.log(socketData.pizzaId);
      //todo: rest call to golang pizza api var args = CartPutArgs;
      
    io.sockets.emit('addQuantity', { pizzaId: socketData.pizzaId,pizzaName : socketData.pizzaName,user : socket.id});   
   
    });

  socket.on('removePizza', function (socketData) {
      //console.log(socketData.pizzaId);
      //todo: rest call to golang pizza 
      io.sockets.emit('removePizza', { pizzaId: socketData.pizzaId,pizzaName : socketData.pizzaName,user : socket.id });
   
    });

  socket.on('reduceQuantity', function (socketData) {
      //console.log(socketData.pizzaId);
      //todo: rest call to golang pizza api
      io.sockets.emit('reduceQuantity', { pizzaId: socketData.pizzaId,pizzaName : socketData.pizzaName,user : socket.id });
   
    });

  socket.on('lookingAt',function(socketData){
      //todo: rest call to golang pizza api
      io.sockets.emit('lookingAt', { pizzaId: socketData.pizzaId,pizzaName : socketData.pizzaName,user : socket.id });
  });

  socket.on('getCatalog',function(){
    catalog = [{id:'0',name: "Pepperoni pizza", price : "$12",img_url:"img/product/1.jpg",desc:"This is a medium spicy Pepperoni Pizza with Tomato sauce, triple Pepperoni and mozzarella cheese."},
    {id:'1',name: "Pizza 2", price : "$10",img_url:"img/product/2.jpg",desc:"This is a medium spicy pizza 2."},
    {id:'2',name: "Pizza 3", price : "$14",img_url:"img/product/3.jpg",desc:"This is a medium spicy pizza 3."},
    {id:'3',name: "Pizza 4", price : "$15",img_url:"img/product/4.jpg",desc:"This is a medium spicy pizza 4."},
    {id:'4',name: "Pizza 5", price : "$18",img_url:"img/product/5.jpg",desc:"This is a medium spicy pizza 5."}];
    socket.emit('catalog',{catalog:catalog,user : socket.id});
    console.log('sent catalog');
  });

  socket.on('closeConnection',function(){
      console.log('Client disconnects'  + socket.id);
      socket.disconnect();
      removePlayer(socket.id);
      io.sockets.emit('left', {user : socket.id });
  });

  socket.on('disconnect', function() {
      console.log('Got disconnected!'  + socket.id);
      socket.disconnect();
      io.sockets.emit('left', {user : socket.id });
      removePlayer(socket.id);
   });
});

function removePlayer(item)
{
var index = connections.indexOf(item);
connections.splice(index, 1);
}


function sendRestGetRequest(url,callback,socketData){
// direct way 
url = url.replace('{SerialNumber}',socketData.SerialNumber);
console.log(url);
client.get(url, function (data, response) {
    // parsed response body as js object 
    //console.log(data);
    callback(data,socketData);
    // raw response 
    //console.log(response);
});

};

function sendRestPostRequest(url,args,callback){
// direct way 
console.log(url);
console.log(args);
client.post(url,args, function (data, response) {
    // parsed response body as js object 
    //console.log(data);
    //replace with uuid varibale returned in post json response
    callback(data);
    // raw response 
    //console.log(response);
});

};

function sendRestPutRequest(url,args){
// direct way 
console.log(url);
console.log(args);
client.put(url,args, function (data, response) {
    // parsed response body as js object 
    console.log(data);
    //replace with uuid varibale returned in post json response
    //callback(data,socketData);
    // raw response 
    //console.log(response);
});

};
