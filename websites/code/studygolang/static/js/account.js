// 账号（登录/注册/忘记密码等）相关js功能
(function(){
	SG.Register = function(){}
	SG.Register.prototype = new SG.Publisher();
	
	jQuery(document).ready(function($) {
		// 同步提交
		$('.submit').on('click', function(evt){
			evt.preventDefault();
			var $form = $(this).parents('.validate-form');
			var validator = $form.validate();
			if (!validator.form()) {
				return false;
			}

			$form.submit();
		});

		var origSrc = '';
		$('#captcha_img').on('click', function(evt){
			evt.preventDefault();

			if (origSrc == '') {
				origSrc = $(this).attr("src");
			}
			$(this).attr("src", origSrc+"?reload=" + (new Date()).getTime());
		});
		
		// 异步提交
		// $('#register-submit').on('click', function(evt){
		// 	evt.preventDefault();
		// 	var validator = $('.validate-form').validate();
		// 	if (!validator.form()) {
		// 		return false;
		// 	}

		// 	new SG.Register().publish(this);
		// });
	});
}).call(this)
