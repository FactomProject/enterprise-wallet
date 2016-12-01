$(window).load(function() {
    LoadTransactions()
});

CurrentCount = 0
ContentLen = 0
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

		// Load past 50 transactions, then stop. Only load more if they scroll
		loopstop = 50
		if(ContentLen < 50) {
			loopstop = ContentLen
		}
		for(; CurrentCount < loopstop; CurrentCount++) {
			AppendNewTransaction(Transactions[CurrentCount], CurrentCount)
		}

	})
}

function AppendNewTransaction(trans, index){
	console.log(trans)
	inputs = ""
	outputs = ""
	token = ""
	personalGain = 0
	personalLoss = 0
	personalECGain = 0

	// Transactions with sent and received to addresses will appear as 2 seperate ones

	for(var i = 0; i < trans.Inputs.length; i++) {
		if(trans.Inputs[i].Name != "") {
			inputs += '<div class="nick red">' + trans.Inputs[i].Name + '<pre class="show-for-large"> (' + trans.Inputs[i].Address + ')</pre></div>'
			personalLoss += trans.Inputs[i].Amount * -1
		}
	}	
	amountLost = (trans.TotalInput / 1e8) * -1
	token = "FCT"

	for(var i = 0; i < trans.Outputs.length; i++) {
		if(trans.Outputs[i].Name != "") {
			if(trans.Outputs[i].Address.startsWith("FA")){
				outputs += '<div class="nick green">' + trans.Outputs[i].Name + '<pre class="show-for-large percent95"> (' + trans.Outputs[i].Address + ')</pre></div>'
				personalGain += trans.Outputs[i].Amount
			} else {
				outputs += '<div class="nick green">' + trans.Outputs[i].Name + '<pre class="show-for-large percent95"> (' + trans.Outputs[i].Address + ')</pre></div>'
				personalECGain += trans.Outputs[i].Amount
			}
		}
	}
	amountGained = trans.TotalFCTOutput / 1e8
	token = "FCT"

	personalAmt = (personalGain + personalLoss) / 1e8
	pic = "sent"
	if(personalAmt >= 0) {
		pic = "received"
	}

	if(personalECGain > 0) {
		pic = "converted"
		trans.Action = "converted"
	}

	$("#transaction-list").append(
       '<tr>' +
            '<td><a data-toggle="transDetails"><i class="transIcon ' + trans.Action + '"><img src="img/transaction_' + pic + '.svg" class="svg"></i></a></td>' +
            '<td>' + trans.Date + ' : <a value="' + CurrentCount + '" id="transaction-link" data-toggle="transDetails">' + trans.Action.capitalize() + '</a>' +
            inputs + outputs + '</td>' +
            '<td style="word-wrap: break-word;">' + Number(personalAmt.toFixed(4)) + token + '</td>' +
        '</tr>'
	)
}


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