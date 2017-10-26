var socket = io.connect('http://localhost:8080');
var app = angular.module("myShoppingList", []); 
app.controller("myCtrl", function($scope) {
    $scope.quantity = 1;
    $scope.totalPrice = "$10";
    $scope.products = [{name: "Hound's Chicken", quantity: 1, price : "$10"},{name: "Fire and Blood spicy pizza", quantity: 1, price : "$12"},{name: "Spices are coming", quantity: 1, price : "$16"},{name: "Winterfell special", quantity: 1, price : "$20"}];
    //$scope.products =  ["Hound's Chicken", "Fire and Blood spicy pizza", "Spices are coming", "Winterfell special"];
  
    $scope.removeItem = function (x) {
        console.log('emit pizza' + x)
        socket.emit('removePizza', { pizza: x });
       }
    $scope.addQuantity = function(x){
       console.log('add quantity' + x)
       socket.emit('addQuantity', { pizza: x }); 
    }
  /*  $scope.calculateTotal = function{
        for(var pizza in $scope.products)
        {

        }
    }*/
    socket.on('removePizza', function (data) {
        console.log('message from server' + data.pizza);
        $scope.products.splice(data.pizza, 1);
        $scope.$apply();
       });

    $scope.closeConnection = function () {
        console.log('close connection');
        socket.emit('closeConnection');
    }

});



  