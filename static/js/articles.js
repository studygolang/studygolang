// 文章相关js功能
(function(){
	SG.Articles = function(){}
	SG.Articles.prototype = new SG.Publisher();
	SG.Articles.prototype.parseContent = function(selector) {
		var markdownString = selector.text();
		// 配置 marked 语法高亮
		marked = SG.markSetting();

		var contentHtml = marked(markdownString);
		contentHtml = SG.replaceCodeChar(contentHtml);
		
		selector.html(contentHtml);

		// emoji 表情解析
		emojify.run(selector.get(0));
	}

	jQuery(document).ready(function($) {
		$('#submit').on('click', function(evt){
			evt.preventDefault();
			var validator = $('.validate-form').validate();
			if (!validator.form()) {
				return false;
			}

			$('#myeditor').text(CKEDITOR.instances.myeditor.getData());
			if (window.localStorage) {
				localStorage.removeItem('autosaveKey');
			}

			$('#txt').text(CKEDITOR.instances.myeditor.document.getBody().getText());

			var articles = new SG.Articles();
			articles.publish(this);
		});

		$(document).keypress(function(evt){
			if (evt.ctrlKey && (evt.which == 10 || evt.which == 13)) {
				$('#submit').click();
			}
		});
	});
}).call(this)
