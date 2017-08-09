// 初始化
(function(){
	window.App = {
		
		notifier:null,
		
		loading:function(){
			return console.log("loading...")
		},
		
		fixUrlDash:function(url){return url.replace(/\/\//g,"/").replace(/:\//,"://")},
		
		// 警告信息显示, to 显示在那个dom前(可以用 css selector)
		alert:function(msg, to){
			$(".alert").remove();
			$(to).before("<div class='alert'><a class='close' href='#' data-dismiss='alert'>X</a>"+msg+"</div>")
		},
		
		// 成功信息显示, to 显示在那个dom前(可以用 css selector)
		notice:function(msg, to){
			$(".alert").remove();
			$(to).before("<div class='alert alert-success'><a class='close' data-dismiss='alert' href='#'>X</a>"+msg+"</div>");
		},
		
		openUrl:function(url){window.open(url)},
		
		// Use this method to redirect so that it can be stubbed in test
		gotoUrl:function(url){window.location=url},
		
		likeable:function(dom){
			var obj = $(dom),
				type = obj.data('type'),
				id = obj.data('id'),
				state = obj.data('state'),
				count = parseInt(obj.data('count'));
			// 没有【喜欢过】
			if (state !== 'liked') {
				$.ajax({url:"/like",type:"POST",data:{type:type,id:id}});
				count++;
				obj.data('count', count);
				App.likeableAsLiked(dom);
			} else {
				$.ajax({url:"/like",type:"DELETE",data:{type:type}});
				if (count > 0) {
					count--;
				}
				obj.data('count', count).data('state', '').attr('title', '喜欢');
				if (count === 0) {
					$('span', dom).text('喜欢');
				} else {
					$('span', dom).text(count+'人喜欢');
				}
				$('i.icon', dom).attr('class', 'icon small_like');
			}
			return false;
		},
		
		likeableAsLiked:function(dom){
			var obj = $(dom),
				count = parseInt(obj.data('count'));
			obj.data('state', 'liked').attr('title', '取消喜欢');
			$('span', dom).text(count+'人喜欢');
			$('i.icon', dom).attr('class', 'icon small_liked');
		},
		
		// 绑定 @ 回复功能（输入框支持@自动提示）
		atReplyable:function(dom, users){
			if(users.length === 0 ) return;
			$(dom).atwho("@", {data:users, tpl:"<li data-value='${username}'>${username} <small>${name}</small></li>"});
		},
		
		// 支持 http://www.emoji-cheat-sheet.com/ 表情
		tipEmojis: function(dom) {
			$(dom).atwho(':', {
				data: window.emojis,
				tpl:"<li data-value='${key}'>${name} <img src='http://a248.e.akamai.net/assets.github.com/images/icons/emoji/${name}.png'  height='20' width='20' /></li>"
			});
		},

		initForDesktopView:function(){
			if(typeof app_mobile != "undefined")
				return;
			$("a[rel=twipsy]").tooltip();

			// CommentAble @ 回复功能
			var users = App.scanLogins($(".cell_comments .comment .info .name a"));
			var result = [];
			for (var username in users) {
				var user = {uid: users[username].uid, username: username, name: users[username].name};
				result.push(user);
			}
			App.atReplyable(".cell_comments_new textarea", result);
		},
		
		// scan logins in jQuery collection and returns as a object,
		// which key is username(账号）, and value is the object of {uid:xxx, name（姓名/昵称）:xxx}.
		scanLogins:function(query){
			var result = {};
			query.each(function(){
				result[$(this).text()] = {
					'uid':$(this).data('uid'),
					'name':$(this).data('name')
				};
			});
			return result;
		},
		
		initNotificationSubscribe:function(){
			// FAYE：Simple pub/sub messaging for the web
		}
	};
	
	$(document).ready(function(){
		//App.initForDesktopView();
		
		// 【时间轴】插件
		$("abbr.timeago").timeago();
		
		$(".alert").alert();
		
		$(".dropdown-toggle").dropdown();
		
		// Web订阅（使用http://faye.jcoglan.com/，TODO：暂不支持）
		if (typeof FAYE_SERVER_URL !="undefined" && FAYE_SERVER_URL !== null) {
			// App.initNotificationSubscribe();
		}
		
		$("form.new_topic,form.new_reply,form.new_note,form.new_page,form.new_resource").sisyphus({timeout:2});
		
		$("form a.reset").click(function(){
			return $.sisyphus().manuallyReleaseData()
		});
		
		/*
		// 绑定评论框（回复） Ctrl+Enter 提交事件
		$(".cell_comments_new textarea").on("keydown","ctrl+return",function(env){
			var tg = $(env.target);
			if (tg.val().trim().length>0) {
				tg.parents("form").submit();
			}
			return false;
		});
		*/

		// Choose 样式（美化），需要http://davidwalsh.name/demo/jquery-chosen.php插件
		// $("select").chosen();

		// 回到顶部
		$("a.go_top").click(function(){
			$("html, body").animate({scrollTop:0},300);
			return false;
		});
		$(window).bind("scroll resize",function(){
			var isscroll = $(window).scrollTop();
			if (isscroll >= 1) {
				$("a.go_top").show();
			} else {
				$("a.go_top").hide();
			}
		});
	});
}).call(this)