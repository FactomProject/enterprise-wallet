

$(window).load(function() {
    LoadAddresses()
});

function LoadAddresses(){
	resp = getRequest("addresses",function(resp){
		obj = JSON.parse(resp)
		console.log(obj)

		obj.FactoidAddresses.List.forEach(function(address){
			$('#factoid-addresses-table tbody').append(
				'<tr>' +
				'<td><a href="receive-factoids"><i class="qr"><img src="img/icon_qr.svg" class="svg"></i></a></td>' +
				'<td>' + address.Name + ' <a href="edit-address-factoid"><i class="edit"><img src="img/icon_edit.svg" class="svg"></i></a></td>' +
				'<td>' + address.Address + '</td>' +
				'<td>' + address.Balance + '</td>' +
				'</tr>');
		})

		obj.EntryCreditAddresses.List.forEach(function(address){
			console.log(address)
			$('#credit-addresses-table tbody').append(
				'<tr>' +
				'<td><a href="receive-factoids"><i class="qr"><img src="img/icon_qr.svg" class="svg"></i></a></td>' +
				'<td>' + address.Name + ' <a href="edit-address-factoid"><i class="edit"><img src="img/icon_edit.svg" class="svg"></i></a></td>' +
				'<td>' + address.Address + '</td>' +
				'<td>' + address.Balance + '</td>' +
				'</tr>');
		})

		obj.ExternalAddresses.List.forEach(function(address){
			$('#external-addresses-table tbody').append(
				'<tr>' +
				'<td><a href="receive-factoids"><i class="qr"><img src="img/icon_qr.svg" class="svg"></i></a></td>' +
				'<td>' + address.Name + ' <a href="edit-address-factoid"><i class="edit"><img src="img/icon_edit.svg" class="svg"></i></a></td>' +
				'<td>' + address.Address + '</td>' +
				'<td>' + address.Balance + '</td>' +
				'</tr>');
		})
 	})
}
//<tr>
//  <td><a href="receive-factoids"><i class="qr"><img src="img/icon_qr.svg" class="svg"></i></a></td>
//  <td>factoid1 <a href="edit-address-factoid"><i class="edit"><img src="img/icon_edit.svg" class="svg"></i></a></td>
//  <td>FA3cih2o2tjEUsnnFR4jX1tQXPpSXFwsp3rhVp6odL5PNCHWvZV1</td>
//	<td>240.82 FCT</td>
//</tr>