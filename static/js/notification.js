// 桌面通知功能
(function(){
	// 把func中的this指向obj对象（或者理解为obj拥有了func方法）
	var applyTool = function(func, obj){
		return function(){
			// apply：改变func中的this指向obj，并立即执行func
			return func.apply(obj, arguments);
		}
	};
	
	// chrome 通知功能Notifications对象
	var Notifier = function(){
		function Notify(){
			this.checkOrRequirePermission = applyTool(this.checkOrRequirePermission, this);
			this.setPermission = applyTool(this.setPermission, this);
			this.enableNotification = false;
			this.checkOrRequirePermission();
		}

		Notify.prototype.hasSupport = function(){return window.webkitNotifications != null;}

		Notify.prototype.requestPermission = function(cb){return window.webkitNotifications.requestPermission(cb);}
		
		Notify.prototype.setPermission = function(){
			if(this.hasPermission()) {
				$("#notification-alert a.close").click();
				this.enableNotification = true;
				if(window.webkitNotifications.checkPermission() === 2) {
					$("#notification-alert a.close").click()
				}
			}
		}

		Notify.prototype.hasPermission = function(){return window.webkitNotifications.checkPermission() === 1;}
		
		Notify.prototype.checkOrRequirePermission = function(){
			if(!this.hasSupport()) {
				console.log("Desktop notifications are not supported for this Browser/OS version yet.");
			}
			if(this.hasPermission()) {
				this.enableNotification = true;
				if(window.webkitNotifications.checkPermission() !== 2) {
					this.showTooltip();
				}
			}
		}
		
		Notify.prototype.showTooltip = function(){
			var self = this;
			$(".breadcrumb").before("<div class='alert alert-info' id='notification-alert'><a href='#' id='link_enable_notifications' style='color:green'>点击这里</a> 开启桌面提醒通知功能。 <a class='close' data-dismiss='alert' href='#'>×</a></div>");
			$("#notification-alert").alert();
			$("#notification-alert").on("click", "a#link_enable_notifications", function(env){
				env.preventDefault();
				self.requestPermission(self.setPermission);
			});
		}
		
		Notify.prototype.visitUrl = function(url) {return window.location.href = url}

		Notify.prototype.notify = function(avatar, title, content, url) {
			var obj, notification;
			if(this.enableNotification && window.Notification) {
				obj = {
					body: content,
					onclick: function() {
						window.parent.focus();
						$.notifier.visitUrl(url);
					}
				};
				notification = new window.Notification(title, obj);
			} else {
				notification = window.webkitNotifications.createNotification(avatar, title, content);
				if (url) {
					notification.onclick = function(){
						return window.parent.focus();
						$.notifier.visitUrl(url);
					}
				}
			}
			notification.show();
		}
		
		return Notify;
	}();
	jQuery.notifier = new Notifier();
}).call(this)