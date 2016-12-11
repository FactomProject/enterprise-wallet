$("#save-changes").on('click', function(){	
	theme = $("#darkTheme").is(":checked")
	exportKeys = $("#export-keys").is(":checked")
	coinControl = $("#coin-control").is(":checked")

	var SettingsStruct = {
    	Values:[]
	}

	SettingsStruct.Values.push(theme)
	SettingsStruct.Values.push(exportKeys)
	SettingsStruct.Values.push(coinControl)

	j = JSON.stringify(SettingsStruct)
	postRequest("adjust-settings", j, function(resp){
		console.log(resp)
	    obj = JSON.parse(resp)
	    if(obj.Error == "none") {
	      location.reload();
	    } else {
	      SetGeneralError("Error: " + obj.Error)
	    }
	})
})