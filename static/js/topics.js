// 话题（帖子）相关js功能
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
	
	SG.Topics = function(){}
	SG.Topics.prototype = new SG.Publisher();
	SG.Topics.prototype.parseContent = function(selector) {
		// 配置 marked 语法高亮
		marked = SG.markSettingNoHightlight();

		selector.each(function() {
			var markdownString = $(this).text();

			var contentHtml = marked(markdownString);
			contentHtml = SG.replaceCodeChar(contentHtml);
			
			$(this).html(contentHtml);

			// emoji 表情解析
			emojify.run(this);
		});
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
			topics.publish(this, function(data) {
				purgeComposeDraft(uid, 'topic');

				setTimeout(function(){
					if (data.tid) {
						window.location.href = '/topics/'+data.tid;
					} else {
						window.location.href = '/topics';
					}
				}, 1000);
			});
		});

		$(document).keypress(function(evt){
			if (evt.ctrlKey && (evt.which == 10 || evt.which == 13)) {
				$('#submit').click();
			}
		});

		SG.registerAtEvent();
	});
}).call(this);
