// 话题（帖子）相关js功能
(function(){
	window.Topics = {
		parseContent: function(selector) {
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
	}
}).call(this)
