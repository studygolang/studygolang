(function(){
	window.SocialShareButton = {
		openUrl: function(url) {
			window.open(url);
			return false;
		},
		share: function(el) {
			var $el = $(el),
				site = $el.data('site'),
				title = encodeURIComponent($el.parent().data('title') || ''),
				img = encodeURIComponent($el.parent().data("img") || ''),
				url = encodeURIComponent($el.parent().data("url") || '');
			if (url.length == 0) {
				// 当前网址
				url = encodeURIComponent(location.href);
			}
			switch (site) {
			case 'weibo':
				SocialShareButton.openUrl('http://service.weibo.com/share/share.php?url='+url+'&type=3&pic='+img+'&title='+title);
				break;
			case 'douban':
				SocialShareButton.openUrl('http://shuo.douban.com/!service/share?href='+url+'&image='+img+'&name='+title);
				break;
			case 'facebook':
				SocialShareButton.openUrl("http://www.facebook.com/sharer.php?t="+title+"&u="+url);
				break;
			case "qq":
				SocialShareButton.openUrl("http://sns.qzone.qq.com/cgi-bin/qzshare/cgi_qzshare_onekey?url="+url+"&title="+title+"&pics="+img);
				break;
			case "tqq":
				SocialShareButton.openUrl("http://share.v.t.qq.com/index.php?c=share&a=index&url="+url+"&title="+title+"&pic="+img);
				break;
			case "baidu":
				SocialShareButton.openUrl("http://hi.baidu.com/pub/show/share?url="+url+"&title="+title+"&content=");
				break;
			case "kaixin001":
				SocialShareButton.openUrl("http://www.kaixin001.com/rest/records.php?url="+url+"&content="+title+"&style=11&pic="+img);
				break;
			case "renren":
				SocialShareButton.openUrl("http://widget.renren.com/dialog/share?resourceUrl="+url+"&title="+title+"&description=");
				break;
			case "google_plus":
				SocialShareButton.openUrl("https://plus.google.com/share?url="+url+"&t="+title);
				break;
			case "google_bookmark":
				SocialShareButton.openUrl("https://www.google.com/bookmarks/mark?op=edit&output=popup&bkmk="+url+"&title="+title);
				break;
			case "delicious":
				SocialShareButton.openUrl("http://www.delicious.com/save?url="+url+"&title="+title+"&jump=yes&pic="+img);
				break;
			}
			return false;
		}
	}
}).call(this)