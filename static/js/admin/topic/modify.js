jQuery(function($){
	var topicConverter = Markdown.getSanitizingConverter();

	topicConverter.hooks.chain("preBlockGamut", function (text, rbg) {
		return text.replace(/^ {0,3}""" *\n((?:.*?\n)+?) {0,3}""" *$/gm, function (whole, inner) {
			return "<blockquote>" + rbg(inner) + "</blockquote>\n";
		});
	});

	var topicEditor = new Markdown.Editor(topicConverter);

	topicEditor.run();

	// 修改提交后的回调
	window.formSuccCallback = function(data) {}

	$('.comment_content').on('blur', function(){
		var cid = $(this).data('cid'),
			content = $(this).text();
		
		$.ajax({
			"url": '/admin/community/comment/del',
			"type": "post",
			"data" : {format:'json', cid:cid, content:content},
			"dataType" : "json",
			"error" : function (jqXHR, textStatus, errorThrown) {
				var errMsg = errorThrown == 'Forbidden' ? "亲，没权限呢!" : "亲，服务器忙!";
				showToast(errMsg);
			},
			"success" : function (data) {
				if (data.ok == 1) {
					showToast("修改成功！");
				} else {
					showToast(data.error);
				}
			}
		});
		
	});

	var showToast = function(content) {
		$('#toast').cftoaster({
			content: content,
			animationTime: 500,
			showTime: 1000,
			maxWidth: 250,
			backgroundColor: '#1a1a1a',
			fontColor: '#eaeaea',
			bottomMargin: 250
		});
	}
});