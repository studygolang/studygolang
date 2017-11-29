// 文章相关js功能
(function(){
	SG.Articles = function(){}
	SG.Articles.prototype = new SG.Publisher();
	SG.Articles.prototype.parseContent = function(selector) {
		var markdownString = selector.text();
		marked = SG.markSettingNoHightlight();

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

			if ($('input[type=radio]:checked').val() == 0) {
				$('#content').val(CKEDITOR.instances.myeditor.getData());
				if (window.localStorage) {
					localStorage.removeItem('autosaveKey');
				}

				$('#txt').val(CKEDITOR.instances.myeditor.document.getBody().getText());
			} else {
				$('#content').val($('#markdown-content').val());
			}

			var articles = new SG.Articles();
			articles.publish(this, function(data) {
				if (typeof cacheKey == "undefined") {
					cacheKey = 'article';
				}
				purgeComposeDraft(uid, cacheKey);

				setTimeout(function(){
					if (data.id) {
						window.location.href = '/articles/'+data.id;
					} else {
						window.location.href = '/articles';
					}
				}, 1000);
			});
		});

		$(document).keypress(function(evt){
			if (evt.ctrlKey && (evt.which == 10 || evt.which == 13)) {
				$('#submit').click();
			}
		});
	});
}).call(this)
