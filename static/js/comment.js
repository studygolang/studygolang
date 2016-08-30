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

		$('.page-comment .md-toolbar .edit').on('click', function(evt){
			evt.preventDefault();
			
			$(this).addClass('cur');
			$('.page-comment .md-toolbar .preview').removeClass('cur');

			$('.page-comment .content-preview').hide();
			$('.page-comment #commentForm').show();
		});
		$('.page-comment .md-toolbar .preview').on('click', function(evt){
			evt.preventDefault();

			// 配置 marked 语法高亮
			marked.setOptions({
				highlight: function (code) {
					return hljs.highlightAuto(code).value;
				}
			});

			$(this).addClass('cur');
			$('.page-comment .md-toolbar .edit').removeClass('cur');

			$('.page-comment #commentForm').hide();
			var content = $('.page-comment #commentForm textarea').val();
			$('.page-comment .content-preview').html(marked(content));
			$('.page-comment .content-preview').show();
		});

		emojify.setConfig({
			// emojify_tag_type : 'span',
			only_crawl_id    : null,
			img_dir          : 'http://hassankhan.me/emojify.js/images/emoji',
			ignored_tags     : { //忽略以下几种标签内的emoji识别
				'SCRIPT'  : 1,
				'TEXTAREA': 1,
				'A'       : 1,
				'PRE'     : 1,
				'CODE'    : 1
			}
		});

		// 异步加载 评论
		window.loadComments = function() {
			var objid = $('.page-comment').data('objid'),
				objtype = $('.page-comment').data('objtype');
			
			var params = {
				'objid': objid,
				'objtype': objtype
			};
			$.getJSON('/object/comments', params, function(data){
				if (data.ok) {
					data = data.data;
					var comments = data.comments;

					var content = '';
					for(var i in comments) {
						var comment = comments[i],
							user = data[comment.uid];

						var avatar = user.avatar;
						if (avatar == "") {
							avatar = 'http://gravatar.duoshuo.com/avatar/'+md5(user.email)+"?s=48";
						} else {
							avatar = 'http://studygolang.qiniudn.com/avatar/'+avatar+'?imageView2/2/w/40';
						}
						
						var cmtTime = SG.timeago(comments[i].ctime);
						if (cmtTime == comments[i].ctime) {
							var cmtTimes = cmtTime.split(" ");
							cmtTime = cmtTimes[0];
						}
						content += contructOneCmt(comment.floor, user.uid, user.username, avatar, comment.content, comment.ctime, cmtTime);
					}

					if (content != "") {
						$('.page-comment .words ul').html(content);
						$('.page-comment .words').removeClass('hide');
					}

					// emoji 表情解析
					emojify.run($('.page-comment .words ul').get(0));
					// twitter emoji 表情解析
					/*
					var result = twemoji.parse($('.page-comment .words ul').get(0), {
						callback: function(icon, options, variant) {
							return ''.concat(options.base, options.size, '/', icon, options.ext);
						},
						size: 16
					});
					*/

					if ($("#is_login_status").val() == 1) {
						SG.registerAtEvent(true, true, $('.page-comment textarea'));
					}
				} else {
					comTip("评论加载失败");
				}
			});
		}

		var contructOneCmt = function(floor, uid, username, avatar, content, ctime, cmtTime, needLight) {
			var oneCmt = '<li id="comment'+floor+'">';
			if (typeof needLight !== "undefined") {
				oneCmt = '<li id="comment'+floor+'" class="light">';
			}
			
			// 配置 marked 语法高亮
			marked.setOptions({
				highlight: function (code) {
					code = code.replace(/&#34;/g, '"');
					code = code.replace(/&#39;/g, "'");
					code = code.replace(/&lt;/g, '<');
					code = code.replace(/&gt;/g, '>');
					code = code.replace(/&amp;/g, '&');
					return hljs.highlightAuto(code).value;
				}
			});
			content = marked(content);
			content = SG.replaceCodeChar(content);
			return oneCmt+
				'<div class="pull-left face">'+
					'<a href="/user/'+username+'" target="_blank"><img src="'+avatar+'" width="48px" height="48px" alt="'+username+'"></a>'+
				'</div>'+
				'<div class="cmt-body">'+
					'<div class="cmt-content">'+
						'<a href="/user/'+username+'" class="name replyName" target="_blank" data-floor="'+floor+'楼" data-uid="'+uid+'">'+username+'</a>：'+
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
							'<textarea class="need-autogrow reply-content" name="content"></textarea>'+
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
				// 隐藏所有回复表单，
				$(".page-comment .respond-submit").hide(10); 
				// 设置回复框的样式并默认选中
				var $cmtBody = thiss.parents(".cmt-body");
				tipLength(thiss, function(){
					var tipWid = $cmtBody.find(".respond-submit .text .tip").outerWidth();
					var cbWid = $cmtBody.width();
					$cmtBody.find(".respond-submit .reply-content").css({'width':cbWid,'padding-left':tipWid+12, 'padding-top':10});
					$cmtBody.find(".respond-submit .reply-content").focus();
				})

				// 显示当前表单
				setTimeout(function(){
					$cmtBody.find(".respond-submit").slideDown(300);
					// 文本框自动伸缩
					$('.need-autogrow').autoGrow();
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
			var	content	= $text.children(".reply-content").val();
			
			if(content == ""){
				alert("其实你想说点什么...");
			} else {
				var that = $(this)
				content = replyTo + content;
				postComment(that, content, function(){
					comTip("回复成功！");
					that.parent(".sub").prev(".text").children(".reply-content").val("");
					that.parents(".respond-submit").slideUp(200);
				});
			}
		});

		var postComment = function(thiss, content, callback){
			thiss.text("稍等").addClass("disabled").attr({"title":'稍等',"disabled":"disabled"});

			var objid = $('.page-comment').data('objid'),
				objtype = $('.page-comment').data('objtype');

			var usernames = SG.analyzeAt(content);
			
			$.ajax({
				type:"post",
				url: '/comment/'+objid,
				data: {
					"objtype":objtype,
					"content":content,
					"usernames": usernames.join(',')
				},
				dataType: 'json',
				success: function(data){
					if(data.ok){
						var comment = data.data;
						var $pageComment = $('.page-comment'),
						username = $pageComment.data('username'),
						uid = $pageComment.data('uid'),
						avatar = $pageComment.data('avatar'),
						cmtTime = SG.timeago(comment.ctime);
						var oneCmt = contructOneCmt(comment.floor, uid, username, avatar, comment.content, comment.ctime, cmtTime, true);

						$('.page-comment .words ul').append(oneCmt);
						$('.page-comment .words').removeClass('hide');

						// emoji 表情解析
						emojify.run($('.page-comment .words ul li:last').get(0));

						// twitter emoji 表情解析
						/*
						twemoji.parse($('.page-comment .words ul li:last').get(0));
						*/
						
						// 注册@
						SG.registerAtEvent(true, true, $('.page-comment textarea'));

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
}).call(this)