// 用户相关js功能
(function(){
	SG.User = function(){}
	SG.User.prototype = {
		edit: function(that) {
			var btnTxt = $(that).text();
			$(that).text("稍等").addClass("disabled").attr({"title":'稍等',"disabled":"disabled"});

			var $form = $(that).parents('form'),
				data = $form.serialize(),
				url = $form.attr('action');

			$.ajax({
				type:"post",
				url: url,
				data: data,
				dataType: 'json',
				success: function(data){
					if(data.ok){
						comTip("修改成功！");
						setTimeout(function(){
							window.location.reload();
						}, 1000);
					}else{
						comTip(data.error);
					}
				},
				complete:function(xmlReq, textStatus){
					$(that).text(btnTxt).removeClass("disabled").removeAttr("disabled").attr({"title":btnTxt});
				},
				error:function(xmlReq, textStatus, errorThrown){
					$(that).text(btnTxt).removeClass("disabled").removeAttr("disabled").attr({"title":btnTxt});
					if (xmlReq.status == 403) {
						comTip("没有编辑权限");
					}
				}
			});
		},
		parseCmtContent: function(selector) {
			selector.each(function(){
				var markdownString = $(this).html();
				// 配置 marked 语法高亮
				marked = SG.markSettingNoHightlight();

				$(this).html(marked(markdownString));

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
				
				// emoji 表情解析
				emojify.run(this);
			});
		}
	}

	jQuery(document).ready(function($) {
		var user = new SG.User();
		user.parseCmtContent($('.recent-comments ul li .content'));

		// 提交
		$('.submit').on('click', function(evt){
			evt.preventDefault();
			var validator = $(this).parents('.validate-form').validate();
			if (!validator.form()) {
				return false;
			}

			user.edit(this);
		});

		$('#active_email').on('click', function(evt){
			evt.preventDefault();

			$.ajax({
				type:"post",
				url: "/account/send_activate_email",
				dataType: 'json',
				success: function(data){
					if(data.ok){
						comTip("激活邮件已发到您邮箱，请查收！");
					}else{
						comTip(data.error);
					}
				},
				error:function(xmlReq, textStatus, errorThrown){
					if (xmlReq.status == 403) {
						comTip("没有操作权限");
					}
				}
			});

			return false;
		});
		
		$('#avatar-tab a').click(function (evt) {
			evt.preventDefault();
			$(this).tab('show');
		});

		$('.btn-gravatar').on('click', function(evt){
			evt.preventDefault();

			var url = $(this).attr('href');
			$.ajax({
				type:"post",
				url: url,
				data: {avatar: ''},
				dataType: 'json',
				success: function(data){
					if(data.ok){
						comTip("操作成功！");
						setTimeout(function(){
							window.location.reload();
						}, 1000);
					}else{
						comTip(data.error);
					}
				},
				error:function(xmlReq, textStatus, errorThrown){
					if (xmlReq.status == 403) {
						comTip("没有操作权限");
					}
				}
			});
		});

		// 实例化一个plupload上传对象
		var uploader = new plupload.Uploader({
			browse_button : 'btn-upload-avatar', // 触发文件选择对话框的按钮，为那个元素id
			url : '/image/upload', // 服务器端的上传页面地址
			filters: {
				mime_types : [ // 只允许上传图片
					{ title : "图片文件", extensions : "jpg,png" }
				],
				max_file_size : '500k', // 最大只能上传 500kb 的文件
				prevent_duplicates : true // 不允许选取重复文件
			},
			multipart_params:{
				avatar: '1'	// 上传的是头像
			},
			multi_selection: false,
			file_data_name: 'img',
			resize: {
				width: 600
			}
		});

		// 在实例对象上调用init()方法进行初始化
		uploader.init();

		uploader.bind('FilesAdded',function(uploader, files){
			// 调用实例对象的start()
			uploader.start();
		});
		uploader.bind('UploadProgress',function(uploader,file){
			// 上传进度
		});
		uploader.bind('FileUploaded',function(uploader,file,responseObject){
			if (responseObject.status == 200) {
				var data = $.parseJSON(responseObject.response);
				if (data.ok) {
					var path = data.data.uri;
					var url = data.data.url;
					var $img = $('#img-preview').find('img');
					$img.attr('src', url);
					$img.attr('alt', file.name);
					$('#img-preview').show();

					$('#upload-avatar').val(path.substr(7));

					$('#upload-btn').removeAttr("disabled");
				} else {
					comTip("上传失败："+data.error);
				}
			} else {
				comTip("上传失败：HTTP状态码："+responseObject.status);
			}
		});
		uploader.bind('Error',function(uploader,errObject){
			comTip("上传出错了："+errObject.message);
		});
	});
}).call(this);
