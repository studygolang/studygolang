jQuery.noConflict();

jQuery(document).ready(function($) {

	function ajaxSubmitCallback(event) {
		event.preventDefault();
		var target = event.target;
		var submitHint = $(target).attr("ajax-hint");
		jConfirm(submitHint, "提示", function(answer) {
			if (!answer) {
				return false;
			}
			var action = $(target).attr("ajax-action"),
				data = 'format=json',
				id = $(target).data('id');
			if (id != "") {
				data += "&id="+id;
			}
			
			$.ajax({
				url : action,
				type : 'post',
				data : data,
				dataType : 'json',
				success : function(data) {
					if (data.ok == 1) {
						var successhint = $(target).attr("success-hint");
						if (successhint != null && successhint != ""){
							jAlert(successhint, "提示");
						}
						if ($(target).attr("callback")) {
							var callback = $(target).attr("callback");
							window[callback](target);
						} else {
							if ($(target).attr("submit-redirect")) {
								if ($(target).attr("submit-redirect") == "#") {
									location.reload();
								}
								location.href = $(target).attr("submit-redirect");
							} else {
								location.href = document.referrer;
							}
						}
					} else {
						jAlert(data.message, "提示");
					}
				}
			});
		});
		return false;
	}

	$('#query_result').on('click', 'a[data-type=ajax-submit]', ajaxSubmitCallback);

	// 列表删除操作的回调函数
	window.delCallback = function(target){
		$(target).parents('tr').remove();
	}
	
	// 表单ajax提交
	$("form[data-type=form-submit]").on('submit', function(event){
		event.preventDefault();
		var target = event.target;
		var submitHint = $(target).attr("submit-hint");
		jConfirm(submitHint, "提示", function(answer) {
			if (!answer) {
				return false;
			}

			var action = $(target).attr("submit-action");
			$.ajax({
				url : action,
				data : $(target).serialize(),
				type : 'post',
				dataType : 'json',
				success : function(data) {
					if (data.code == 0) {
						var successhint = $(target).attr("success-hint");
						if (successhint != null && successhint != ""){
							jAlert(successhint, "提示");
						}
						if ($(target).attr("submit-redirect")) {
							if ($(target).attr("submit-redirect") == "#") {
								location.reload();
							}
							location.href = $(target).attr("submit-redirect");
						} else if ($(target).attr('close')) {
							$.colorbox.close();
						} else {
							// 回退到上一个页面
							//location.href = document.referrer;
						}
					} else {
						jAlert(data.message, "提示");
					}
				}
			});
		});
	});

	var showProgress = function() {
		$('#loaders').show();
	}
	var hideProgress = function() {
		$('#loaders').hide();
	}
	
	var getParams = function() {
		var queryParams = GLOBAL_CONF['query_params'],
			params = {};
		for( var k in queryParams) {
			params[k] = $.trim($(queryParams[k]).val());
		}
		return params;
	}
	
	$('#queryform').on('submit', function(evt) {
		evt.preventDefault();
		
		var url = GLOBAL_CONF['action_query'],
			params = getParams();
		
		showProgress();

		$.ajax({
			"url": url,
			"type": "post",
			"data" : params,
			"dataType" : "html",
			"error" : function (jqXHR, textStatus, errorThrown) {
				hideProgress();
				var errMsg = errorThrown == 'Forbidden' ? "亲，没权限呢!" : "亲，服务器忙!"; jAlert(errMsg, "提示");
			},
			"success" : function (data) {
				hideProgress();
				$("#query_result").html(data);
				bindEvt(true);
			}
		});

		return false;
	});

	// 查询结果(page为0表示当前页)
	var queryResult = function(page) {
		if (!page) {
			page = $('#cur_page').val();
		}
		var params = getParams();
		params["page"] = page;
		params["limit"] = $('#limit').val();

		showProgress();

		var url = GLOBAL_CONF['action_query'];
		$.ajax({
			"url": url,
			"type": "post",
			"data" : params,
			"dataType" : "html",
			"error" : function (jqXHR, textStatus, errorThrown) {
				hideProgress();
				var errMsg = errorThrown == 'Forbidden' ? "亲，没权限呢!" : "亲，服务器忙!"; jAlert(errMsg, "提示");
			},
			"success" : function (data) {
				$("#query_result").html(data);
				hideProgress();
				bindEvt(true);
			}
		});
	}

	// bind分页及其他事件
	var bindEvt = function(needUniform) {
		// 对bind的页面样式处理
		if (needUniform) {
			$("#query_result").find('input:checkbox, input:radio, select.uniformselect').uniform();
		}

		// 分页
		$('.pagination').jqPagination({
			current_page: $('#cur_page').val(),
			max_page: $('#totalPages').val(),
			page_string: '当前页 {current_page} 共 {max_page} 页', 
			paged: function(page) {
				// do something with the page variable
				queryResult(page);
			}
		});
	};
	
	bindEvt(false);
});
