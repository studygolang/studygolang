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

		// 编辑 tab
		$('.page').on('click', '.comment-edit-tab', function(evt){
			evt.preventDefault();

			var $this = $(this);
			var $tabMenu = $this.parent()
			var commentGroup = $tabMenu.data('comment-group')
			$this.addClass('cur');
			$tabMenu.children('.comment-preview-tab').removeClass('cur')

			$('.comment-content-preview[data-comment-group="' + commentGroup + '"]').hide();
			$('.comment-content-text[data-comment-group="' + commentGroup + '"]').show();
		});
		// 点击预览 tab
		$('.page').on('click', '.comment-preview-tab', function(evt){
			evt.preventDefault();

			var marked = SG.markSettingNoHightlight();

			var $this = $(this).addClass('cur');
			var $tabMenu = $this.parent();
			var commentGroup = $tabMenu.data('comment-group')
			var $preview = $('.comment-content-preview[data-comment-group="' + commentGroup + '"]')
			var $text = $('.comment-content-text[data-comment-group="' + commentGroup + '"]')
			$tabMenu.children('.comment-edit-tab').removeClass('cur');

			$text.hide();
			var content = $text.children('textarea').val();
			$preview.html(marked(content));
			// emoji 表情解析
			emojify.run($preview.get(0));
			$preview.show();

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

		// 切换显示评论和编辑评论
		function toggleCommentShowOrEdit(floor, show) {
			var $markdown = $('.markdown[data-floor="' + floor + '"]')
			var $content = $markdown.children('.content')
			var $editWrapper = $markdown.children('.edit-wrapper')
			if (show) {
				$content.show()
				$editWrapper.hide()
			} else {
				$content.hide()
				$editWrapper.show()
				var $textarea = $editWrapper.children('textarea')
				$textarea.val($textarea.data('raw-content')).focus()
			}
		}

		// 点击编辑评论按钮
		$('#replies').on('click', '.btn-edit', function(evt) {
			evt.preventDefault()
			var floor = $(this).data('floor')
			var $markdown = $('.markdown[data-floor="' + floor + '"]')
			var $editWrapper = $markdown.children('.edit-wrapper')
			var $textarea = $editWrapper.children('textarea')
			toggleCommentShowOrEdit(floor, false)

			var $uploadBtn = $('.upload-img[data-floor="' + floor + '"]')

			// 复制上传
			// 防止重复上传
			var pasteUpload = $textarea.data('paste-uploader')
			if (!pasteUpload) {
				pasteUpload = $textarea.pasteUploadImage('/image/paste_upload')
				$textarea.data('paste-uploader', pasteUpload)
			}

			// 点击按钮上传
			// 防止重复上传
			var uploader = $uploadBtn.data('uploader')
			if (!uploader) {
				uploader = window.initPLUpload({
					ele: $uploadBtn[0]
				})
				$uploadBtn.data('uploader', uploader)
			}
		});

		// 点击取消编辑评论按钮
		$('#replies').on('click', '.btn.cancel', function(evt) {
			evt.stopPropagation();
			var floor = $(this).data('floor');
			toggleCommentShowOrEdit(floor, true)
		})

		// 点击提交编辑后的评论
		$('#replies').on('click', '.btn.submit', function(evt) {
			evt.stopPropagation();
			var floor = $(this).data('floor');
			var $markdown = $('.markdown[data-floor="' + floor + '"]')
			var $submitBtn = $(this)
			var $editWrapper = $markdown.children('.edit-wrapper')
			var $textarea = $editWrapper.find('textarea')
			var $content = $markdown.children('.content')
			var content = $textarea.val()
			var cid = $submitBtn.data("cid")

			editComment($submitBtn, cid, content, function() {
				$textarea.data('raw-content', content)
				$content.html(parseCmtContent(content))
				toggleCommentShowOrEdit(floor, true)
			})
		})

		// 点击回复某人
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
		window.loadComments = function(p) {
			// 默认取最后一页
			p = p || 0;

			var objid = $('.comment-list').data('objid'),
				objtype = $('.comment-list').data('objtype');

			var params = {
				'objid': objid,
				'objtype': objtype,
				'p': p
			};
			$.getJSON('/object/comments', params, function(data){
				if (data.ok) {
					data = data.data;
					var comments = data.comments,
						replyComments = data.reply_comments;

					var content = '';
					for(var i in comments) {
						var comment = comments[i],
							meUid = $('[name="me-uid"]').val(),
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
							var replyComment = replyComments[comment.reply_floor]
							comment.reply_user = data[replyComment.uid];
							comment.reply_content = replyComment.content;
						}
						comment.rawContent = comment.content
						comment.content = parseCmtContent(comment.content);
						content += $.templates('#one-comment').render({comment: comment, user: user, me: {uid: meUid}});
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

		var editComment = function(thiss, cid, content, callback) {
			thiss.text("稍等").addClass("disabled").attr({"title":'稍等',"disabled":"disabled"});

			$.ajax({
				type:"post",
				url: '/object/comments/' + cid,
				data: {
					"content": content,
				},
				dataType: 'json',
				success: function(data){
					if(data.ok) {
						comTip("修改成功！");
						callback()
						thiss.text("提交").removeClass("disabled").removeAttr("disabled").attr({"title":'提交'});
					} else {
						alert(data.error);
					}
				},
				error: function() {
					thiss.text("提交").removeClass("disabled").removeAttr("disabled").attr({"title":'提交'});
				}
			})
		}

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
							meUid = $('[name="me-uid"]').val(),
							user = {};

						user.username = $pageComment.data('username'),
						user.uid = $pageComment.data('uid'),
						user.avatar = $pageComment.data('avatar'),
						comment.cmt_time = SG.timeago(comment.ctime);
						if (comment.reply_floor > 0) {
							comment.content = content.substr(1);
						}
						comment.reply_floor = 0;
						comment.rawContent = comment.content
						comment.content = parseCmtContent(comment.content);

						var oneCmt = $.templates('#one-comment').render({comment: comment, user: user, is_new: true, me: {uid: meUid}});

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

	////// 评论翻页 ///////////
	$('.page_input').on('keydown', function(event) {
		if (event.keyCode == 13) {
			var p = $(this).val();
			$('.cmt-page .page-num a:nth-child('+p+')').trigger('click');
		}
	});

	$('.ctrl-page button').on('click', function() {
		var p = $('.cmt-page .page_input').val();

		if ($(this).hasClass('prev-page')) {
			p--;
		} else {
			p++;
		}

		$('.cmt-page .page-num a:nth-child('+p+')').trigger('click');
	});

	$('.ctrl-page button').on('mouseover', function() {
		if (!$(this).hasClass('disable_now')) {
			$(this).addClass('hover_now');
		}
	});

	$('.ctrl-page button').on('mousedown', function() {
		$(this).addClass('active_now');
	});

	$('.ctrl-page button').on('mouseleave', function() {
		$(this).removeClass('hover_now');
		$(this).removeClass('active_now');
	});

	$('.cmt-page .page-num a').on('click', function(evt) {
		evt.preventDefault();
		$('.page-num .page_current').removeClass('page_current').addClass('page_normal');

		var p = $(this).data('page'),
			pageMax = $('.cmt-page .page_input').attr("max");

		$('.cmt-page .page-num a:nth-child('+p+')').removeClass('page_normal').addClass('page_current')
		$('.page-num .page_input').val(p);

		$('.cmt-page .ctrl-page button')
			.removeClass('disable_now')
			.removeAttr("disabled");

		if (p == 1) {
			$('.cmt-page .prev-page')
				.removeClass('hover_now')
				.removeClass('active_now')
				.addClass('disable_now')
				.attr("disabled", "disabled");
		} else if (p == pageMax) {
			$('.cmt-page .next-page')
				.removeClass('hover_now')
				.removeClass('active_now')
				.addClass('disable_now')
				.attr("disabled", "disabled");
		}

		loadComments(p);

		return false;
	});
	/////////// 评论翻页 end //////////////

}).call(this);
