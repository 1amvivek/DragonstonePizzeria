var pop = document.getElementById('item_modal');
var btn = document.getElementById("myBtn");
var span = document.getElementsByClassName("close")[0];
btn.onclick = function () {
  pop.style.display = "block";
}

// click on "x"
span.onclick = function () {
  pop.style.display = "none";
}

// click outside pop-up window
window.onclick = function (event) {
  if (event.target == pop) {
    pop.style.display = "none";
  }
}