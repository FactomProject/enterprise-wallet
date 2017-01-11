FCTDecminalLength = 8 // Number of decimals places

function getRequest(item, func) {
  var req = new XMLHttpRequest()

  req.onreadystatechange = function() {
    if(req.readyState == 4) {
      func(req.response)
    }
  }
  req.open("GET", "/GET?request=" + item, true)
  req.send()
}

function postRequest(request, jsonObj, func) {
  var req = new XMLHttpRequest()

  req.onreadystatechange = function() {
    if(req.readyState == 4) {
      func(req.response)
    }
  }

  var formData = new FormData();
  formData.append("request", request)
  formData.append("json", jsonObj)

  req.open("POST", "/POST")
  req.send(formData)
}

// Jquery on all pages
$(window).load(function() {
    updateBalances()
});
setInterval(updateBalances,5000);

// Updates total balances on the page
function updateBalances() {
  getRequest("balances", function(resp){
        obj = JSON.parse(resp)
        if (obj.Error != "none") {
          return
        } 

        $("#ec-balance").text(obj.Content.EC)
        fcBal = formatFC(obj.Content.FC)
        $("#factoid-balance").text(fcBal[0] + ".")
        if(fcBal.length > 1) {
          $("#factoid-balance-trailing").text(fcBal[1])
        } else {
          $("#factoid-balance-trailing").text(0)
        }
  })
}

function formatFC(fcBalance){
  dec = FCTNormalize(fcBalance)
  decStr = dec.toString()
  decSplit = decStr.split(".")

  return decSplit
}

function FCTNormalize(fct) {
  return Number((fct/1e8).toFixed(FCTDecminalLength))
}

Synced = false
// On most pages
checkSynced()
setInterval(checkSynced,3000);
function checkSynced(){
  getRequest("synced", function(resp){
    obj = JSON.parse(resp)
    // console.log(obj)
    if (!Synced && obj.Content.Synced == true) {
      $("#synced-indicator").slideUp(100)
      Synced = true
    }
  })
}

function SetGeneralError(err) {
  $("#success-zone").slideUp(100)
  $("#error-zone").text(err)
  $("#error-zone").slideDown(100)
}

function SetGeneralSuccess(mes) {
  $("#error-zone").slideUp(100)
  $("#success-zone").text(mes)
  $("#success-zone").slideDown(100)
}

function saveTextAsFile(text, filename) {
    var textToWrite = text
    var textFileAsBlob = new Blob([textToWrite], { type: 'text/plain' })
    var fileNameToSaveAs = filename

    var downloadLink = document.createElement("a");
    downloadLink.download = fileNameToSaveAs;
    window.URL = window.URL || window.webkitURL;
    downloadLink.href = window.URL.createObjectURL(textFileAsBlob);
    downloadLink.style.display = "none";
    document.body.appendChild(downloadLink);
    downloadLink.click();
}