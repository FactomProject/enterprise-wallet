
$(document).ready(function () {
    fixUp()
});

function fixUp() {
    // Convert any img.svg to inline svg
    $('img.svg').each(function () {
        var $img = $(this),
            imgID = $img.attr('id'),
            imgClass = $img.attr('class'),
            imgURL = $img.attr('src');
        $.get(imgURL, function (data) {
            var $svg = $(data).find('svg');
            if (typeof imgID !== 'undefined') {
                $svg = $svg.attr('id', imgID);
            }
            if (typeof imgClass !== 'undefined') {
                $svg = $svg.attr('class', imgClass + ' replaced-svg');
            }
            $svg = $svg.removeAttr('xmlns:a');
            $img.replaceWith($svg);
        }, 'xml');
    });
    
    // CTA click
    $('.newCTA a.transaction').click(function(){
        $(this).parent().toggleClass('active');
        return false;
    });
    $(document).mouseup(function(e) {
        var subject = $('.newCTA.active');
        if(e.target.id != subject.attr('id') && !subject.has(e.target).length) {
            subject.removeClass('active');
        }
    });
}

function reload_js(src) {
    $('script[src="' + src + '"]').remove();
    $('<script>').attr('src', src).appendTo('head');
}
// Dynamic page loading
$(function() {
    if(Modernizr.history){
        var newHash      = "",
            $mainContent = $("#dynamic-content"),
            $pageWrap    = $("body"),
            baseHeight   = 0,
            $el;
            
        $pageWrap.height($pageWrap.height());
        baseHeight = $pageWrap.height() - $mainContent.height();
        
        $("body").delegate("a[nav-click='true']", "click", function() {
            if($(this).attr("nav-click") != "true") {
                return
            }
            _link = $(this).attr("href");
            history.pushState(null, null, _link);
            loadContent(_link);
            return false;
        });

        function loadContent(href){
            $mainContent
                    .find("#guts")
                    .fadeOut(200, function() {
                        $mainContent.hide().load(href + " #guts", function() {
                            $mainContent.fadeIn(200, function() {
                                $pageWrap.animate({
                                    height: baseHeight + $mainContent.height() + "px"
                                });
                            });
                            $("#nav-list [class='active'").removeClass("active")
                            console.log(href);
                            reload_js('js/all.js');
                            switch(href) {
                                case "/":
                                    ChangeNav("transactions", 1)
                                    //loadScript("index")
                                    LoadTransactions()
                                    break;
                                case "AddressBook":
                                    ChangeNav("address-book", 2)
                                    //loadScript("addressbook")
                                    LoadAddresses()
                                    break;
                                case "Settings":
                                    console.log("what?")
                                    ChangeNav("settings", 3)
                                    break;
                                case "send-factoids":
                                    ChangeNav("send-factoids", 1)
                                    LoadAddressesSendConvert()
                                    break;
                                case "create-entry-credits":
                                    ChangeNav("create-entry-credits", 1)
                                    LoadAddressesSendConvert()
                                    break;
                                case "import-export-transaction":
                                    ChangeNav("send-factoids", 1)
                                    LoadAddressesSendConvert()
                                    break;
                                case "new-address":
                                    ChangeNav("send-factoids", 2)
                                    break;
                                case "notFound":
                                    ChangeNav("notFound", 1)
                                    break;
                                case "receive-factoids":
                                    ChangeNav("receive-factoids", 2)
                                    LoadRecAddresses()
                                    break;
                                default:
                                    ChangeNav("", 1)
                                    break;
                            }
                        });
                    });
        }
        
        $(window).bind('popstate', function(){
           _link = location.pathname.replace(/^.*[\\\/]/, ''); //get filename only
           loadContent(_link);
        });

    } // otherwise, history is not supported, so nothing fancy here.    
});

function ChangeNav(mainClass, activeWindow) {
    $("main").removeClass()

    $("main").addClass(mainClass)
    if(activeWindow == 1) {
        $("#transactions-nav").addClass("active")
        console.log("1")
    } else if(activeWindow == 2) {
        $("#address-book-nav").addClass("active")
        console.log("2")
    } else {
        $("#settings-nav").addClass("active")
        console.log("3")
    }
    fixUp();
    $(document).foundation();
}

function loadScript(script) {
    $.getScript( "js/" + script + ".js", function( data, textStatus, jqxhr ) {
      console.log( data ); // Data returned
      console.log( textStatus ); // Success
      console.log( jqxhr.status ); // 200
      console.log( "Load was performed." );
    });
}