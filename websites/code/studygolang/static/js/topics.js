// 话题（帖子）相关js功能
(function(){
	window.Topics = {
		replies_per_page: 50,
		
		// 往话题编辑器里面插入图片代码
		appendImageFromUpload:function(srcs){
			var txtBox = $(".topic_editor"),
				caret_pos = txtBox.caretPos(),
				src_merged = "";
			for (src in srcs) {
				src_merged = "![]("+src+")\n";
			}
			var source = txtBox.val();
			before_text = source.slice(0, caret_pos);
			txtBox.val(before_text + src_merged + source.slice(caret_pos+1, source.count));
			txtBox.caretPos(caret_pos+src_merged.length);
			txtBox.focus();
		},
		
		// 上传图片
		initUploader:function(){
			$("#topic_add_image").bind("click", function(){
				$("#topic_upload_images").click();
				return false;
			});

			var opts = {
				url:"/photos",
				type:"POST",
				beforeSend:function(){
					$("#topic_add_image").hide();
					$("#topic_add_image").before("<span class='loading'>上传中...</span>");
				},
				success:function(result, status, xhr){
					$("#topic_add_image").parent().find("span").remove();
					$("#topic_add_image").show();
					Topics.appendImageFromUpload([result]);
				}
			};
			$("#topic_upload_images").fileUpload(opts);
			return false;
		},
		
		// 回复
		reply:function(floor, username) {
			var reply_body = $("#wmd-input"),
				new_text = "#"+floor+"楼 @"+username+" ";
			new_text += reply_body.val();
			reply_body.focus().val(new_text);
			return false;
		},
		
		// Given floor, calculate which page this floor is in
		pageOfFloor:function(floor) {
			return Math.floor((floor - 1) / Topics.replies_per_page) + 1
		},
		
		// 跳到指定楼。如果楼层在当前页，高亮该层，否则跳转到楼层所在页面并添
		// 加楼层的 anchor。返回楼层 DOM Element 的 jQuery 对象
		// 
		//  -   floor: 回复的楼层数，从1开始
		gotoFloor:function(floor){
			var replyEl = $("#reply"+floor);
			if (replyEl.length > 0) {
				// 高亮
				Topics.highlightReply(replyEl);
			} else {
				page = Topics.pageOfFloor(floor);
			}
			// TODO: merge existing query string
			url = window.location.pathname + "?p="+page+ "#reply"+floor;
			App.gotoUrl( url );
			return replyEl;
		},
		
		// 高亮指定楼。取消其它楼的高亮
		// 
		// replyEl: 需要高亮的 DOM Element，需要 jQuery 对象
		highlightReply:function(replyEl) {
			$("#replies .reply").removeClass("light");
			replyEl.addClass("light");
		},
		
		// 异步更改用户 like 过的回复的 like 按钮的状态
		checkRepliesLikeStatus:function(user_liked_reply_ids){
			for (id in user_liked_reply_ids) {
				el = $("#replies a.likeable[data-id=#{id}]");
				App.likeableAsLiked(el);
			}
		},
		
		// Ajax 回复后的事件
		replyCallback:function(success, msg) {
			$("#main .alert-message").remove();
			if (success) {
				$("abbr.timeago", $("#replies .reply").last()).timeago();
				$("abbr.timeago", $("#replies .total")).timeago();
				$("#new_reply textarea").val('');
				$("#preview").text('');
				App.notice(msg, '#reply');
			} else {
				App.alert(msg, '#reply');
				$("#new_reply textarea").focus();
				$('#btn_reply').button('reset');
			}
		},
		
		// 预览（不需要）
		preview:function(body){
			$("#preview").text("Loading...");
			$.post("/topics/preview", {body:body}, function(data){
				return $("#preview").html(data.body)
			}, "json");
		},
		
		hookPreview:function(switcher, textarea){
			// put div#preview after textarea
			var preview_box = $(document.createElement("div")).attr("id", "preview");
			preview_box.addClass("body");
			$(textarea).after(preview_box);
			preview_box.hide();

			$(".edit a",switcher).click(function(){
				$(".preview", switcher).removeClass("active")
				$(this).parent().addClass("active");
				$(preview_box).hide();
				$(textarea).show();
				return false;
			});
			$(".preview a", switcher).click(function(){
				$(".edit", switcher).removeClass("active");
				$(this).parent().addClass("active");
				$(preview_box).show();
				$(textarea).hide();
				Topics.preview($(textarea).val());
				return false;
			});
		},
			
		initCloseWarning:function(el, msg){
			if (el.length == 0) {
				return false;
			}
			if (!msg) {
				msg = "离开本页面将丢失未保存页面!";
			}
			$("input[type=submit]").click(function(){
				$(window).unbind("beforeunload");
			});
			el.change(function(){
				if (el.val().length > 0) {
					$(window).bind("beforeunload",function(e){
						if ($.browser.msie) {
							e.returnValue = msg;
						} else {
							return msg;
						}
					});
				} else {
					 $(window).unbind("beforeunload");
				}
			});
		},
		
		// 收藏
		favorite:function(el){
			var $el = $(el),
				topic_id = $el.data("id");
			if ($el.hasClass("small_bookmarked")) {
				var data = {type: "unfavorite"};
				$.ajax({url: "/topics/"+topic_id+"/favorite", data: data, type: "POST"});
				$el.attr("title","收藏");
				$el.attr("class","icon small_bookmark");
			} else {
				$.post("/topics/"+topic_id+"/favorite");
				$el.attr("title","取消收藏");
				$el.attr("class","icon small_bookmarked");
			}
			return false;
		},
		
		// 关注
		follow:function(el){
			var $el = $(el),
				topic_id = $el.data("id"),
				followed = $el.data("followed");
			if (followed) {
				$.ajax({url:"/topics/"+topic_id+"/unfollow",type:"POST"});
				$el.data("followed", false);
				$("i", el).attr("class", "icon small_follow");
			} else {
				$.ajax({url:"/topics/"+topic_id+"/follow",type:"POST"});
				$el.data("followed", true);
				$("i",el).attr("class", "icon small_followed");
			}
			return false;
		},
		analyzeAt: function(text) {
			var result = {
				floor:0,
				username:[]
			};
			
			String(text).replace(/^#(\d+)楼|[^@]*@([^\s@]{4,20})\s*/g, function (match, floor, username) {
				floor && (result.floor = floor) || result.username.push(username);
			});
			
			return result;
		}
	}
		
	$(document).ready(function(){
		var bodyEl = $("body");
		$("textarea").on("keydown","ctrl+return", function(el){
			if ($(el.target).val().trim().length > 0) {
				$("#reply > form").submit()
			}
			return false;
		});

		Topics.initCloseWarning($("textarea.closewarning"));
		
		// 不自动伸缩
		// $("textarea").autoGrow();
		
		// TODO 图片上传
		//Topics.initUploader();
		
		$("a.at_floor").on("click", function(a){
			var floor = $(this).data('floor');
			Topics.gotoFloor(floor);
		});

		// also highlight if hash is reply
		var matchResult = window.location.hash.match(/^#reply(\d+)$/);
		if (matchResult != null) {
			Topics.highlightReply($("#reply"+matchResult[1]));
		}
		
		// 绑定回复按钮
		$("a.small_reply").on("click", function(){
			return Topics.reply($(this).data("floor"), $(this).data("login"));
		});

		Topics.hookPreview($(".editor_toolbar"), $(".topic_editor"));
		/*bodyEl.on("keydown.m", function(el){
			$("#markdown_help_tip_modal").modal({keyboard: true, backdrop:true, show:true});
		});*/

		// @ Reply
		var usersMap = App.scanLogins($("#topic_show .leader a[data-author]"));
		$.extend(usersMap, App.scanLogins($("#replies span.name a")));
		var users = function(){
			var result = [];
			for (var username in usersMap) {
				// @不出自己
				if (username == MeUsername) {
					continue;
				}
				var user = {uid: usersMap[username].uid, username: username, name: usersMap[username].name};
				result.push(user);
			}
			return result;
		}();
		// console.log(users);
		
		App.atReplyable(".cell_comments_new textarea", users);
		// 回复表单提交
		$("#new_reply").submit(function(env){
			env.preventDefault();
			var text = $('#wmd-input').val().trim();
			if (text == "") {
				$('#reply_content').addClass('error');
				$('#alert_info').show().html("回复内容不能为空")
				return false;
			}
			var result = Topics.analyzeAt(text),
				usernameArr = result['username']
				uidArr = [];
			for (var i in usernameArr) {
				// TODO:只支持参与过当前页面的用户
				if (usernameArr[i] in usersMap) {
					uidArr.push(usersMap[usernameArr[i]].uid);
				}
			}
			$('#uid').val(uidArr.join(','));
			var self = $(this);
			$.post(self.attr('action'), self.serialize(), function(data){
				if (data.errno) {
					$('#reply_content').addClass('error');
					$('#alert_info').show().html(data.error)
				} else {
					/*
					var content = $('#wmd-input').val(),
						floor = parseInt($('.small_reply').last().data('floor')) + 1;
					var replyHtml = '<div class="reply" id="reply1">'+
					'<div class="pull-left face"><a href="/user/{{.me.username}}"><img alt="{{.me.username}}" class="uface" src="{{gravatar .me.email 48}}" style="width:48px;height:48px;"></a></div>'+
				  '<div class="infos">'+
					'<div class="info">'+
					  '<span class="name">'+
						'<a href="/user/{{.me.username}}" data-name="{{.me.username}}">{{.me.username}}</a>'+floor+'楼, <abbr class="timeago" title="">'+$.timeago(new Date())+'</abbr>'+
					  '</span>'+
					  '<span class="opts">'+
						'<a href="/topics/9173/replies/86742/edit" class="edit icon small_edit" data-uid="744" title="修改回帖"></a>'+
						'<a href="#" class="icon small_reply" data-floor="1" data-login="{{.me.username}}" title="回复此楼"></a>'+
					  '</span>'+
					'</div>'+
					'<div class="body reply-content"><p>'+content+'</p></div>'+
				  '</div>'+
				'</div>';
				  $('.items').append(replyHtml);
				  */
				  // TODO:
				  location.reload();
				}
			}, 'json');
			return false;
		});

		// App.tipEmojis("textarea");

		// Focus title field in new-topic page
		$("body.topics-controller.new-action #topic_title").focus();

		$('#wmd-input').on('focus', function(){
			$(this).parents('.control-group').removeClass('error');
			$('#alert_info').hide();
		});
	});
}).call(this)
