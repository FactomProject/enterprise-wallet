/*$(window).load(function() {
    LoadAddresses()
});*/

function LoadInitialAddresses(){
	resp = getRequest("addresses-no-bal",function(resp){
		obj = JSON.parse(resp)
		
		if(obj.FactoidAddresses.List != null) {
			obj.FactoidAddresses.List.forEach(function(address){
				$('#factoid-addresses-table tbody').append(addressTableRow(address, "factoid", true));
			})
		}
		if(obj.EntryCreditAddresses.List != null) {
			obj.EntryCreditAddresses.List.forEach(function(address){
				$('#credit-addresses-table tbody').append(addressTableRow(address, "entry-credits", true));
			})
		}
		if(obj.ExternalAddresses.List != null) {
			obj.ExternalAddresses.List.forEach(function(address){
				$('#external-addresses-table tbody').append(addressTableRow(address, "external", true));
			})
		}
		sortNames(true)
 	})
}

function LoadAddresses(){
	LoadInitialAddresses()
	resp = getRequest("addresses",function(resp){
		obj = JSON.parse(resp)
		//console.log(resp)
		
		if(obj.FactoidAddresses.List != null) {
			$('#factoid-addresses-table tbody').html("")
			obj.FactoidAddresses.List.forEach(function(address){
				$('#factoid-addresses-table tbody').append(addressTableRow(address, "factoid", false));
			})
		}
		if(obj.EntryCreditAddresses.List != null) {
			$('#credit-addresses-table tbody').html("")
			obj.EntryCreditAddresses.List.forEach(function(address){
				$('#credit-addresses-table tbody').append(addressTableRow(address, "entry-credits", false));
			})
		}
		if(obj.ExternalAddresses.List != null) {
			$('#external-addresses-table tbody').html("")
			obj.ExternalAddresses.List.forEach(function(address){
				$('#external-addresses-table tbody').append(addressTableRow(address, "external", false));
			})
		}
		sortNames(true)
 	})
}

function addressTableRow(address, type, loading) {
	if(address.Address.startsWith("FA")){
		token = " FCT"
		address.Balance = Number(address.Balance.toFixed(4))
	} else {
		token = " EC"
	}

	if(loading) {
		address.Balance = "..."
		token = ""
	}

	star = '<small><span id="star" class="fa fa-star-o" aria-hidden="true" value="0"></span></small>'
	if(address.Seeded) {
		star = '<small><span id="star" class="fa fa-star" aria-hidden="true" value="1"></span></small>'
	}

	if(type == "external") {
		star = ""
	}

	shortAddr = address.Address
	// Potential to shorten address
	//shortAddr = shortAddr.substring(0, 52)

	return		'<tr>' +
				'<td><a nav-click="true" href="receive-factoids?address=' + address.Address + '&name=' + address.Name + '"><i class="qr"><img src="img/icon_qr.svg" class="svg"></i></a></td>' +
				'<td><span id="name">' + address.Name + '</span> <a nav-click="true" href="edit-address-' + type + '?address=' + address.Address + '&name=' + address.Name + '"><i class="edit"><img src="img/icon_edit.svg" class="svg"></i></a></td>' +
				'<td><pre>' + star + " " + shortAddr + '</pre></td>' +
				'<td><span id="balance">' + address.Balance  + "</span>" + token + '</td>' +
				'</tr>'
}


function ResetNotMe(me) {
	if(me != 0) {
		resetNameSort()
	}
	if(me != 1) {
		resetBalancesSort()
	}
	if(me != 2) {
		resetSeededSort()
	}
}

// Sort Names
$("table").on('click', "#sort-names", function(e){	
	ResetNotMe(0)
	if($("#sort-names-icon").hasClass("fa-sort")) { // norm to desc
		$("th #sort-names-icon").removeClass("fa-sort")
		$("th #sort-names-icon").addClass("fa-sort-desc")
		sortNames(true) 

	} else if($("#sort-names-icon").hasClass("fa-sort-asc")){ // asc to desc
		$("th #sort-names-icon").removeClass("fa-sort-asc")
		$("th #sort-names-icon").addClass("fa-sort-desc")
		sortNames(true)

	} else if($("#sort-names-icon").hasClass("fa-sort-desc")){ // desc to asc
		$("th #sort-names-icon").removeClass("fa-sort-desc")
		$("th #sort-names-icon").addClass("fa-sort-asc")
		sortNames(false)
	}
})

function resetNameSort() {
	$("th #sort-names-icon").removeClass("fa-sort-asc")
	$("th #sort-names-icon").removeClass("fa-sort-desc")

	$("th #sort-names-icon").addClass("fa-sort")
}
function sortNames(order) {
	array = $("#factoid-addresses-table tbody tr").get()
	valArray = $("#factoid-addresses-table tbody tr").find("#name").get()
	array = generalSort(stringLessThan, array, valArray, order, order)
	$("#factoid-addresses-table tbody").html(array)

	array = $("#credit-addresses-table tbody tr").get()
	valArray = $("#credit-addresses-table tbody tr").find("#name").get()
	array = generalSort(stringLessThan, array, valArray, order, order)
	$("#credit-addresses-table tbody").html(array)

	array = $("#external-addresses-table tbody tr").get()
	valArray = $("#external-addresses-table tbody tr").find("#name").get()
	array = generalSort(stringLessThan, array, valArray, order, order)
	$("#external-addresses-table tbody").html(array)
}

// Sort Balances
$("table").on('click', "#sort-balances", function(e){
	ResetNotMe(1)
	if($("#sort-balances-icon").hasClass("fa-sort-amount")) { // norm to desc
		$("th #sort-balances-icon").removeClass("fa-sort-amount")
		$("th #sort-balances-icon").addClass("fa-sort-amount-desc")
		sortBalances(false) 

	} else if($("#sort-balances-icon").hasClass("fa-sort-amount-asc")){ // asc to desc
		$("th #sort-balances-icon").removeClass("fa-sort-amount-asc")
		$("th #sort-balances-icon").addClass("fa-sort-amount-desc")
		sortBalances(false)

	} else if($("#sort-balances-icon").hasClass("fa-sort-amount-desc")){ // desc to asc
		$("th #sort-balances-icon").removeClass("fa-sort-amount-desc")
		$("th #sort-balances-icon").addClass("fa-sort-amount-asc")
		sortBalances(true)
	}
})

function resetBalancesSort() {
	$("th #sort-balances-icon").removeClass("fa-sort-amount-asc")
	$("th #sort-balances-icon").removeClass("fa-sort-amount-desc")

	$("th #sort-balances-icon").addClass("fa-sort-amount")
}

function sortBalances(order) {
	array = $("#factoid-addresses-table tbody tr").get()
	valArray = $("#factoid-addresses-table tbody tr").find("#balance").get()
	array = generalSort(isLessThan, array, valArray, order, order)
	$("#factoid-addresses-table tbody").html(array)

	array = $("#credit-addresses-table tbody tr").get()
	valArray = $("#credit-addresses-table tbody tr").find("#balance").get()
	array = generalSort(isLessThan, array, valArray, order, order)
	$("#credit-addresses-table tbody").html(array)

	array = $("#external-addresses-table tbody tr").get()
	valArray = $("#external-addresses-table tbody tr").find("#balance").get()
	array = generalSort(isLessThan, array, valArray, order, order)
	$("#external-addresses-table tbody").html(array)
}

// Sort Seeded
$("table").on('click', "#sort-seeded", function(e){
	ResetNotMe(2)
	if($("#sort-seeded-icon").hasClass("fa-sort")) { // norm to desc
		$("th #sort-seeded-icon").removeClass("fa-sort")
		$("th #sort-seeded-icon").addClass("fa-sort-desc")
		sortSeeded(false) 

	} else if($("#sort-seeded-icon").hasClass("fa-sort-asc")){ // asc to desc
		$("th #sort-seeded-icon").removeClass("fa-sort-asc")
		$("th #sort-seeded-icon").addClass("fa-sort-desc")
		sortSeeded(false)

	} else if($("#sort-seeded-icon").hasClass("fa-sort-desc")){ // desc to asc
		$("th #sort-seeded-icon").removeClass("fa-sort-desc")
		$("th #sort-seeded-icon").addClass("fa-sort-asc")
		sortSeeded(true)
	}
})

function resetSeededSort() {
	$("th #sort-seeded-icon").removeClass("fa-sort-asc")
	$("th #sort-seeded-icon").removeClass("fa-sort-desc")

	$("th #sort-seeded-icon").addClass("fa-sort")
}

function sortSeeded(order) {
	array = $("#factoid-addresses-table tbody tr").get()
	valArray = $("#factoid-addresses-table tbody tr").get()
	array = generalSort(starLessThan, array, valArray, order, true)
	$("#factoid-addresses-table tbody").html(array)

	array = $("#credit-addresses-table tbody tr").get()
	valArray = $("#credit-addresses-table tbody tr").get()
	array = generalSort(starLessThan, array, valArray, order, true)
	$("#credit-addresses-table tbody").html(array)
}

// Valuearray is the actual values to be compares
// Array is the corrosponding array of original data.
// order is 'true' for descending, 'false' for ascending
// 		returned: sorted array
function generalSort(lessThanFunction, array, valueArray, order, keep) {
  peerLen = valueArray.length
  for(index = 0; index < peerLen; index++) {
    tmpVal = valueArray[index]
    tmp = array[index]

    if(order == true || keep) {
      for(j = index - 1; j > -1 && !lessThanFunction(valueArray[j], tmpVal, order); j--) {
        valueArray[j+1] = valueArray[j]
        array[j+1] = array[j]
      }
    } else {
      for(j = index - 1; j > -1 && lessThanFunction(valueArray[j], tmpVal , order); j--) {
        valueArray[j+1] = valueArray[j]
        array[j+1] = array[j]
      }
    }

    valueArray[j+1] = tmpVal
    array[j+1] = tmp
  }
  return array
}

// 1 is true (a < b)
// 0 is false (a >= b)
function starLessThan(a, b, o) {
	as = $(a).find("#star")
	bs = $(b).find("#star")

	aN = $(as).attr("value")
	bN = $(bs).attr("value")

	if(!o) {
		aN = flipInt(aN)
		bN = flipInt(bN)
	}

	aM = aN + $(a).find("#name").text()
	bM = bN + $(b).find("#name").text()
	return (String(aM)<String(bM)?1:0) 
}


function flipInt(i) {return (Number(i)==1?0:1)}

function stringLessThan(a, b, o){
    return (String(a.innerText)<String(b.innerText)?1:0) 
}

function isLessThan(a, b, o) {
	return (Number(a.innerText)<Number(b.innerText)?1:0)
}
