var interval = 5000;	//两次滚动之间的时间间隔 
var stepsize = 28;		//滚动一次的长度，必须是行高的倍数,这样滚动的时候才不会断行 
var objInterval = null;

$(document).ready( function(){
	// 顶部动态获取
	$.getJSON('/dymanics/recent.json', function(data){
		if (data.ok) {
			data = data.data;

			var content = '';
			for (var i in data) {
				content += '<li>'+
					'<a href="'+data[i].url+'" title="'+data[i].content+'" target="_blank">'+data[i].content+'</a>'+
					'</li>';
			}

			$('#top ul').html(content);

			// 用上部的内容填充下部 
			$("#bottom").html($("#top").html());

			// 给显示的区域绑定鼠标事件
			$("#content").bind("mouseover",function(){StopScroll();});
			$("#content").bind("mouseout",function(){StartScroll();});

			// 启动定时器 
			StartScroll();
		}
	});

	// 登录
	$('.login').submit(function(evt) {
		evt.preventDefault();
		$.post('/account/login.json', $(this).serialize(), function(data) {
			if (data.ok) {
				location.reload();
			} else {
				comTip(data.error);
			}
		});
	});
	
}); 

// 启动定时器，开始滚动 
function StartScroll(){ 
	objInterval=setInterval("verticalloop()", interval);
} 

// 清除定时器，停止滚动 
function StopScroll(){ 
	window.clearInterval(objInterval); 
} 

// 控制滚动 
function verticalloop(){ 
	// 判断是否上部内容全部移出显示区域 
	// 如果是，从新开始;否则，继续向上移动 
	if($("#content").scrollTop()>=$("#top").outerHeight()){ 
		$("#content").scrollTop($("#content").scrollTop()-$("#top").outerHeight()); 
	} 
	// 使用jquery创建滚动时的动画效果 
	$("#content").animate( 
	{"scrollTop" : $("#content").scrollTop()+stepsize +"px"},600,function(){ 
		// 这里用于显示滚动区域的scrollTop，实际应用中请删除 
		// $("#foot").html("scrollTop:"+$("#content").scrollTop()); 
	}); 
}