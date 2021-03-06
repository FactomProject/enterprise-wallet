var HideSyncError = false

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
                                    height: "100%"
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
                                    ChangeNav("settings", 3)
                                    break;
                                case "Backup":
                                    ChangeNav("backup-main", 4, "address-book", true)
                                    LoadBackup0()
                                    break;
                                case "backup1":
                                    ChangeNav("backup-main", 4, true)
                                    break;
                                case "backup2":
                                    ChangeNav("backup-main", 4, true)
                                    break;
                                case "backup3":
                                    ChangeNav("backup-main", 4, true)
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
                                case "import-seed":
                                    ChangeNav("import-seed", 3, true)
                                    break;
                                case "success-screen-import":
                                    ChangeNav("success-screen", 3, true)
                                    $("#backup-message").hide()
                                    $("#backup-nav-item").removeClass("never-backedup")
                                    $("#backup-import-success-message").text("Your seed has been successfully restored!")
                                    break;
                                case "success-screen-backup":
                                    ChangeNav("success-screen", 4, true)
                                    $("#backup-message").hide()
                                    $("#backup-nav-item").removeClass("never-backedup")
                                    $("#backup-import-success-message").text("Your seed has been successfully backed up!")
                                    break;
                                default:
                                    if(href.indexOf("receive-factoids?address") == 0){
                                        ChangeNav("receive-factoids", 2)
                                        LoadFixedAddress()
                                    } else if(href.indexOf("edit-address") == 0){
                                        ChangeNav("receive-factoids", 2)
                                        GetDefaultData()
                                    } else {
                                        ChangeNav("", 1)
                                    }
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

function ChangeNav(mainClass, activeWindow, extraClass, hideBalances) {
    $("main").removeClass()

    if(extraClass === undefined) {
        hideBalances = false
    }

    if(extraClass === true || extraClass === false) {
        hideBalances = extraClass
        extraClass = undefined
    }

    if(hideBalances) {
        HideSyncError = true
        $(".balances").hide()
        $("#synced-indicator").hide()
    } else {
        $(".balances").show()
        HideSyncError = false
        //$("#synced-indicator").show()
    }

    $("main").addClass(mainClass)
    if(extraClass !== undefined) {
       $("main").addClass(extraClass) 
    }
    if(activeWindow == 1) {
        $("#transactions-nav").addClass("active")
    } else if(activeWindow == 2) {
        $("#address-book-nav").addClass("active")
    } else if(activeWindow == 3){
        $("#settings-nav").addClass("active")
    } else {
        $("#backup-nav").addClass("active")

    }
    fixUp();
    $('#guts').foundation();
}

function loadScript(script) {
    $.getScript( "js/" + script + ".js", function( data, textStatus, jqxhr ) {
      console.log( data ); // Data returned
      console.log( textStatus ); // Success
      console.log( jqxhr.status ); // 200
      console.log( "Load was performed." );
    });
}