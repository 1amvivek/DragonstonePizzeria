var socket = io.connect('http://localhost:8080');
var app = angular.module("myShoppingList", []); 
app.controller("myCtrl", function($scope,$http) {
    $scope.quantity = 1;
    $scope.totalPrice = "$0";
    $scope.products = [];
    $scope.logs =  [];
    $scope.shoppingCartText = "Your Cart";
    
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

    $scope.createGroupCart = function(x){
         var groupName = document.getElementById("groupName").value;
         var yourName = document.getElementById("yourName").value;
         sessionStorage["myName"] = yourName;
         socket.emit('createGroupCart', { groupName: groupName, yourName:sessionStorage.myName }); 
    }

    $scope.joinGroupCart = function(x){
         var SerialNumber = document.getElementById("groupId").value;
         var yourName = document.getElementById("newUserName").value;
         sessionStorage["myName"] = yourName;
         socket.emit('joinGroupCart', { SerialNumber: SerialNumber, yourName:sessionStorage.myName }); 
    }

    $scope.addItem = function () {
        console.log('add pizza' + $scope.modalPizzaId);
        socket.emit('addPizza',{'pizzaId' : $scope.modalPizzaId,'pizzaName' : $scope.catalog[$scope.modalPizzaId].name});
       }

    $scope.removeItem = function (x) {
        console.log('emit pizza'+ x)
        socket.emit('removePizza', { pizzaId: x });
       }
    $scope.addQuantity = function(x){
       console.log('add quantity' + x)
       socket.emit('addQuantity', { pizzaId: x,'pizzaName' : $scope.catalog[$scope.modalPizzaId].name }); 
    }

    $scope.reduceQuantity = function(x){
       console.log('reduce quantity' + x)
       if($scope.products[x].quantity>1)
         socket.emit('reduceQuantity', { pizzaId: x }); 
    }

    socket.on('createGroupCart', function (data) {
        sessionStorage.groupSerialNumber = data.response.SerialNumber;
        sessionStorage.CartSerialnumber = data.response.CartSerialNumber;
        sessionStorage.LogsSerialnumber = data.response.LogsSerialNumber;
        $scope.logs.push(data.logs);
        $scope.shoppingCartText = "Your Cart ID:" + sessionStorage.groupSerialNumber;
        $scope.$apply();
       });


    socket.on('catalog', function (data) {
      console.log('received catalog');
        $scope.catalog = data.catalog;
        $scope.$apply();
       });

    socket.on('join', function (data) {
        console.log('new user joined');
        $scope.logs.push(data.logs);
        $scope.$apply();
       });

    socket.on('joinGroupCart', function (data) {
        sessionStorage.groupSerialNumber = data.response.SerialNumber;
        sessionStorage.CartSerialnumber = data.response.CartSerialNumber;
        sessionStorage.LogsSerialnumber = data.response.LogsSerialNumber;
        });

    socket.on('left', function (data) {
      console.log('new user joined');
        $scope.logs.push(''+data.user+' left the group cart');
        $scope.$apply();
       });

    socket.on('addPizza', function (data) {
        console.log('message from server' + data.pizzaId + "price: " + data.totalPrice);
        $scope.products.push({name: data.pizzaName, quantity: 1, price : $scope.catalog[data.pizzaId].price});
        $scope.logs.push(data.logs);
        calculateTotal();
        $scope.$apply();
        //sessionStorage.cartUuid = data.cartUuid;
        //console.log(sessionStorage.cartUuid);
       });  

    socket.on('removePizza', function (data) {
        console.log('message from server' + data.pizzaId + "price: " + data.totalPrice);
        $scope.products.splice(data.pizzaId, 1);
        calculateTotal();
        $scope.logs.push(''+data.user+' removed ' + $scope.catalog[data.pizzaId].name + ' from the cart');
        $scope.$apply();
       });

    socket.on('addQuantity', function (data) {
        console.log('message from server - addQuantity:' + data.pizzaId);
        $scope.products[data.pizzaId].quantity++;
        calculateTotal();
        $scope.logs.push(''+data.user+' added quantity for ' + $scope.catalog[data.pizzaId].name); 
        $scope.$apply();
       });

    socket.on('reduceQuantity', function (data) {
        console.log('message from server - reduceQuantity:' + data.pizzaId);
        $scope.products[data.pizzaId].quantity--;
        calculateTotal();
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

    calculateTotal = function(){
      var products = $scope.products;
      var totalPrice = 0;
      for(var i=0;i<products.length;i++){ 
        var price = Number((products[i].price).replace("$", ""));
        var amount =  price * products[i].quantity;
        totalPrice = totalPrice + amount;
      }
      $scope.totalPrice = "$" + totalPrice;
    }
});

