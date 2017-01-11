FCTDecminalLength=8;function getRequest(a,c){var d=new XMLHttpRequest;d.onreadystatechange=function(){4==d.readyState&&c(d.response)};d.open("GET","/GET?request="+a,!0);d.send()}function postRequest(a,c,d){var b=new XMLHttpRequest;b.onreadystatechange=function(){4==b.readyState&&d(b.response)};var e=new FormData;e.append("request",a);e.append("json",c);b.open("POST","/POST");b.send(e)}$(window).load(function(){updateBalances()});setInterval(updateBalances,5E3);
function updateBalances(){getRequest("balances",function(a){obj=JSON.parse(a);"none"==obj.Error&&($("#ec-balance").text(obj.Content.EC),fcBal=formatFC(obj.Content.FC),$("#factoid-balance").text(fcBal[0]+"."),1<fcBal.length?$("#factoid-balance-trailing").text(fcBal[1]):$("#factoid-balance-trailing").text(0))})}function formatFC(a){dec=FCTNormalize(a);decStr=dec.toString();return decSplit=decStr.split(".")}function FCTNormalize(a){return Number((a/1E8).toFixed(FCTDecminalLength))}checkSynced();
setInterval(checkSynced,3E3);
function checkSynced(){getRequest("synced",function(a){obj=JSON.parse(a);switch(obj.Content.Stage){case 0:$("#load-message").text("Setting up...");break;case 1:$("#load-message").text("Gathering new transactions...");break;case 2:$("#load-message").text("Checking any new addresses...");break;case 3:$("#load-message").text("Sorting transactions...")}eBlockPercent=obj.Content.EntryHeight/obj.Content.LeaderHeight;eBlockPercent=HelperFunctionForPercent(eBlockPercent,50);fBlockPercent=obj.Content.FblockHeight/
obj.Content.LeaderHeight;fBlockPercent=HelperFunctionForPercent(fBlockPercent,50);percent=eBlockPercent+fBlockPercent;98<percent?$("#sync-bar").removeClass("alert"):$("#sync-bar").addClass("alert");$("#load-percent").text(percent.toFixed(2));1==obj.Content.Synced&&$("#synced-indicator").slideUp(100)})}function HelperFunctionForPercent(a,c){if(void 0==a||NaN==a)a=0;a*=c;a>c&&(a=c);return a}
function SetGeneralError(a){$("#success-zone").slideUp(100);$("#error-zone").text(a);$("#error-zone").slideDown(100)}function SetGeneralSuccess(a){$("#error-zone").slideUp(100);$("#success-zone").text(a);$("#success-zone").slideDown(100)}function saveTextAsFile(a,c){var d=new Blob([a],{type:"text/plain"}),b=document.createElement("a");b.download=c;window.URL=window.URL||window.webkitURL;b.href=window.URL.createObjectURL(d);b.style.display="none";document.body.appendChild(b);b.click()};
