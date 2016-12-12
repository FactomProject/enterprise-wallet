$(window).load(function() {
    LoadAddresses()
});

function LoadAddresses(){
	resp = getRequest("addresses",function(resp){
		obj = JSON.parse(resp)
		//console.log(resp)
		
		if(obj.FactoidAddresses.List != null) {
			obj.FactoidAddresses.List.forEach(function(address){
				$('#factoid-addresses-table tbody').append(addressTableRow(address, "factoid"));
			})
		}
		if(obj.EntryCreditAddresses.List != null) {
			obj.EntryCreditAddresses.List.forEach(function(address){
				$('#credit-addresses-table tbody').append(addressTableRow(address, "entry-credits"));
			})
		}
		if(obj.ExternalAddresses.List != null) {
			obj.ExternalAddresses.List.forEach(function(address){
				$('#external-addresses-table tbody').append(addressTableRow(address, "external"));
			})
		}
 	})
}

function addressTableRow(address, type) {
	console.log(address)
	if(address.Address.startsWith("FA")){
		token = " FCT"
		address.Balance = Number(address.Balance.toFixed(4))
	} else {
		token = " EC"
	}

	star = '<small><span class="fa fa-star-o" aria-hidden="true"></span></small>'
	if(address.Seeded) {
		star = '<small><span class="fa fa-star" aria-hidden="true"></span></small>'
	}

	shortAddr = address.Address
	// Potential to shorten address
	//shortAddr = shortAddr.substring(0, 52)

	return		'<tr>' +
				'<td><a href="receive-factoids?address=' + address.Address + '&name=' + address.Name + '"><i class="qr"><img src="img/icon_qr.svg" class="svg"></i></a></td>' +
				'<td>' + address.Name + ' <a href="edit-address-' + type + '?address=' + address.Address + '&name=' + address.Name + '"><i class="edit"><img src="img/icon_edit.svg" class="svg"></i></a></td>' +
				'<td><pre>' + star + " " + shortAddr + '</pre></td>' +
				'<td>' + address.Balance + token + '</td>' +
				'</tr>'
}