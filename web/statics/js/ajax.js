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