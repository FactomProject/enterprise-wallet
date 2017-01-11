var $jscomp={scope:{},checkStringArgs:function(a,c,b){if(null==a)throw new TypeError("The 'this' value for String.prototype."+b+" must not be null or undefined");if(c instanceof RegExp)throw new TypeError("First argument to String.prototype."+b+" must not be a regular expression");return a+""}};
$jscomp.defineProperty="function"==typeof Object.defineProperties?Object.defineProperty:function(a,c,b){if(b.get||b.set)throw new TypeError("ES3 does not support getters and setters.");a!=Array.prototype&&a!=Object.prototype&&(a[c]=b.value)};$jscomp.getGlobal=function(a){return"undefined"!=typeof window&&window===a?a:"undefined"!=typeof global&&null!=global?global:a};$jscomp.global=$jscomp.getGlobal(this);
$jscomp.polyfill=function(a,c,b,d){if(c){b=$jscomp.global;a=a.split(".");for(d=0;d<a.length-1;d++){var e=a[d];e in b||(b[e]={});b=b[e]}a=a[a.length-1];d=b[a];c=c(d);c!=d&&null!=c&&$jscomp.defineProperty(b,a,{configurable:!0,writable:!0,value:c})}};
$jscomp.polyfill("String.prototype.startsWith",function(a){return a?a:function(a,b){var c=$jscomp.checkStringArgs(this,a,"startsWith");a+="";for(var e=c.length,f=a.length,h=Math.max(0,Math.min(b|0,c.length)),g=0;g<f&&h<e;)if(c[h++]!=a[g++])return!1;return g>=f}},"es6-impl","es3");$jscomp.findInternal=function(a,c,b){a instanceof String&&(a=String(a));for(var d=a.length,e=0;e<d;e++){var f=a[e];if(c.call(b,f,e,a))return{i:e,v:f}}return{i:-1,v:void 0}};
$jscomp.polyfill("Array.prototype.find",function(a){return a?a:function(a,b){return $jscomp.findInternal(this,a,b).v}},"es6-impl","es3");
function LoadAddresses(){resp=getRequest("addresses",function(a){obj=JSON.parse(a);null!=obj.FactoidAddresses.List&&obj.FactoidAddresses.List.forEach(function(a){$("#factoid-addresses-table tbody").append(addressTableRow(a,"factoid"))});null!=obj.EntryCreditAddresses.List&&obj.EntryCreditAddresses.List.forEach(function(a){$("#credit-addresses-table tbody").append(addressTableRow(a,"entry-credits"))});null!=obj.ExternalAddresses.List&&obj.ExternalAddresses.List.forEach(function(a){$("#external-addresses-table tbody").append(addressTableRow(a,
"external"))});sortNames(!0)})}
function addressTableRow(a,c,b){a.Address.startsWith("FA")?(token=" FCT",a.Balance=Number(a.Balance.toFixed(4))):token=" EC";b='<small><span id="star" class="fa fa-star-o" aria-hidden="true" value="0"></span></small>';a.Seeded&&(b='<small><span id="star" class="fa fa-star" aria-hidden="true" value="1"></span></small>');"external"==c&&(b="");shortAddr=a.Address;return'<tr><td><a nav-click="true" href="receive-factoids?address='+a.Address+"&name="+a.Name+'"><i class="qr"><img src="img/icon_qr.svg" class="svg"></i></a></td><td><span id="name">'+a.Name+
'</span> <a nav-click="true" href="edit-address-'+c+"?address="+a.Address+"&name="+a.Name+'"><i class="edit"><img src="img/icon_edit.svg" class="svg"></i></a></td><td><pre>'+b+" "+shortAddr+'</pre></td><td><span id="balance">'+a.Balance+"</span>"+token+"</td></tr>"}function ResetNotMe(a){0!=a&&resetNameSort();1!=a&&resetBalancesSort();2!=a&&resetSeededSort()}
$("table").on("click","#sort-names",function(a){ResetNotMe(0);$("#sort-names-icon").hasClass("fa-sort")?($("th #sort-names-icon").removeClass("fa-sort"),$("th #sort-names-icon").addClass("fa-sort-desc"),sortNames(!0)):$("#sort-names-icon").hasClass("fa-sort-asc")?($("th #sort-names-icon").removeClass("fa-sort-asc"),$("th #sort-names-icon").addClass("fa-sort-desc"),sortNames(!0)):$("#sort-names-icon").hasClass("fa-sort-desc")&&($("th #sort-names-icon").removeClass("fa-sort-desc"),$("th #sort-names-icon").addClass("fa-sort-asc"),
sortNames(!1))});function resetNameSort(){$("th #sort-names-icon").removeClass("fa-sort-asc");$("th #sort-names-icon").removeClass("fa-sort-desc");$("th #sort-names-icon").addClass("fa-sort")}
function sortNames(a){array=$("#factoid-addresses-table tbody tr").get();valArray=$("#factoid-addresses-table tbody tr").find("#name").get();array=generalSort(stringLessThan,array,valArray,a,a);$("#factoid-addresses-table tbody").html(array);array=$("#credit-addresses-table tbody tr").get();valArray=$("#credit-addresses-table tbody tr").find("#name").get();array=generalSort(stringLessThan,array,valArray,a,a);$("#credit-addresses-table tbody").html(array);array=$("#external-addresses-table tbody tr").get();
valArray=$("#external-addresses-table tbody tr").find("#name").get();array=generalSort(stringLessThan,array,valArray,a,a);$("#external-addresses-table tbody").html(array)}
$("table").on("click","#sort-balances",function(a){ResetNotMe(1);$("#sort-balances-icon").hasClass("fa-sort-amount")?($("th #sort-balances-icon").removeClass("fa-sort-amount"),$("th #sort-balances-icon").addClass("fa-sort-amount-desc"),sortBalances(!1)):$("#sort-balances-icon").hasClass("fa-sort-amount-asc")?($("th #sort-balances-icon").removeClass("fa-sort-amount-asc"),$("th #sort-balances-icon").addClass("fa-sort-amount-desc"),sortBalances(!1)):$("#sort-balances-icon").hasClass("fa-sort-amount-desc")&&
($("th #sort-balances-icon").removeClass("fa-sort-amount-desc"),$("th #sort-balances-icon").addClass("fa-sort-amount-asc"),sortBalances(!0))});function resetBalancesSort(){$("th #sort-balances-icon").removeClass("fa-sort-amount-asc");$("th #sort-balances-icon").removeClass("fa-sort-amount-desc");$("th #sort-balances-icon").addClass("fa-sort-amount")}
function sortBalances(a){array=$("#factoid-addresses-table tbody tr").get();valArray=$("#factoid-addresses-table tbody tr").find("#balance").get();array=generalSort(isLessThan,array,valArray,a,a);$("#factoid-addresses-table tbody").html(array);array=$("#credit-addresses-table tbody tr").get();valArray=$("#credit-addresses-table tbody tr").find("#balance").get();array=generalSort(isLessThan,array,valArray,a,a);$("#credit-addresses-table tbody").html(array);array=$("#external-addresses-table tbody tr").get();
valArray=$("#external-addresses-table tbody tr").find("#balance").get();array=generalSort(isLessThan,array,valArray,a,a);$("#external-addresses-table tbody").html(array)}
$("table").on("click","#sort-seeded",function(a){ResetNotMe(2);$("#sort-seeded-icon").hasClass("fa-sort")?($("th #sort-seeded-icon").removeClass("fa-sort"),$("th #sort-seeded-icon").addClass("fa-sort-desc"),sortSeeded(!1)):$("#sort-seeded-icon").hasClass("fa-sort-asc")?($("th #sort-seeded-icon").removeClass("fa-sort-asc"),$("th #sort-seeded-icon").addClass("fa-sort-desc"),sortSeeded(!1)):$("#sort-seeded-icon").hasClass("fa-sort-desc")&&($("th #sort-seeded-icon").removeClass("fa-sort-desc"),$("th #sort-seeded-icon").addClass("fa-sort-asc"),
sortSeeded(!0))});function resetSeededSort(){$("th #sort-seeded-icon").removeClass("fa-sort-asc");$("th #sort-seeded-icon").removeClass("fa-sort-desc");$("th #sort-seeded-icon").addClass("fa-sort")}
function sortSeeded(a){array=$("#factoid-addresses-table tbody tr").get();valArray=$("#factoid-addresses-table tbody tr").get();array=generalSort(starLessThan,array,valArray,a,!0);$("#factoid-addresses-table tbody").html(array);array=$("#credit-addresses-table tbody tr").get();valArray=$("#credit-addresses-table tbody tr").get();array=generalSort(starLessThan,array,valArray,a,!0);$("#credit-addresses-table tbody").html(array)}
function generalSort(a,c,b,d,e){peerLen=b.length;for(index=0;index<peerLen;index++){tmpVal=b[index];tmp=c[index];if(1==d||e)for(j=index-1;-1<j&&!a(b[j],tmpVal,d);j--)b[j+1]=b[j],c[j+1]=c[j];else for(j=index-1;-1<j&&a(b[j],tmpVal,d);j--)b[j+1]=b[j],c[j+1]=c[j];b[j+1]=tmpVal;c[j+1]=tmp}return c}
function starLessThan(a,c,b){as=$(a).find("#star");bs=$(c).find("#star");aN=$(as).attr("value");bN=$(bs).attr("value");b||(aN=flipInt(aN),bN=flipInt(bN));aM=aN+$(a).find("#name").text();bM=bN+$(c).find("#name").text();return String(aM)<String(bM)?1:0}function flipInt(a){return 1==Number(a)?0:1}function stringLessThan(a,c,b){return String(a.innerText)<String(c.innerText)?1:0}function isLessThan(a,c,b){return Number(a.innerText)<Number(c.innerText)?1:0};
