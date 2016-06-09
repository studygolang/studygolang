jQuery(function($){
	var allmenu2 = $.parseJSON(ALL_MENU2);
	var optionHtml = '<option value="0">请选择</option>';
	var menu1 = $('#menu1').val();
	var curMenu2 = allmenu2[menu1];
	for(var i in curMenu2) {
		if (curMenu2[i][0] == menu2) {
			optionHtml += '<option value="'+curMenu2[i][0]+'" selected>'+curMenu2[i][1]+'</option>';
		} else {
			optionHtml += '<option value="'+curMenu2[i][0]+'">'+curMenu2[i][1]+'</option>';
		}
	}
	$('#menu2').html(optionHtml);
	$.uniform.update("#menu2");

	window.formSuccCallback = function(data) {}
});