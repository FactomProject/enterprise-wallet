

$(window).load(function() {
    LoadAddresses()
});

function LoadAddresses(){
	resp = getRequest("addresses",function(resp){
		console.log(resp)
		obj = JSON.parse(resp)
		
		obj.FactoidAddresses.List.forEach(function(address){
			$('#factoid-addresses-table tbody').append(addressTableRow(address, "factoid"));
		})

		obj.EntryCreditAddresses.List.forEach(function(address){
			console.log("ASD")
			$('#credit-addresses-table tbody').append(addressTableRow(address, "entry-credits"));
		})

		obj.ExternalAddresses.List.forEach(function(address){
			$('#external-addresses-table tbody').append(addressTableRow(address, "external"));
		})
 	})
}

function addressTableRow(address, type) {
	return		'<tr>' +
				'<td><a href="receive-factoids?address=' + address.Address + '&name=' + address.Name + '"><i class="qr"><img src="img/icon_qr.svg" class="svg"></i></a></td>' +
				'<td>' + address.Name + ' <a href="edit-address-' + type + '?address=' + address.Address + '&name=' + address.Name + '"><i class="edit"><img src="img/icon_edit.svg" class="svg"></i></a></td>' +
				'<td>' + address.Address + '</td>' +
				'<td>' + address.Balance + '</td>' +
				'</tr>'
}
//<tr>
//  <td><a href="receive-factoids"><i class="qr"><img src="img/icon_qr.svg" class="svg"></i></a></td>
//  <td>factoid1 <a href="edit-address-factoid"><i class="edit"><img src="img/icon_edit.svg" class="svg"></i></a></td>
//  <td>FA3cih2o2tjEUsnnFR4jX1tQXPpSXFwsp3rhVp6odL5PNCHWvZV1</td>
//	<td>240.82 FCT</td>
//</tr>