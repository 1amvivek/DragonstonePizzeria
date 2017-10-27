var app = require('express')();
var server = require('http').Server(app);
var io = require('socket.io')(server);
var totalPrice = "$40";

server.listen(8080, function() {
  console.log('Server running at http://127.0.0.1:8080/');
});

io.sockets.on('connection', function(socket) {
  
  console.log('new client:' + socket.id);
   
  socket.on('removePizza', function (data) {
      console.log(data.pizza);
      //todo: rest call to golang pizza api
      io.sockets.emit('removePizza', { pizza: data.pizza,totalPrice:totalPrice });
   
    });

  socket.on('addQuantity', function (data) {
      console.log(data.pizza);
      //todo: rest call to golang pizza api
      io.sockets.emit('addQuantity', { pizza: data.pizza,totalPrice:totalPrice });
   
    });

  socket.on('reduceQuantity', function (data) {
      console.log(data.pizza);
      //todo: rest call to golang pizza api
      io.sockets.emit('reduceQuantity', { pizza: data.pizza,totalPrice:totalPrice });
   
    });


  socket.on('closeConnection',function(){
      console.log('Client disconnects'  + socket.id);
      socket.disconnect();
  });

  socket.on('disconnect', function() {
      console.log('Got disconnected!'  + socket.id);
      socket.disconnect();
   });
});