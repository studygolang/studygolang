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

		$('#comment-content').on('change', function() {
			var content = $(this).val();

			var objdata = {content: content};

			saveReplyDraft(uid, keyprefix, objid, objdata);
		});

		(function() {
			if (typeof keyprefix === "undefined") {
				return;
			}
			var draft = loadReplyDraft(uid, keyprefix, objid);
			if (draft) {
				$('#comment-content').val(draft.content);
			}
		})();

		$('.page-comment .md-toolbar .edit').on('click', function(evt){
			evt.preventDefault();
			
			$(this).addClass('cur');
			$('.page-comment .md-toolbar .preview').removeClass('cur');

			$('.page-comment .content-preview').hide();
			$('.page-comment #commentForm .text').show();
		});
		$('.page-comment .md-toolbar .preview').on('click', function(evt){
			evt.preventDefault();

			var marked = SG.markSettingNoHightlight();

			$(this).addClass('cur');
			$('.page-comment .md-toolbar .edit').removeClass('cur');

			$('.page-comment #commentForm .text').hide();
			var content = $('.page-comment #commentForm textarea').val();
			$('.page-comment .content-preview').html(marked(content));
			// emoji 表情解析
			emojify.run($('.page-comment .content-preview').get(0));
			$('.page-comment .content-preview').show();

			Prism.highlightAll();
		});

		$('#replies').on('mouseenter', '.reply', function(evt) {
			$(this).find('.op-reply').removeClass('hideable');
		});
		$('#replies').on('mouseleave', '.reply', function(evt) {
			$(this).find('.op-reply').addClass('hideable');
		});

		$('#replies').on('click', '.reply_user', function(evt) {
			if ($(evt.target).hasClass('reply_user')) {
				$(this).parents('.reply-to-block').find('.markdown').toggleClass('dn');
			}
		});

		$('#replies').on('click', '.btn-reply', function(evt) {
			evt.preventDefault();

			var floor = $(this).data('floor'),
				username = $(this).data('username');
			var $replyTo = $('.md-toolbar .reply-to');

			$replyTo.data('floor', floor).data('username', username);

			var title = '回复#'+floor+'楼';
			$replyTo.children('.fa-mail-reply').attr('title', title);
			$replyTo.children('.user').attr('title', title).attr('href', '#reply'+floor).text(username+' #'+floor);
			$replyTo.removeClass('dn');

			$('#commentForm textarea').focus();
		});

		$('.md-toolbar .reply-to .close').on('click', function(evt) {
			evt.preventDefault();
			$(this).parents('.reply-to').addClass('dn').data('floor', '').data('username', '');
		});

		// 支持粘贴上传图片
		$('#comment-content').pasteUploadImage('/image/paste_upload');

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

		// 异步加载 评论
		window.loadComments = function() {
			var objid = $('.comment-list').data('objid'),
				objtype = $('.comment-list').data('objtype');
			
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
							if (isHttps) {
								user.avatar = 'https://secure.gravatar.com/avatar/'+md5(user.email)+"?s=48";
							} else {
								user.avatar = 'http://gravatar.com/avatar/'+md5(user.email)+"?s=48";
							}
						} else if (avatar.indexOf('http') === -1) {
							user.avatar = cdnDomain+'avatar/'+avatar+'?imageView2/2/w/48';
						}

						var cmtTime = SG.timeago(comment.ctime);
						if (cmtTime == comment.ctime) {
							var cmtTimes = cmtTime.split(" ");
							comment.cmt_time = cmtTimes[0];
						} else {
							comment.cmt_time = cmtTime;
						}

						if (comment.reply_floor > 0) {
							var replyComment = comments[comment.reply_floor-1]
							comment.reply_user = data[replyComment.uid];
							comment.reply_content = replyComment.content;
						}

						comment.content = parseCmtContent(comment.content);

						content += $.templates('#one-comment').render({comment: comment, user: user});
					}

					if (content != '') {
						$('.comment-list .words').html(content);

						// 链接，add target=_blank
						$('.comment-list .words .markdown').on('mousedown', 'a', function(evt){
							var url = $(this).attr('href');
							$(this).attr('target', '_blank');
						});

						$('.comment-list .markdown img').attr('data-action', 'zoom');

						$('.comment-list .markdown img').on('click', function() {
							$(this).parents('.box_white').css('overflow', 'visible');
						});
					}
					$('.comment-list .words').removeClass('hide');
					$('.comment-list .words').find('code[class*="language-"]').parent('pre').addClass('line-numbers');
					Prism.highlightAll();

					// emoji 表情解析
					emojify.run($('.comment-list .words').get(0));

					if ($("#is_login_status").val() == 1) {
						SG.registerAtEvent(true, true, $('.page-comment textarea'));
					}
				} else {
					comTip("回复加载失败");
				}
			});
		}

		var parseCmtContent = function(content) {
			var marked = SG.markSettingNoHightlight();
			content = SG.preProcess(content);
			content = marked(content);
			return SG.replaceCodeChar(content);
		};

		// 回复提交
		$('#comment-submit').on('click', function(){
			var content = $('#commentForm textarea').val();

			if(content == ""){
				alert("其实你想说点什么...");
			} else {
				var floor = $('.md-toolbar .reply-to').data('floor');
				if (parseInt(floor, 10) > 0) {
					var username = $('.md-toolbar .reply-to').data('username');
					content = '#'+floor+'楼 @'+username+' '+content;
				}
				postComment($(this), content, function(comment) {
					comTip("回复成功！");
					purgeReplyDraft(uid, keyprefix, objid);

					$('#commentForm textarea').val('');

					$('.md-toolbar .reply-to .close').click();
				});
			}
		});

		var postComment = function(thiss, content, callback){
			thiss.text("稍等").addClass("disabled").attr({"title":'稍等',"disabled":"disabled"});

			var objid = $('.comment-list').data('objid'),
				objtype = $('.comment-list').data('objtype');

			var usernames = SG.analyzeAt(content);
			
			$.ajax({
				type:"post",
				url: '/comment/'+objid,
				data: {
					"objtype": objtype,
					"content": content,
					"usernames": usernames.join(',')
				},
				dataType: 'json',
				success: function(data){
					if(data.ok){
						var comment = data.data;

						var $pageComment = $('.comment-list'),
							user = {};
						
						user.username = $pageComment.data('username'),
						user.uid = $pageComment.data('uid'),
						user.avatar = $pageComment.data('avatar'),
						comment.cmt_time = SG.timeago(comment.ctime);
						if (comment.reply_floor > 0) {
							comment.content = content.substr(1);
						}
						comment.reply_floor = 0;
						comment.content = parseCmtContent(comment.content);

						var oneCmt = $.templates('#one-comment').render({comment: comment, user: user, is_new: true});

						var $cmtNumObj = $('#replies .cmtnum'),
							cmtNum = parseInt($cmtNumObj.text(), 10);
						if (cmtNum == 0) {
							$('.comment-list .words').html('');
						}

						$('.comment-list .words').append(oneCmt).removeClass('hide');
						Prism.highlightAll();

						// emoji 表情解析
						emojify.run($('.comment-list .words .reply:last').get(0));
						
						// 注册@
						SG.registerAtEvent(true, true, $('.page-comment textarea'));

						cmtNum++;

						$cmtNumObj.text(cmtNum);
						
						setTimeout(function(){
							$('.comment-list .words .reply').removeClass('light');
						}, 2000);
						callback();
					} else {
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
}).call(this);