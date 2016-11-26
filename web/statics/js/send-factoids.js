/*
<div class="row" id="all-outputs">

<div class="row" id="single-output">
    <div class="small-12 medium-7 large-8 columns">
        <div class="input-group">
            <input id="output-factoid-address" type="text" name="output1" class="input-group-field" placeholder="Type factoid address">
            <a data-toggle="addressbook" class="input-group-button button"><i class="fa fa-book"></i></a>
        </div>
    </div>
    <div class="small-10 medium-4 large-3 columns">
        <div class="input-group">
            <input id="output-factoid-amount" type="text" class="input-group-field" name="output1-num" placeholder="Amount of factoids">
            <span class="input-group-label">FCT</span>
        </div>
    </div>
    <div class="small-2 medium-1 columns">
            <a id="append-new-output" class="button expanded newPlus">&nbsp;</a>
    </div>
</div>
*/

// Add new output address
$("#append-new-output").click(function(){
	$("#all-outputs").append(
'<div class="row" id="single-output">' +
'    <div class="small-12 medium-7 large-8 columns">' +
'        <div class="input-group input-group-error" id="output-factoid-address-container">' +
'            <input id="output-factoid-address" type="text" name="output1" class="input-group-field" placeholder="Type factoid address">' +
'            <a data-toggle="addressbook" class="input-group-button button"><i class="fa fa-book"></i></a>' +
'        </div>' +
'    </div>' +
'    <div class="small-10 medium-4 large-3 columns">' +
'        <div class="input-group">' +
'            <input id="output-factoid-amount" type="text" class="input-group-field" name="output1-num" placeholder="Amount of factoids">' +
'            <span class="input-group-label">FCT</span>' +
'        </div>' +
'    </div>' +
'    <div class="small-2 medium-1 columns">' +
'            <a id="remove-new-output" class="button expanded newMinus">&nbsp;</a>' +
'    </div>' +
'</div>')
})

// Remove output address
$("#all-outputs").on('click','#remove-new-output', function(){
	jQuery(this).parent().parent().remove()
})

// Ensure factoids being sent are valid, this is not a security feature, but an ease of use
// feature
$("#all-outputs").on("keypress", "#output-factoid-amount", function(evt) {
  var self = $(this);
  self.val(self.val().replace(/[^0-9\.]/g, ''));
  if ((evt.which != 46 || self.val().indexOf('.') != -1) && (evt.which < 48 || evt.which > 57)) {
    evt.preventDefault();
  }

  decSplit = $(this).val().split(".")
  if(decSplit.length > 2) {
    evt.preventDefault();
  }
});

// Update Fee
$("#all-outputs").on('change', '#output-factoid-amount', function(){
	// Need to determine new fee
})

$("#all-outputs").on('click', '#output-factoid-address-container', function(){
	$(this).removeClass("input-group-error")
})

$("#send-entire-transaction").on('click', function(){
	// var transObject = new Object()
	var transObject = {
		OutputAddresses:[],
		OutputAmounts:[]
	}

	$("#all-outputs #single-output").each(function(){
		add = $(this).find("#output-factoid-address").val()
		amt = $(this).find("#output-factoid-amount").val()
		
		transObject.OutputAddresses.push(add)
		transObject.OutputAmounts.push(amt)
	})

	j = JSON.stringify(transObject)
	postRequest("send-transaction", j, function(resp){
		console.log(resp)
	})
})