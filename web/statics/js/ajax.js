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