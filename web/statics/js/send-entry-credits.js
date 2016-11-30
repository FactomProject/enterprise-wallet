counter = 1
function addNewOutputAddress(defaultVal, error) {
  eClass = ""
  if(error){
    eClass = "input-group-error"
  }
  $("#all-outputs").append(
  '<div class="row" id="single-output">' +
  '    <div class="small-12 medium-7 large-8 columns">' +
  '        <div class="input-group ' + eClass + '" id="output-factoid-address-container">' +
  '            <pre><input id="output-factoid-address" type="text" name="output1" class="input-group-field percent95" placeholder="Type entry credit address" value="' + defaultVal + '"></pre>' +
  '            <!-- <a data-toggle="addressbook" class="input-group-button button" id="addressbook-' + counter + '"><i class="fa fa-book"></i></a> -->' +
  '        </div>' +
  '    </div>' +
  '    <div class="small-10 medium-4 large-3 columns">' +
  '        <div class="input-group">' +
  '            <input id="output-factoid-amount" type="text" class="input-group-field" name="output1-num" placeholder="Amount of entry credits">' +
  '            <span class="input-group-label">EC</span>' +
  '        </div>' +
  '    </div>' +
  '    <div class="small-2 medium-1 columns">' +
  '            <a id="remove-new-output" class="button expanded newMinus">&nbsp;</a>' +
  '    </div>' +
  '</div>')
}

// Add new output address
$("#append-new-output").click(function(){
  addNewOutputAddress("", true)
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
  if ((evt.which < 48 || evt.which > 57)) {
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
    TransType:"ec",
		OutputAddresses:[],
		OutputAmounts:[]
	}

  errMessage = ""
  faErr = false
  amtErr = false

	$("#all-outputs #single-output").each(function(){
		add = $(this).find("#output-factoid-address").val()
    if(!add.startsWith("EC")) {
      $(this).find("#output-factoid-address-container").addClass("input-group-error")
      faErr = true
      err = true
    }

		amt = $(this).find("#output-factoid-amount").val()
		if(amt == 0 || amt == undefined) {
      amtErr = true
      $(this).find("#output-factoid-amount").addClass("input-group-error")
    }

		transObject.OutputAddresses.push(add)
		transObject.OutputAmounts.push(amt)
	})

  if(err){
    if(faErr){errMessage += "Addresses must start with 'EC'. "}
    if(amtErr){errMessage += "Amounts should not be 0. "}
    SetGeneralError("Error(s): " + errMessage)
    return
  }

	j = JSON.stringify(transObject)
	postRequest("send-transaction", j, function(resp){
		console.log(resp)
    obj = JSON.parse(resp)
    if(obj.Error == "none") {
      SetGeneralSuccess("Transaction sent, still working on better confirmation")
    } else {
      SetGeneralError("Error: " + obj.Error)
    }
	})
})

// Load the Reveal
$(window).load(function() {
    LoadAddresses()
});

function LoadAddresses(){
  resp = getRequest("addresses",function(resp){
    obj = JSON.parse(resp)
    
    /*obj.FactoidAddresses.List.forEach(function(address){
      $('#addresses-reveal').append(factoidECRadio(address, "factoid"));
    })*/

    obj.EntryCreditAddresses.List.forEach(function(address){
      $('#addresses-reveal').append(factoidECRadio(address, "entry-credits"));
    })

    obj.ExternalAddresses.List.forEach(function(address){
      if(address.Address.startsWith("EC")){
        $('#addresses-reveal').append(factoidECRadio(address, "external"));
      }
    })
  })
}

function factoidECRadio(address, type){
return '<pre>' +
'  <input type="radio" name="address" id="address" value="' + address.Address + '"> <span id="address-name" name="' + address.Name + '">' + address.Name + '</span>' +
'</pre>'
}

$('#addresses-reveal').on("mouseover", "#address-name", function(){
  $(this).css("font-size", "90%")
  $(this).text($(this).parent().find("#address").val());
})

$('#addresses-reveal').on("mouseout", "#address-name", function(){
  $(this).text($(this).attr("name"));
  $(this).css("font-size", "100%")
})

done = false

$("#addresses-reveal-button").on("click", function(){
  newAddress = $("input[name='address']:checked").val()
  if(newAddress == undefined) {
    return
  }

  done = false
  $("#all-outputs #single-output").each(function(){
    if(!done) {
      addressDOM = $(this).find("#output-factoid-address")
      add = addressDOM.val()
      if(add == "") {
        addressDOM.val(newAddress)
        done = true
      }
    }
  })

  // No empty slot found
  if(!done){
    addNewOutputAddress(newAddress, false)
  }
})