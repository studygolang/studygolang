// 话题（帖子）相关js功能
(function(){
	SG.Topics = function(){}
	SG.Topics.prototype = new SG.Publisher();
	SG.Topics.prototype.parseContent = function(selector) {
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

		// 注册 @ 和 表情
		var registerAtEvent = function() {
			// @ 本站其他人
			$('form textarea').atwho({
				at: "@",
				data: "/at/users.json"
			}).atwho({
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
