var pop = document.getElementById('item_modal');
var span = document.getElementsByClassName("close")[0];

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
