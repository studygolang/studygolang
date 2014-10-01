$(function(){
	$('.sidebar .top ul li').on('mouseenter', function(evt){
		
		if (evt.target.tagName != 'LI') {
			return;
		}
		$(this).parent().find('a').removeClass('cur');
		$(this).children('a').addClass('cur');

		var sbContent = $(this).parents('.top').next();
		var left = 0;

		sbContent.children().removeClass('hidden').hide();
		switch ($(this).attr('class')) {
		case 'first':
			sbContent.children('.first').show();
			left = "18px";
			break;
		case 'second':
			sbContent.children('.second').show();
			left = "114px";
			break;
		case 'last':
			sbContent.children('.last').show();
			left = "210px";
			break;
		}
		$(this).parents('.top').children('.bar').animate({left: left}, "fast");
	});
	
	// º”‘ÿ≤‡±ﬂ¿∏
	$.getJSON("/topics/recent.json", function(data){
		if (data.ok) {
			data = data.data;

			var content = '';
			for(var i in data) {
				content += '<li>'+
						'<a href="/topics/'+data[i].tid+'" title="'+data[i].title+'">'+data[i].title+'</a>'+
						'</li>'
			}
			$('.sb-content .first ul').html(content);
		}
	});

	// º”‘ÿ≤‡±ﬂ¿∏
	$.getJSON("/articles/recent.json", function(data){
		if (data.ok) {
			data = data.data;

			var content = '';
			for(var i in data) {
				content += '<li>'+
						'<a href="/articles/'+data[i].id+'" title="'+data[i].title+'">'+data[i].title+'</a>'+
						'</li>'
			}
			$('.sb-content .second ul').html(content);
		}
	});

	// º”‘ÿ≤‡±ﬂ¿∏
	$.getJSON("/comment/recent.json", function(data){
		if (data.ok) {
			data = data.data;

			var content = '';
			for(var i in data) {
				content += '<li>'+
						'<a href="/articles/'+data[i].id+'" title="'+data[i].title+'">'+data[i].title+'</a>'+
						'</li>'
			}
			$('.sb-content .second ul').html(content);
		}
	});
});