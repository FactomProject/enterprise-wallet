var Name = ""
var Address = ""

$(window).load(function() {
    GetDefaultData()
});

function GetDefaultData(){
	Name = $("#address-name").val()
	Address = $("#address-field").val()
}

$("#display-private-key").click(function(){
	console.log("ASD")
	jsonOBJ = '{"Address":"' + Address + '"}'
	postRequest("display-private-key", jsonOBJ, function(resp){
		obj = JSON.parse(resp)
		if (obj.Error != "none") {
			$("#private-key-field").val(obj.Error)
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
			console.log("NameChangeResponse:" + resp)
		})
	}
})

