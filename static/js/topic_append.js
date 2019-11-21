// 主题附言相关js功能
(function(){
	emojify.setConfig({
		// emojify_tag_type : 'span',
		only_crawl_id    : null,
		img_dir          : SG.EMOJI_DOMAIN,
		ignored_tags     : { //忽略以下几种标签内的emoji识别
			'SCRIPT'  : 1,
			'TEXTAREA': 1,
			'A'       : 1,
			'PRE'     : 1,
			'CODE'    : 1
		}
	});
	
	SG.TopicAppend = function(){}
	SG.TopicAppend.prototype = new SG.Publisher();

	jQuery(document).ready(function($) {
		// 文本框自动伸缩
		$('.need-autogrow').autoGrow();
		
		$('#content').on('keydown', function(e) {
			if (e.keyCode == 9) {
				e.preventDefault();
				var indent = "\t";
				var start = this.selectionStart;
				var end = this.selectionEnd;
				var selected = window.getSelection().toString();
				selected = indent + selected.replace(/\n/g, '\n' + indent);
				this.value = this.value.substring(0, start) + selected
						+ this.value.substring(end);
				this.setSelectionRange(start + indent.length, start
						+ selected.length);
			}
		});

		$('#content').on('input propertychange', function() {
			var markdownString = $(this).val();

			// 配置 marked 语法高亮
			marked = SG.markSettingNoHightlight();

			var contentHtml = marked(markdownString);
			contentHtml = SG.replaceCodeChar(contentHtml);
			
			$('#content-preview').html(contentHtml);

			// emoji 表情解析
			emojify.run($('#content-preview').get(0));
		});

		// 提交附言
		$('#submit').on('click', function(evt){
			evt.preventDefault();
			var validator = $('.validate-form').validate();
			if (!validator.form()) {
				return false;
			}

			// if ($('.usernames').length != 0) {
			// 	var usernames = SG.analyzeAt($('#content').val());
			// 	$('.usernames').val(usernames);
			// }

			var topicAppend = new SG.TopicAppend();
			topicAppend.publish(this);
		});

		$(document).keypress(function(evt){
			if (evt.ctrlKey && (evt.which == 10 || evt.which == 13)) {
				$('#submit').click();
			}
		});

		SG.registerAtEvent();
	});
}).call(this);
