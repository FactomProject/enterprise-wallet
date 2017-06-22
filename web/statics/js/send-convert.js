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


  defAmt = '<input id="output-factoid-amount" type="text" class="input-group-field" name="output1-num" placeholder="Amount of ' + PageToken + '">'
  if(defaultAmt != 0) {
     defAmt = '<input id="output-factoid-amount" type="text" class="input-group-field" name="output1-num" placeholder="Amount of ' + PageToken + '" value="' + defaultAmt + '">'
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

  defAmt = '<input id="input-factoid-amount" type="text" class="input-group-field" name="output1-num" placeholder="Amount of factoids">'
  if(defaultAmt != 0) {
    defAmt = '<input id="input-factoid-amount" type="text" class="input-group-field" name="output1-num" placeholder="Amount of factoids" value="' + defaultAmt + '">'
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
    if(PageTokenABR == "FCT") {
      if(evt.which < 48 || (evt.which > 57 && (evt.which != 190))) evt.preventDefault();
    } else {
        if(evt.which < 48 || evt.which > 57) evt.preventDefault();
      }
    });
//   if(PageTokenABR == "FCT") {
//     var self = $(this);
//     // self.val(self.val().replace(/[^0-9\.]/g, ''));
//     if ((evt.which != 46 || self.val().indexOf('.') != -1) && (evt.which < 48 || evt.which > 57) && evt.which != 8) {
//       //evt.preventDefault();
//     }

//     decSplit = $(this).val().split(".")
//     if(decSplit.length > 2) {
//       evt.preventDefault();
//     }
//   } else {
//     var self = $(this);
//     self.val(self.val().replace(/[^0-9\.]/g, ''));
//     if ((evt.which < 48 || evt.which > 57) && evt.which != 8) {
//       evt.preventDefault();
//     }
//   }
// });


// Update Fee
$("#all-outputs").on('change', '#output-factoid-amount', function(){
	// Need to determine new fee
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

function MakeTransaction(sig) {
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
      disableInput()
      ShowNewButtons()
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
      if(importexport) {
        setExportDownload(obj.Content.Json)
        SetGeneralSuccess('Click "Export Transaction" to download, or go back to editing it')
      } else {
        SetGeneralSuccess('Click "Send Transaction" to send, or go back to editing it')
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
  '  <input type="radio" name="' + name + '" id="address" value="' + address.Address + 
  '"> <span for="' + address.Address + '" id="address-name" name="(' + FCTNormalize(address.Balance)  + " FCT) "  + address.Name + '">(' + FCTNormalize(address.Balance)  + " FCT) " + address.Name + '</span>' +
  '</pre><br />'
}

$('#addresses-reveal').on("mouseover", "#address-name", function(){
  $(this).css("font-size", "90%")
  $(this).text($(this).parent().find("#address").val());
})

$('#addresses-reveal').on("mouseout", "#address-name", function(){
  $(this).text($(this).attr("name"));
  $(this).css("font-size", "100%")
})

$('#fee-addresses-reveal').on("mouseover", "#address-name", function(){
  $(this).css("font-size", "90%")
  $(this).text($(this).parent().find("#address").val());
})

$('#fee-addresses-reveal').on("mouseout", "#address-name", function(){
  $(this).text($(this).attr("name"));
  $(this).css("font-size", "100%")
})

function factoidECRadio(address, type){
  return '<pre>' +
  '  <input type="radio" name="address" id="address" value="' + address.Address + '"> <span id="address-name" name="' + address.Name + '">' + address.Name + '</span>' +
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
