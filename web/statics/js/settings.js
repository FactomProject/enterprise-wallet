$("#save-changes").on('click', function(){	
	theme = $("#darkTheme").is(":checked")
	exportKeys = $("#export-keys").is(":checked")
	coinControl = $("#coin-control").is(":checked")
	importExport = $("#import-export").is(":checked")
    fd = $("input[name='factomd']:checked").val()
    if (fd == "custom") {
        fd = $("#factomd-location").val()
    }

	var SettingsStruct = {
        Values:[],
    	FactomdLocation:""
	}

	SettingsStruct.Values.push(theme)
	SettingsStruct.Values.push(exportKeys)
	SettingsStruct.Values.push(coinControl)
    SettingsStruct.Values.push(importExport)
	SettingsStruct.FactomdLocation = fd

	j = JSON.stringify(SettingsStruct)
	postRequest("adjust-settings", j, function(resp){
	    obj = JSON.parse(resp)
	    if(obj.Error == "none") {
	    	if((window.location.href).includes("success")){
	    		window.location.href = window.location.href
	    	} else {
	    		window.location.href = window.location.href + "?success=true"
	    	}
	      	//location.reload();
	    } else {
	    	SetGeneralError("Error: " + obj.Error)
	    }
	})
})

$("#customFactomd").on('click', function(){
	if($("#customFactomd").is(":checked")){
		$("#factomd-location-container").removeClass("hide")
	} else {
		$("#factomd-location-container").addClass("hide")
	}
})

$("#export-seed").on('click', function(){
	postRequest("get-seed", "", function(resp){
	    obj = JSON.parse(resp)
	    if(obj.Error == "none") {
	    	saveTextAsFile(obj.Content, "WalletSeed.txt")
	    } else {
	    	SetGeneralError("Error: " + obj.Error)
	    }
	})
})

//selected = false
// Import/Export
// $("#settings-import-file").on('click', function(e){
// 	document.getElementById('settings-uploaded-file').click()
// })


// $("#settings-uploaded-file").on('change', function(){
// 	input = document.getElementById('settings-uploaded-file');
// 	if (!input) {
// 		SetGeneralError("Error: Couldn't find the fileinput element.")
// 	}
// 	else if (!input.files) {
// 		SetGeneralError("This browser doesn't seem to support the `files` property of file inputs.")
// 	}
// 	else if (!input.files[0]) {
// 		SetGeneralError("Please select a file before clicking 'Import From File'")
// 	}
// 	else {
// 	file = input.files[0];
// 	fr = new FileReader();
// 	fr.onload = receivedText;
// 	fr.readAsText(file);
// 	//fr.readAsDataURL(file);
// 	}
// })

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
