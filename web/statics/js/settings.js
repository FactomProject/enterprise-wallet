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
	      location.reload();
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