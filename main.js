var textElement;

$(document).ready(function() {
	textElement = $("#text");
	update();
});

function update() {
	$.ajax({
		url: "/update",
		method: "GET",
		dataType: "json",
	}).done(function(resp) {
		console.log(resp);
		textElement.text(resp.text);
		update();
	}).fail(function(resp) {
		alert(resp.responseText);
	});
}
