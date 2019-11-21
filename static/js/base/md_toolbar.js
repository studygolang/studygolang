// markdown tool bar 相关功能
(function(){
	jQuery(document).ready(function($) {
		$('form .md-toolbar .edit').on('click', function(evt){
			evt.preventDefault();
			
			$(this).addClass('cur');

			var $mdToobar = $(this).parents('.md-toolbar');
			$mdToobar.find('.preview').removeClass('cur');

			$mdToobar.nextAll('.content-preview').hide();
			$mdToobar.next().show();
		});
		
		$('form .md-toolbar .preview').on('click', function(evt){
			evt.preventDefault();

			// 配置 marked 语法高亮
			marked = SG.markSettingNoHightlight();

			$(this).addClass('cur');
			var $mdToobar = $(this).parents('.md-toolbar');
			$mdToobar.find('.edit').removeClass('cur');

			var $textarea = $mdToobar.next();
			$textarea.hide();
			var content = $textarea.val();
			var $contentPreview = $mdToobar.nextAll('.content-preview');
			$contentPreview.html(marked(content));
			$contentPreview.show();
		});

		$('form .preview_btn').on('click', function(evt) {
			evt.preventDefault();

			// 配置 marked 语法高亮
			marked = SG.markSettingNoHightlight();

			var $mdToobar = $('form .md-toolbar');
			$mdToobar.find('.preview').addClass('cur');
			$mdToobar.find('.edit').removeClass('cur');

			var $textarea = $mdToobar.next();
			$textarea.hide();
			var content = $textarea.val();
			var $contentPreview = $mdToobar.nextAll('.content-preview');
			$contentPreview.html(marked(content));
			$contentPreview.show();
		});
	});
}).call(this);
