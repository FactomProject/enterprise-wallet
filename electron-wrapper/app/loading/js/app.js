//const ipc = require('electron').ipcRenderer

function sendChoiceToMain(secure) {
	if(secure) {
		var dom = document.getElementById("secure-password-input")
		if(checkValidPassword(dom.value)) {
			ipc.send('submitForm', dom.value)
			dom.classList.remove("has-error")
			document.getElementById("error-text").innerHTML = ""
		} else {
			dom.classList.add("has-error")
		}
	} else {
		ipc.send('submitForm', "");
	}
	return false
}

function updateCheckbox() {
	var c = document.getElementById('checkbox');
	if(c.checked) {
		document.getElementById("proceed-button").disabled = false;
	} else {
		document.getElementById("proceed-button").disabled = true;
	}
}

function checkValidPassword(pass) {
	if(pass.length < 8) {
		document.getElementById("error-text").innerHTML = "Password must be at least 8 characters in length"
		return false
	}
	return true
}