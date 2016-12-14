var $jscomp={scope:{},findInternal:function(a,b,c){a instanceof String&&(a=String(a));for(var e=a.length,d=0;d<e;d++){var f=a[d];if(b.call(c,f,d,a))return{i:d,v:f}}return{i:-1,v:void 0}}};$jscomp.defineProperty="function"==typeof Object.defineProperties?Object.defineProperty:function(a,b,c){if(c.get||c.set)throw new TypeError("ES3 does not support getters and setters.");a!=Array.prototype&&a!=Object.prototype&&(a[b]=c.value)};
$jscomp.getGlobal=function(a){return"undefined"!=typeof window&&window===a?a:"undefined"!=typeof global&&null!=global?global:a};$jscomp.global=$jscomp.getGlobal(this);$jscomp.polyfill=function(a,b,c,e){if(b){c=$jscomp.global;a=a.split(".");for(e=0;e<a.length-1;e++){var d=a[e];d in c||(c[d]={});c=c[d]}a=a[a.length-1];e=c[a];b=b(e);b!=e&&null!=b&&$jscomp.defineProperty(c,a,{configurable:!0,writable:!0,value:b})}};
$jscomp.polyfill("Array.prototype.find",function(a){return a?a:function(a,c){return $jscomp.findInternal(this,a,c).v}},"es6-impl","es3");$jscomp.checkStringArgs=function(a,b,c){if(null==a)throw new TypeError("The 'this' value for String.prototype."+c+" must not be null or undefined");if(b instanceof RegExp)throw new TypeError("First argument to String.prototype."+c+" must not be a regular expression");return a+""};
$jscomp.polyfill("String.prototype.startsWith",function(a){return a?a:function(a,c){var b=$jscomp.checkStringArgs(this,a,"startsWith");a+="";for(var d=b.length,f=a.length,h=Math.max(0,Math.min(c|0,b.length)),g=0;g<f&&h<d;)if(b[h++]!=a[g++])return!1;return g>=f}},"es6-impl","es3");importexport=!1;PageTokenABR="FCT";PageToken="factoids";AddressPrefix="FA";PageTransType="factoid";
"1"==$("#token-header").attr("value")?(PageTokenABR="EC",PageToken="entry credits",AddressPrefix="EC",PageTransType="ec"):"2"==$("#token-header").attr("value")&&(importexport=!0);counter=2;
function addNewOutputAddress(a,b){eClass="";b&&(eClass="input-group-error");str="factoid";"FCT"!=PageTokenABR&&(str="entry credit");$("#all-outputs").append('<div class="row single-output-'+counter+'" id="single-output">    <div class="small-12 medium-7 large-8 columns">        <div class="input-group '+eClass+'" id="output-factoid-address-container">            <pre><input id="output-factoid-address" type="text" name="output1" class="input-group-field percent95" placeholder="Type '+str+' address"></pre>        <a id="addressbook-button" data-toggle="addressbook" class="input-group-button button input-group-field" id="addressbook" value="'+
counter+'"><i class="fa fa-book"></i></a>        </div>    </div>    <div class="small-10 medium-4 large-3 columns">        <div class="input-group">            <input id="output-factoid-amount" type="text" class="input-group-field" name="output1-num" placeholder="Amount of '+PageToken+'">            <span class="input-group-label">'+PageTokenABR+'</span>        </div>    </div>    <div class="small-2 medium-1 columns">            <a id="remove-new-output" class="button expanded newMinus">&nbsp;</a>    </div></div>');
counter+=1}$("#append-new-output").click(function(){$(this).hasClass("disabled-input")||addNewOutputAddress("",!0)});
function addNewInputAddress(a,b){eClass="";b&&(eClass="input-group-error");$("#all-inputs").append('<div class="row single-input-'+counter+'" id="single-input">    <div class="small-12 medium-7 large-8 columns">        <div class="input-group '+eClass+'" id="input-factoid-address-container">        <pre><input id="input-factoid-address" type="text" name="input1" class="input-group-field percent95 disabled-input" placeholder="Choose factoid address" disabled></pre>        <a id="addressbook-button" data-toggle="fee-addressbook" class="input-group-button button input-group-field" id="addressbook" value="'+counter+
'"><i class="fa fa-book"></i></a>        </div>    </div>    <div class="small-10 medium-4 large-3 columns">        <div class="input-group">            <input id="input-factoid-amount" type="text" class="input-group-field" name="output1-num" placeholder="Amount of factoids">            <span class="input-group-label">FCT</span>        </div>    </div>    <div class="small-2 medium-1 columns">            <a id="remove-new-input" class="button expanded newMinus">&nbsp;</a>    </div></div>');counter+=
1}$("#append-new-input").click(function(){$(this).hasClass("disabled-input")||addNewInputAddress("",!0)});$("#all-outputs").on("click","#remove-new-output",function(){jQuery(this).parent().parent().remove()});$("#all-inputs").on("click","#remove-new-input",function(){jQuery(this).parent().parent().remove()});
$("#all-outputs").on("keypress","#output-factoid-amount",function(a){if("FCT"==PageTokenABR){var b=$(this);b.val(b.val().replace(/[^0-9\.]/g,""));46==a.which&&b.val().indexOf(".");decSplit=$(this).val().split(".");2<decSplit.length&&a.preventDefault()}else b=$(this),b.val(b.val().replace(/[^0-9\.]/g,"")),(48>a.which||57<a.which)&&8!=a.which&&a.preventDefault()});$("#all-outputs").on("change","#output-factoid-amount",function(){});$("#all-outputs").on("click","#output-factoid-address-container",function(){$(this).removeClass("input-group-error")});
$("#all-inputs").on("click","#input-factoid-amount",function(){$(this).removeClass("input-group-error")});$("#all-outputs").on("click","#output-factoid-amount",function(){$(this).removeClass("input-group-error")});$("#make-entire-transaction").on("click",function(){Input&&("1"==$(this).attr("value")?$("#sign-transaction").prop("checked")?MakeTransaction(!1):MakeTransaction(!0):MakeTransaction(!0))});
function MakeTransaction(a){transObject=getTransactionObject(!0);a||(transObject.TransType="nosig");null!=transObject&&(j=JSON.stringify(transObject),postRequest("make-transaction",j,function(a){obj=JSON.parse(a);"none"==obj.Error?(disableInput(),ShowNewButtons(),totalInput=obj.Content.Total/1E8,feeFact=obj.Content.Fee/1E8,total=totalInput+feeFact,$("#transaction-total").attr("value",total),$("#transaction-fee").attr("value",feeFact),importexport?(setExportDownload(obj.Content.Json),SetGeneralSuccess('Click "Export Transaction" to download, or go back to editing it')):
SetGeneralSuccess('Click "Send Transaction" to send, or go back to editing it')):SetGeneralError("Error: "+obj.Error)}))}function setExportDownload(a){obj=JSON.parse(a);fileExt=parseInt(obj.millitimestamp);$("#export-transaction").click(function(){$(this).attr("href","data:text/plain;charset=UTF-8,"+encodeURIComponent(a));$(this).attr("download","Exported-"+fileExt)})}$("#send-entire-transaction").on("click",function(){SendTransaction()});
function getTransactionObject(a){var b={TransType:PageTransType,OutputAddresses:[],OutputAmounts:[],InputAddresses:[],InputAmounts:[],FeeAddress:""};errMessage="";feeErr=amtErr=faErr=!1;$("#all-outputs #single-output").each(function(){err=!1;add=$(this).find("#output-factoid-address").val();add.startsWith(AddressPrefix)||($(this).find("#output-factoid-address-container").addClass("input-group-error"),err=faErr=!0);amt=$(this).find("#output-factoid-amount").val();if(0==amt||void 0==amt)$(this).find("#output-factoid-amount").addClass("input-group-error"),
err=amtErr=!0;b.OutputAddresses.push(add);b.OutputAmounts.push(amt)});a&&!$("#coin-control").hasClass("coin-control")&&($("#all-inputs #single-input").each(function(){err=!1;add=$(this).find("#input-factoid-address").val();add.startsWith("FA")||($(this).find("#input-factoid-address-container").addClass("input-group-error"),err=faErr=!0);amt=$(this).find("#input-factoid-amount").val();if(0==amt||void 0==amt)$(this).find("#input-factoid-amount").addClass("input-group-error"),err=amtErr=!0;b.InputAddresses.push(add);
b.InputAmounts.push(amt)}),b.FeeAddress=$("#fee-factoid-address").val(),30>b.FeeAddress.length&&($("#fee-factoid-address").addClass("input-group-error"),err=feeErr=!0),b.TransType="custom");return err?(faErr&&(errMessage+="Addresses must start with '"+AddressPrefix+"'. "),amtErr&&(errMessage+="Amounts should not be 0. "),feeErr&&(errMessage+="Fee Address must be given. "),SetGeneralError("Error(s): "+errMessage),null):b}$("#needed-input-button").on("click",function(a){Input?GetNeededInput():a.preventDefault()});
function GetNeededInput(){transObject=getTransactionObject(!1);null!=transObject&&(j=JSON.stringify(transObject),postRequest("get-needed-input",j,function(a){obj=JSON.parse(a);"none"==obj.Error?($("#input-needed-amount").val(FCTNormalize(obj.Content)),HideMessages()):SetGeneralError("Error: "+obj.Error)}))}
function SendTransaction(){transObject=getTransactionObject(!0);null!=transObject&&(j=JSON.stringify(transObject),postRequest("send-transaction",j,function(a){obj=JSON.parse(a);"none"==obj.Error?(disableInput(),HideNewButtons(),SetGeneralSuccess("Transaction Sent, transaction ID: "+obj.Content),ShowNewTransaction()):(enableInput(),HideNewButtons(),$("#transaction-fee").attr("value","----"),$("#transaction-total").attr("value","----"),SetGeneralError("Error: "+obj.Error))}))}
$("#edit-transaction").on("click",function(){enableInput();HideNewButtons();$("#transaction-fee").attr("value","----");$("#transaction-total").attr("value","----");HideMessages()});$(window).load(function(){LoadAddresses();$("#coin-control").hasClass("coin-control")&&$("#fee-address-input").css("display","none")});
function LoadAddresses(){resp=getRequest("addresses",function(a){obj=JSON.parse(a);null!=obj.FactoidAddresses.List&&obj.FactoidAddresses.List.forEach(function(a){$("#fee-addresses-reveal").append(factoidAddressRadio(a,"fee-address"))});"FCT"==PageTokenABR?(null!=obj.FactoidAddresses.List&&obj.FactoidAddresses.List.forEach(function(a){$("#addresses-reveal").append(factoidAddressRadio(a,"address"))}),null!=obj.ExternalAddresses.List&&obj.ExternalAddresses.List.forEach(function(a){a.Address.startsWith("FA")&&
$("#addresses-reveal").append(factoidAddressRadio(a,"address"))})):(null!=obj.EntryCreditAddresses.List&&obj.EntryCreditAddresses.List.forEach(function(a){$("#addresses-reveal").append(factoidECRadio(a,"address"))}),null!=obj.ExternalAddresses.List&&obj.ExternalAddresses.List.forEach(function(a){a.Address.startsWith("EC")&&$("#addresses-reveal").append(factoidECRadio(a,"address"))}))})}
function factoidAddressRadio(a,b){return'<pre>  <input type="radio" name="'+b+'" id="address" value="'+a.Address+'"> <span id="address-name" name="'+a.Name+'">'+a.Name+"</span></pre><br />"}$("#addresses-reveal").on("mouseover","#address-name",function(){$(this).css("font-size","90%");$(this).text($(this).parent().find("#address").val())});$("#addresses-reveal").on("mouseout","#address-name",function(){$(this).text($(this).attr("name"));$(this).css("font-size","100%")});
function factoidECRadio(a,b){return'<pre>  <input type="radio" name="address" id="address" value="'+a.Address+'"> <span id="address-name" name="'+a.Name+'">'+a.Name+"</span></pre> <br />"}$("#addresses-reveal").on("mouseover","#address-name",function(){$(this).css("font-size","90%");$(this).text($(this).parent().find("#address").val())});$("#addresses-reveal").on("mouseout","#address-name",function(){$(this).text($(this).attr("name"));$(this).css("font-size","100%")});done=!1;
$("#addresses-reveal-button").on("click",function(){newAddress=$("input[name='address']:checked").val();void 0!=newAddress&&($(".single-output-"+toChange+" #output-factoid-address").val(newAddress),$(".single-output-"+toChange+" #output-factoid-address-container").removeClass("input-group-error"))});toChange="-1";$("#all-outputs").on("click","#addressbook-button",function(){toChange=$(this).attr("value");$("input[type=radio]").attr("checked",!1)});
$("#all-inputs").on("click","#addressbook-button",function(){toChange=$(this).attr("value");$("input[type=radio]").attr("checked",!1)});$("#fee-address-input").on("click","#addressbook-button",function(){toChange=$(this).attr("value");$("input[type=radio]").attr("checked",!1)});
$("#fee-addresses-reveal-button").on("click",function(){newAddress=$("input[name='fee-address']:checked").val();void 0!=newAddress&&("-1"==toChange?($("#fee-factoid-address").val(newAddress),$("#fee-factoid-address").removeClass("input-group-error")):($(".single-input-"+toChange+" #input-factoid-address").val(newAddress),$(".single-input-"+toChange+" #input-factoid-address-container").removeClass("input-group-error")))});function HideNewButtons(){$("#second-stage-buttons").slideUp(100)}
function ShowNewButtons(){$("#second-stage-buttons").slideDown(100)}function ShowNewTransaction(){$("#new-transaction").slideDown(100)}$("#new-transaction").on("click",function(){location.reload()});Input=!0;
function disableInput(){Input=!1;$(".input-group").each(function(){$(this).addClass("disabled-input");$(this).prop("disabled",!0)});$(".input-group-field").each(function(){$(this).addClass("disabled-input");$(this).prop("disabled",!0)});$("#needed-input-button").addClass("disabled-input");$("#needed-input-button").prop("disabled",!0);$("#addressbook-button").addClass("disabled-input");$("#addressbook-button").prop("disabled",!0);$("#make-entire-transaction").addClass("disabled-input");$("#make-entire-transaction").prop("disabled",
!0);$("#first-stage-buttons").slideUp(100)}
function enableInput(){Input=!0;$(".input-group").each(function(){$(this).removeClass("disabled-input");$(this).prop("disabled",!1)});$(".input-group-field").each(function(){$(this).removeClass("disabled-input");$(this).prop("disabled",!1)});$("#transaction-fee").prop("disabled",!0);$("#transaction-total").prop("disabled",!0);$("#needed-input-button").removeClass("disabled-input");$("#needed-input-button").prop("disabled",!1);$("#addressbook-button").removeClass("disabled-input");$("#addressbook-button").prop("disabled",
!1);$("#make-entire-transaction").removeClass("disabled-input");$("#make-entire-transaction").prop("disabled",!1);$("#first-stage-buttons").slideDown(100);keepFeeDisabled()}function keepFeeDisabled(){$("#input-factoid-address").prop("disabled",!0);$("#input-factoid-address").addClass("disabled-input")}function HideMessages(){$("#error-zone").slideUp(100);$("#success-zone").slideUp(100)}
$("#import-file").on("click",function(){(input=document.getElementById("uploaded-file"))?input.files?input.files[0]?(file=input.files[0],fr=new FileReader,fr.onload=receivedText,fr.readAsDataURL(file)):alert("Please select a file before clicking 'Load'"):alert("This browser doesn't seem to support the `files` property of file inputs."):alert("Um, couldn't find the fileinput element.")});function receivedText(){console.log(fr.result)};
