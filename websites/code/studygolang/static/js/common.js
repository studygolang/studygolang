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