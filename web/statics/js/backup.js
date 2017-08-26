$("#backup-input-verify-button").on('click', function(){
	$("#backup-input-verify-button").addClass("backup-btn-checking")

	var seedSingle = $("#given-seed").attr("value");
	var inputedSeedSplit = new Array(12);

	$(".backup-input").each(function(){
		var i = Number($(this).attr("index"))
		inputedSeedSplit[i-1] = $(this).val()
	})

	var clear = function() {
		$("#backup-input-verify-button").removeClass("backup-btn-failed")
		$("#backup-input-verify-button").removeClass("backup-btn-checking")
		$("#backup-input-verify-button").removeClass("backup-btn-verified")
	}

	if(inputedSeedSplit.join(" ") === seedSingle) {
		clear()
		$("#backup-input-verify-button").addClass("backup-btn-verified")
	} else {
		clear()
		$("#backup-input-verify-button").addClass("backup-btn-failed")
		setTimeout(clear, 3000);
	}
})

$("#backup-html-form").on("submit",function(event){event.preventDefault()})

function LoadBackup0() {
	LoadAddresses(false)
}