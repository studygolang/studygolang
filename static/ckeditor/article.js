$(function(){
	CKEDITOR.plugins.addExternal('prism', '/static/ckeditor/plugins/prism/', 'plugin.js');
	$('#edit').on('click', function(){
		var txt = $(this).text();
		if (txt == '编辑') {
			$('#myeditor').attr('contenteditable', true);
			$('#myeditor').html($('#content_tpl').html());
			if (!CKEDITOR.instances.myeditor) {
				MyEditorConfig.extraPlugins = MyEditorExtraPlugins+',prism,sourcedialog';
				MyEditorConfig.toolbarGroups = [
					{ name: 'undo' },
					{ name: 'basicstyles', groups: [ 'basicstyles', 'cleanup' ] },
					{ name: 'paragraph', groups: [ 'list', 'indent', 'blocks', 'align' ] },
					{ name: 'links' },
					{ name: 'insert' },
					{ name: 'styles' },
					{ name: 'document', groups: [ 'mode', 'document' ] }
				];
				MyEditorConfig.removeButtons = 'Anchor,SpecialChar,HorizontalRule,Table,Styles,Subscript,Superscript';
				CKEDITOR.inline( 'myeditor', MyEditorConfig );
			}

			$(this).text('完成');
		} else {
			if (CKEDITOR.instances.myeditor) {
				var content = CKEDITOR.instances.myeditor.getData();
				modify(content);

				CKEDITOR.instances.myeditor.destroy();

				Prism.highlightAll();
			}

			$('#myeditor').attr('contenteditable', false);
			$(this).text('编辑');
		}
	});

	CKEDITOR.on('instanceReady', function(evt, editor) {
		$('#myeditor').find('.cke_widget_element').each(function(){
			$(this).addClass('line-numbers').css('background-color', '#000');
		});
	});

	function modify(content)
	{
		var url = '/articles/modify',
			data = { id: $('#title').data('id'), content:content };

		$.ajax({
			type: "post",
			url: url,
			data: data,
			dataType: 'json',
			success: function(data){
				if(data.ok){
					if (typeof data.msg != "undefined") {
						comTip(data.msg);
					} else {
						comTip("修改成功！");
					}
				}else{
					comTip(data.error);
				}
			},
			complete:function(xmlReq, textStatus){
			},
			error:function(xmlReq, textStatus, errorThrown){
				if (xmlReq.status == 403) {
					comTip("没有修改权限");
				}
			}
		});
	}
});