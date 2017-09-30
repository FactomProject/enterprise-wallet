const ipc = require('electron').ipcRenderer


const {shell} = require('electron')
function openBlog(link) {
	shell.openExternal(link)
}

function sendChoiceToMain(secure) {
	if(secure) {
		var dom = document.getElementById("secure-password-input")
		if(checkValidPassword(dom.value)) {
			ipc.send('submitForm', dom.value)
			clearError()
		} else {
			dom.classList.add("has-error")
		}
	} else {
		ipc.send('submitForm', "");
	}
	return false
}

function clearError() {
	document.getElementById("secure-password-input").classList.remove("has-error")
	document.getElementById("error-text").innerHTML = ""
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

$("#termsbox").bind('scroll', function() {
	if($("#termsbox").scrollTop() + $("#termsbox").innerHeight() >= .95 * $("#termsbox").prop('scrollHeight')) {
		// Toggle class to show proceed
	}
});
