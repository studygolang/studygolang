// 话题（帖子）相关js功能
(function(){
	emojify.setConfig({
		// emojify_tag_type : 'span',
		only_crawl_id    : null,
		img_dir          : 'http://hassankhan.me/emojify.js/images/emoji',
		ignored_tags     : { //忽略以下几种标签内的emoji识别
			'SCRIPT'  : 1,
			'TEXTAREA': 1,
			'A'       : 1,
			'PRE'     : 1,
			'CODE'    : 1
		}
	});
	
	SG.Topics = function(){}
	SG.Topics.prototype = new SG.Publisher();
	SG.Topics.prototype.parseContent = function(selector) {
		var markdownString = SG.preProcess(selector.text());
		// 配置 marked 语法高亮
		marked = SG.markSetting();

		var contentHtml = marked(markdownString);
		contentHtml = SG.replaceCodeChar(contentHtml);
		selector.html(contentHtml);

		// emoji 表情解析
		emojify.run(selector.get(0));
	}

	jQuery(document).ready(function($) {
		// 发布主题
		$('#submit').on('click', function(evt){
			evt.preventDefault();
			var validator = $('.validate-form').validate();
			if (!validator.form()) {
				return false;
			}

			if ($('.usernames').length != 0) {
				var usernames = SG.analyzeAt($('#content').val());
				$('.usernames').val(usernames);
			}

			var topics = new SG.Topics();
			topics.publish(this);
		});

		$(document).keypress(function(evt){
			if (evt.ctrlKey && (evt.which == 10 || evt.which == 13)) {
				$('#submit').click();
			}
		});

		SG.registerAtEvent();
	});
}).call(this)
