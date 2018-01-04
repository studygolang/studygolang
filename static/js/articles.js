// 文章相关js功能
(function(){
	SG.Articles = function(){}
	SG.Articles.prototype = new SG.Publisher();
	SG.Articles.prototype.parseContent = function(selector) {
		var markdownString = selector.text();
		marked = SG.markSettingNoHightlight();

		var contentHtml = marked(markdownString);
		contentHtml = SG.replaceCodeChar(contentHtml);
		
		selector.html(contentHtml);

		// emoji 表情解析
		emojify.run(selector.get(0));
	}

	jQuery(document).ready(function($) {
		$('#submit').on('click', function(evt){
			evt.preventDefault();
			var validator = $('.validate-form').validate();
			if (!validator.form()) {
				return false;
			}

			if ($('input[type=radio]:checked').val() == 0) {
				$('#content').val(CKEDITOR.instances.myeditor.getData());
				if (window.localStorage) {
					localStorage.removeItem('autosaveKey');
				}

				$('#txt').val(CKEDITOR.instances.myeditor.document.getBody().getText());
			} else {
				$('#content').val($('#markdown-content').val());
			}

			var articles = new SG.Articles();
			articles.publish(this, function(data) {
				if (typeof cacheKey == "undefined") {
					cacheKey = 'article';
				}
				purgeComposeDraft(uid, cacheKey);

				setTimeout(function(){
					if (data.id) {
						window.location.href = '/articles/'+data.id;
					} else {
						window.location.href = '/articles';
					}
				}, 1000);
			});
		});

		$(document).keypress(function(evt){
			if (evt.ctrlKey && (evt.which == 10 || evt.which == 13)) {
				$('#submit').click();
			}
		});

		///////////////////// 收入专栏相关操作 //////////////////////////
		$('.add-collection').on('click', function(evt) {
			evt.preventDefault();
			
			var articleId = $('#title').data('id');
			$.getJSON('/subject/mine?article_id='+articleId, function(result) {
				if (result.ok) {
					var subjects = result.data.subjects;
					fillSubjects(subjects);
	
					$('body').addClass('modal-open');
					$('.add-self').fadeIn();
				}
			});
		});
	
		$('.add-self .close').on('click', function() {
			$('body').removeClass('modal-open');
			$('.add-self').fadeOut();
		})
	
		var noteListHtml = '';
		$('.add-self .search-btn').on('click', function() {
			var kw = $('.add-self .search-input').val();
			if (kw == "") {
				$('#self-note-list').html(noteListHtml);
				return;
			}
	
			noteListHtml = $('#self-note-list').html();
			$('#self-note-list').html('');
			var placeholder = $('.add-self .modal-collections-placeholder');
			placeholder.show();
			
			var articleId = $('#title').data('id');
			$.getJSON('/subject/mine?kw='+encodeURIComponent(kw)+'&article_id='+articleId, function(result) {
				placeholder.hide();
				
				if (result.ok) {
					var subjects = result.data.subjects;
					if (subjects.length == 0) {
						$('#self-note-list').html('<div class="default">未找到相关专栏</div>');
					} else {
						fillSubjects(subjects);
					}
				} else {
					$('#self-note-list').html('<div class="default">'+result.msg+'</div>');
				}
			})
		})
	
		$('.add-self .search-input').on('change', function() {
			if ($(this).val() == '') {
				$('#self-note-list').html(noteListHtml);
			}
		});
	
		$(document).keypress(function(evt){
			if (evt.which == 10 || evt.which == 13) {
				$('.add-self .search-btn').click();
			}
		});
	
		$('.add-self').on('click', '.action-btn', function() {
			var $collectInfo = $(this).parent(),
				sid = $collectInfo.data('sid'),
				articleId = $('#title').data('id');
			
			var that = this;
	
			if ($(this).hasClass('push')) {
				$.post('/subject/contribute', {sid: sid, article_id: articleId}, function(result) {
					if (result.ok) {
						$(that).removeClass('push').addClass('remove').
							before(' <span class="status has-add">已收入</span>').text('移除');
					} else {
						alert(result.error);
					}
				});
			} else {
				$.post('/subject/remove_contribute', {sid: sid, article_id: articleId}, function(result) {
					if (result.ok) {
						$(that).removeClass('remove').addClass('push').text('收入');
						$collectInfo.children('.status').remove();
					} else {
						alert(result.error);
					}
				});
			}
		});
	
		function fillSubjects(subjects) {
			var listHtml = '';
			for(var i in subjects) {
				var subject = subjects[i];
	
				listHtml += '<li>'+
					'<a href="/subject/'+subject.id+'" class="avatar-collection"><img src="'+subject.cover+'"></a>'+
					'<div class="collection-info" data-sid="'+subject.id+'">'+
						'<a href="/subject/'+subject.id+'" class="collection-name">'+subject.name+'</a>'+
						'<div class="meta">'+subject.username+' 编</div>';
				
				if (subject.had_add) {
					listHtml += ' <span class="status has-add">已收入</span>'+
						'<a class="action-btn remove">移除</a>';
				} else {
					listHtml += '<a class="action-btn push">收入</a>';
				}
	
				listHtml += '</div></li>';
			}
			$('#self-note-list').html(listHtml);
		}
	});
}).call(this);
