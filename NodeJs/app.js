var app = require('express')();
var server = require('http').Server(app);
var io = require('socket.io')(server);
var totalPrice = "$40";


//Handle from post data
var bodyParser = require('body-parser');

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
