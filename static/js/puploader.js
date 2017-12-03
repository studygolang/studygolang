$(function(){
	// 实例化一个plupload上传对象
	var uploader = new plupload.Uploader({
		browse_button : 'upload-img', // 触发文件选择对话框的按钮，为那个元素id
		url : '/image/upload', // 服务器端的上传页面地址
		filters: {
			mime_types : [ //只允许上传图片
				{ title : "图片文件", extensions : "jpg,gif,png,bmp" }
			],
			max_file_size : '5mb', // 最大只能上传 5mb 的文件
			prevent_duplicates : true // 不允许选取重复文件
		},
		multi_selection: false,
		file_data_name: 'img'
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
	uploader.bind('FileUploaded', function(uploader, file, responseObject) {
		window.uploadSuccess(uploader, file, responseObject)
	});
	uploader.bind('Error',function(uploader,errObject){
		comTip("上传出错了："+errObject.message);
	});

	window.uploadSuccess = function(uploader,file,responseObject){
		if (responseObject.status == 200) {
			var data = $.parseJSON(responseObject.response);
			if (data.ok) {
				var text = $('.main-textarea').val();
				text += '!['+file.name+']('+data.data.url+')';
				$('.main-textarea').val(text);
			} else {
				comTip("上传失败："+data.error);
			}
		} else {
			comTip("上传失败：HTTP状态码："+responseObject.status);
		}
	}
});