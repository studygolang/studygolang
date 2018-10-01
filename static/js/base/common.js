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
