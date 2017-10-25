<script src="/static/js/libs/jquery.timeago.zh-CN.js"></script>
<script src="//cdn.bootcss.com/zoom.js/0.0.1/zoom.min.js"></script>
<script src="/static/js/libs/md5.js"></script>
<script type="text/javascript">
var uid = {{.me.Uid}};
var isHttps = {{.is_https}},
	cdnDomain = "{{.cdn_domain}}";
if (isHttps) {
	var wsUrl = 'wss://{{.wshost}}/ws?uid='+uid;
} else {
	var wsUrl = 'ws://{{.wshost}}/ws?uid='+uid;
}
var GLaunchTime = {{timestamp .app.LaunchTime}}*1000;
</script>
<script src="/static/js/common.js"></script>
<script type="text/javascript" src="/static/js/libs/paste-upload-image.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/jsrender/0.9.84/jsrender.min.js"></script>
<script type="text/javascript">
$.views.settings.delimiters("[%", "%]");
// $.views.settings.debugMode(true);
</script>
<script src="//cdn.bootcss.com/emojify.js/1.1.0/js/emojify.min.js"></script>
