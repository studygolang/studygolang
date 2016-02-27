jQuery(function($){
	var allmenu2 = $.parseJSON(ALL_MENU2);
	
	$('#menu1').on('change', function(){
		var optionHtml = '<option value="0">请选择</option>';
		
		var menu1 = $(this).val();
		if (menu1 == 0) {
			$('#menu2').html(optionHtml);
			$('#menu2').get(0).options[0].selected = 0;
			$.uniform.update("#menu2");
			return
		}

		var curMenu2 = allmenu2[menu1];
		for(var i in curMenu2) {
			optionHtml += '<option value="'+curMenu2[i][0]+'">'+curMenu2[i][1]+'</option>';
		}
		
		$('#menu2').html(optionHtml);
		
	});
});