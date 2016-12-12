$("#generate-source").on("change", function(){
	selected = $("#generate-source option:selected").val()
	if(selected == "new-external-address") {
		$("#sec-pub").text("Public")
	} else {
		$("#sec-pub").text("Private")
	}
	if(selected == "import-address" || selected == "new-external-address"){
		$('#private-key-input').prop("disabled", false);

		$('#private-key-input').removeClass("disabled-input");

		$('#nickname-input').addClass("input-group-error");
		$("#private-key-input-container").addClass("input-group-error");
	} else {
		$('#private-key-input').prop("disabled", true);

		$('#private-key-input').addClass("disabled-input")

		$('#nickname-input').addClass("input-group-error");
		$("#private-key-input-container").removeClass("input-group-error");
	}
})

$("#private-key-input-container").on('click', function(){
	$(this).removeClass("input-group-error")
})

$("#nickname-input").on('click', function(){
	$(this).removeClass("input-group-error")
})

$("#generate-source").on('change', function(){
	selected = $("#generate-source option:selected").val()
	if(selected == "import-address"){
		$("#private-key-input").attr("placeholder","Type address private key")
	} else if(selected == "random-ec"){
		$("#private-key-input").attr("placeholder","A new entry credit address will be created")
	} else if(selected == "random-factoid"){
		$("#private-key-input").attr("placeholder","A new factoid address will be created")
	} else if(selected == "new-external-address"){
		$("#private-key-input").attr("placeholder","Type a public address to add to your contacts")
	}
})

$("#add-to-addressbook").on("click", function(){
	$("#error-zone").slideUp(100)
	Name = $("#nickname-input").val()
	if(Name == ""){
		SetError("Need a NickName for the new address")
		$('#nickname-input').addClass("input-group-error");
		return
	}
	selected = $("#generate-source option:selected").val()
	if(selected == "import-address"){
		sec = $("#private-key-input").val()
		if(!(sec.startsWith("Fs") || sec.startsWith("Es"))){
			SetError("Not a valid private key. Should start with 'Fs' for a factoid address or 'Es' for an entry credit address")
			$("#private-key-input-container").addClass("input-group-error");
			return
		}

		postRequest("is-valid-address", sec, function(resp){
			if(resp == "false") { // Not valid
				SetError("Not a valid private key.")
				$("#private-key-input-container").addClass("input-group-error");
				return
			} else { // Is valid, generate off new private key
				var newAddressObj = {
			    	Name:Name,
					Secret:sec
				}

				j = JSON.stringify(newAddressObj)
				postRequest("new-address", j, function(resp){
					obj = JSON.parse(resp)
					if(obj.Error == "none"){
						SetSuccess(obj)
					} else {
						SetError(obj.Error)
					}
				})
			}
		})
	} else if(selected == "random-ec"){
		postRequest("generate-new-address-ec", Name, function(resp){ // Generate new key
			obj = JSON.parse(resp)
			if(obj.Error == "none"){
				SetSuccess(obj)
			} else {
				SetError(obj.Error)
			}
		})
	} else if(selected == "random-factoid"){
		postRequest("generate-new-address-factoid", Name, function(resp){ // Generate new key
			obj = JSON.parse(resp)
			if(obj.Error == "none"){
				SetSuccess(obj)
			} else {
				SetError(obj.Error)
			}
		})
	} else if(selected == "new-external-address"){
		pub = $("#private-key-input").val()
		if(!(pub.startsWith("FA") || pub.startsWith("EC"))){
			SetError("Not a valid public key. Should start with 'FA' for a factoid address or 'EC' for an entry credit address." +
				" This option adds an address to your external addresses for easier use.")
			$("#private-key-input-container").addClass("input-group-error");
			return
		}

		postRequest("is-valid-address", pub, function(resp){
			if(resp == "false") { // Not valid
				SetError("Not a valid public key.")
				$("#private-key-input-container").addClass("input-group-error");
				return
			} else { // Is valid, generate off new private key
				var newAddressObj = {
			    	Name:Name,
					Public:pub
				}

				j = JSON.stringify(newAddressObj)
				postRequest("new-external-address", j, function(resp){
					obj = JSON.parse(resp)
					if(obj.Error == "none"){
						SetSuccess(obj)
					} else {
						SetError(obj.Error)
					}
				})
			}
		})
	} else {
		SetError("An error has occurred. No address type selected, please try selecting from the dropdown menu again, or reload this page.")
	}
})

function SetError(err) {
	$("#success-zone").slideUp(100)
	$("#error-zone").text(err)
	$("#error-zone").slideDown(100)
}

function SetSuccess(obj) {
	$("#error-zone").slideUp(100)
	$("#success-link").attr("href", "/receive-factoids?address=" + obj.Content.Address + "&name=" + obj.Content.Name)
	$("#success-zone > #name").text(obj.Content.Name)
	$("#success-zone > pre").text(obj.Content.Address)


	$("#success-zone").slideDown(100)
}