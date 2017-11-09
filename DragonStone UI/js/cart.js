var socket = io.connect('http://localhost:8080');
var app = angular.module("myShoppingList", []); 
app.controller("myCtrl", function($scope,$http) {
    $scope.quantity = 1;
    $scope.totalPrice = "$58";
    $scope.products = [];
    $scope.logs =  ["Arun Created Group Cart"];
    
    socket.emit('getCatalog');
    $scope.showModal = function (x) {
      $scope.modalPizzaName = $scope.catalog[x].name;
      $scope.modalPizzaPrize = $scope.catalog[x].price;
      $scope.modalPizzaUrl = $scope.catalog[x].img_url;
      $scope.modalPizzaDesc = $scope.catalog[x].desc;
      $scope.modalPizzaId = $scope.catalog[x].id;
    
      var pop = document.getElementById('item_modal');
      pop.style.display = "block";
      socket.emit('lookingAt',{'pizzaId' : $scope.modalPizzaId});
    }

    $scope.addItem = function () {
        console.log('add pizza' + $scope.modalPizzaId);
        socket.emit('addPizza',{'pizzaId' : $scope.modalPizzaId});
       }

    $scope.removeItem = function (x) {
        console.log('emit pizza'+ x)
        socket.emit('removePizza', { pizzaId: x });
       }
    $scope.addQuantity = function(x){
       console.log('add quantity' + x)
       socket.emit('addQuantity', { pizzaId: x }); 
    }

    $scope.reduceQuantity = function(x){
       console.log('reduce quantity' + x)
       if($scope.products[x].quantity>1)
         socket.emit('reduceQuantity', { pizzaId: x }); 
    }

    socket.on('catalog', function (data) {
      console.log('received catalog');
        $scope.catalog = data.catalog;
        $scope.$apply();
       });

    socket.on('join', function (data) {
      console.log('new user joined');
        $scope.logs.push(''+data.user+' joined the group cart');
        $scope.$apply();
       });

    socket.on('left', function (data) {
      console.log('new user joined');
        $scope.logs.push(''+data.user+' left the group cart');
        $scope.$apply();
       });

    socket.on('addPizza', function (data) {
        console.log('message from server' + data.pizzaId + "price: " + data.totalPrice);
        $scope.products.push({name: $scope.catalog[data.pizzaId].name, quantity: 1, price : $scope.catalog[data.pizzaId].price});
        $scope.logs.push(''+data.user+' added ' + $scope.catalog[data.pizzaId].name + ' to the cart');
        $scope.totalPrice = data.totalPrice;
        $scope.$apply();
        sessionStorage.cartUuid = data.cartUuid;
        console.log(sessionStorage.cartUuid);
       });

    socket.on('removePizza', function (data) {
        console.log('message from server' + data.pizzaId + "price: " + data.totalPrice);
        $scope.products.splice(data.pizzaId, 1);
        $scope.totalPrice = data.totalPrice;
        $scope.logs.push(''+data.user+' removed ' + $scope.catalog[data.pizzaId].name + ' from the cart');
        $scope.$apply();
       });

    socket.on('addQuantity', function (data) {
        console.log('message from server - addQuantity:' + data.pizzaId);
        $scope.products[data.pizzaId].quantity++;
        $scope.totalPrice = data.totalPrice;
        $scope.logs.push(''+data.user+' added quantity for ' + $scope.catalog[data.pizzaId].name); 
        $scope.$apply();
       });

    socket.on('reduceQuantity', function (data) {
        console.log('message from server - reduceQuantity:' + data.pizzaId);
        $scope.products[data.pizzaId].quantity--;
        $scope.totalPrice = data.totalPrice;
        $scope.logs.push(''+data.user+' reduced quantity for ' + $scope.catalog[data.pizzaId].name);
        $scope.$apply();
       });

    socket.on('lookingAt',function(data){
      //todo: rest call to golang pizza api
        $scope.logs.push(''+data.user+' is looking at ' + $scope.catalog[data.pizzaId].name);
        $scope.$apply();     
      });

    $scope.closeConnection = function () {
        console.log('close connection');
        socket.emit('closeConnection');
    }

});