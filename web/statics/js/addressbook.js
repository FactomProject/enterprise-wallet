

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
	} else {
		token = " EC"
	}

	star = '<small><span class="fa fa-star-o" aria-hidden="true"></span></small>'
	if(address.Seeded) {
		star = '<small><span class="fa fa-star" aria-hidden="true"></span></small>'
	}
	return		'<tr>' +
				'<td><a href="receive-factoids?address=' + address.Address + '&name=' + address.Name + '"><i class="qr"><img src="img/icon_qr.svg" class="svg"></i></a></td>' +
				'<td>' + address.Name + ' <a href="edit-address-' + type + '?address=' + address.Address + '&name=' + address.Name + '"><i class="edit"><img src="img/icon_edit.svg" class="svg"></i></a></td>' +
				'<td><pre>' + star + " " + address.Address + '</pre></td>' +
				'<td>' + address.Balance + token + '</td>' +
				'</tr>'
}
//<tr>
//  <td><a href="receive-factoids"><i class="qr"><img src="img/icon_qr.svg" class="svg"></i></a></td>
//  <td>factoid1 <a href="edit-address-factoid"><i class="edit"><img src="img/icon_edit.svg" class="svg"></i></a></td>
//  <td>FA3cih2o2tjEUsnnFR4jX1tQXPpSXFwsp3rhVp6odL5PNCHWvZV1</td>
//	<td>240.82 FCT</td>
//</tr>