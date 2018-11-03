// studygolang 全局对象（空间）
var SG = {};

SG.EMOJI_DOMAIN = 'https://cdnjs.cloudflare.com/ajax/libs/emojify.js/1.1.0/images/basic';

function goTop()
{
	$(window).scroll(function(e) {
		// 若滚动条离顶部大于100元素
		if($(window).scrollTop() > 100)
			$("#gotop").fadeIn(500);// 以1秒的间隔渐显id=gotop的元素
		else
			$("#gotop").fadeOut(500);// 以1秒的间隔渐隐id=gotop的元素
	});
};

// 通用的发布功能
SG.Publisher = function(){}
SG.Publisher.prototype = {
	publish: function(that, callback) {
		var btnTxt = $(that).text();
		$(that).text("稍等").addClass("disabled").attr({"title":'稍等',"disabled":"disabled"});

		var $form = $(that).parents('form'),
			data = $form.serialize(),
			url = $form.attr('action');

		$.ajax({
			type:"post",
			url: url,
			data: data,
			dataType: 'json',
			success: function(data){
				if(data.ok){
					$form.get(0).reset();

					if (typeof data.msg != "undefined") {
						comTip(data.msg);
					} else {
						comTip("发布成功！");
					}

					if (typeof callback != "undefined") {
						callback(data.data);
						return;
					}

					setTimeout(function(){
						var redirect = $form.data('redirect');
						if (redirect) {
							window.location.href = redirect;
						}
					}, 1000);
				} else {
					comTip(data.error);
				}
			},
			complete:function(xmlReq, textStatus){
				$(that).text(btnTxt).removeClass("disabled").removeAttr("disabled").attr({"title":btnTxt});
			},
			error:function(xmlReq, textStatus, errorThrown){
				$(that).text(btnTxt).removeClass("disabled").removeAttr("disabled").attr({"title":btnTxt});
				if (xmlReq.status == 403) {
					comTip("没有修改权限");
				}
			}
		});
	}
}

SG.replaceSpecialChar = function(str) {
	str = str.replace(/&#34;/g, '"');
	str = str.replace(/&#39;/g, "'");
	str = str.replace(/&lt;/g, '<');
	str = str.replace(/&gt;/g, '>');
	str = str.replace(/&amp;/g, '&');
	return str;
}

SG.markSetting = function() {
	var renderer = new marked.Renderer();

	// 对 html 进行处理
	renderer.html = function(html) {
		if (html.indexOf('<script') != -1) {
			return html.replace(/</g, '&lt;');
		} else if (html.indexOf('<input') != -1) {
			return html.replace(/</g, '&lt;');
		} else if (html.indexOf('<select') != -1) {
			return html.replace(/</g, '&lt;');
		} else if (html.indexOf('<textarea') != -1) {
			return html.replace(/</g, '&lt;');
		} else {
			return html;
		}
	};

	marked.setOptions({
		renderer: renderer,
		// 配置 marked 语法高亮
		highlight: function (code) {
			code = SG.replaceSpecialChar(code);
			return hljs.highlightAuto(code).value;
		}
	});

	return marked;
}

SG.markSettingNoHightlight = function() {
	var renderer = new marked.Renderer();

	// 对 html 进行处理
	renderer.html = function(html) {
		if (html.indexOf('<script') != -1) {
			return html.replace(/</g, '&lt;');
		} else if (html.indexOf('<input') != -1) {
			return html.replace(/</g, '&lt;');
		} else if (html.indexOf('<select') != -1) {
			return html.replace(/</g, '&lt;');
		} else if (html.indexOf('<textarea') != -1) {
			return html.replace(/</g, '&lt;');
		} else {
			return html;
		}
	};

	marked.setOptions({
		renderer: renderer,
		highlight: function (code) {
			code = SG.replaceSpecialChar(code);
			return code;
		}
	});

	return marked;
}

// 替换 `` 代码块中的 "<>& 等字符
SG.replaceCodeChar = function(code) {
	code = code.replace(/<code class="lang-/g, '<code class="language-');
	return code.replace(/<code>.*<\/code>/g, function(matched, index, origin) {
		return SG.replaceSpecialChar(matched);
	});
}

// marked 处理之前进行预处理
SG.preProcess = function(content) {
	// 对引用进行处理
	content = content.replace(/&gt;/g, '>');
	return content;
}

// 分析 @ 的用户
SG.analyzeAt = function(text) {
	var usernames = [];

	String(text).replace(/[^@]*@([^\s@]{4,20})\s*/g, function (match, username) {
		usernames.push(username);
	});

	return usernames;
}

// registerAtEvent
// 注册 @ 和 表情
SG.registerAtEvent = function(isAt, isEmoji, selector) {
	if (typeof isAt == "undefined") {
		isAt = true;
	}

	if (typeof isEmoji == "undefined") {
		isEmoji = true;
	}

	if (typeof selector == "undefined") {
		selector = $('form textarea');
	}

	if (isAt) {
		var cachequeryMentions = {}, itemsMentions;
		// @ 本站其他人
		selector.atwho({
			at: "@",
			tpl: "<li data-value='${atwho-at}${username}'><img src='${avatar}' height='20' width='20' /> ${username}</li>",
			search_key: "username",
			callbacks: {
				remote_filter: function (query, render_view) {
					var thisVal = query,
					self = $(this);
					if( !self.data('active') ){
						self.data('active', true);
						itemsMentions = cachequeryMentions[thisVal]
						if(typeof itemsMentions == "object"){
							render_view(itemsMentions);
						} else {
							if (self.xhr) {
								self.xhr.abort();
							}
							self.xhr = $.getJSON("/at/users",{
								term: thisVal
							}, function(data) {
								cachequeryMentions[thisVal] = data
								render_view(data);
							});
						}
						self.data('active', false);
					}
				}
			}
		});
	}

	if (isEmoji) {
		selector.atwho({
			at: ":",
			data: window.emojis,
			tpl:"<li data-value='${key}'><img src='"+SG.EMOJI_DOMAIN+"/${name}.png' height='20' width='20' /> ${name}</li>"
		});
	}
}

jQuery(document).ready(function($) {
	// timeago：100 天之内才显示 timeago
	$.timeago.settings.cutoff = 1000*60*60*24*100;

	// 历史原因，其他 js 使用了。（当时版本 timeago 不支持 cutoff）
	// time 的格式 2014-10-02 11:40:01
	SG.timeago = function(time) {
		return $.timeago(time);
	};

	$('.timeago').timeago();

	// tooltip
	$('.tool-tip').tooltip();

	// 点击回到顶部的元素
	$("#gotop").click(function(e) {
		// 以1秒的间隔返回顶部
		$('body,html').animate({scrollTop:0}, 100);
	});
	/*
	$("#gotop").mouseover(function(e) {
		$(this).css("background","url(/static/img/top.gif) no-repeat 0px 0px");
	});
	$("#gotop").mouseout(function(e) {
		$(this).css("background","url(/static/img/top.gif) no-repeat -70px 0px");
	});
	*/

	goTop();// 实现回到顶部元素的渐显与渐隐

	//全局淡入淡出提示框 comTip
	window.comTip = function(msg){
		$("<div>").addClass("comTip").text(msg).appendTo("body");
		var timer = setInterval(function(){
			if($(".comTip").width()){
				clearInterval(timer);
				var	l = ($(window).width()-$(".comTip").outerWidth())/2;
				var	t = ($(window).height()-$(".comTip").outerHeight())/2;
				t = (t<0?0:t)+$(window).scrollTop();
				$(".comTip").css({left:l,top:t}).fadeIn(500);
				setTimeout(function(){
					$(".comTip").fadeOut(1000);
				},1800)
				setTimeout(function(){
					$(".comTip").remove()
				},3000)
			}
		},500)
	}

	// 全局公用弹出层方法
	// 弹层
	window.openPop = function(popid)
	{
		if (hadPop) {
			return;
		}

		hadPop = true;
		var pop = $(popid);
		var l = ($(window).width() - pop.outerWidth())/2;
		var t = ($(window).height() - pop.outerHeight())/2;
		t = (t<0 ? 0 : t) + $(window).scrollTop();
		pop.css({left:l,top:$(window).scrollTop(),opacity:0,display:'block'}).animate({left:l,top:t,opacity:1},500);
		$("#sg-overlay").css({width:$(document).width(),height:$(document).height()}).fadeIn(300);
	}

	// 关闭弹层
	window.closePop = function()
	{
		hadPop = false;
		$(".pop").hide();
		$("#sg-overlay").fadeOut(300);
	}

	$("#sg-overlay").click(function(){closePop()});

	// 弹窗异步登录
	$('#login-pop .login-form form').on('submit', function(evt){
		evt.preventDefault();

		var username = $('#form_username').val(),
			passwd = $('#form_passwd').val();

		if (username == "") {
			$('#form_username').parent().addClass('has-error');
			return;
		}
		if (passwd == "") {
			$('#form_passwd').parent().addClass('has-error');
			return;
		}

		$.post('/account/login', $(this).serialize(), function(data){
			if (data.ok) {
				location.reload();
			} else {
				$('#login-pop .login-form .error').text(data.error).show();
			}
		});
	});

	$('#username, #passwd').on('focus', function(){$('#login-pop .login-form .error').hide();});

	// 发送喜欢(取消喜欢)
	var postLike = function(that, callback){
		if ($('#is_login_status').val() != 1) {
			openPop("#login-pop");
			return;
		}

		var objid = $(that).data('objid'),
			objtype = $(that).data('objtype'),
			likeFlag = parseInt($(that).data('flag'), 10);

		if (likeFlag) {
			likeFlag = 0;
		} else {
			likeFlag = 1;
		}

		$.post('/like/'+objid, {objtype:objtype, flag:likeFlag}, function(data){
			if (data.ok) {

				$(that).data('flag', likeFlag);

				var likeNum = parseInt($(that).children('.likenum').text(), 10);
				// 已喜欢
				if (likeFlag) {
					comTip("感谢赞！");
					$(that).attr('title', '取消赞').text('取消赞');
					likeNum++;
				} else {
					comTip("已取消赞！");
					$(that).attr('title', '赞').text('赞');
					likeNum--;
				}

				$(that).children('.likenum').text(likeNum);

				callback(likeNum, likeFlag);
			} else {
				alert(data.error);
			}
		});
	}

	// 详情页喜欢(取消喜欢)
	$('.page #content-thank a').on('click', function(evt){
		evt.preventDefault();

		var that = this;
		postLike(that, function(likeNum, likeFlag){
			// $('.page .meta .p-comment .like .likenum').text(likeNum);
		});
	});

	// 列表页直接点喜欢(取消喜欢)
	$('.article .metatag .like').on('click', function(evt){
		evt.preventDefault();

		var that = this;
		postLike(that, function(likeNum, likeFlag){
			if (likeFlag) {
				$(that).children('i').removeClass('glyphicon-heart-empty').addClass('glyphicon-heart');
			} else {
				$(that).children('i').removeClass('glyphicon-heart').addClass('glyphicon-heart-empty');
			}
		});
	});

	// 收藏(取消收藏)
	var postFavorite = function(that, callback) {

		if ($('#is_login_status').val() != 1) {
			openPop("#login-pop");
			return;
		}

		var objid = $(that).data('objid'),
			objtype = $(that).data('objtype'),
			hadCollect = parseInt($(that).data('collect'), 10);

		if (hadCollect) {
			hadCollect = 0;
		} else {
			hadCollect = 1;
		}

		$.post('/favorite/'+objid, {objtype:objtype, collect:hadCollect}, function(data){
			if (data.ok) {
				callback(hadCollect);
			} else {
				alert(data.error);
			}
		});
	};

	// 详情页收藏(取消收藏)
	$('.page .collect').on('click', function(evt){
		evt.preventDefault();

		var that = this;
		postFavorite(that, function(hadCollect){
			$('.page .collect').data('collect', hadCollect);

			if (hadCollect) {
				comTip("感谢收藏！");
				$('.page .collect').attr('title', '取消收藏').text('取消收藏');
			} else {
				$('.page .collect').attr('title', '稍后再读').text('加入收藏');
				comTip("已取消收藏！");
			}
		});
	});

	// 收藏页 取消收藏
	$('.article .metatag .collect').on('click', function(evt){
		evt.preventDefault();

		var that = this;
		postFavorite(that, function(){
			$(that).parents('article').fadeOut();
		});
	});

	window.saveComposeDraft = function(uid, keyprefix, objdata) {
		var key = keyprefix+':compose:by:' + uid;
		lscache.set(key, objdata, 525600);
		console.log('Compose draft for UID ' + uid + ' is saved');
	};

	window.loadComposeDraft = function(uid, keyprefix) {
		var key = keyprefix+":compose:by:" + uid;
		var draft = lscache.get(key);
		console.log("Loaded compose draft for UID " + uid);

		return draft;
	}

	window.purgeComposeDraft = function(uid, keyprefix) {
		var key = keyprefix+":compose:by:" + uid;
		lscache.remove(key);
		console.log("Purged compose draft for UID " + uid);
	}

	window.saveReplyDraft = function(uid, keyprefix, objid, objdata) {
		var key = keyprefix+':'+objid+':reply:by:' + uid;
		lscache.set(key, objdata, 525600);
		console.log('Reply draft for ' + keyprefix + ':' + objid + ' is saved');
	};

	window.loadReplyDraft = function(uid, keyprefix, objid) {
		var key = keyprefix+':'+objid+':reply:by:' + uid;
		var draft = lscache.get(key);
		console.log('Loaded reply draft for ' + keyprefix + ':' + objid);

		return draft;
	}

	window.purgeReplyDraft = function(uid, keyprefix, objid) {
		var key = keyprefix+':'+objid+':reply:by:' + uid;
		lscache.remove(key);
		console.log('Purged reply draft for ' + keyprefix + ':' + objid);
	}

	// 图片响应式
	setTimeout(function(){
		$('.page .content img').each(function(){
			if ($(this).hasClass('emoji')) {
				return;
			}

			if ($(this).hasClass('no-zoom')) {
				return;
			}

			$(this).addClass('img-responsive').attr('data-action', 'zoom');
		})

		$('.page .content img').on('click', function() {
			$(this).parents('.box_white').css('overflow', 'visible');
		});
	}, 1000);

	var origSrc = '';
	$('#reload-captcha').on('click', function(evt){
		evt.preventDefault();

		if (origSrc == '') {
			origSrc = $(this).attr("src");
		}
		$(this).attr("src", origSrc+"?reload=" + (new Date()).getTime());
	});

	// 表格响应式
	setTimeout(function() {
		$('.page .content table').addClass('table').wrap('<div class="table-responsive"></div>');
	}, 2000);

});

// 在线人数统计
window.WebSocket = window.WebSocket || window.MozWebSocket;
if (window.WebSocket) {
	var websocket = new WebSocket(wsUrl);

	websocket.onopen = function(evt){
		// console.log("open");
		// console.log(evt);
	}

	websocket.onclose = function(evt){
		// console.log("close");
		// console.log(evt);
	}

	websocket.onmessage = function(msgEvent){
		data = JSON.parse(msgEvent.data);
		switch (data.type) {
		case 0:
			var $badge = $('#user_message_count .badge'),
				curVal = parseInt($badge.text(), 10);
			totalVal = parseInt(data.body) + curVal;
			if (totalVal > 0) {
				$badge.addClass('badge-warning').text(totalVal);
			} else {
				$badge.removeClass('badge-warning').text(0);
			}
			break;
		case 1:
			$('#onlineusers').text(data.body.online);
			if (data.body.maxonline) {
				$('#maxonline').text(data.body.maxonline);
			}
			break;
		}
	}

	websocket.onerror = function(evt) {
		// console.log(evt);
	}
}

var hadPop = false;

$(function(){
	$(window).scroll(function() {
		// 滚动条所在位置的高度
		var totalheight = parseFloat($(window).height()) + parseFloat($(window).scrollTop());
		// 当前文档高度   小于或等于   滚动条所在位置高度  则是页面底部
		if(($(document).height()) <= totalheight) {
			if($("#is_login_status").val() != 1){
				// openPop("#login-pop");
			}
		}

		// 控制导航栏
		$('.navbar').css('position', $(window).scrollTop() > 0 ? 'fixed' : 'relative')

		if ($(window).scrollTop() > 0) {
			$('#wrapper').css('margin-top', '52px');
		} else {
			$('#wrapper').css('margin-top', '-20px');
		}
	});

	$('#login-pop .close').on('click', function() {
		closePop();
	});
});

// markdown tool bar 相关功能
(function(){
	jQuery(document).ready(function($) {
		$('form .md-toolbar .edit').on('click', function(evt){
			evt.preventDefault();
			
			$(this).addClass('cur');

			var $mdToobar = $(this).parents('.md-toolbar');
			$mdToobar.find('.preview').removeClass('cur');

			$mdToobar.nextAll('.content-preview').hide();
			$mdToobar.next().show();
		});
		
		$('form .md-toolbar .preview').on('click', function(evt){
			evt.preventDefault();

			// 配置 marked 语法高亮
			marked = SG.markSettingNoHightlight();

			$(this).addClass('cur');
			var $mdToobar = $(this).parents('.md-toolbar');
			$mdToobar.find('.edit').removeClass('cur');

			var $textarea = $mdToobar.next();
			$textarea.hide();
			var content = $textarea.val();
			var $contentPreview = $mdToobar.nextAll('.content-preview');
			$contentPreview.html(marked(content));
			$contentPreview.show();
		});

		$('form .preview_btn').on('click', function(evt) {
			evt.preventDefault();

			// 配置 marked 语法高亮
			marked = SG.markSettingNoHightlight();

			var $mdToobar = $('form .md-toolbar');
			$mdToobar.find('.preview').addClass('cur');
			$mdToobar.find('.edit').removeClass('cur');

			var $textarea = $mdToobar.next();
			$textarea.hide();
			var content = $textarea.val();
			var $contentPreview = $mdToobar.nextAll('.content-preview');
			$contentPreview.html(marked(content));
			$contentPreview.show();
		});
	});
}).call(this);

window.initPLUpload = function (options) {
	options = options || {}
	options.ele = options.ele || 'upload-img'
	options.fileUploaded = options.fileUploaded || function(file, data) {
		var $textarea = $(options.ele).parents('.md-toolbar').next().children('textarea');
		if ($textarea.length == 0) {
			$textarea = $('.main-textarea');
		}
		var text = $textarea.val();
		text += '!['+file.name+']('+data.data.url+')';
		$textarea.val(text);
	}
	
	// 实例化一个plupload上传对象
	var uploader = new plupload.Uploader({
		browse_button : options.ele, // 触发文件选择对话框的按钮，为那个元素id
		url : '/image/upload', // 服务器端的上传页面地址
		filters: {
			mime_types : [ //只允许上传图片
				{ title : "图片文件", extensions : "jpg,gif,png,bmp" }
			],
			max_file_size : '5mb', // 最大只能上传 5mb 的文件
			prevent_duplicates : true // 不允许选取重复文件
		},
		multi_selection: false,
		file_data_name: 'img'
	});

	// 在实例对象上调用init()方法进行初始化
	uploader.init();

	uploader.bind('FilesAdded',function(uploader, files){
		// 调用实例对象的start()
		uploader.start();
	});
	uploader.bind('UploadProgress',function(uploader,file){
		// 上传进度
	});
	uploader.bind('FileUploaded', function(uploader, file, responseObject) {
		if (responseObject.status == 200) {
			var data = $.parseJSON(responseObject.response);
			if (data.ok) {
				options.fileUploaded(file, data)
			} else {
				comTip("上传失败："+data.error);
			}
		} else {
			comTip("上传失败：HTTP状态码："+responseObject.status);
		}
	});
	uploader.bind('Error',function(uploader,errObject){
		comTip("上传出错了："+errObject.message);
	});

	return uploader;
}

$(function(){
	initPLUpload()
});
jQuery(document).ready(function(){
	
	$('.upload_img_single').Huploadify({
		auto: true,
		fileTypeExts: '*.png;*.jpg;*.JPG;*.bmp;*.gif',// 不限制上传文件请修改成'*.*'
		multi:false,
		fileSizeLimit: 5*1024*1024, // 大小限制
		uploader : '/image/upload', // 文件上传目标地址
		buttonText : '上传',
		fileObjName : 'img',
		showUploadedPercent:true,
		onUploadSuccess : function(file, data) {
			data = $.parseJSON(data);
			if (data.ok) {
				var url = data.data.url;
				$('.img_url').val(url);
				$('img.show_img').attr('src', url);
				$('a.show_img').attr('href', url);
			} else {
				if (window.jAlert) {
					jAlert(data.error, '错误');
				} else {
					alert(data.error);
				}
			}
		}
	});
});
// 评论相关js
(function(){
	window.Comment = {};

	$(document).ready(function(){
		// 文本框事件
		$(".page-comment #commentForm textarea").on('click', function(){
			// 没有登录
			if($("#is_login_status").val() != 1){
				openPop("#login-pop");
			}
		});

		$('#comment-content').on('change', function() {
			var content = $(this).val();

			var objdata = {content: content};

			saveReplyDraft(uid, keyprefix, objid, objdata);
		});

		(function() {
			if (typeof keyprefix === "undefined") {
				return;
			}
			var draft = loadReplyDraft(uid, keyprefix, objid);
			if (draft) {
				$('#comment-content').val(draft.content);
			}
		})();

		// 编辑 tab
		$('.page').on('click', '.comment-edit-tab', function(evt){
			evt.preventDefault();

			var $this = $(this);
			var $tabMenu = $this.parent()
			var commentGroup = $tabMenu.data('comment-group')
			$this.addClass('cur');
			$tabMenu.children('.comment-preview-tab').removeClass('cur')

			$('.comment-content-preview[data-comment-group="' + commentGroup + '"]').hide();
			$('.comment-content-text[data-comment-group="' + commentGroup + '"]').show();
		});
		// 点击预览 tab
		$('.page').on('click', '.comment-preview-tab', function(evt){
			evt.preventDefault();

			var marked = SG.markSettingNoHightlight();

			var $this = $(this).addClass('cur');
			var $tabMenu = $this.parent();
			var commentGroup = $tabMenu.data('comment-group')
			var $preview = $('.comment-content-preview[data-comment-group="' + commentGroup + '"]')
			var $text = $('.comment-content-text[data-comment-group="' + commentGroup + '"]')
			$tabMenu.children('.comment-edit-tab').removeClass('cur');

			$text.hide();
			var content = $text.children('textarea').val();
			$preview.html(marked(content));
			// emoji 表情解析
			emojify.run($preview.get(0));
			$preview.show();

			Prism.highlightAll();
		});

		$('#replies').on('mouseenter', '.reply', function(evt) {
			$(this).find('.op-reply').removeClass('hideable');
		});
		$('#replies').on('mouseleave', '.reply', function(evt) {
			$(this).find('.op-reply').addClass('hideable');
		});

		$('#replies').on('click', '.reply_user', function(evt) {
			if ($(evt.target).hasClass('reply_user')) {
				$(this).parents('.reply-to-block').find('.markdown').toggleClass('dn');
			}
		});

		// 切换显示评论和编辑评论
		function toggleCommentShowOrEdit(floor, show) {
			var $markdown = $('.markdown[data-floor="' + floor + '"]')
			var $content = $markdown.children('.content')
			var $editWrapper = $markdown.children('.edit-wrapper')
			if (show) {
				$content.show()
				$editWrapper.hide()
			} else {
				$content.hide()
				$editWrapper.show()
				var $textarea = $editWrapper.children('textarea')
				$textarea.val($textarea.data('raw-content')).focus()
			}
		}

		// 点击编辑评论按钮
		$('#replies').on('click', '.btn-edit', function(evt) {
			evt.preventDefault()
			var floor = $(this).data('floor')
			var $markdown = $('.markdown[data-floor="' + floor + '"]')
			var $editWrapper = $markdown.children('.edit-wrapper')
			var $textarea = $editWrapper.children('textarea')
			toggleCommentShowOrEdit(floor, false)

			var $uploadBtn = $('.upload-img[data-floor="' + floor + '"]')

			// 复制上传
			// 防止重复上传
			var pasteUpload = $textarea.data('paste-uploader')
			if (!pasteUpload) {
				pasteUpload = $textarea.pasteUploadImage('/image/paste_upload')
				$textarea.data('paste-uploader', pasteUpload)
			}

			// 点击按钮上传
			// 防止重复上传
			var uploader = $uploadBtn.data('uploader')
			if (!uploader) {
				uploader = window.initPLUpload({
					ele: $uploadBtn[0]
				})
				$uploadBtn.data('uploader', uploader)
			}
		});

		// 点击取消编辑评论按钮
		$('#replies').on('click', '.btn.cancel', function(evt) {
			evt.stopPropagation();
			var floor = $(this).data('floor');
			toggleCommentShowOrEdit(floor, true)
		})

		// 点击提交编辑后的评论
		$('#replies').on('click', '.btn.submit', function(evt) {
			evt.stopPropagation();
			var floor = $(this).data('floor');
			var $markdown = $('.markdown[data-floor="' + floor + '"]')
			var $submitBtn = $(this)
			var $editWrapper = $markdown.children('.edit-wrapper')
			var $textarea = $editWrapper.find('textarea')
			var $content = $markdown.children('.content')
			var content = $textarea.val()
			var cid = $submitBtn.data("cid")

			editComment($submitBtn, cid, content, function() {
				$textarea.data('raw-content', content)
				$content.html(parseCmtContent(content))
				toggleCommentShowOrEdit(floor, true)
			})
		})

		// 点击回复某人
		$('#replies').on('click', '.btn-reply', function(evt) {
			evt.preventDefault();

			var floor = $(this).data('floor'),
				username = $(this).data('username');
			var $replyTo = $('.md-toolbar .reply-to');

			$replyTo.data('floor', floor).data('username', username);

			var title = '回复#'+floor+'楼';
			$replyTo.children('.fa-mail-reply').attr('title', title);
			$replyTo.children('.user').attr('title', title).attr('href', '#reply'+floor).text(username+' #'+floor);
			$replyTo.removeClass('dn');

			$('#commentForm textarea').focus();
		});

		$('.md-toolbar .reply-to .close').on('click', function(evt) {
			evt.preventDefault();
			$(this).parents('.reply-to').addClass('dn').data('floor', '').data('username', '');
		});

		// 支持粘贴上传图片
		$('#comment-content').pasteUploadImage('/image/paste_upload');

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

		// 异步加载 评论
		window.loadComments = function(p) {
			// 默认取最后一页
			p = p || 0;

			var objid = $('.comment-list').data('objid'),
				objtype = $('.comment-list').data('objtype');

			var params = {
				'objid': objid,
				'objtype': objtype,
				'p': p
			};
			$.getJSON('/object/comments', params, function(data){
				if (data.ok) {
					data = data.data;
					var comments = data.comments,
						replyComments = data.reply_comments;

					var content = '';
					for(var i in comments) {
						var comment = comments[i],
							meUid = $('[name="me-uid"]').val(),
							user = data[comment.uid];

						var avatar = user.avatar;
						if (avatar == "") {
							if (isHttps) {
								user.avatar = 'https://secure.gravatar.com/avatar/'+md5(user.email)+"?s=48";
							} else {
								user.avatar = 'http://gravatar.com/avatar/'+md5(user.email)+"?s=48";
							}
						} else if (avatar.indexOf('http') === -1) {
							user.avatar = cdnDomain+'avatar/'+avatar+'?imageView2/2/w/48';
						}

						var cmtTime = SG.timeago(comment.ctime);
						if (cmtTime == comment.ctime) {
							var cmtTimes = cmtTime.split(" ");
							comment.cmt_time = cmtTimes[0];
						} else {
							comment.cmt_time = cmtTime;
						}

						if (comment.reply_floor > 0) {
							var replyComment = replyComments[comment.reply_floor]
							comment.reply_user = data[replyComment.uid];
							comment.reply_content = replyComment.content;
						}
						comment.rawContent = comment.content
						comment.content = parseCmtContent(comment.content);
						content += $.templates('#one-comment').render({comment: comment, user: user, me: {uid: meUid}});
					}

					if (content != '') {
						$('.comment-list .words').html(content);

						// 链接，add target=_blank
						$('.comment-list .words .markdown').on('mousedown', 'a', function(evt){
							var url = $(this).attr('href');
							$(this).attr('target', '_blank');
						});

						$('.comment-list .markdown img').attr('data-action', 'zoom');

						$('.comment-list .markdown img').on('click', function() {
							$(this).parents('.box_white').css('overflow', 'visible');
						});
					}
					$('.comment-list .words').removeClass('hide');
					$('.comment-list .words').find('code[class*="language-"]').parent('pre').addClass('line-numbers');
					Prism.highlightAll();

					// emoji 表情解析
					emojify.run($('.comment-list .words').get(0));

					if ($("#is_login_status").val() == 1) {
						SG.registerAtEvent(true, true, $('.page-comment textarea'));
					}
				} else {
					comTip("回复加载失败");
				}
			});
		}

		var parseCmtContent = function(content) {
			var marked = SG.markSettingNoHightlight();
			content = SG.preProcess(content);
			content = marked(content);
			return SG.replaceCodeChar(content);
		};

		// 回复提交
		$('#comment-submit').on('click', function(){
			var content = $('#commentForm textarea').val();

			if(content == ""){
				alert("其实你想说点什么...");
			} else {
				var floor = $('.md-toolbar .reply-to').data('floor');
				if (parseInt(floor, 10) > 0) {
					var username = $('.md-toolbar .reply-to').data('username');
					content = '#'+floor+'楼 @'+username+' '+content;
				}
				postComment($(this), content, function(comment) {
					comTip("回复成功！");
					purgeReplyDraft(uid, keyprefix, objid);

					$('#commentForm textarea').val('');

					$('.md-toolbar .reply-to .close').click();
				});
			}
		});

		var editComment = function(thiss, cid, content, callback) {
			thiss.text("稍等").addClass("disabled").attr({"title":'稍等',"disabled":"disabled"});

			$.ajax({
				type:"post",
				url: '/object/comments/' + cid,
				data: {
					"content": content,
				},
				dataType: 'json',
				success: function(data){
					if(data.ok) {
						comTip("修改成功！");
						callback()
						thiss.text("提交").removeClass("disabled").removeAttr("disabled").attr({"title":'提交'});
					} else {
						alert(data.error);
					}
				},
				error: function() {
					thiss.text("提交").removeClass("disabled").removeAttr("disabled").attr({"title":'提交'});
				}
			})
		}

		var postComment = function(thiss, content, callback){
			thiss.text("稍等").addClass("disabled").attr({"title":'稍等',"disabled":"disabled"});

			var objid = $('.comment-list').data('objid'),
				objtype = $('.comment-list').data('objtype');

			var usernames = SG.analyzeAt(content);

			$.ajax({
				type:"post",
				url: '/comment/'+objid,
				data: {
					"objtype": objtype,
					"content": content,
					"usernames": usernames.join(',')
				},
				dataType: 'json',
				success: function(data){
					if(data.ok){
						var comment = data.data;

						var $pageComment = $('.comment-list'),
							meUid = $('[name="me-uid"]').val(),
							user = {};

						user.username = $pageComment.data('username'),
						user.uid = $pageComment.data('uid'),
						user.avatar = $pageComment.data('avatar'),
						comment.cmt_time = SG.timeago(comment.ctime);
						if (comment.reply_floor > 0) {
							comment.content = content.substr(1);
						}
						comment.reply_floor = 0;
						comment.rawContent = comment.content
						comment.content = parseCmtContent(comment.content);

						var oneCmt = $.templates('#one-comment').render({comment: comment, user: user, is_new: true, me: {uid: meUid}});

						var $cmtNumObj = $('#replies .cmtnum'),
							cmtNum = parseInt($cmtNumObj.text(), 10);
						if (cmtNum == 0) {
							$('.comment-list .words').html('');
						}

						$('.comment-list .words').append(oneCmt).removeClass('hide');
						Prism.highlightAll();

						// emoji 表情解析
						emojify.run($('.comment-list .words .reply:last').get(0));

						// 注册@
						SG.registerAtEvent(true, true, $('.page-comment textarea'));

						cmtNum++;

						$cmtNumObj.text(cmtNum);

						setTimeout(function(){
							$('.comment-list .words .reply').removeClass('light');
						}, 2000);
						callback();
					} else {
						alert(data.error);
					}
				},
				complete:function(){
					thiss.text("提交").removeClass("disabled").removeAttr("disabled").attr({"title":'提交'});
				},
				error:function(){
					thiss.text("提交").removeClass("disabled").removeAttr("disabled").attr({"title":'提交'});
				}
			});
		}
	});

	////// 评论翻页 ///////////
	$('.page_input').on('keydown', function(event) {
		if (event.keyCode == 13) {
			var p = $(this).val();
			$('.cmt-page .page-num a:nth-child('+p+')').trigger('click');
		}
	});

	$('.ctrl-page button').on('click', function() {
		var p = $('.cmt-page .page_input').val();

		if ($(this).hasClass('prev-page')) {
			p--;
		} else {
			p++;
		}

		$('.cmt-page .page-num a:nth-child('+p+')').trigger('click');
	});

	$('.ctrl-page button').on('mouseover', function() {
		if (!$(this).hasClass('disable_now')) {
			$(this).addClass('hover_now');
		}
	});

	$('.ctrl-page button').on('mousedown', function() {
		$(this).addClass('active_now');
	});

	$('.ctrl-page button').on('mouseleave', function() {
		$(this).removeClass('hover_now');
		$(this).removeClass('active_now');
	});

	$('.cmt-page .page-num a').on('click', function(evt) {
		evt.preventDefault();
		$('.page-num .page_current').removeClass('page_current').addClass('page_normal');

		var p = $(this).data('page'),
			pageMax = $('.cmt-page .page_input').attr("max");

		$('.cmt-page .page-num a:nth-child('+p+')').removeClass('page_normal').addClass('page_current')
		$('.page-num .page_input').val(p);

		$('.cmt-page .ctrl-page button')
			.removeClass('disable_now')
			.removeAttr("disabled");

		if (p == 1) {
			$('.cmt-page .prev-page')
				.removeClass('hover_now')
				.removeClass('active_now')
				.addClass('disable_now')
				.attr("disabled", "disabled");
		} else if (p == pageMax) {
			$('.cmt-page .next-page')
				.removeClass('hover_now')
				.removeClass('active_now')
				.addClass('disable_now')
				.attr("disabled", "disabled");
		}

		loadComments(p);

		return false;
	});
	/////////// 评论翻页 end //////////////

}).call(this);
