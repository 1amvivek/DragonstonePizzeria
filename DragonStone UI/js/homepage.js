var increment = 0;
var word_array = ["<i>Fiery</i>", "<i>Hot</i>", "<i>Sizzling</i>", "<i>Cheesy</i>"];
var domElement;
function nextWord() {
  increment++;
  domElement.style.opacity = 0;
  if (increment > (word_array.length - 1)) {
    increment = 0;
  }
  setTimeout('slide()', 1000);
}
function slide() {
  domElement.innerHTML = word_array[increment];
  domElement.style.opacity = 1;
  setTimeout('nextWord()', 1000);
}
