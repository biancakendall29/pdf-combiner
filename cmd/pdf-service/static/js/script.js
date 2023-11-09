let dropZone = document.getElementById("drop-zone");
let fileInput = document.getElementById("file-input");

dropZone.addEventListener("click", function () {
  fileInput.click();
});

dropZone.addEventListener("dragover", function (e) {
  e.preventDefault();
  dropZone.classList.add("dragging");
});

dropZone.addEventListener("dragleave", function (e) {
  dropZone.classList.remove("dragging");
});

dropZone.addEventListener("drop", function (e) {
  e.preventDefault();
  dropZone.classList.remove("dragging");

  let files = e.dataTransfer.files;
  fileInput.files = files;
});
