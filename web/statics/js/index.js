$(window).load(function() {
    LoadTransactions()
});

CurrentCount = 0
ContentLen = 0
LoopstopIncrement = 10 // Amount to load on scroll
Loopstop = 15
var Transactions
Done = false

function LoadTransactions() {
	getRequest("related-transactions", function(resp){
		$("#loading-container").remove()
		obj = JSON.parse(resp)
		console.log(obj)

		if(obj.Error != "none"){
			SetGeneralError(obj.Error)
			return
		}

		ContentLen = obj.Content.length
		Transactions = obj.Content

		// Load past x transactions, then stop. Only load more if they scroll
		if(ContentLen < Loopstop) {
			Loopstop = ContentLen
		}
		for(; CurrentCount < Loopstop; CurrentCount++) {
			AppendNewTransaction(Transactions[CurrentCount], CurrentCount)
		}

	})
}

function LoadCached() {
		// Load past x transactions, then stop. Only load more if they scroll
		if(ContentLen < Loopstop) {
			Loopstop = ContentLen
		}
		for(; CurrentCount < Loopstop; CurrentCount++) {
			AppendNewTransaction(Transactions[CurrentCount], CurrentCount)
		}
}

function AppendNewTransaction(trans, index){
	console.log(trans)

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
        '<td><a data-toggle="transDetails"><i class="transIcon ' + pic + '"><img src="img/transaction_' + pic + '.svg" class="svg"></i></a></td>' +
        '<td>' + date + ' : <a value="' + index + '" id="transaction-link" data-toggle="transDetails">' + pic.capitalize() + '</a>' +
        addrs + '</td>' +
        '<td style="word-wrap: break-word;">' + Number(amt.toFixed(4)) + ' ' + token + '</td>' +
    '</tr>'
)
}

$("main").bind('scroll', function() {
	//console.log($("main").outerHeight(), $("main").scrollTop(), $("main").innerHeight(), $("main").prop('scrollHeight'), $("main").prop('offsetHeight'))
	// Total Height
	// $("main").prop('scrollHeight')
	if($("main").scrollTop() + $("main").innerHeight() >= .8 * $("main").prop('scrollHeight')) {
		Loopstop += LoopstopIncrement
		LoadCached()
	}
});

$("#transaction-list").on('click', '#transaction-link', function(){
	$("#transDetails #details").html(getTransDetails($(this).attr("value")))
	$("#transDetails #link").attr("href", "http://explorer.factom.org/tx/" + Transactions[$(this).attr("value")].TxID)
	// TODO: Remove local link or correct port
	$("#transDetails #local-link").attr("href", "http://localhost:8090/search?input=" + Transactions[$(this).attr("value")].TxID + "&type=facttransaction")
})

function getTransDetails(index){
	trans = Transactions[index]
	inputs =  ""
	outputs = ""
	ecOutputs = ""

	for(var i = 0; i < trans.Inputs.length; i++) {
		inputs += "<div>" + trans.Inputs[i].Name + "(<pre>" + trans.Inputs[i].Address + "</pre>)</div>"
	} 

	for(var i = 0; i < trans.Outputs.length; i++) {
		if(trans.Outputs[i].Address.startsWith("FA")) {
			outputs += "<div>" + trans.Outputs[i].Name + "(<pre>" + trans.Outputs[i].Address + "</pre>)</div>"
		} else {
			ecOutputs += "<div>" + trans.Outputs[i].Name + "(<pre>" + trans.Outputs[i].Address + "</pre>)</div>"
		}
	} 

	htmlBody = '' +
	'<div>Input Total: ' + (trans.TotalInput / 1e8).toFixed(4) + '</div>' +
	inputs +
	'<div>Output Total: ' + (trans.TotalFCTOutput / 1e8).toFixed(4) + ' FCT</div>' +
	outputs +
	'<div>EC Output Total: ' + (trans.TotalECOutput / 1e8).toFixed(4) + ' FCT</div>' +
	ecOutputs

	return htmlBody
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