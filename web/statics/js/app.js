
$(document).ready(function () {
   
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
    
});
