(function(){
	window.Projects = {
		publish: function(that){
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
						$form.get(0).reset();
						
						comTip("发布成功！");
						
						setTimeout(function(){
							window.location.href = "/projects";
						}, 3000);
					}else{
						alert(data.error);
					}
				},
				complete:function(){
					$(that).text(btnTxt).removeClass("disabled").removeAttr("disabled").attr({"title":btnTxt});
				},
				error:function(){
					$(that).text(btnTxt).removeClass("disabled").removeAttr("disabled").attr({"title":btnTxt});
				}
			});
		},

		parseDesc: function(){
			var markdownString = $('.project .desc').html();
			marked.setOptions({
				highlight: function (code) {
					return hljs.highlightAuto(code).value;
				}
			});

			$('.project .desc').html(marked(markdownString));
		}
	};
	
	jQuery(document).ready(function($) {
		var IS_PREVIEW = false;
		$('.preview').on('click', function(){
			// console.log(hljs.listLanguages());
			if (IS_PREVIEW) {
				$('.preview-div').hide();
				$('#desc').show();
				IS_PREVIEW = false;
			} else {
				var markdownString = $('#desc').val();
				marked.setOptions({
					highlight: function (code) {
						return hljs.highlightAuto(code).value;
					}
				});

				$('#desc').hide();
				$('.preview-div').html(marked(markdownString)).show();
				IS_PREVIEW = true;
			}
		});

		// 发布项目
		$('#submit').on('click', function(evt){
			evt.preventDefault();
			var validator = $('.validate-form').validate();
			if (!validator.form()) {
				return false;
			}
			Projects.publish(this);
		});
	});
}).call(this)