/*$(window).load(function() {
    LoadAddresses()
});*/

function LoadInitialAddresses(showSeeded){
	if(showSeeded === undefined) {
		showSeeded = true
	}
	resp = getRequest("addresses-no-bal",function(resp){
		var count = 0
		obj = JSON.parse(resp)
		
		if(obj.FactoidAddresses.List != null) {
			obj.FactoidAddresses.List.forEach(function(address){
				if(!address.Seeded ) { count++ }
				if(!address.Seeded || showSeeded) {
					$('#factoid-addresses-table tbody').append(addressTableRow(address, "factoid", true));
				}
			})
		}
		if(obj.EntryCreditAddresses.List != null) {
			obj.EntryCreditAddresses.List.forEach(function(address){
				if(!address.Seeded ) { count++ }
				if(!address.Seeded || showSeeded) {
					$('#credit-addresses-table tbody').append(addressTableRow(address, "entry-credits", true));
				}
			})
		}
		if(obj.ExternalAddresses.List != null) {
			obj.ExternalAddresses.List.forEach(function(address){
				$('#external-addresses-table tbody').append(addressTableRow(address, "external", true));
			})
		}
		sortNames(true)

		if(!showSeeded && count === 0) {
			// This means we only need to show if there are not seeded addresses
			$(".not-all-backed-up").hide()
			$(".all-backed-up").show()
		}
 	})
}

function LoadAddresses(showSeeded){
	if(showSeeded === undefined) {
		showSeeded = true
	}
	LoadInitialAddresses(showSeeded)
	resp = getRequest("addresses",function(resp){
		obj = JSON.parse(resp)
		//console.log(resp)
		
		if(obj.FactoidAddresses.List != null) {
			$('#factoid-addresses-table tbody').html("")
			obj.FactoidAddresses.List.forEach(function(address){
				if(!address.Seeded || showSeeded) {
					$('#factoid-addresses-table tbody').append(addressTableRow(address, "factoid", false));
				}
			})
		}
		if(obj.EntryCreditAddresses.List != null) {
			$('#credit-addresses-table tbody').html("")
			obj.EntryCreditAddresses.List.forEach(function(address){
				if(!address.Seeded || showSeeded) {
					$('#credit-addresses-table tbody').append(addressTableRow(address, "entry-credits", false));
				}
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
		address.Balance = ShrinkFixedPoint(FCTNormalize(address.Balance),4) //Number(address.Balance.toFixed(4))
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
/*$(window).load(function() {
    LoadTransactions()
});*/

CurrentCount = 0
ContentLen = 0
LoopstopIncrement = 15 // Amount to load on scroll
Loopstop = 20
var Transactions = new Array()
Done = false

function LoadTransactions() {
	getRequest("related-transactions", function(resp){
		obj = JSON.parse(resp)
		if(obj.Error == "none" && obj.Content == null){
			setTimeout(function(){
			  LoadTransactions()
			}, 5000);
			return
		}
		$("#loading-container").remove()

		if(obj.Error != "none"){
			SetGeneralError(obj.Error)
			return
		}

		ContentLen = obj.Content.length
		Transactions = obj.Content

		if(obj.Content.length > 0 && obj.Content[0].TxID == "empty") {
			SetGeneralError("No transactions found for your addresses.")
			return
		} else if(obj.Content.length == 0) {
			return
		}

		// Load past x transactions, then stop. Only load more if they scroll
		if(ContentLen < Loopstop) {
			Loopstop = ContentLen
		}
		for(; CurrentCount < Loopstop; CurrentCount++) {
			AppendNewTransaction(Transactions[CurrentCount], CurrentCount)
		}
	})
}

// Load past x transactions, then stop. Only load more if they scroll
Empty = false
function LoadCached() {
	if (Empty){return}
	if(ContentLen < Loopstop * 2) {
		// Request more
		var requestObject = {
			Current:ContentLen,
		    More:Loopstop*5
		}
		j = JSON.stringify(requestObject)
		postRequest("more-cached-transaction", j, function(resp){
			obj = JSON.parse(resp)
			if(obj.Error != "none"){
				return
			}
			if(obj.Content == null){
				return
			}
			if(obj.Content.length == 0) {
				Empty = true
				return
			}
			ContentLen = ContentLen + obj.Content.length
			Transactions = Transactions.concat(obj.Content)
		})
	}
	if(ContentLen < Loopstop) {
		Loopstop = ContentLen
	}
	for(; CurrentCount < Loopstop; CurrentCount++) {
		AppendNewTransaction(Transactions[CurrentCount], CurrentCount)
	}
}

function AppendNewTransaction(trans, index){
	// Transactions are split into 3 transactions if sent/recieve/converted is all happening.
	// function appendTrans(pic, index, amt, token, date, addrs)
	if(trans.Action[0] == true) { // Sent
		pic = "sent"
		amt = 0
		token = "FCT"
		addrs = ""

		for(var i = 0; i < trans.Inputs.length; i++) {
			if(trans.Inputs[i].Name != "") {
				addrs += '<div class="nick">' + trans.Inputs[i].Name + '<pre class="show-for-large"> (' + trans.Inputs[i].Address + ')</pre></div>'
				amt += trans.Inputs[i].Amount / 1e8
			}
		}	

		appendTrans(pic, index, amt*-1, token, trans.Date, addrs)
	}

	if(trans.Action[1] == true) { // Received
		pic = "received"
		amt = 0
		token = "FCT"
		addrs = ""

		for(var i = 0; i < trans.Outputs.length; i++) {
			if(trans.Outputs[i].Name != "") {
				if(trans.Outputs[i].Address.startsWith("FA")){
					addrs += '<div class="nick">' + trans.Outputs[i].Name + '<pre class="show-for-large percent95"> (' + trans.Outputs[i].Address + ')</pre></div>'
					amt += trans.Outputs[i].Amount / 1e8
				}
			}
		}

		appendTrans(pic, index, amt, token, trans.Date, addrs)
	}

	if(trans.Action[2] == true) { // Converted
		pic = "converted"
		amt = 0
		token = "FCT"
		addrs = ""

		for(var i = 0; i < trans.Outputs.length; i++) {
			if(trans.Outputs[i].Name != "") {
				if(trans.Outputs[i].Address.startsWith("EC")){
					addrs += '<div class="nick">' + trans.Outputs[i].Name + '<pre class="show-for-large percent95"> (' + trans.Outputs[i].Address + ')</pre></div>'
					amt += trans.Outputs[i].Amount / 1e8
				}
			}
		}

		appendTrans(pic, index, amt, token, trans.Date, addrs)
	}
}

function appendTrans(pic, index, amt, token, date, addrs) {
	$("#transaction-list").append(
   '<tr>' +
        '<td><a id="transaction-link" data-toggle="transDetails" value="' + index + '"><i class="transIcon ' + pic + '"><img src="img/transaction_' + pic + '.svg" class="svg"></i></a></td>' +
        '<td>' + date + ' : <a value="' + index + '" id="transaction-link" data-toggle="transDetails">' + pic.capitalize() + '</a>' +
        addrs + '</td>' +
        '<td style="word-wrap: break-word;">' + ShrinkFixedPoint(amt,4) + ' ' + token + '</td>' +
    '</tr>'
)
}

$("main").bind('scroll', function() {
	//console.log($("main").outerHeight(), $("main").scrollTop(), $("main").innerHeight(), $("main").prop('scrollHeight'), $("main").prop('offsetHeight'))
	// Total Height
	// $("main").prop('scrollHeight')
	//console.log("scroll", $("#transaction-list").scrollTop(), $("html").innerHeight(), $("html").prop('scrollHeight'))
	//console.log("scroll", $("main").scrollTop(), $("main").innerHeight(), $("main").prop('scrollHeight'))
	//console.log("scroll", $("body").scrollTop(), $("body").innerHeight(), $("body").prop('scrollHeight'))
	if($("main").scrollTop() + $("main").innerHeight() >= .8 * $("main").prop('scrollHeight')) {
		Loopstop += LoopstopIncrement
		LoadCached()
	}
});

LOCAL_EXPLORER = true
$("#transaction-list").on('click', '#transaction-link', function(){
	port = $("#controlpanel-port").text()
	factomd = $("#factomd-location").text()
	setTransDetails($(this).attr("value"))
	$("#transDetails #link").attr("href", "http://explorer.factom.org/tx/" + Transactions[$(this).attr("value")].TxID)
	
	if(!(factomd.includes("localhost") || factomd.includes("127.0.0.1"))) {
		LOCAL_EXPLORER = false
		$("#transDetails #local-link").addClass("disabled-input")
		$("#transDetails #local-link").prop("disabled", true)
		$("#transDetails #local-link").attr("href", "")
	} else {
		$("#transDetails #local-link").attr("href", "http://localhost:" + port + "/search?input=" + Transactions[$(this).attr("value")].TxID + "&type=facttransaction")
	}
})

$('#transDetails #local-link').click(function(e) {
	if (!LOCAL_EXPLORER) {
		e.preventDefault();
	}
});

function setTransDetails(index) {
	trans = Transactions[index]
	$("#trans-detail-txid").text(trans.TxID)

	$("#trans-details-inputs").html("")
	for(var i = 0; i < trans.Inputs.length; i++) {
		$("#trans-details-inputs").append('<tr>' +
			'<td>' + trans.Inputs[i].Address + '</td>' +
			'<td>' + FCTNormalize(trans.Inputs[i].Amount) + ' FCT</td>' +
			'</tr>')
	}

	$("#trans-details-outputs").html("")
	for(var i = 0; i < trans.Outputs.length; i++) {
		$("#trans-details-outputs").append('<tr>' +
			'<td>' + trans.Outputs[i].Address + '</td>' +
			'<td>' + FCTNormalize(trans.Outputs[i].Amount) + ' FCT</td>' +
			'</tr>')
	}

	$("#trans-details-outputs").append('<tr class="total">' +
	'<td> Total </td>' +
	'<td>' + FCTNormalize(trans.TotalInput) + ' FCT</td>' +
	'</tr>')

	$("#total-transacted").text(FCTNormalize(trans.TotalECOutput + trans.TotalFCTOutput))
	$("#trans-date").text(trans.Date + " at " + trans.Time)
}

String.prototype.capitalize = function() {
    return this.charAt(0).toUpperCase() + this.slice(1);
}
/*
        <tr>
            <td><a data-toggle="transDetails"><i class="transIcon sent"><img src="img/transaction_sent.svg" class="svg"></i></a></td>
            <td>1/10/2016 : <a data-toggle="transDetails">Sent</a> <div class="nick">factoid1<span class="show-for-large"> (FA3cih2o2tjEUsnnFR4jX1tQXPpSXFwsp3rhVp6odL5PNCHWvZV1)</span></div></td>
            <td>-510.04 FCT</td>
        </tr>
*/
// Import/Export page acts differently
importexport = false

// Load the Reveal
/*$(window).load(function() {
    LoadAddresses()
    if($("#coin-control").hasClass("coin-control")) {
      $("#fee-address-input").css("display", "none")
    }
});*/

// Used for sending factoids or converting to entry credits
PageTokenABR = "FCT"
PageToken = "factoids"
AddressPrefix = "FA"
PageTransType = "factoid"

if($("#token-header").attr("value") == "1") {
  PageTokenABR = "EC"
  PageToken = "entry credits"
  AddressPrefix = "EC"
  PageTransType = "ec"
} else if ($("#token-header").attr("value") == "2"){
  importexport = true
}

counter = 2
function addNewOutputAddress(defaultAdd, defaultAmt, error, first) {
  eClass = ""
  if(error){
    eClass = "input-group-error"
  }
  
  str = "factoid"
  if(PageTokenABR != "FCT") {
    str = "entry credit"
  }

  defAdd = '<pre><input id="output-factoid-address" type="text" name="output1" class="input-group-field percent95" placeholder="Type ' + str + ' address"></pre>'
  if(defaultAdd != "") {
    defAdd = '<pre><input id="output-factoid-address" type="text" name="output1" class="input-group-field percent95" placeholder="Type ' + str + ' address" value="' + defaultAdd + '"></pre>'
  }


  defAmt = '<input id="output-factoid-amount" type="text" class="input-group-field" name="output1-num" placeholder="Amount">'
  if(defaultAmt != 0) {
     defAmt = '<input id="output-factoid-amount" type="text" class="input-group-field" name="output1-num" placeholder="Amount" value="' + defaultAmt + '">'
  }

  button = '<a id="remove-new-output" class="button expanded newMinus">&nbsp;</a>'
  if(first) {
    button = '<a id="append-new-output" class="button expanded newPlus">&nbsp;</a>'
  }


  $("#all-outputs").append(
  '<div class="row single-output-' + counter + '" id="single-output">' +
  '    <div class="small-12 medium-7 large-8 columns">' +
  '        <div class="input-group ' + eClass + '" id="output-factoid-address-container">' +
  '        ' + defAdd +
  '        <a id="addressbook-button" data-toggle="addressbook" class="input-group-button button input-group-field" id="addressbook" value="' + counter + '"><i class="fa fa-book"></i></a>' +
  '        </div>' +
  '    </div>' +
  '    <div class="small-10 medium-4 large-3 columns">' +
  '        <div class="input-group">' +
  '            ' + defAmt + 
  '            <span class="input-group-label">' + PageTokenABR + '</span>' +
  '        </div>' +
  '    </div>' +
  '    <div class="small-2 medium-1 columns">' +
  '            ' + button + 
  '    </div>' +
  '</div>')

  counter = counter + 1
}

// Add new output address
$("#append-new-output").click(function(){
  if($(this).hasClass("disabled-input")) {
    return
  }
  addNewOutputAddress("", 0, true, false)
})

function addNewInputAddress(defaultAdd, defaultAmt, error, first) {
  eClass = ""
  if(error){
    eClass = "input-group-error"
  }

  defAmt = '<input id="input-factoid-amount" type="text" class="input-group-field" name="output1-num" placeholder="Amount">'
  if(defaultAmt != 0) {
    defAmt = '<input id="input-factoid-amount" type="text" class="input-group-field" name="output1-num" placeholder="Amount" value="' + defaultAmt + '">'
  }


  defAdd = '<pre><input id="input-factoid-address" type="text" name="input1" class="input-group-field percent95 disabled-input" placeholder="Choose factoid address" disabled></pre>'
  if(defaultAdd != "") {
     defAdd = '<pre><input id="input-factoid-address" type="text" name="input1" class="input-group-field percent95 disabled-input" placeholder="Choose factoid address" disabled value="' + defaultAdd  + '"></pre>'
  }

  button = '<a id="remove-new-input" class="button expanded newMinus">&nbsp;</a>'
  if(first) {
    button = '<a id="append-new-input" class="button expanded newPlus">&nbsp;</a>'
  }

  $("#all-inputs").append(
  '<div class="row single-input-' + counter + '" id="single-input">' +
  '    <div class="small-12 medium-7 large-8 columns">' +
  '        <div class="input-group ' + eClass + '" id="input-factoid-address-container">' +
  '        ' + defAdd +
  '        <a id="addressbook-button" data-toggle="fee-addressbook" class="input-group-button button input-group-field" id="addressbook" value="' + counter + '"><i class="fa fa-book"></i></a>' +
  '        </div>' +
  '    </div>' +
  '    <div class="small-10 medium-4 large-3 columns">' +
  '        <div class="input-group">' +
  '            ' + defAmt +
  '            <span class="input-group-label">FCT</span>' +
  '        </div>' +
  '    </div>' +
  '    <div class="small-2 medium-1 columns">' +
  '            ' + button +
  '    </div>' +
  '</div>')
  counter = counter + 1
}

$("#append-new-input").click(function(){
  if($(this).hasClass("disabled-input")) {
    return
  }
  addNewInputAddress("", 0, true, false)
})

// Remove output address
$("#all-outputs").on('click','#remove-new-output', function(){
	jQuery(this).parent().parent().remove()
})

$("#all-inputs").on('click','#remove-new-input', function(){
  jQuery(this).parent().parent().remove()
})

// Ensure factoids/ec being sent are valid, this is not a security feature, but an ease of use
// feature
$("#all-outputs").on("keypress", "#output-factoid-amount", function(evt) {
  if((evt.which < 48 || evt.which > 57) && !(evt.Key == "." || evt.which == 46)) evt.preventDefault();
});


function mockTransaction() {
    if($("#make-entire-transaction").attr("value") == "1") {
        if($("#sign-transaction").prop('checked')) {
            MakeTransaction(true, true)
        } else {
            MakeTransaction(false, true)
        }
    } else {
        MakeTransaction(true, true)
    }
}


// Update Fee
$("#all-outputs").on('keyup', 'input[type=text]', function(){
    // Need to determine new fee
    mockTransaction()
})

$("#all-outputs").on('click', '#output-factoid-address-container', function(){
	$(this).removeClass("input-group-error")
})

$("#all-inputs").on('click', '#input-factoid-amount', function(){
  $(this).removeClass("input-group-error")
})

$("#all-outputs").on('click', '#output-factoid-amount', function(){
  $(this).removeClass("input-group-error")
})

$("#make-entire-transaction").on('click', function(){
  //$("#sign-transaction").prop('checked')
  if(Input) {
    if($(this).attr("value") == "1") {
      if($("#sign-transaction").prop('checked')) {
        MakeTransaction(true)
      } else {
        MakeTransaction(false)
      }
    } else {
      MakeTransaction(true)
    }
  }
})

function MakeTransaction(sig, checkonly = false) {
  transObject = getTransactionObject(true)

  if(transObject == null || transObject == undefined) {
    return
  }
  if(!sig) {
    transObject.TransType = "nosig"
  }

  if(transObject == null) {
    return
  }

  j = JSON.stringify(transObject)
  postRequest("make-transaction", j, function(resp){
    //console.log(resp)
    obj = JSON.parse(resp)
    //console.log(obj)
    if(obj.Error == "none") {
        if (!checkonly) {
            disableInput()
            ShowNewButtons()
        }
      totalInput = obj.Content.Total / 1e8
      feeFact = obj.Content.Fee / 1e8
      
      total = totalInput + feeFact 

      if(transObject.FeeAddress != "") {
        for(var i = 0; i < transObject.OutputAddresses.length; i++) {
          if(transObject.OutputAddresses[i] == transObject.FeeAddress) {
            total = totalInput - feeFact
            break
          }
        }           
      }

      $("#transaction-total").attr("value", total)
      $("#transaction-fee").attr("value", feeFact)
      if (checkonly) {
        HideMessages()
      } else {
        if(importexport) {
            setExportDownload(obj.Content.Json)
            SetGeneralSuccess('Click "Export Transaction" to download, or go back to editing it')
        } else {
            SetGeneralSuccess('Click "Send Transaction" to send, or go back to editing it')
        }
      }
    } else {
      SetGeneralError("Error: " + obj.Error)
    }
  })
}


function setExportDownload(json) {
  obj = JSON.parse(json)
  console.log(obj.params.transaction)
  fileExt = Date.now()
  $("#export-transaction").attr("value", obj.params.transaction)
  $("#export-transaction").attr("fileExt", fileExt)
  $("#export-transaction").click(function() {
    saveTextAsFile($(this).attr("value"), "Exported-" + $(this).attr("fileExt"))
    // $(this).attr("href", "data:text/plain;charset=UTF-8," + encodeURIComponent(obj.params.transaction))
    // $(this).attr("download", "Exported-" + fileExt)
  })
}

$("#send-entire-transaction").on('click', function(){
    disableInput()
  SendTransaction()
})

function getTransactionObject(checkInput) {
  var transObject = {
    TransType:PageTransType,
    OutputAddresses:[],
    OutputAmounts:[],

    InputAddresses:[],
    InputAmounts:[],
    FeeAddress:""
  }

  errMessage = ""
  faErr = false
  amtErr = false
  feeErr = false


  $("#all-outputs #single-output").each(function(){
    err = false
    add = $(this).find("#output-factoid-address").val()
    if(!add.startsWith(AddressPrefix) && !importexport) {
      $(this).find("#output-factoid-address-container").addClass("input-group-error")
      faErr = true
      err = true
    }

    amt = $(this).find("#output-factoid-amount").val()
    if(Number(amt) == 0 || amt == undefined || amt == "") {
      $(this).find("#output-factoid-amount").addClass("input-group-error")
      amtErr = true
      err = true
    }

    transObject.OutputAddresses.push(add)
    transObject.OutputAmounts.push(amt)
  })

  if(checkInput) {
    // Only FCT for inputs
    if(!$("#coin-control").hasClass("coin-control")) {
      $("#all-inputs #single-input").each(function(){
        add = $(this).find("#input-factoid-address").val()
        if(!add.startsWith("FA")) {
          $(this).find("#input-factoid-address-container").addClass("input-group-error")
          faErr = true
          err = true
        }

        amt = $(this).find("#input-factoid-amount").val()
        if(Number(amt) == 0 || amt == undefined || amt == "") {
          $(this).find("#input-factoid-amount").addClass("input-group-error")
          amtErr = true
          err = true
        }

        transObject.InputAddresses.push(add)
        transObject.InputAmounts.push(amt)
      })

      transObject.FeeAddress = $("#fee-factoid-address").val()
      if(transObject.FeeAddress.length <  52) {
        $("#fee-factoid-address").addClass("input-group-error")
        feeErr = true
        err = true
      }
      transObject.TransType = "custom"
    }
  }

  if(err){
    if(faErr){errMessage += "Addresses must start with '" + AddressPrefix + "' for output and 'FCT' for input. "}
    if(amtErr){errMessage += "Amounts should not be 0. "}
    if(feeErr){errMessage += "Fee Address must be given. "}
    SetGeneralError("Error(s): " + errMessage)
    return null
  }

  return transObject
}

$("#needed-input-button").on('click', function(evt){
  if(!Input) {
    evt.preventDefault()
  } else {
    GetNeededInput()
  }
})

CurrentInput = 0
TotalNeeded = 0
InputLeft = 0
function GetNeededInput() {
  transObject = getTransactionObject(false)

  if(transObject == null) {
    return
  }

  CurrentInput = 0
  // console.log(transObject)
  for(var i = 0; i < transObject.InputAmounts.length; i++) {
    if(transObject.InputAmounts[i] != undefined) {
      CurrentInput += Number(transObject.InputAmounts[i])
    }
  }

  j = JSON.stringify(transObject)
  postRequest("get-needed-input", j, function(resp){
    obj = JSON.parse(resp)
    if(obj.Error == "none") {
      if(obj.Content != null || obj.Content != undefined) {
        $("#input-needed-amount").val(FCTNormalize(obj.Content))
        console.log(CurrentInput)
        TotalNeeded = FCTNormalize(obj.Content)
        InputLeft = TotalNeeded - CurrentInput
      }
      HideMessages()
    } else {
      SetGeneralError("Error: " + obj.Error)
    }
  })
}

function SendTransaction() {
  transObject = getTransactionObject(true)

  if(transObject == null) {
    return
  }

  j = JSON.stringify(transObject)
  postRequest("send-transaction", j, function(resp){
    obj = JSON.parse(resp)
    if(obj.Error == "none") {
      disableInput()
      HideNewButtons()
      
      SetGeneralSuccess('Transaction Sent, transaction ID: ' + obj.Content )
      ShowNewTransaction()
    } else {
      enableInput()
      HideNewButtons()
      $("#transaction-fee").attr("value", "----")
      $("#transaction-total").attr("value", "----")
      SetGeneralError("Error: " + obj.Error)
    }
  })
}

$("#edit-transaction").on('click', function(){
  enableInput()
  HideNewButtons()
  $("#transaction-fee").attr("value", "----")
  $("#transaction-total").attr("value", "----")
  HideMessages()
})

function LoadAddressesSendConvert(){
  resp = getRequest("addresses-no-bal",function(resp){
    obj = JSON.parse(resp)

    if(obj.FactoidAddresses.List  != null) {
      obj.FactoidAddresses.List.forEach(function(address){
        $('#fee-addresses-reveal').append(factoidAddressRadio(address, "fee-address"));
      })
    }

    if(PageTokenABR == "FCT") {
      if(obj.FactoidAddresses.List  != null) {
        obj.FactoidAddresses.List.forEach(function(address){
          $('#addresses-reveal').append(factoidAddressRadio(address, "address"));
        })
      }

      if(obj.ExternalAddresses.List  != null) {
        obj.ExternalAddresses.List.forEach(function(address){
          if(address.Address.startsWith("FA")){
            $('#addresses-reveal').append(factoidAddressRadio(address, "address"));
          }
        })          
      }
    } else {
      if(obj.EntryCreditAddresses.List  != null) {
        obj.EntryCreditAddresses.List.forEach(function(address){
          $('#addresses-reveal').append(factoidECRadio(address, "address"));
        })
      }
      if(obj.ExternalAddresses.List  != null) {
        obj.ExternalAddresses.List.forEach(function(address){
          if(address.Address.startsWith("EC")){
            $('#addresses-reveal').append(factoidECRadio(address, "address"));
          }
        })
      }
    }
  })

  if($("#coin-control").hasClass("coin-control")) {
    $("#fee-address-input").css("display", "none")
  }
}

function factoidAddressRadio(address, name){
return '<pre>' +
  '  <input type="radio" name="' + name + '" id="address-'+address.Address+'" value="' + address.Address + 
  '"> <span id="address-name" data-address="'+address.Address+'" data-name="(' + FCTNormalize(address.Balance)  + " FCT) "  + address.Name + '">(' + FCTNormalize(address.Balance)  + " FCT) " + address.Name + '</span>' +
  '</pre><br />'
}

$('#addresses-reveal,#fee-addresses-reveal').on("mouseover", "#address-name", function(){
  $(this).css("font-size", "90%")
  $(this).text($(this).data("address"));
})

$('#addresses-reveal,#fee-addresses-reveal').on("mouseout", "#address-name", function(){
  $(this).text($(this).data("name"));
  $(this).css("font-size", "100%")
})

function factoidECRadio(address, type){
  return '<pre>' +
  '  <input type="radio" name="address" id="address" value="' + address.Address + '"> <span id="address-name" data-address="'+address.Address+'" data-name="' + address.Name + '">' + address.Name + '</span>' +
  '</pre> <br />'
}

done = false

$("#addresses-reveal-button").on("click", function(){
  newAddress = $("input[name='address']:checked").val()
  if(newAddress == undefined) {
    return
  }

  $(".single-output-" + toChange + " #output-factoid-address").val(newAddress)
  $(".single-output-" + toChange + " #output-factoid-address-container").removeClass("input-group-error")
  mockTransaction()
})

toChange = "-1"

$("#all-outputs").on('click', "#addressbook-button", function(){
  toChange = $(this).attr("value")
  $("input[type=radio]").attr('checked', false)
})

$("#all-inputs").on('click', "#addressbook-button", function(){
  toChange = $(this).attr("value")
  $("input[type=radio]").attr('checked', false)
})

$("#fee-address-input").on('click', "#addressbook-button", function(){
  toChange = $(this).attr("value")
  $("input[type=radio]").attr('checked', false)
})

$("#fee-addresses-reveal-button").on("click", function(){
  newAddress = $("input[name='fee-address']:checked").val()
  if(newAddress == undefined) {
    return
  }

  if(toChange == "-1") {
    $("#fee-factoid-address").val(newAddress)
    $("#fee-factoid-address").removeClass("input-group-error")
  } else {
    $(".single-input-" + toChange + " #input-factoid-address").val(newAddress)
    $(".single-input-" + toChange + " #input-factoid-address-container").removeClass("input-group-error")
  }
  mockTransaction()
})

function HideNewButtons() {
  $("#second-stage-buttons").slideUp(100)
  //$("#send-entire-transaction").slideUp(100)
}

function ShowNewButtons() {
  $("#export-transaction").slideDown(100)
  $("#broadcast-transaction").slideUp(100)

  $("#second-stage-buttons").slideDown(100)
  //$("#send-entire-transaction").slideDown(100)
}

function ShowNewTransaction() {
  $("#new-transaction").slideDown(100)
}

$("#new-transaction").on('click', function(){
  location.reload()
})

Input = true

function disableInput() {
  Input = false
  $(".input-group").each(function(){
    $(this).addClass("disabled-input")
    $(this).prop("disabled", true)
  })

  $(".input-group-field").each(function(){
    $(this).addClass("disabled-input")
    $(this).prop("disabled", true)
  })

  $("#needed-input-button").addClass("disabled-input")
  $("#needed-input-button").prop("disabled", true)
  $("#addressbook-button").addClass("disabled-input")
  $("#addressbook-button").prop("disabled", true)
  $("#make-entire-transaction").addClass("disabled-input")
  $("#make-entire-transaction").prop("disabled", true)
  $("#send-entire-transaction").addClass("disabled-input")
  $("#send-entire-transaction").prop("disabled", true)
  $("#first-stage-buttons").slideUp(100)
  $("#import-file").addClass("disabled-input")
  $("#import-file").prop("disabled", true)
}

function enableInput() {
  Input = true
  $(".input-group").each(function(){
    $(this).removeClass("disabled-input")
    $(this).prop("disabled", false)
  })

  $(".input-group-field").each(function(){
    $(this).removeClass("disabled-input")
    $(this).prop("disabled", false)
  })

  $("#transaction-fee").prop("disabled", true)
  $("#transaction-total").prop("disabled", true)

  $("#needed-input-button").removeClass("disabled-input")
  $("#needed-input-button").prop("disabled", false)

  $("#addressbook-button").removeClass("disabled-input")
  $("#addressbook-button").prop("disabled", false)
  $("#make-entire-transaction").removeClass("disabled-input")
  $("#make-entire-transaction").prop("disabled", false)
  $("#send-entire-transaction").removeClass("disabled-input")
  $("#send-entire-transaction").prop("disabled", false)
  $("#first-stage-buttons").slideDown(100)

  $("#import-file").removeClass("disabled-input")
  $("#import-file").prop("disabled", false)

  keepFeeDisabled()
}

// We need to keep these disabled.
function keepFeeDisabled() {
  //$("#fee-factoid-address").prop("disabled", true)
  //$("#fee-factoid-address").addClass("disabled-input")

  $("#input-factoid-address").prop("disabled", true)
  $("#input-factoid-address").addClass("disabled-input")
}

function HideMessages(){
  $("#error-zone").slideUp(100)
  $("#success-zone").slideUp(100)
}

// Import/Export
$("#import-file").on('click', function(){
  document.getElementById('uploaded-file').click()
})

$("#uploaded-file").on('change', function(){
  input = document.getElementById('uploaded-file');
  if (!input) {
    SetGeneralError("Error with upload file javascript.")
  }
  else if (!input.files) {
    SetGeneralError("This browser doesn't seem to support the `files` property of file inputs.")
  }
  else if (!input.files[0]) {
    SetGeneralError("No file found")             
  }
  else {
    file = input.files[0];
    fr = new FileReader();
    fr.onload = importTrans;
    fr.readAsText(file);
  }
})

// Do action with imported transaction
function importTrans() {
  x = fr.result
  postRequest("import-transaction", fr.result, function(resp){
    obj = JSON.parse(resp)
    if(obj.Error == "none") {
      total = 0
      $("#all-inputs").html("")
      for(var i = 0; i < obj.Content.InputAddresses.length; i++) {
        first = false
        if(i == 0) {
          first = true
        }
        addNewInputAddress(obj.Content.InputAddresses[i], Number(obj.Content.InputAmounts[i]), false, first)
        total += obj.Content.InputAmounts[i]
      }
      $("#all-outputs").html("")
      for(var i = 0; i < obj.Content.OutputAddresses.length; i++) {
        first = false
        if(i == 0) {
          first = true
        }
        addNewOutputAddress(obj.Content.OutputAddresses[i], Number(obj.Content.OutputAmounts[i]), false, first)
      }

      if(obj.Content.Signature) {
        $("#sign-transaction").attr('checked', true)
      } else {
        $("#sign-transaction").attr('checked', false)
      }

      disableInput()
      $("#transaction-total").val(Number(total))
      ShowNewButtons()
      $("#export-transaction").slideUp(1)
      $("#broadcast-transaction").slideDown(1)
    } else {
      SetGeneralError(obj.Error)
    }
  })
  console.log(fr.result)
}

$("#broadcast-transaction").on('click', function(){
  if(!$("#sign-transaction").prop('checked')){
    SetGeneralError("Transaction is not signed. Click the sign button if you contain the private keys to the inputs.")
    return
  }
  postRequest("broadcast-transaction", null, function(resp){
    obj = JSON.parse(resp)
    if(obj.Error == "none") {
      SetGeneralSuccess('Transaction Sent, transaction ID: ' + obj.Content )
      $("#broadcast-transaction").addClass("disabled-input")
      $("#broadcast-transaction").prop("disabled", true)
    } else {
      SetGeneralError(obj.Error)
    }
  })
})

/* http://stackoverflow.com/questions/12281775/get-data-from-file-input-in-jquery
<script>        
  function handleFileSelect()
  {               
    if (!window.File || !window.FileReader || !window.FileList || !window.Blob) {
      alert('The File APIs are not fully supported in this browser.');
      return;
    }   

    input = document.getElementById('fileinput');
    if (!input) {
      alert("Um, couldn't find the fileinput element.");
    }
    else if (!input.files) {
      alert("This browser doesn't seem to support the `files` property of file inputs.");
    }
    else if (!input.files[0]) {
      alert("Please select a file before clicking 'Load'");               
    }
    else {
      file = input.files[0];
      fr = new FileReader();
      fr.onload = receivedText;
      //fr.readAsText(file);
      fr.readAsDataURL(file);
    }
  }

  function receivedText() {
    document.getElementById('editor').appendChild(document.createTextNode(fr.result));
  }           

</script>


<input type="file" id="fileinput"/>
<input type='button' id='btnLoad' value='Load' onclick='handleFileSelect();'>
<div id="editor"></div>

*/
$("#save-changes").on('click', function(){	
	theme = $("#darkTheme").is(":checked")
	exportKeys = $("#export-keys").is(":checked")
	coinControl = $("#coin-control").is(":checked")
	importExport = $("#import-export").is(":checked")
	fd = $("#factomd-location").val()

	if(!$("#customFactomd").is(":checked")){
		fd = "localhost:8088"
	}

	var SettingsStruct = {
    	Values:[],
    	FactomdLocation:""
	}

	SettingsStruct.Values.push(theme)
	SettingsStruct.Values.push(exportKeys)
	SettingsStruct.Values.push(coinControl)
	SettingsStruct.Values.push(importExport)
	SettingsStruct.FactomdLocation = fd

	j = JSON.stringify(SettingsStruct)
	postRequest("adjust-settings", j, function(resp){
	    obj = JSON.parse(resp)
	    if(obj.Error == "none") {
	    	if((window.location.href).includes("success")){
	    		window.location.href = window.location.href
	    	} else {
	    		window.location.href = window.location.href + "?success=true"
	    	}
	      	//location.reload();
	    } else {
	    	SetGeneralError("Error: " + obj.Error)
	    }
	})
})

$("#customFactomd").on('click', function(){
	if($("#customFactomd").is(":checked")){
		$("#factomd-location-container").removeClass("hide")
	} else {
		$("#factomd-location-container").addClass("hide")
	}
})

$("#export-seed").on('click', function(){
	postRequest("get-seed", "", function(resp){
	    obj = JSON.parse(resp)
	    if(obj.Error == "none") {
	    	saveTextAsFile(obj.Content, "WalletSeed.txt")
	    } else {
	    	SetGeneralError("Error: " + obj.Error)
	    }
	})
})

//selected = false
// Import/Export
// $("#settings-import-file").on('click', function(e){
// 	document.getElementById('settings-uploaded-file').click()
// })


// $("#settings-uploaded-file").on('change', function(){
// 	input = document.getElementById('settings-uploaded-file');
// 	if (!input) {
// 		SetGeneralError("Error: Couldn't find the fileinput element.")
// 	}
// 	else if (!input.files) {
// 		SetGeneralError("This browser doesn't seem to support the `files` property of file inputs.")
// 	}
// 	else if (!input.files[0]) {
// 		SetGeneralError("Please select a file before clicking 'Import From File'")
// 	}
// 	else {
// 	file = input.files[0];
// 	fr = new FileReader();
// 	fr.onload = receivedText;
// 	fr.readAsText(file);
// 	//fr.readAsDataURL(file);
// 	}
// })

// Do action with imported transaction
function receivedText() {
	is = fr.result
	len = is.split(" ")
	if(len.length != 12) {
		SetGeneralError("Seed must be 12 words");
		return
	}
	document.getElementById('data-expand').click()
	$("#import-seed-reveal-text").text(fr.result)
	$("#import-seed-reveal-cancel").click()
}

$("#import-seed-reveal-confirm").on('click', function(){
	seed = $("#import-seed-reveal-text").text()
	var SeedStruct  = {
    	Seed:seed,
  	}
  	j = JSON.stringify(SeedStruct)
	postRequest("import-seed", j, function(resp) {
		obj = JSON.parse(resp)
		if(obj.Error == "none") {
	    	SetGeneralSuccess("Seed has been changed to: " + obj.Content)
	    } else {
	    	SetGeneralError("Error: " + obj.Error)
	    }
	})
})
// Load drop down if we were not directed with a specific link
function LoadRecAddresses(){
	resp = getRequest("addresses-no-bal",function(resp){
		obj = JSON.parse(resp)
		
		if(obj.FactoidAddresses.List != null){
			obj.FactoidAddresses.List.forEach(function(address){
				$('#receiving-address').append(dropDownOption(address));
			})
		}

		if(obj.EntryCreditAddresses.List != null){
			obj.EntryCreditAddresses.List.forEach(function(address){
				$('#receiving-address').append(dropDownOption(address));
			})
		}

		if(obj.ExternalAddresses.List != null) {
			obj.ExternalAddresses.List.forEach(function(address){
			//if(address.Address.startsWith("FA") == true) {
				$('#receiving-address').append(dropDownOption(address));
			//}
			})
		}

		UpdateSelectedInfo()
 	})
}

// If we only have 1 address, then we don't want to load dropdown
function LoadFixedAddress(){
	add = $("#receiving-address-fixed small").text()
	updateWithGivenAddress(add)
}

function dropDownOption(address) {
	return		'<option value="' + address.Name + '">' + address.Name + ' (' + address.Address + ')</option>'
}

// Copy to clipboard
$("#rec-copy-to-clipboard").on('click', function(){
	var aux = document.createElement("input");
	aux.setAttribute("value", $('#selected-address-info').val());
	document.body.appendChild(aux);
	aux.select();
	document.execCommand("copy");
	document.body.removeChild(aux);
})

// When new selected
$("#receiving-address").on('change', function(){
	UpdateSelectedInfo()
})

function UpdateSelectedInfo() {
	add = $("#receiving-address option:selected").text()
	updateWithGivenAddress(add)
}

function updateWithGivenAddress(address){
	// Get only the address
	splits = address.split("(")
	splits = splits[1].split(")")
	$("#selected-address-info").val(splits[0])
	$("#selected-address-info").text(splits[0])

	jsonOBJ = '{"Address":"' + splits[0] + '"}'
	postRequest("get-address", jsonOBJ, function(resp){
		obj = JSON.parse(resp)
		if (obj.Error != "none") {
			$("#balance").val("Error")
		} else {
			if(obj.Content.Address.startsWith("FA")) {
				$("#balance").val(FCTNormalize(obj.Content.Balance))
				$("#balance-type").text("FCT")
			} else {
				$("#balance").val(obj.Content.Balance)
				$("#balance-type").text("EC")
			}
		}
	})

	// Remove last QR code
	$('#qrcode').text("")
	// Add new QR code
	$('#qrcode').qrcode({
		width: 256,
		height: 256,
		text: splits[0]
	})
}

// QR Code from: https://github.com/jeromeetienne/jquery-qrcode
// MIT License
(function(r){r.fn.qrcode=function(h){var s;function u(a){this.mode=s;this.data=a}function o(a,c){this.typeNumber=a;this.errorCorrectLevel=c;this.modules=null;this.moduleCount=0;this.dataCache=null;this.dataList=[]}function q(a,c){if(void 0==a.length)throw Error(a.length+"/"+c);for(var d=0;d<a.length&&0==a[d];)d++;this.num=Array(a.length-d+c);for(var b=0;b<a.length-d;b++)this.num[b]=a[b+d]}function p(a,c){this.totalCount=a;this.dataCount=c}function t(){this.buffer=[];this.length=0}u.prototype={getLength:function(){return this.data.length},
write:function(a){for(var c=0;c<this.data.length;c++)a.put(this.data.charCodeAt(c),8)}};o.prototype={addData:function(a){this.dataList.push(new u(a));this.dataCache=null},isDark:function(a,c){if(0>a||this.moduleCount<=a||0>c||this.moduleCount<=c)throw Error(a+","+c);return this.modules[a][c]},getModuleCount:function(){return this.moduleCount},make:function(){if(1>this.typeNumber){for(var a=1,a=1;40>a;a++){for(var c=p.getRSBlocks(a,this.errorCorrectLevel),d=new t,b=0,e=0;e<c.length;e++)b+=c[e].dataCount;
for(e=0;e<this.dataList.length;e++)c=this.dataList[e],d.put(c.mode,4),d.put(c.getLength(),j.getLengthInBits(c.mode,a)),c.write(d);if(d.getLengthInBits()<=8*b)break}this.typeNumber=a}this.makeImpl(!1,this.getBestMaskPattern())},makeImpl:function(a,c){this.moduleCount=4*this.typeNumber+17;this.modules=Array(this.moduleCount);for(var d=0;d<this.moduleCount;d++){this.modules[d]=Array(this.moduleCount);for(var b=0;b<this.moduleCount;b++)this.modules[d][b]=null}this.setupPositionProbePattern(0,0);this.setupPositionProbePattern(this.moduleCount-
7,0);this.setupPositionProbePattern(0,this.moduleCount-7);this.setupPositionAdjustPattern();this.setupTimingPattern();this.setupTypeInfo(a,c);7<=this.typeNumber&&this.setupTypeNumber(a);null==this.dataCache&&(this.dataCache=o.createData(this.typeNumber,this.errorCorrectLevel,this.dataList));this.mapData(this.dataCache,c)},setupPositionProbePattern:function(a,c){for(var d=-1;7>=d;d++)if(!(-1>=a+d||this.moduleCount<=a+d))for(var b=-1;7>=b;b++)-1>=c+b||this.moduleCount<=c+b||(this.modules[a+d][c+b]=
0<=d&&6>=d&&(0==b||6==b)||0<=b&&6>=b&&(0==d||6==d)||2<=d&&4>=d&&2<=b&&4>=b?!0:!1)},getBestMaskPattern:function(){for(var a=0,c=0,d=0;8>d;d++){this.makeImpl(!0,d);var b=j.getLostPoint(this);if(0==d||a>b)a=b,c=d}return c},createMovieClip:function(a,c,d){a=a.createEmptyMovieClip(c,d);this.make();for(c=0;c<this.modules.length;c++)for(var d=1*c,b=0;b<this.modules[c].length;b++){var e=1*b;this.modules[c][b]&&(a.beginFill(0,100),a.moveTo(e,d),a.lineTo(e+1,d),a.lineTo(e+1,d+1),a.lineTo(e,d+1),a.endFill())}return a},
setupTimingPattern:function(){for(var a=8;a<this.moduleCount-8;a++)null==this.modules[a][6]&&(this.modules[a][6]=0==a%2);for(a=8;a<this.moduleCount-8;a++)null==this.modules[6][a]&&(this.modules[6][a]=0==a%2)},setupPositionAdjustPattern:function(){for(var a=j.getPatternPosition(this.typeNumber),c=0;c<a.length;c++)for(var d=0;d<a.length;d++){var b=a[c],e=a[d];if(null==this.modules[b][e])for(var f=-2;2>=f;f++)for(var i=-2;2>=i;i++)this.modules[b+f][e+i]=-2==f||2==f||-2==i||2==i||0==f&&0==i?!0:!1}},setupTypeNumber:function(a){for(var c=
j.getBCHTypeNumber(this.typeNumber),d=0;18>d;d++){var b=!a&&1==(c>>d&1);this.modules[Math.floor(d/3)][d%3+this.moduleCount-8-3]=b}for(d=0;18>d;d++)b=!a&&1==(c>>d&1),this.modules[d%3+this.moduleCount-8-3][Math.floor(d/3)]=b},setupTypeInfo:function(a,c){for(var d=j.getBCHTypeInfo(this.errorCorrectLevel<<3|c),b=0;15>b;b++){var e=!a&&1==(d>>b&1);6>b?this.modules[b][8]=e:8>b?this.modules[b+1][8]=e:this.modules[this.moduleCount-15+b][8]=e}for(b=0;15>b;b++)e=!a&&1==(d>>b&1),8>b?this.modules[8][this.moduleCount-
b-1]=e:9>b?this.modules[8][15-b-1+1]=e:this.modules[8][15-b-1]=e;this.modules[this.moduleCount-8][8]=!a},mapData:function(a,c){for(var d=-1,b=this.moduleCount-1,e=7,f=0,i=this.moduleCount-1;0<i;i-=2)for(6==i&&i--;;){for(var g=0;2>g;g++)if(null==this.modules[b][i-g]){var n=!1;f<a.length&&(n=1==(a[f]>>>e&1));j.getMask(c,b,i-g)&&(n=!n);this.modules[b][i-g]=n;e--; -1==e&&(f++,e=7)}b+=d;if(0>b||this.moduleCount<=b){b-=d;d=-d;break}}}};o.PAD0=236;o.PAD1=17;o.createData=function(a,c,d){for(var c=p.getRSBlocks(a,
c),b=new t,e=0;e<d.length;e++){var f=d[e];b.put(f.mode,4);b.put(f.getLength(),j.getLengthInBits(f.mode,a));f.write(b)}for(e=a=0;e<c.length;e++)a+=c[e].dataCount;if(b.getLengthInBits()>8*a)throw Error("code length overflow. ("+b.getLengthInBits()+">"+8*a+")");for(b.getLengthInBits()+4<=8*a&&b.put(0,4);0!=b.getLengthInBits()%8;)b.putBit(!1);for(;!(b.getLengthInBits()>=8*a);){b.put(o.PAD0,8);if(b.getLengthInBits()>=8*a)break;b.put(o.PAD1,8)}return o.createBytes(b,c)};o.createBytes=function(a,c){for(var d=
0,b=0,e=0,f=Array(c.length),i=Array(c.length),g=0;g<c.length;g++){var n=c[g].dataCount,h=c[g].totalCount-n,b=Math.max(b,n),e=Math.max(e,h);f[g]=Array(n);for(var k=0;k<f[g].length;k++)f[g][k]=255&a.buffer[k+d];d+=n;k=j.getErrorCorrectPolynomial(h);n=(new q(f[g],k.getLength()-1)).mod(k);i[g]=Array(k.getLength()-1);for(k=0;k<i[g].length;k++)h=k+n.getLength()-i[g].length,i[g][k]=0<=h?n.get(h):0}for(k=g=0;k<c.length;k++)g+=c[k].totalCount;d=Array(g);for(k=n=0;k<b;k++)for(g=0;g<c.length;g++)k<f[g].length&&
(d[n++]=f[g][k]);for(k=0;k<e;k++)for(g=0;g<c.length;g++)k<i[g].length&&(d[n++]=i[g][k]);return d};s=4;for(var j={PATTERN_POSITION_TABLE:[[],[6,18],[6,22],[6,26],[6,30],[6,34],[6,22,38],[6,24,42],[6,26,46],[6,28,50],[6,30,54],[6,32,58],[6,34,62],[6,26,46,66],[6,26,48,70],[6,26,50,74],[6,30,54,78],[6,30,56,82],[6,30,58,86],[6,34,62,90],[6,28,50,72,94],[6,26,50,74,98],[6,30,54,78,102],[6,28,54,80,106],[6,32,58,84,110],[6,30,58,86,114],[6,34,62,90,118],[6,26,50,74,98,122],[6,30,54,78,102,126],[6,26,52,
78,104,130],[6,30,56,82,108,134],[6,34,60,86,112,138],[6,30,58,86,114,142],[6,34,62,90,118,146],[6,30,54,78,102,126,150],[6,24,50,76,102,128,154],[6,28,54,80,106,132,158],[6,32,58,84,110,136,162],[6,26,54,82,110,138,166],[6,30,58,86,114,142,170]],G15:1335,G18:7973,G15_MASK:21522,getBCHTypeInfo:function(a){for(var c=a<<10;0<=j.getBCHDigit(c)-j.getBCHDigit(j.G15);)c^=j.G15<<j.getBCHDigit(c)-j.getBCHDigit(j.G15);return(a<<10|c)^j.G15_MASK},getBCHTypeNumber:function(a){for(var c=a<<12;0<=j.getBCHDigit(c)-
j.getBCHDigit(j.G18);)c^=j.G18<<j.getBCHDigit(c)-j.getBCHDigit(j.G18);return a<<12|c},getBCHDigit:function(a){for(var c=0;0!=a;)c++,a>>>=1;return c},getPatternPosition:function(a){return j.PATTERN_POSITION_TABLE[a-1]},getMask:function(a,c,d){switch(a){case 0:return 0==(c+d)%2;case 1:return 0==c%2;case 2:return 0==d%3;case 3:return 0==(c+d)%3;case 4:return 0==(Math.floor(c/2)+Math.floor(d/3))%2;case 5:return 0==c*d%2+c*d%3;case 6:return 0==(c*d%2+c*d%3)%2;case 7:return 0==(c*d%3+(c+d)%2)%2;default:throw Error("bad maskPattern:"+
a);}},getErrorCorrectPolynomial:function(a){for(var c=new q([1],0),d=0;d<a;d++)c=c.multiply(new q([1,l.gexp(d)],0));return c},getLengthInBits:function(a,c){if(1<=c&&10>c)switch(a){case 1:return 10;case 2:return 9;case s:return 8;case 8:return 8;default:throw Error("mode:"+a);}else if(27>c)switch(a){case 1:return 12;case 2:return 11;case s:return 16;case 8:return 10;default:throw Error("mode:"+a);}else if(41>c)switch(a){case 1:return 14;case 2:return 13;case s:return 16;case 8:return 12;default:throw Error("mode:"+
a);}else throw Error("type:"+c);},getLostPoint:function(a){for(var c=a.getModuleCount(),d=0,b=0;b<c;b++)for(var e=0;e<c;e++){for(var f=0,i=a.isDark(b,e),g=-1;1>=g;g++)if(!(0>b+g||c<=b+g))for(var h=-1;1>=h;h++)0>e+h||c<=e+h||0==g&&0==h||i==a.isDark(b+g,e+h)&&f++;5<f&&(d+=3+f-5)}for(b=0;b<c-1;b++)for(e=0;e<c-1;e++)if(f=0,a.isDark(b,e)&&f++,a.isDark(b+1,e)&&f++,a.isDark(b,e+1)&&f++,a.isDark(b+1,e+1)&&f++,0==f||4==f)d+=3;for(b=0;b<c;b++)for(e=0;e<c-6;e++)a.isDark(b,e)&&!a.isDark(b,e+1)&&a.isDark(b,e+
2)&&a.isDark(b,e+3)&&a.isDark(b,e+4)&&!a.isDark(b,e+5)&&a.isDark(b,e+6)&&(d+=40);for(e=0;e<c;e++)for(b=0;b<c-6;b++)a.isDark(b,e)&&!a.isDark(b+1,e)&&a.isDark(b+2,e)&&a.isDark(b+3,e)&&a.isDark(b+4,e)&&!a.isDark(b+5,e)&&a.isDark(b+6,e)&&(d+=40);for(e=f=0;e<c;e++)for(b=0;b<c;b++)a.isDark(b,e)&&f++;a=Math.abs(100*f/c/c-50)/5;return d+10*a}},l={glog:function(a){if(1>a)throw Error("glog("+a+")");return l.LOG_TABLE[a]},gexp:function(a){for(;0>a;)a+=255;for(;256<=a;)a-=255;return l.EXP_TABLE[a]},EXP_TABLE:Array(256),
LOG_TABLE:Array(256)},m=0;8>m;m++)l.EXP_TABLE[m]=1<<m;for(m=8;256>m;m++)l.EXP_TABLE[m]=l.EXP_TABLE[m-4]^l.EXP_TABLE[m-5]^l.EXP_TABLE[m-6]^l.EXP_TABLE[m-8];for(m=0;255>m;m++)l.LOG_TABLE[l.EXP_TABLE[m]]=m;q.prototype={get:function(a){return this.num[a]},getLength:function(){return this.num.length},multiply:function(a){for(var c=Array(this.getLength()+a.getLength()-1),d=0;d<this.getLength();d++)for(var b=0;b<a.getLength();b++)c[d+b]^=l.gexp(l.glog(this.get(d))+l.glog(a.get(b)));return new q(c,0)},mod:function(a){if(0>
this.getLength()-a.getLength())return this;for(var c=l.glog(this.get(0))-l.glog(a.get(0)),d=Array(this.getLength()),b=0;b<this.getLength();b++)d[b]=this.get(b);for(b=0;b<a.getLength();b++)d[b]^=l.gexp(l.glog(a.get(b))+c);return(new q(d,0)).mod(a)}};p.RS_BLOCK_TABLE=[[1,26,19],[1,26,16],[1,26,13],[1,26,9],[1,44,34],[1,44,28],[1,44,22],[1,44,16],[1,70,55],[1,70,44],[2,35,17],[2,35,13],[1,100,80],[2,50,32],[2,50,24],[4,25,9],[1,134,108],[2,67,43],[2,33,15,2,34,16],[2,33,11,2,34,12],[2,86,68],[4,43,27],
[4,43,19],[4,43,15],[2,98,78],[4,49,31],[2,32,14,4,33,15],[4,39,13,1,40,14],[2,121,97],[2,60,38,2,61,39],[4,40,18,2,41,19],[4,40,14,2,41,15],[2,146,116],[3,58,36,2,59,37],[4,36,16,4,37,17],[4,36,12,4,37,13],[2,86,68,2,87,69],[4,69,43,1,70,44],[6,43,19,2,44,20],[6,43,15,2,44,16],[4,101,81],[1,80,50,4,81,51],[4,50,22,4,51,23],[3,36,12,8,37,13],[2,116,92,2,117,93],[6,58,36,2,59,37],[4,46,20,6,47,21],[7,42,14,4,43,15],[4,133,107],[8,59,37,1,60,38],[8,44,20,4,45,21],[12,33,11,4,34,12],[3,145,115,1,146,
116],[4,64,40,5,65,41],[11,36,16,5,37,17],[11,36,12,5,37,13],[5,109,87,1,110,88],[5,65,41,5,66,42],[5,54,24,7,55,25],[11,36,12],[5,122,98,1,123,99],[7,73,45,3,74,46],[15,43,19,2,44,20],[3,45,15,13,46,16],[1,135,107,5,136,108],[10,74,46,1,75,47],[1,50,22,15,51,23],[2,42,14,17,43,15],[5,150,120,1,151,121],[9,69,43,4,70,44],[17,50,22,1,51,23],[2,42,14,19,43,15],[3,141,113,4,142,114],[3,70,44,11,71,45],[17,47,21,4,48,22],[9,39,13,16,40,14],[3,135,107,5,136,108],[3,67,41,13,68,42],[15,54,24,5,55,25],[15,
43,15,10,44,16],[4,144,116,4,145,117],[17,68,42],[17,50,22,6,51,23],[19,46,16,6,47,17],[2,139,111,7,140,112],[17,74,46],[7,54,24,16,55,25],[34,37,13],[4,151,121,5,152,122],[4,75,47,14,76,48],[11,54,24,14,55,25],[16,45,15,14,46,16],[6,147,117,4,148,118],[6,73,45,14,74,46],[11,54,24,16,55,25],[30,46,16,2,47,17],[8,132,106,4,133,107],[8,75,47,13,76,48],[7,54,24,22,55,25],[22,45,15,13,46,16],[10,142,114,2,143,115],[19,74,46,4,75,47],[28,50,22,6,51,23],[33,46,16,4,47,17],[8,152,122,4,153,123],[22,73,45,
3,74,46],[8,53,23,26,54,24],[12,45,15,28,46,16],[3,147,117,10,148,118],[3,73,45,23,74,46],[4,54,24,31,55,25],[11,45,15,31,46,16],[7,146,116,7,147,117],[21,73,45,7,74,46],[1,53,23,37,54,24],[19,45,15,26,46,16],[5,145,115,10,146,116],[19,75,47,10,76,48],[15,54,24,25,55,25],[23,45,15,25,46,16],[13,145,115,3,146,116],[2,74,46,29,75,47],[42,54,24,1,55,25],[23,45,15,28,46,16],[17,145,115],[10,74,46,23,75,47],[10,54,24,35,55,25],[19,45,15,35,46,16],[17,145,115,1,146,116],[14,74,46,21,75,47],[29,54,24,19,
55,25],[11,45,15,46,46,16],[13,145,115,6,146,116],[14,74,46,23,75,47],[44,54,24,7,55,25],[59,46,16,1,47,17],[12,151,121,7,152,122],[12,75,47,26,76,48],[39,54,24,14,55,25],[22,45,15,41,46,16],[6,151,121,14,152,122],[6,75,47,34,76,48],[46,54,24,10,55,25],[2,45,15,64,46,16],[17,152,122,4,153,123],[29,74,46,14,75,47],[49,54,24,10,55,25],[24,45,15,46,46,16],[4,152,122,18,153,123],[13,74,46,32,75,47],[48,54,24,14,55,25],[42,45,15,32,46,16],[20,147,117,4,148,118],[40,75,47,7,76,48],[43,54,24,22,55,25],[10,
45,15,67,46,16],[19,148,118,6,149,119],[18,75,47,31,76,48],[34,54,24,34,55,25],[20,45,15,61,46,16]];p.getRSBlocks=function(a,c){var d=p.getRsBlockTable(a,c);if(void 0==d)throw Error("bad rs block @ typeNumber:"+a+"/errorCorrectLevel:"+c);for(var b=d.length/3,e=[],f=0;f<b;f++)for(var h=d[3*f+0],g=d[3*f+1],j=d[3*f+2],l=0;l<h;l++)e.push(new p(g,j));return e};p.getRsBlockTable=function(a,c){switch(c){case 1:return p.RS_BLOCK_TABLE[4*(a-1)+0];case 0:return p.RS_BLOCK_TABLE[4*(a-1)+1];case 3:return p.RS_BLOCK_TABLE[4*
(a-1)+2];case 2:return p.RS_BLOCK_TABLE[4*(a-1)+3]}};t.prototype={get:function(a){return 1==(this.buffer[Math.floor(a/8)]>>>7-a%8&1)},put:function(a,c){for(var d=0;d<c;d++)this.putBit(1==(a>>>c-d-1&1))},getLengthInBits:function(){return this.length},putBit:function(a){var c=Math.floor(this.length/8);this.buffer.length<=c&&this.buffer.push(0);a&&(this.buffer[c]|=128>>>this.length%8);this.length++}};"string"===typeof h&&(h={text:h});h=r.extend({},{render:"canvas",width:256,height:256,typeNumber:-1,
correctLevel:2,background:"#ffffff",foreground:"#000000"},h);return this.each(function(){var a;if("canvas"==h.render){a=new o(h.typeNumber,h.correctLevel);a.addData(h.text);a.make();var c=document.createElement("canvas");c.width=h.width;c.height=h.height;for(var d=c.getContext("2d"),b=h.width/a.getModuleCount(),e=h.height/a.getModuleCount(),f=0;f<a.getModuleCount();f++)for(var i=0;i<a.getModuleCount();i++){d.fillStyle=a.isDark(f,i)?h.foreground:h.background;var g=Math.ceil((i+1)*b)-Math.floor(i*b),
j=Math.ceil((f+1)*b)-Math.floor(f*b);d.fillRect(Math.round(i*b),Math.round(f*e),g,j)}}else{a=new o(h.typeNumber,h.correctLevel);a.addData(h.text);a.make();c=r("<table></table>").css("width",h.width+"px").css("height",h.height+"px").css("border","0px").css("border-collapse","collapse").css("background-color",h.background);d=h.width/a.getModuleCount();b=h.height/a.getModuleCount();for(e=0;e<a.getModuleCount();e++){f=r("<tr></tr>").css("height",b+"px").appendTo(c);for(i=0;i<a.getModuleCount();i++)r("<td></td>").css("width",
d+"px").css("background-color",a.isDark(e,i)?h.foreground:h.background).appendTo(f)}}a=c;jQuery(a).appendTo(this)})}})(jQuery);

$("#generate-source").on("change", function(){
	selected = $("#generate-source option:selected").val()
	if(selected == "new-external-address") {
		$("#sec-pub").text("Public")
	} else {
		$("#sec-pub").text("Private")
	}
	if(selected == "import-address" || selected == "new-external-address" || selected == "import-koinify"){
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
	$("#private-key-input").val("")
	if(selected == "import-address"){
		$("#private-key-input").attr("placeholder","Type address private key")
	} else if(selected == "random-ec"){
		$("#private-key-input").attr("placeholder","A new entry credit address will be created")
	} else if(selected == "random-factoid"){
		$("#private-key-input").attr("placeholder","A new factoid address will be created")
	} else if(selected == "new-external-address"){
		$("#private-key-input").attr("placeholder","Type a public address to add to your contacts")
	} else if(selected == "import-koinify"){
		$("#private-key-input").attr("placeholder","Type in your Koinify phrase")
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
	} else if(selected == "import-koinify") {
		koinify = $("#private-key-input").val()

		var newAddressObj = {
	    	Name:Name,
			Koinify:koinify
		}

		j = JSON.stringify(newAddressObj)
		postRequest("import-koinify", j, function(resp){
			obj = JSON.parse(resp)
			if(obj.Error == "none"){
				SetSuccess(obj)
			} else {
				SetError(obj.Error)
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
var Name = ""
var Address = ""

/*$(window).load(function() {
    GetDefaultData()
});*/

function GetDefaultData(){
	Name = $("#address-name").val()
	Address = $("#address-field").val()

	jsonOBJ = '{"Address":"' + Address + '"}'
	postRequest("get-address", jsonOBJ, function(resp){
		obj = JSON.parse(resp)
		if (obj.Error != "none") {
			$("#balance-container").text("Can not find the addresses in address book")
		} else {
			if(obj.Content.Address.startsWith("FA")) {
				$("#balance").text(FCTNormalize(obj.Content.Balance))
			} else {
				$("#balance").text(obj.Content.Balance)
			}
		}
	})
}

$("#display-private-key").click(function(){
	jsonOBJ = '{"Address":"' + Address + '"}'
	postRequest("display-private-key", jsonOBJ, function(resp){
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
$("#backup-input-verify-button").on('click', function(){
	$("#backup-input-verify-button").addClass("backup-btn-checking")

	var seedSingle = $("#given-seed").attr("value");
	var inputedSeedSplit = new Array(12);

	$(".backup-input").each(function(){
		var i = Number($(this).attr("index"))
		inputedSeedSplit[i-1] = $(this).val()
	})

	var clear = function() {
		$("#backup-input-verify-button").removeClass("backup-btn-failed")
		$("#backup-input-verify-button").removeClass("backup-btn-checking")
		$("#backup-input-verify-button").removeClass("backup-btn-verified")
		$("#backup-error-message").removeClass("active")
	}

	if(inputedSeedSplit.join(" ") === seedSingle) {
		clear()
		$("#backup-input-verify-button").addClass("backup-btn-verified")
		setTimeout(function() {
			document.getElementById("link-to-success").click()
		}, 1000);
	} else {
		clear()
		$("#backup-input-verify-button").addClass("backup-btn-failed")
		$("#backup-error-message").addClass("active")
		setTimeout(clear, 3000);
	}
})

$(".external-link").on('click', function(e){
	e.preventDefault()
	require('electron').shell.openExternal($(this).attr("href"))
})

$("#backup-html-form").on("submit",function(event){event.preventDefault()})

function LoadBackup0() {
	LoadAddresses(false)
}

function updateBackupConfirmCheckbox() {
	var c = $("#wrote-down-confirm-checkbox")
	if(!c.is(':checked')) {
		c.prop('checked', true);
	} else {
		c.prop('checked', false);
	}

	if(c.is(':checked')) {
		document.getElementById("wrote-down-confirm").disabled = false;
	} else {
		document.getElementById("wrote-down-confirm").disabled = true;
	}
}

// Importing

$("#import-input-confirm-button").on('click', function(){
	$("#import-input-confirm-button").addClass("backup-btn-checking")

	var seedSingle = $("#given-seed").attr("value");
	var inputedSeedSplit = new Array(12);

	$(".import-input").each(function(){
		var i = Number($(this).attr("index"))
		inputedSeedSplit[i-1] = $(this).val()
	})

	var clear = function() {
		$("#import-input-confirm-button").removeClass("backup-btn-failed")
		$("#import-input-confirm-button").removeClass("backup-btn-checking")
		$("#import-input-confirm-button").removeClass("backup-btn-verified")
		$("#backup-error-message").removeClass("active")
	}


	seed = inputedSeedSplit.join(" ")
	var SeedStruct  = {
    	Seed:seed,
  	}
  	j = JSON.stringify(SeedStruct)
	postRequest("import-seed", j, function(resp) {
		obj = JSON.parse(resp)
		if(obj.Error == "none") {
			clear()
			$("#import-input-confirm-button").addClass("backup-btn-verified")
	    	SetGeneralSuccess("Seed has been changed to: " + obj.Content)
	    	setTimeout(function() {
	    		document.getElementById("link-to-success").click()
	    	}, 1000);
	    } else {
	    	clear()
			$("#import-input-confirm-button").addClass("backup-btn-failed")
	    	$("#backup-error-message").addClass("active")
	    	setTimeout(clear, 3000);
	    }
	})

	/*if(inputedSeedSplit.join(" ") === seedSingle) {
		clear()
		$("#import-input-confirm-button").addClass("backup-btn-verified")
	} else {
		clear()
		$("#import-input-confirm-button").addClass("backup-btn-failed")
		setTimeout(clear, 3000);
	}*/
})


// Do action with imported transaction
function receivedText() {
	is = fr.result
	len = is.split(" ")
	if(len.length != 12) {
		SetGeneralError("Seed must be 12 words");
		return
	}
	document.getElementById('data-expand').click()
	$("#import-seed-reveal-text").text(fr.result)
	$("#import-seed-reveal-cancel").click()
}

$("#import-seed-reveal-confirm").on('click', function(){
	seed = $("#import-seed-reveal-text").text()
	var SeedStruct  = {
    	Seed:seed,
  	}
  	j = JSON.stringify(SeedStruct)
	postRequest("import-seed", j, function(resp) {
		obj = JSON.parse(resp)
		if(obj.Error == "none") {
	    	SetGeneralSuccess("Seed has been changed to: " + obj.Content)
	    } else {
	    	SetGeneralError("Error: " + obj.Error)
	    }
	})
})