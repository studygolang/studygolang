$(function(){
	var safeStr = function(str) {
		return str.replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, "&quot;").replace(/'/g, "&#039;");
	}

	$('.sidebar .top ul li').on('mouseenter', function(evt){
		
		if (evt.target.tagName != 'LI') {
			// return;
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

	var gravatar = function(avatar, email, size) {
		if (avatar == "") {
			if (isHttps) {
				avatar = 'https://secure.gravatar.com/avatar/'+md5(email)+"?s="+size;
			} else {
				avatar = 'http://gravatar.com/avatar/'+md5(email)+"?s="+size;
			}
		} else {
			if (avatar.indexOf('http') == 0) {
				avatar += '&s='+size;
			} else {
				avatar = cdnDomain+'avatar/'+avatar+'?imageView2/2/w/'+size;
			}
		}

		return avatar;
	}

	// 侧边栏——最新帖子
	var topicRecent = function(data) {
		if (data.ok) {
			data = data.data;

			var content = '';
			for(var i in data) {
				var title = safeStr(data[i].title);
				content += '<li class="truncate">'+
						'<i></i><a href="/topics/'+data[i].tid+'?fr=sidebar" title="'+title+'">'+title+'</a>'+
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
				var title = safeStr(data[i].title);
				content += '<li class="truncate">'+
						'<i></i><a href="/articles/'+data[i].id+'?fr=sidebar" title="'+title+'">'+title+'</a>'+
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
				var logo = data[i].logo;

				title = safeStr(title);
				content += '<li>'+
					'<a href="/p/'+uri+'">'+
						'<div class="logo"><img src="'+logo+'" alt="'+data[i].name+'" width="48px"/></div>'+
					'</a>'+
					'<div class="title">'+
						'<h4>'+
							'<a href="/p/'+uri+'?fr=sidebar" title="'+title+'">'+title+'</a>'+
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
				var title = safeStr(data[i].title);
				content += '<li class="truncate">'+
						'<i></i><a href="/resources/'+data[i].id+'?fr=sidebar" title="'+title+'">'+title+'</a>'+
						'</li>'
			}
			$('.sb-content .resource-list ul').html(content);
		}
	}

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
	
	// 侧边栏——最新评论
	var commentRecent = function(data){
		if (data.ok) {
			data = data.data;
			var comments = data.comments;

			var content = '';
			for(var i in comments) {
				var url = comments[i].objinfo.uri+comments[i].objid;

				var user = data[comments[i].uid];
				var avatar = gravatar(user.avatar, user.email, 40);
				
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
							'<span>在 <a href="'+url+'#commentForm" title="'+comments[i].objinfo.title+'">'+comments[i].objinfo.title+'</a> 中评论</span>'+
						'</div>'+
						'<div class="w-comment">'+
							'<span>'+comments[i].content+'</span>'+
						'</div>'+
					'</div>'+
				'</li>';
			}
			$('.sb-content .cmt-list ul').html(content);
			emojify.run($('.sb-content .cmt-list ul').get(0));
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
				var avatar = gravatar(data[i].avatar, data[i].email, 48);
				
				content	+= '<li	class="pull-left">'+
					'<div class="avatar">'+
					'<a href="/user/'+data[i].username+'" title="'+data[i].username+'"><img alt="'+data[i].username+'" class="img-circle" src="'+avatar+'" width="48px" height="48px"></a>'+
					'</div>'+
		  			'<div class="name" style="white-space: nowrap;"><a style="word-break: normal;" href="/user/'+data[i].username+'" title="'+data[i].username+'">'+data[i].username+'</a></div>'+
		  		'</li>';
			}
			$('.sb-content '+id+' ul').html(content);
		}
	}

	var websiteStat = function(data) {
		if (data.ok) {
			data = data.data;

			var content = '<li>会员数: <span>'+data.user+'</span> 人</li>';
			if (data.topic > 0) {
				content += '<li>主题数: <span>'+data.topic+'</span> 个</li>';
			}
			if (data.article > 0) {
				content += '<li>文章数: <span>'+data.article+'</span> 篇</li>';
			}
			if (data.comment > 0) {
				content += '<li>回复数: <span>'+data.comment+'</span> 条</li>';
			}
			if (data.resource > 0) {
				content += '<li>资源数: <span>'+data.resource+'</span> 个</li>';
			}
			if (data.project > 0) {
				content += '<li>项目数: <span>'+data.project+'</span> 个</li>';
			}
			if (data.book > 0) {
				content += '<li>图书数: <span>'+data.book+'</span> 本</li>';
			}

			$('.sb-content .stat-list ul').html(content);
		}
	}

	var readingRecent = function(data) {
		if (data.ok) {
			data = data.data;

			if (!data || data.length == 0) {
				$('.sb-content .reading-list').parents('.sidebar').hide();
				return;
			}

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

	var hotNodes = function(data) {
		if (data.ok) {
			data = data.data;
			if (data == null) {
				return;
			}

			var content = '';
			for(var i in data) {
				content += '<li><a href="/go/'+data[i].ename+'?fr=sidebar" title="'+data[i].name+'">'+data[i].name+'</a></li>';
			}
			
			$('.sb-content .node-list ul').html(content);
		}
	}

	var friendLinks = function(data) {
		if (data.ok) {
			data = data.data;
			if (data == null) {
				return;
			}

			var content = '';
			for(var i in data) {
				content += '<li style="margin-left:15px; margin-bottom:5px;">'+
							'<a href="'+data[i].url+'" target="_blank" title="'+data[i].name+'">'+data[i].name+'</a>'+
						'</li>';
			}
			
			$('.sb-content .friendslink-list ul').html(content);
		}
	}

	// 侧边栏——排行榜
	var rankList = function(result, dataKeys){
		if (result.ok) {
			data = result.data;
			var list = data.list;

			var content = '';
			for(var i in list) {
				var path = data.path + list[i].id,
					title = list[i].title;
				switch (data.objtype) {
				case 0:
					path = data.path + list[i].tid;
					break;
				case 4:
					title = list[i].category + list[i].name;
					if (list[i].uri != '') {
						path = data.path + list[i].uri;
					}
					break;
				case 5:
					title = list[i].name;
					break;
				}
				title = safeStr(title);

				var pos = parseInt(i, 10) + 1;

				var rankFlag = '';
				if (pos < 4) {
					rankFlag = '<img src="'+cdnDomain+'static/img/rank_medal'+pos+'.png" width="20px">';
				} else {
					rankFlag = '<em>'+pos+'</em>';
				}

				content += '<li>'+
						rankFlag+'<a href="'+path+'?fr=sidebar" title="'+title+'">'+title+'</a> - '+list[i].rank_view+' 阅读'+
						'</li>'
			}

			$('.sb-content .rank-list').each(function(index) {
				if ($(this).data('objtype') == data.objtype && $(this).data('rank_type') == data.rank_type) {
					$(this).children().html(content);
				}
			});
		}
	}

	var sidebar_callback = {
		"/topics/recent": {"func": topicRecent, "class": ".topic-list"},
		"/articles/recent": {"func": articleRecent, "class": ".article-list"},
		"/projects/recent": {"func": projectRecent, "class": ".project-list"},
		"/resources/recent": {"func": resourceRecent, "class": ".resource-list"},
		"/comments/recent": {"func": commentRecent, "class": ".cmt-list"},
		"/users/active": {"func": userActive, "class": "#active-list"},
		"/users/newest": {"func": userNewest, "class": "#newest-list"},
		"/websites/stat": {"func": websiteStat, "class": ".stat-list"},
		"/readings/recent": {"func": readingRecent, "class": ".reading-list"},
		"/nodes/hot": {"func": hotNodes, "class": ".node-list"},
		"/friend/links": {"func": friendLinks, "class": ".friendslink-list"},
		"/rank/view": {"func": rankList, "class": ".rank-list", data_keys:["objtype", "rank_type"]},
	};
	
	if (typeof SG.SIDE_BARS != "undefined") {

		for (var i in SG.SIDE_BARS) {
			if (typeof sidebar_callback[SG.SIDE_BARS[i]] != "undefined") {
				var sbObj = sidebar_callback[SG.SIDE_BARS[i]],
					$dataSelector = $('.sidebar .sb-content '+sbObj['class']);

				if ($dataSelector.length == 0) {
					continue;
				}

				if (!sbObj.data_keys) {
					var limit = $dataSelector.data('limit');
					if (limit == "") {
						limit = 10;
					}
					$.ajax({
						type:"get",
						url: SG.SIDE_BARS[i],
						data: {limit: limit},
						dataType: 'json',
						success: sbObj['func'],
						ifModified: true
					});
					
					continue;
				}

				$dataSelector.each(function(index) {
					var limit = $(this).data('limit');
					var params = {limit: limit};

					for (var j in sbObj.data_keys) {
						var k = sbObj.data_keys[j];
						params[k] = $(this).data(k);
					}
					
					$.ajax({
						type:"get",
						url: SG.SIDE_BARS[i],
						data: params,
						dataType: 'json',
						success: sbObj['func'],
						ifModified: true
					});
				});
			}
		}
	}
	
});