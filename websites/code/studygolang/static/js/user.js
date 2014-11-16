// 用户相关js功能
(function(){
	SG.User = function(){}
	SG.User.prototype.parseCmtContent = function(selector) {
		selector.each(function(){
			var markdownString = $(this).html();
			// 配置 marked 语法高亮
			marked = SG.markSetting();

			$(this).html(marked(markdownString));

			emojify.setConfig({
				// emojify_tag_type : 'span',
				only_crawl_id    : null,
				img_dir          : 'http://www.emoji-cheat-sheet.com/graphics/emojis',
				ignored_tags     : { //忽略以下几种标签内的emoji识别
					'SCRIPT'  : 1,
					'TEXTAREA': 1,
					'A'       : 1,
					'PRE'     : 1,
					'CODE'    : 1
				}
			});
			
			// emoji 表情解析
			emojify.run(this);
		});
	}

	jQuery(document).ready(function($) {
		new SG.User().parseCmtContent($('.recent-comments ul li .content'));
	});
}).call(this)
