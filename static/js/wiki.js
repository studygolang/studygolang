(function(){
	SG.Wiki = function() {}
	SG.Wiki.prototype = new SG.Publisher();
	SG.Wiki.prototype.parseDesc = function(){
		var markdownString = $('.page .content').text();
		// 配置 marked 语法高亮
		marked = SG.markSettingNoHightlight();

		var contentHtml = marked(markdownString);
		contentHtml = SG.replaceCodeChar(contentHtml);

		$('.page .content').html(contentHtml);
	}

	jQuery(document).ready(function($) {
		// 发布 Wiki
		$('#submit').on('click', function(evt){
			evt.preventDefault();
			var validator = $('.validate-form').validate();
			if (!validator.form()) {
				return false;
			}

			var Wiki = new SG.Wiki()
			Wiki.publish(this);
		});

		$(document).keypress(function(evt){
			if (evt.ctrlKey && (evt.which == 10 || evt.which == 13)) {
				$('#submit').click();
			}
		});
	});
}).call(this);
