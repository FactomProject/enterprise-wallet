$("#backup-input-verify-button").on('click', function(){
	$("#backup-input-verify-button").addClass("backup-btn-checking")

	var seedSingle = $("#given-seed").attr("value");
	var inputedSeedSplit = new Array(12);

	$(".backup-input").each(function(){
		var i = Number($(this).attr("index"))
		inputedSeedSplit[i-1] = $(this).val()
	})

	var clear = function() {
		$("#backup-input-verify-button").removeClass("backup-btn-failed")
		$("#backup-input-verify-button").removeClass("backup-btn-checking")
		$("#backup-input-verify-button").removeClass("backup-btn-verified")
	}

	if(inputedSeedSplit.join(" ") === seedSingle) {
		clear()
		$("#backup-input-verify-button").addClass("backup-btn-verified")
	} else {
		clear()
		$("#backup-input-verify-button").addClass("backup-btn-failed")
		setTimeout(clear, 3000);
	}
})

$("#backup-html-form").on("submit",function(event){event.preventDefault()})

function LoadBackup0() {
	LoadAddresses(false)
}

function updateBackupConfirmCheckbox() {
	var c = $("#wrote-down-confirm-checkbox")
	if(!c.is(':checked')) {
		c.prop('checked', true);
	} else {
		c.prop('checked', false);
	}

	if(c.is(':checked')) {
		document.getElementById("wrote-down-confirm").disabled = false;
	} else {
		document.getElementById("wrote-down-confirm").disabled = true;
	}
}

// Importing

$("#import-input-confirm-button").on('click', function(){
	$("#import-input-confirm-button").addClass("backup-btn-checking")

	var seedSingle = $("#given-seed").attr("value");
	var inputedSeedSplit = new Array(12);

	$(".import-input").each(function(){
		var i = Number($(this).attr("index"))
		inputedSeedSplit[i-1] = $(this).val()
	})

	var clear = function() {
		$("#import-input-confirm-button").removeClass("backup-btn-failed")
		$("#import-input-confirm-button").removeClass("backup-btn-checking")
		$("#import-input-confirm-button").removeClass("backup-btn-verified")
	}


	seed = inputedSeedSplit.join(" ")
	var SeedStruct  = {
    	Seed:seed,
  	}
  	j = JSON.stringify(SeedStruct)
	postRequest("import-seed", j, function(resp) {
		obj = JSON.parse(resp)
		if(obj.Error == "none") {
			clear()
			$("#import-input-confirm-button").addClass("backup-btn-verified")
	    	SetGeneralSuccess("Seed has been changed to: " + obj.Content)
	    } else {
	    	clear()
			$("#import-input-confirm-button").addClass("backup-btn-failed")
	    	SetGeneralError("Error: " + obj.Error)
	    	setTimeout(clear, 3000);
	    }
	})

	/*if(inputedSeedSplit.join(" ") === seedSingle) {
		clear()
		$("#import-input-confirm-button").addClass("backup-btn-verified")
	} else {
		clear()
		$("#import-input-confirm-button").addClass("backup-btn-failed")
		setTimeout(clear, 3000);
	}*/
})


// Do action with imported transaction
function receivedText() {
	is = fr.result
	len = is.split(" ")
	if(len.length != 12) {
		SetGeneralError("Seed must be 12 words");
		return
	}
	document.getElementById('data-expand').click()
	$("#import-seed-reveal-text").text(fr.result)
	$("#import-seed-reveal-cancel").click()
}

$("#import-seed-reveal-confirm").on('click', function(){
	seed = $("#import-seed-reveal-text").text()
	var SeedStruct  = {
    	Seed:seed,
  	}
  	j = JSON.stringify(SeedStruct)
	postRequest("import-seed", j, function(resp) {
		obj = JSON.parse(resp)
		if(obj.Error == "none") {
	    	SetGeneralSuccess("Seed has been changed to: " + obj.Content)
	    } else {
	    	SetGeneralError("Error: " + obj.Error)
	    }
	})
})