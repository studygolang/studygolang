jQuery(function($){
	$(document).delegate(".js-checkbox", "click", function(){
		var parent = $(this).closest("tr");
		var flag = $(this).is(":checked");
		$("option", parent).each(function(){
			this.selected = flag;
		});
	});

	$(".js-select-mul").click(function(){
		var inputThis = $(this).closest("tr").find(".js-checkbox");
		if($(this).val()){
			inputThis[0].checked = true;
			inputThis.closest("span").addClass("checked");
		} else {
			inputThis[0].checked = false;
			inputThis.closest("span").removeClass("checked");
		}
	})

	$('#role_form').on('submit', function(evt) {
		evt.preventDefault();

		$('#loaders').show();

		// 获取被选项
		var authorities = [];
		$(".js-select-mul").each(function(){
			var value = $(this).val();
			if(value){
				authorities = authorities.concat(value);
			}
		});
		$(".js-checkbox").each(function(){
			var value = $(this).val();
			if(value && $(this).is(":checked")){
				authorities = authorities.concat(value);
			}
		});

		var roleid   = $('#roleid').val();
		var name = $('#name').val();
		if (!name) {
			$('#loaders').hide();
			jAlert("(角色名称)不能为空!", "提示");
			return false;
		}

		var params     = {
			"name" : name,
			"authorities" : authorities,
			"submit" : 1,
			"format" : "json"
		};

		if (roleid) {
			params["roleid"] = roleid;
		}

		var url = $(this).attr('action');

		$.ajax({
			"url": url,
			"type": "post",
			"data" : params,
			"dataType" : "json",
			"error" : function (jqXHR, textStatus, errorThrown) {
				$('#loaders').hide();
				var errMsg = errorThrown == 'Forbidden' ? "亲，没权限呢!" : "亲，服务器忙!"; jAlert(errMsg, "提示");
			},
			"success" : function (data) {
				$('#loaders').hide();
				if (data['ok'] != 1) {
					jAlert(data['error'], "错误");
				} else {
					jAlert("提交成功!", "提示", function() {
						if (roleid) {
							window.close();
						} else {
							$('#name').val("");
						}
					});
				}
			}
		});
		return false;
	});

})