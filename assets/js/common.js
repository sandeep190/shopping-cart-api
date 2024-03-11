function addCard(productId){
	$.ajax({
		url: "/users/addtoCarts/"+productId,
		type: "POST",
		success: function (data) {
			if (data.status) {
				window.location.reload();
				return;
			} else {
				alert("Somthings went wrong!!!");
			}
		},
	});
}
