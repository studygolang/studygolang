jQuery(document).ready(function(){
	
	$('.upload_img_single').Huploadify({
		auto: true,
		fileTypeExts: '*.png;*.jpg;*.JPG;*.bmp;*.gif',// 不限制上传文件请修改成'*.*'
		multi:false,
		fileSizeLimit: 5*1024*1024, // 大小限制
		uploader : '/image/upload', // 文件上传目标地址
		buttonText : '上传',
		fileObjName : 'img',
		showUploadedPercent:true,
		onUploadSuccess : function(file, data) {
			data = $.parseJSON(data);
			if (data.ok) {
				var url = data.data.url;
				$('.img_url').val(url);
				$('img.show_img').attr('src', url);
				$('a.show_img').attr('href', url);
			} else {
				if (window.jAlert) {
					jAlert(data.error, '错误');
				} else {
					alert(data.error);
				}
			}
		}
	});
});