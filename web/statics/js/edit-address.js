var Name = ""
var Address = ""

$(window).load(function() {
    GetDefaultData()
});

function GetDefaultData(){
	Name = $("#address-name").val()
	Address = $("#address-field").val()
	console.log(Name, Address)
}

$("#display-private-key").click(function(){
	jsonOBJ = '{"Address":"' + Address + '"}'
	postRequest("display-private-key", jsonOBJ, function(resp){
				console.log(resp)
		obj = JSON.parse(resp)
		if (obj.Error != "none") {
			$("#private-key-field").val(obj.Error)
			SetGeneralError("Error: " + obj.Error)
		} else {
			$("#private-key-field").val(obj.Content)
		}
	})
})

$("#save-name-change").click(function(){
	NewName = $("#address-name").val()
	jsonOBJ = '{"Address":"' + Address + '", "Name":"' + NewName + '"}'

	if (NewName != Name) {
		postRequest("address-name-change", jsonOBJ, function(resp){
			obj = JSON.parse(resp)
			if (obj.Error != "none") {
				SetGeneralError("Error: " + obj.Error)
			} else {
				SetGeneralSuccess(obj.Content + ": The name has been changed")
			}
		})
	} else {
		SetGeneralError("Newname is the same as the original")
	}
})

$("#delete-address").on('click', function(){
	name = $("#address-name").val()
	jsonOBJ = '{"Address":"' + Address + '", "Name":"' + name + '"}'
	postRequest("delete-address", jsonOBJ, function(resp){
		obj = JSON.parse(resp)
		if (obj.Error != "none") {
			SetGeneralError("Error: " + obj.Error)
		} else {
			SetGeneralSuccess(obj.Content + ": The name has been changed")
		}
	})
})

$("#copy-to-clipboard").on('click', function(){
	var aux = document.createElement("input");
	//console.log($('#selected-address-info').val())
	aux.setAttribute("value", $('#private-key-field').val());
	document.body.appendChild(aux);
	aux.select();
	document.execCommand("copy");
	document.body.removeChild(aux);
})
