// 资源相关js功能
(function(){
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

		// 注册 @ 和 表情
		var registerAtEvent = function() {
			// @ 本站其他人
			$('form textarea')/*.atwho({
				at: "@",
				data: "/at/users.json"
			})*/.atwho({
				at: ":",
				data: window.emojis,
				tpl:"<li data-value='${key}'><img src='http://www.emoji-cheat-sheet.com/graphics/emojis/${name}.png' height='20' width='20' /> ${name}</li>"
			})/*.atwho({
				at: "\\",
				data: window.twemojis,
				tpl:"<li data-value='${name}'><img src='https://twemoji.maxcdn.com/16x16/${key}.png' height='16' width='16' /> ${name}</li>"
			})*/;
		}

		registerAtEvent();
	});
}).call(this)
