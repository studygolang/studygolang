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
	
	SG.Message = function(){}
	SG.Message.prototype = new SG.Publisher();
	SG.Message.prototype.parseContent = function(selector) {
		var markdownString = selector.text();
		// 配置 marked 语法高亮
		marked.setOptions({
			highlight: function (code) {
				code = code.replace(/&#34;/g, '"');
				code = code.replace(/&lt;/g, '<');
				code = code.replace(/&gt;/g, '>');
				return hljs.highlightAuto(code).value;
			}
		});

		selector.html(marked(markdownString));

		// emoji 表情解析
		emojify.run(selector.get(0));
	}

	jQuery(document).ready(function($) {

		// 发送消息
		$('#submit').on('click', function(evt){
			evt.preventDefault();
			var validator = $('.validate-form').validate();
			if (!validator.form()) {
				return false;
			}

			var message = new SG.Message();
			message.publish(this);
		});

		$(document).keypress(function(evt){
			if (evt.ctrlKey && (evt.which == 10 || evt.which == 13)) {
				$('#submit').click();
			}
		});

		SG.registerAtEvent(false, true);
	});
}).call(this);
