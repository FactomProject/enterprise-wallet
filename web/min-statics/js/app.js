var $jscomp={scope:{},findInternal:function(a,b,c){a instanceof String&&(a=String(a));for(var e=a.length,d=0;d<e;d++){var f=a[d];if(b.call(c,f,d,a))return{i:d,v:f}}return{i:-1,v:void 0}}};$jscomp.defineProperty="function"==typeof Object.defineProperties?Object.defineProperty:function(a,b,c){if(c.get||c.set)throw new TypeError("ES3 does not support getters and setters.");a!=Array.prototype&&a!=Object.prototype&&(a[b]=c.value)};
$jscomp.getGlobal=function(a){return"undefined"!=typeof window&&window===a?a:"undefined"!=typeof global&&null!=global?global:a};$jscomp.global=$jscomp.getGlobal(this);$jscomp.polyfill=function(a,b,c,e){if(b){c=$jscomp.global;a=a.split(".");for(e=0;e<a.length-1;e++){var d=a[e];d in c||(c[d]={});c=c[d]}a=a[a.length-1];e=c[a];b=b(e);b!=e&&null!=b&&$jscomp.defineProperty(c,a,{configurable:!0,writable:!0,value:b})}};
$jscomp.polyfill("Array.prototype.find",function(a){return a?a:function(a,c){return $jscomp.findInternal(this,a,c).v}},"es6-impl","es3");$(document).ready(function(){fixUp()});
function fixUp(){$("img.svg").each(function(){var a=$(this),b=a.attr("id"),c=a.attr("class"),e=a.attr("src");$.get(e,function(d){d=$(d).find("svg");"undefined"!==typeof b&&(d=d.attr("id",b));"undefined"!==typeof c&&(d=d.attr("class",c+" replaced-svg"));d=d.removeAttr("xmlns:a");a.replaceWith(d)},"xml")});$(".newCTA a.transaction").click(function(){$(this).parent().toggleClass("active");return!1});$(document).mouseup(function(a){var b=$(".newCTA.active");a.target.id==b.attr("id")||b.has(a.target).length||
b.removeClass("active")})}function reload_js(a){$('script[src="'+a+'"]').remove();$("<script>").attr("src",a).appendTo("head")}
$(function(){if(Modernizr.history){var a=function(a){b.find("#guts").fadeOut(200,function(){b.hide().load(a+" #guts",function(){b.fadeIn(200,function(){c.animate({height:e+b.height()+"px"})});$("#nav-list [class='active'").removeClass("active");console.log(a);reload_js("js/all.js");switch(a){case "/":ChangeNav("transactions",1);LoadTransactions();break;case "AddressBook":ChangeNav("address-book",2);LoadAddresses();break;case "Settings":ChangeNav("settings",3);break;case "send-factoids":ChangeNav("send-factoids",
1);LoadAddressesSendConvert();break;case "create-entry-credits":ChangeNav("create-entry-credits",1);LoadAddressesSendConvert();break;case "import-export-transaction":ChangeNav("send-factoids",1);LoadAddressesSendConvert();break;case "new-address":ChangeNav("send-factoids",2);break;case "notFound":ChangeNav("notFound",1);break;case "receive-factoids":ChangeNav("receive-factoids",2);LoadRecAddresses();break;default:0==a.indexOf("receive-factoids?address")?(ChangeNav("receive-factoids",2),LoadFixedAddress()):
0==a.indexOf("edit-address")?(ChangeNav("receive-factoids",2),GetDefaultData()):ChangeNav("",1)}})})},b=$("#dynamic-content"),c=$("body"),e=0;c.height(c.height());e=c.height()-b.height();$("body").delegate("a[nav-click='true']","click",function(){if("true"==$(this).attr("nav-click"))return _link=$(this).attr("href"),history.pushState(null,null,_link),a(_link),!1});$(window).bind("popstate",function(){_link=location.pathname.replace(/^.*[\\\/]/,"");a(_link)})}});
function ChangeNav(a,b){$("main").removeClass();$("main").addClass(a);1==b?$("#transactions-nav").addClass("active"):2==b?$("#address-book-nav").addClass("active"):$("#settings-nav").addClass("active");fixUp();$(document).foundation()}function loadScript(a){$.getScript("js/"+a+".js",function(a,c,e){console.log(a);console.log(c);console.log(e.status);console.log("Load was performed.")})};
