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

	// 文本框事件
	$(".page-comment #commentForm textarea").click(function(){
		// 没有登录
		if($("#is_login_status").val() != 1){
			openPop("#login-pop");
		}
	})

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
		
		$.post('/account/login.json', $(this).serialize(), function(data){
			if (data.ok) {
				location.reload();
			} else {
				$('#login-pop .login-form .error').text(data.error).show();
			}
		});
	});

	$('#username, #passwd').on('focus', function(){$('#login-pop .login-form .error').hide();});

	// 异步加载 评论
	window.loadComments = function() {
		var objid = $('.page-comment').data('objid'),
			objtype = $('.page-comment').data('objtype');
		
		var params = {
			'objid': objid,
			'objtype': objtype
		};
		$.getJSON('/object/comments.json', params, function(data){
			if (data.ok) {
				data = data.data;
				var comments = data.comments;

				var content = '';
				for(var i in comments) {
					var comment = comments[i],
						user = data[comment.uid];

					var avatar = user.avatar;
					if (avatar == "") {
						avatar = 'http://www.gravatar.com/avatar/'+md5(user.email)+"?s=48";
					}
					
					var cmtTime = SG.timeago(comments[i].ctime);
					if (cmtTime == comments[i].ctime) {
						var cmtTimes = cmtTime.split(" ");
						cmtTime = cmtTimes[0];
					}
					content += contructOneCmt(comment.floor, user.username, avatar, comment.content, comment.ctime, cmtTime);
				}

				if (content != "") {
					$('.page-comment .words ul').html(content);
					$('.page-comment .words').removeClass('hide');
				}
			} else {
				comTip("评论加载失败");
			}
		});
	}

	var contructOneCmt = function(floor, username, avatar, content, ctime, cmtTime, needLight) {
		var oneCmt = '<li id="comment'+floor+'">';
		if (typeof needLight !== "undefined") {
			oneCmt = '<li id="comment'+floor+'" class="light">';
		}
		return oneCmt+
			'<div class="pull-left face">'+
				'<a href="/user/'+username+'" target="_blank"><img src="'+avatar+'" width="48px" height="48px" alt="'+username+'"></a>'+
			'</div>'+
			'<div class="cmt-body">'+
				'<div class="cmt-content">'+
					'<a href="/user/'+username+'" class="name replyName" target="_blank" data-floor="'+floor+'楼">'+username+'</a>：'+
					'<span>'+content+'</span>'+
					'<!--'+
					'<span>'+
						'<a href="" onclick="return confirm(\'确定删除该条评论?\');" title="删除">删除</a>'+
					'</span>'+
					'-->'+
				'</div>'+
				'<div class="cmt-meta">'+
					floor+'楼, <span title="'+ctime+'">'+cmtTime+'</span>'+
					'<a href="#" class="small_reply" data-floor="'+floor+'" title="回复此楼"><i class="glyphicon glyphicon-comment"></i> 回复</a>'+
				'</div>'+

				'<!--回复开始-->'+
				'<div class="respond-submit">'+
					'<div class="text">'+
						'<input type="text" name="content" value="">'+
						'<div class="tip"></div>'+
					'</div>'+
					'<div class="sub clr">'+
						'<button>提交</button>'+
					'</div>'+
				'</div>'+
			'</div>'+
		'</li>';
	}

	var tipLength = function(thiss, callback){// 先输入文本 宽度计算完成 回调
		var $reply = thiss.parent(".cmt-meta").prev(".cmt-content").children(".replyName");
		var	replyName =	'#'+$reply.data('floor')+' @'+$reply.text()+' ';

		var $tip = thiss.parents(".cmt-body").find(".respond-submit .text .tip");
		$tip.text(replyName);
		var	timer =	setInterval(function(){
			if($tip.outerWidth()){
				callback();
				clearInterval(timer);
			}
		},100)
	}

	// 回复交互表单
	$(".page-comment").on('click', '.small_reply', function(event){
		event.preventDefault();
		
		var thiss = $(this);
		if($("#is_login_status").val() == 1){// 如果登录
			// 隐藏所依回复表单，
			$(".page-comment .respond-submit").hide(10); 
			// 设置input的样式并默认选中
			var $cmtBody = thiss.parents(".cmt-body");
			tipLength(thiss, function(){
				var tipWid = $cmtBody.find(".respond-submit .text .tip").outerWidth();
				var cbWid = $cmtBody.width();
				$cmtBody.find(".respond-submit input").css({'width':cbWid,'padding-left':tipWid+10});
				$cmtBody.find(".respond-submit input").focus();
			})
			
			// 显示当前表单
			setTimeout(function(){
				$cmtBody.find(".respond-submit").slideDown(300);
			},150)

		} else {//未登录
			openPop("#login-pop");
		}
		event.stopPropagation();
	});

	// 点击其他地方收起回复框
	$(".page-comment").on('click', '.respond-submit', function(event){event.stopPropagation();});
	$(document).click(function(){$(".respond-submit").slideUp(200);});

	// 评论提交
	$('.page-comment #commentForm .sub button').on('click', function(){
		var content = $('.page-comment #commentForm textarea').val();

		if(content == ""){
			alert("其实你想说点什么...");
		} else {
			postComment($(this), content, function(comment){
				comTip("评论成功！");
				$('#commentForm textarea').val('');
			});
		}
	});
	// 回复表单提交
	$(".page-comment").on('click', '.cmt-body .sub button', function(){
		var $text = $(this).parent(".sub").prev(".text");
		var replyTo = $text.children('.tip').text();
		var	content	= $text.children("input").val();
		
		if(content == ""){
			alert("其实你想说点什么...");
		} else {
			var that = $(this)
			content = replyTo + content;
			postComment(that, content, function(){
				comTip("回复成功！");
				that.parent(".sub").prev(".text").children("input").val("");
				that.parents(".respond-submit").slideUp(200);
			});
		}
	});

	var postComment = function(thiss, content, callback){
		thiss.text("稍等").addClass("disabled").attr({"title":'稍等',"disabled":"disabled"});

		var objid = $('.page-comment').data('objid'),
			objtype = $('.page-comment').data('objtype');
		$.ajax({
			type:"post",
			url: '/comment/'+objid+'.json',
			data: {
				"objtype":objtype,
				"content":content
			},
			dataType: 'json',
			success: function(data){
				if(data.errno == 0){
					var comment = data.data;
					var $pageComment = $('.page-comment'),
					username = $pageComment.data('username'),
					avatar = $pageComment.data('avatar'),
					cmtTime = SG.timeago(comment.ctime);
					var oneCmt = contructOneCmt(comment.floor, username, avatar, comment.content, comment.ctime, cmtTime, true);

					$('.page-comment .words ul').append(oneCmt);
					$('.page-comment .words').removeClass('hide');

					var $cmtNumObj = $('.page-comment .words h3 .cmtnum'),
						cmtNum = parseInt($cmtNumObj.text(), 10) + 1;

					$cmtNumObj.text(cmtNum);
					$('.page .meta .p-comment .cmtnum').text(cmtNum);
					
					setTimeout(function(){
						$('.page-comment .words ul li').removeClass('light');
					}, 2000);
					callback();
				}else{
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