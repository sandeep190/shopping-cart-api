function addCard(productId,quantity) {
	if(quantity ==0){
		quantity = parseInt($('#quantity_details').val());
	}
	$.ajax({
		url: "/users/addtoCarts",
		type: "POST",
		data: JSON.stringify({ "quantity": quantity, "product_id": productId }),
		success: function (data) {
			if (data.status) {
				alert(data.message)
				//window.location.reload();
				return;
			} else {
				alert(data.message);
			}
		},
	});
}

function updateCart(productId, type) {
	$.ajax({
		url: "/users/carts/",
		type: "POST",
		data: JSON.stringify({ "type": type, "product_id": productId }),
		success: function (data) {
			if (data.status) {
				//alert(data.message)
				window.location.reload();
				return;
			} else {
				alert(data.message);
			}
		},
	});
}
