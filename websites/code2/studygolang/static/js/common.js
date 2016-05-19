// studygolang 全局对象（空间）
var SG = {};

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
	publish: function(that) {
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
					
					setTimeout(function(){
						var redirect = $form.data('redirect');
						if (redirect) {
							window.location.href = redirect;
						}
					}, 1000);
				}else{
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
	// 配置 marked 语法高亮
	marked.setOptions({
		highlight: function (code) {
			code = SG.replaceSpecialChar(code);
			return hljs.highlightAuto(code).value;
		}
	});

	return marked;
}

// 替换 `` 代码块中的 "<>& 等字符
SG.replaceCodeChar = function(code) {
	return code.replace(/<code>.*<\/code>/g, function(matched, index, origin) {
		return SG.replaceSpecialChar(matched);
	});
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
			tpl:"<li data-value='${key}'><img src='http://www.emoji-cheat-sheet.com/graphics/emojis/${name}.png' height='20' width='20' /> ${name}</li>"
		})/*.atwho({
			at: "\\",
			data: window.twemojis,
			tpl:"<li data-value='${name}'><img src='https://twemoji.maxcdn.com/16x16/${key}.png' height='16' width='16' /> ${name}</li>"
		})*/;
	}
}

jQuery(document).ready(function($) {
	// timeago：3 天之内才显示 timeago

	// time 的格式 2014-10-02 11:40:01
	SG.timeago = function(time) {
		var ago = new Date(time),
			now = new Date();

		if (now - ago > 3 * 86400 * 1000) {
			return time;
		}

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
		closePop();
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
		$(".pop").hide();
		$("#sg-overlay").fadeOut(300);
	}

	$("#sg-overlay").click(function(){closePop()});

	// 弹窗异步登录
	$('#login-pop .login-form form').on('submit', function(evt){
		evt.preventDefault();

		var username = $('#username').val(),
			passwd = $('#passwd').val();

		if (username == "") {
			$('#username').parent().addClass('has-error');
			return;
		}
		if (passwd == "") {
			$('#passwd').parent().addClass('has-error');
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
					comTip("感谢喜欢！");
					$(that).addClass('hadlike').attr('title', '取消喜欢');
					likeNum++;
				} else {
					comTip("已取消喜欢！");
					$(that).removeClass('hadlike').attr('title', '我喜欢');
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
	$('.page .like-btn').on('click', function(evt){
		evt.preventDefault();

		var that = this;
		postLike(that, function(likeNum, likeFlag){
			$('.page .meta .p-comment .like .likenum').text(likeNum);
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
				$('.page .collect').addClass('hadlike').attr('title', '取消收藏');
			} else {
				$('.page .collect').removeClass('hadlike').attr('title', '稍后再读');
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
	
});

// 在线人数统计
window.WebSocket = window.WebSocket || window.MozWebSocket;
if (window.WebSocket) {
	var websocket = new WebSocket(wsUrl);

	websocket.onopen = function(){
		// console.log("open");
	}

	websocket.onclose = function(){
		// console.log("close");
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
	
	// websocket.onerror = onError;
}

$(function(){
	if (Math.random()*50 <= 1) {
		$('.ad').each(function(){
			var url = $(this).attr('href');

			var adImg = new Image();
			adImg.src = url;
		});
	}
});