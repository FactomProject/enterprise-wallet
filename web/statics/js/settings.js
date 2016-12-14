$("#save-changes").on('click', function(){	
	theme = $("#darkTheme").is(":checked")
	exportKeys = $("#export-keys").is(":checked")
	coinControl = $("#coin-control").is(":checked")
	importExport = $("#import-export").is(":checked")

	var SettingsStruct = {
    	Values:[]
	}

	SettingsStruct.Values.push(theme)
	SettingsStruct.Values.push(exportKeys)
	SettingsStruct.Values.push(coinControl)
	SettingsStruct.Values.push(importExport)

	j = JSON.stringify(SettingsStruct)
	postRequest("adjust-settings", j, function(resp){
	    obj = JSON.parse(resp)
	    if(obj.Error == "none") {
	    	window.location.href = window.location.href + "?success=true";
	      	//location.reload();
	    } else {
	    	SetGeneralError("Error: " + obj.Error)
	    }
	})
})

$("#export-seed").on('click', function(){
	postRequest("get-seed", "", function(resp){
	    obj = JSON.parse(resp)
	    if(obj.Error == "none") {
	    	var link = document.createElement("a");
			link.download = "WalletSeed.txt";
			link.href = "data:text/plain;charset=UTF-8," + encodeURIComponent(obj.Content);
			link.click();
	    } else {
	    	SetGeneralError("Error: " + obj.Error)
	    }
	})
})

// Import/Export
$("#import-file").on('click', function(e){
  input = document.getElementById('uploaded-file');
  if (!input) {
  	SetGeneralError("Error: Couldn't find the fileinput element.")
  	e.stopPropagation()
  }
  else if (!input.files) {
  	SetGeneralError("This browser doesn't seem to support the `files` property of file inputs.")
  	e.stopPropagation()
  }
  else if (!input.files[0]) {
  	SetGeneralError("Please select a file before clicking 'Import From File'")
  	e.stopPropagation()  
  }
  else {
    file = input.files[0];
    fr = new FileReader();
    fr.onload = receivedText;
    fr.readAsText(file);
    //fr.readAsDataURL(file);
  }
})

// Do action with imported transaction
function receivedText() {
	is = fr.result
	len = is.split(" ")
	if(len.length != 12) {
		SetGeneralError("Seed must be 12 words");
		return
	}
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
		console.log(obj)
		if(obj.Error == "none") {
	    	SetGeneralSuccess("Seed has been changed to: " + obj.Content)
	    } else {
	    	SetGeneralError("Error: " + obj.Error)
	    }
	})
})
