$(function(){
	$('.sidebar .top ul li').on('mouseenter', function(evt){
		
		if (evt.target.tagName != 'LI') {
			return;
		}
		$(this).parent().find('a').removeClass('cur');
		$(this).children('a').addClass('cur');

		var sbContent = $(this).parents('.top').next();
		var left = 0;

		sbContent.children().removeClass('hidden').hide();
		switch ($(this).attr('class')) {
		case 'first':
			sbContent.children('.first').show();
			left = "18px";
			break;
		case 'second':
			sbContent.children('.second').show();
			left = "114px";
			break;
		case 'last':
			sbContent.children('.last').show();
			left = "210px";
			break;
		}
		$(this).parents('.top').children('.bar').animate({left: left}, "fast");
	});

	// 侧边栏——最新帖子
	var topicRecent = function(data) {
		if (data.ok) {
			data = data.data;

			var content = '';
			for(var i in data) {
				content += '<li class="truncate">'+
						'<i></i><a href="/topics/'+data[i].tid+'" title="'+data[i].title+'">'+data[i].title+'</a>'+
						'</li>'
			}
			$('.sb-content .topic-list ul').html(content);
		}
	}

	// 侧边栏——最新博文
	var articleRecent = function(data){
		if (data.ok) {
			data = data.data;

			var content = '';
			for(var i in data) {
				content += '<li class="truncate">'+
						'<i></i><a href="/articles/'+data[i].id+'" title="'+data[i].title+'">'+data[i].title+'</a>'+
						'</li>'
			}
			$('.sb-content .article-list ul').html(content);
		}
	}

	// 侧边栏——最新开源项目
	var projectRecent = function(data){
		if (data.ok) {
			data = data.data;

			var content = '';
			for(var i in data) {
				var uri = data[i].id;
				if (data[i].uri != '') {
					uri = data[i].uri;
				}

				var title = data[i].category + ' ' + data[i].name;

				var logo = 'http://studygolang.qiniudn.com/gopher/default_project.jpg?imageView2/2/w/48';
				if (data[i].logo != '') {
					logo = data[i].logo;
				}
				content += '<li>'+
					'<a href="/p/'+uri+'">'+
						'<div class="logo"><img src="'+logo+'" alt="'+data[i].name+'" width="48px"/></div>'+
					'</a>'+
					'<div class="title">'+
						'<h4>'+
							'<a href="/p/'+uri+'" title="'+title+'">'+title+'</a>'+
						'</h4>'+
					'</div>'+
				'</li>';
			}
			$('.sb-content .project-list ul').html(content);
		}
	}

	// 侧边栏——最新资源
	var resourceRecent = function(data){
		if (data.ok) {
			data = data.data;

			var content = '';
			for(var i in data) {
				content += '<li class="truncate">'+
						'<i></i><a href="/resources/'+data[i].id+'" title="'+data[i].title+'">'+data[i].title+'</a>'+
						'</li>'
			}
			$('.sb-content .resource-list ul').html(content);
		}
	}

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
	
	// 侧边栏——最新评论
	var commentRecent = function(data){
		if (data.ok) {
			data = data.data;
			var comments = data.comments;

			var content = '';
			for(var i in comments) {
				var url = '';
				switch(comments[i].objtype) {
				case 0:
					url = '/topics/';
					break;
				case 1:
					url = '/articles/';
					break;
				case 2:
					url = '/resources/';
					break;
				case 4:
					url = '/p/';
					break;
				}
				url += comments[i].objid;

				var user = data[comments[i].uid];

				var avatar = user.avatar;
				if (avatar == "") {
					avatar = 'http://www.gravatar.com/avatar/'+md5(user.email)+"?s=40";
				}
				
				var cmtTime = SG.timeago(comments[i].ctime);
				if (cmtTime == comments[i].ctime) {
					var cmtTimes = cmtTime.split(" ");
					cmtTime = cmtTimes[0];
				}
				
				content += '<li>'+
					'<div class="pic">'+
						'<a href="/user/'+user.username+'" target="_blank">'+
							'<img src="'+avatar+'" alt="'+user.username+'" width="40px" height="40px">'+
						'</a>'+
					'</div>'+
					'<div class="word">'+
						'<div class="w-name">'+
							'<a href="/user/'+user.username+'" target="_blank" title="'+user.username+'">'+user.username+'</a>'+
							'<span title="'+comments[i].ctime+'">'+cmtTime+'</span>'+
						'</div>'+
						'<div class="w-page">'+
							'<span>在<a href="'+url+'#commentForm" title="'+comments[i].objinfo.title+'">'+comments[i].objinfo.title+'  </a>中评论</span>'+
						'</div>'+
						'<div class="w-comment">'+
							'<span>'+comments[i].content+'</span>'+
						'</div>'+
					'</div>'+
				'</li>';
			}
			$('.sb-content .cmt-list ul').html(content);
			emojify.run($('.sb-content .cmt-list ul li .w-comment').get(0));
		}
	}

	var userActive = function(data) {
		userList(data, '#active-list');
	}

	var userNewest = function(data) {
		userList(data, '#newest-list');
	}
	
	var userList = function(data, id) {
		if (data.ok) {
			data = data.data;

			var content = '';
			for(var	i in data) {
				var avatar = data[i].avatar;
				if (avatar == "") {
					avatar = 'http://www.gravatar.com/avatar/'+md5(data[i].email)+"?s=40";
				}
				
				content	+= '<li	class="pull-left">'+
					'<div class="avatar">'+
					'<a href="/user/'+data[i].username+'" title="'+data[i].username+'"><img alt="'+data[i].username+'" class="img-circle" src="'+avatar+'" width="48px" height="48px"></a>'+
					'</div>'+
		  			'<div class="name"><a href="/user/'+data[i].username+'" title="'+data[i].username+'">'+data[i].username+'</a></div>'+
		  		'</li>';
			}
			$('.sb-content '+id+' ul').html(content);
		}
	}

	var websiteStat = function(data) {
		if (data.ok) {
			data = data.data;

			var content = '<li>会员数: <span>'+data.user+'</span> 人</li>'+
				'<li>博文数: <span>'+data.article+'</span> 篇</li>'+
				'<li>话题数: <span>'+data.topic+'</span> 个</li>'+
				'<li>评论数: <span>'+data.comment+'</span> 条</li>'+
				'<li>资源数: <span>'+data.resource+'</span> 个</li>'+
				'<li>项目数: <span>'+data.project+'</span> 个</li>';

			$('.sb-content .stat-list ul').html(content);
		}
	}

	var readingRecent = function(data) {
		if (data.ok) {
			data = data.data;

			var content = '';
			if (data.length == 1) {
				data = data[0];
				content = '<li><a href="/readings/'+data.id+'" target="_blank">'+data.content+'</a></li>';
			} else {
				for(var i in data) {
					content += '<li>'+
						'<a href="/readings/'+data[i].id+'">'+
							'<div class="time"><span>10-25</span></div>'+
						'</a>'+
						'<div class="title">'+
							'<h4>'+
								'<a href="/readings/'+data[i].id+'">'+data[i].content+'</a>'+
							'</h4>'+
						'</div>'+
					'</li>';
				}
			}
			
			$('.sb-content .reading-list ul').html(content);
		}
	}

	var sidebar_callback = {
		"/topics/recent.json": {"func": topicRecent, "class": ".topic-list"},
		"/articles/recent.json": {"func": articleRecent, "class": ".article-list"},
		"/projects/recent.json": {"func": projectRecent, "class": ".project-list"},
		"/resources/recent.json": {"func": resourceRecent, "class": ".resource-list"},
		"/comments/recent.json": {"func": commentRecent, "class": ".cmt-list"},
		"/users/active.json": {"func": userActive, "class": "#active-list"},
		"/users/newest.json": {"func": userNewest, "class": "#newest-list"},
		"/websites/stat.json": {"func": websiteStat, "class": ".stat-list"},
		"/readings/recent.json": {"func": readingRecent, "class": ".reading-list"},
	};
	
	if (typeof SG.SIDE_BARS != "undefined") {

		for (var i in SG.SIDE_BARS) {
			if (typeof sidebar_callback[SG.SIDE_BARS[i]] != "undefined") {
				var sbObj = sidebar_callback[SG.SIDE_BARS[i]];
				var limit = $('.sidebar .sb-content '+sbObj['class']).data('limit');
				if (limit == "") {
					limit = 10;
				}
				
				$.getJSON(SG.SIDE_BARS[i], {limit: limit}, sbObj['func']);
			}
		}
	}
	
});