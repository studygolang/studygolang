// 资源相关js功能
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
	
	SG.Resources = function(){}
	SG.Resources.prototype = new SG.Publisher();
	SG.Resources.prototype.parseContent = function(selector) {
		var markdownString = selector.text();
		// 配置 marked 语法高亮
		marked = SG.markSettingNoHightlight();
		var contentHtml = marked(markdownString);
		contentHtml = SG.replaceCodeChar(contentHtml);
		selector.html(contentHtml);

		// emoji 表情解析
		emojify.run(selector.get(0));
	}

	jQuery(document).ready(function($) {
		// 资源形式选择
		$('.res-form input:radio').on('click', function(){
			var $form = $(this).parents('form');

			var $resUrl = $form.find('.res-url'),
				$resContent = $form.find('.res-content');
			
			if ($(this).val() == '只是链接') {
				$resUrl.show();
				$resContent.hide();
				$('#url').addClass('{required:true,url:true}');
				$('textarea#content').removeClass('required');
			} else {
				$resUrl.hide();
				$resContent.show();
				$('textarea#content').addClass('required');
				$('#url').removeClass('{required:true,url:true}');
			}
		});
		
		// 分享资源
		$('#submit').on('click', function(evt){
			evt.preventDefault();
			var validator = $('.validate-form').validate();
			if (!validator.form()) {
				return false;
			}

			/* 资源暂时不支持 @
			if ($('.usernames').length != 0) {
				var usernames = SG.analyzeAt($('#content').val());
				$('.usernames').val(usernames);
			}
			*/

			var resources = new SG.Resources();
			resources.publish(this);
		});

		$(document).keypress(function(evt){
			if (evt.ctrlKey && (evt.which == 10 || evt.which == 13)) {
				$('#submit').click();
			}
		});
		
		SG.registerAtEvent(false, true);
	});
}).call(this);
