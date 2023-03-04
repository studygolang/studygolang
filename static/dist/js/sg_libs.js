// emoji表情 http://www.emoji-cheat-sheet.com/
var emojis = [
	// people
    "bowtie", "smile", "laughing", "blush", "smiley", "relaxed", "smirk",
	"heart_eyes", "kissing_heart", "kissing_closed_eyes", "flushed", "relieved", "satisfied", "grin",
	"wink", "stuck_out_tongue_winking_eye", "stuck_out_tongue_closed_eyes", "grinning", "kissing", "kissing_smiling_eyes", "stuck_out_tongue", "sleeping",
	"worried", "frowning", "anguished", "open_mouth", "grimacing", "confused", "hushed", "expressionless",
	"unamused", "sweat_smile", "sweat", "disappointed_relieved", "weary", "pensive", "disappointed", "confounded",
	"fearful", "cold_sweat", "persevere", "cry", "sob", "joy", "astonished",
	"scream", "neckbeard", "tired_face", "angry", "rage", "triumph", "sleepy",
	"yum", "mask", "sunglasses", "dizzy_face", "imp", "smiling_imp", "neutral_face",
	"no_mouth", "innocent", "alien", "yellow_heart", "blue_heart", "purple_heart", "heart",
	"green_heart", "broken_heart", "heartbeat", "heartpulse", "two_hearts", "revolving_hearts", "cupid",
	"sparkling_heart", "sparkles", "star", "star2", "dizzy", "boom", "collision",
	"anger", "exclamation", "question", "grey_exclamation", "grey_question", "zzz", "dash",
	"sweat_drops", "notes", "musical_note", "fire", "hankey", "poop", "shit",
	"+1", "thumbsup", "-1", "thumbsdown", "ok_hand", "punch", "facepunch",
	"fist", "v", "wave", "hand", "raised_hand", "open_hands", "point_up",
	"point_down", "point_left", "point_right", "raised_hands", "pray", "point_up_2", "clap",
	"muscle", "metal", "fu", "walking", "runner", "running", "couple",
	"family", "two_men_holding_hands", "two_women_holding_hands", "dancer", "dancers", "ok_woman", "no_good",
	"information_desk_person", "raising_hand", "bride_with_veil", "person_with_pouting_face", "person_frowning", "bow", "couplekiss",
	"couple_with_heart", "massage", "haircut", "nail_care", "boy", "girl", "woman",
	"man", "baby", "older_woman", "older_man", "person_with_blond_hair", "man_with_gua_pi_mao", "man_with_turban",
	"construction_worker", "cop", "angel", "princess", "smiley_cat", "smile_cat", "heart_eyes_cat",
	"kissing_cat", "smirk_cat", "scream_cat", "crying_cat_face", "joy_cat", "pouting_cat", "japanese_ogre",
	"japanese_goblin", "see_no_evil", "hear_no_evil", "speak_no_evil", "guardsman", "skull", "feet",
	"lips", "kiss", "droplet", "ear", "eyes", "nose", "tongue",
	"love_letter", "bust_in_silhouette", "busts_in_silhouette", "speech_balloon", "thought_balloon", "feelsgood", "finnadie",
	"goberserk", "godmode", "hurtrealbad", "rage1", "rage2", "rage3", "rage4",
	"suspect", "trollface",
	// nature
	"sunny", "umbrella", "cloud", "snowflake", "snowman", "zap", "cyclone",
	"foggy", "ocean", "cat", "dog", "mouse", "hamster", "rabbit",
	"wolf", "frog", "tiger", "koala", "bear", "pig", "pig_nose",
	"cow", "boar", "monkey_face", "monkey", "horse", "racehorse", "camel",
	"sheep", "elephant", "panda_face", "snake", "bird", "baby_chick", "hatched_chick",
	"hatching_chick", "chicken", "penguin", "turtle", "bug", "honeybee", "ant",
	"beetle", "snail", "octopus", "tropical_fish", "fish", "whale", "whale2",
	"dolphin", "cow2", "ram", "rat", "water_buffalo", "tiger2", "rabbit2",
	"dragon", "goat", "rooster", "dog2", "pig2", "mouse2", "ox",
	"dragon_face", "blowfish", "crocodile", "dromedary_camel", "leopard", "cat2", "poodle",
	"paw_prints", "bouquet", "cherry_blossom", "tulip", "four_leaf_clover", "rose", "sunflower",
	"hibiscus", "maple_leaf", "leaves", "fallen_leaf", "herb", "mushroom", "cactus",
	"palm_tree", "evergreen_tree", "deciduous_tree", "chestnut", "seedling", "blossom", "ear_of_rice",
	"shell", "globe_with_meridians", "sun_with_face", "full_moon_with_face", "new_moon_with_face", "new_moon", "waxing_crescent_moon",
	"first_quarter_moon", "waxing_gibbous_moon", "full_moon", "waning_gibbous_moon", "last_quarter_moon", "waning_crescent_moon", "last_quarter_moon_with_face",
	"first_quarter_moon_with_face", "moon", "earth_africa", "earth_americas", "earth_asia", "volcano", "milky_way",
	"partly_sunny", "octocat", "squirrel",
	// objects
	"bamboo", "gift_heart", "dolls", "school_satchel", "mortar_board", "flags", "fireworks",
	"sparkler", "wind_chime", "rice_scene", "jack_o_lantern", "ghost", "santa", "christmas_tree",
	"gift", "bell", "no_bell", "tanabata_tree", "tada", "confetti_ball", "balloon",
	"crystal_ball", "cd", "dvd", "floppy_disk", "camera", "video_camera", "movie_camera",
	"computer", "tv", "iphone", "phone", "telephone", "telephone_receiver", "pager",
	"fax", "minidisc", "vhs", "sound", "speaker", "mute", "loudspeaker",
	"mega", "hourglass", "hourglass_flowing_sand", "alarm_clock", "watch", "radio", "satellite",
	"loop", "mag", "mag_right", "unlock", "lock", "lock_with_ink_pen", "closed_lock_with_key",
	"key", "bulb", "flashlight", "high_brightness", "low_brightness", "electric_plug", "battery",
	"calling", "email", "mailbox", "postbox", "bath", "bathtub", "shower",
	"toilet", "wrench", "nut_and_bolt", "hammer", "seat", "moneybag", "yen",
	"dollar", "pound", "euro", "credit_card", "money_with_wings", "e-mail", "inbox_tray",
	"outbox_tray", "envelope", "incoming_envelope", "postal_horn", "mailbox_closed", "mailbox_with_mail", "mailbox_with_no_mail",
	"package", "door", "smoking", "bomb", "gun", "hocho", "pill",
	"syringe", "page_facing_up", "page_with_curl", "bookmark_tabs", "bar_chart", "chart_with_upwards_trend", "chart_with_downwards_trend",
	"scroll", "clipboard", "calendar", "date", "card_index", "file_folder", "open_file_folder",
	"scissors", "pushpin", "paperclip", "black_nib", "pencil2", "straight_ruler", "triangular_ruler",
	"closed_book", "green_book", "blue_book", "orange_book", "notebook", "notebook_with_decorative_cover", "ledger",
	"books", "bookmark", "name_badge", "microscope", "telescope", "newspaper", "football",
	"basketball", "soccer", "baseball", "tennis", "8ball", "rugby_football", "bowling",
	"golf", "mountain_bicyclist", "bicyclist", "horse_racing", "snowboarder", "swimmer", "surfer",
	"ski", "spades", "hearts", "clubs", "diamonds", "gem", "ring",
	"trophy", "musical_score", "musical_keyboard", "violin", "space_invader", "video_game", "black_joker",
	"flower_playing_cards", "game_die", "dart", "mahjong", "clapper", "memo", "pencil",
	"book", "art", "microphone", "headphones", "trumpet", "saxophone", "guitar",
	"shoe", "sandal", "high_heel", "lipstick", "boot", "shirt", "tshirt",
	"necktie", "womans_clothes", "dress", "running_shirt_with_sash", "jeans", "kimono", "bikini",
	"ribbon", "tophat", "crown", "womans_hat", "mans_shoe", "closed_umbrella", "briefcase",
	"handbag", "pouch", "purse", "eyeglasses", "fishing_pole_and_fish", "coffee", "tea",
	"sake", "baby_bottle", "beer", "beers", "cocktail", "tropical_drink", "wine_glass",
	"fork_and_knife", "pizza", "hamburger", "fries", "poultry_leg", "meat_on_bone",
	"spaghetti", "curry", "fried_shrimp", "bento", "sushi", "fish_cake", "rice_ball",
	"rice_cracker", "rice", "ramen", "stew", "oden", "dango", "egg",
	"bread", "doughnut", "custard", "icecream", "ice_cream", "shaved_ice", "birthday",
	"cake", "cookie", "chocolate_bar", "candy", "lollipop", "honey_pot", "apple",
	"green_apple", "tangerine", "lemon", "cherries", "grapes", "watermelon", "strawberry", "peach", "melon",
	"banana", "pear", "pineapple", "sweet_potato", "eggplant", "tomato", "corn",
];

emojis = $.map(emojis, function(value, i) {return {key:':'+value+':', name:value}});
(function($){
$.fn.Huploadify = function(opts){
	var itemTemp = '<div id="${fileID}" class="uploadify-queue-item"><div class="uploadify-progress"><div class="uploadify-progress-bar"></div></div><span class="up_filename">${fileName}</span><span class="uploadbtn">上传</span><span class="delfilebtn">删除</span></div>';
	var defaults = {
		fileTypeExts:'*.*',//允许上传的文件类型，格式'*.jpg;*.doc'
		uploader:'',//文件提交的地址
		auto:false,//是否开启自动上传
		method:'post',//发送请求的方式，get或post
		multi:true,//是否允许选择多个文件
		formData:null,//发送给服务端的参数，格式：{key1:value1,key2:value2}
		fileObjName:'file',//在后端接受文件的参数名称，如PHP中的$_FILES['file']
		fileSizeLimit:2048,//允许上传的文件大小，单位KB
		showUploadedPercent:true,//是否实时显示上传的百分比，如20%
		showUploadedSize:false,//是否实时显示已上传的文件大小，如1M/2M
		buttonText:'选择文件',//上传按钮上的文字
		removeTimeout: 1000,//上传完成后进度条的消失时间
		itemTemplate:itemTemp,//上传队列显示的模板
		onUploadStart:null,//上传开始时的动作
		onUploadSuccess:null,//上传成功的动作
		onUploadComplete:null,//上传完成的动作
		onUploadAllComplete: null, // 批量上传时，所有的都上传完后回调
		onUploadError:null, //上传失败的动作
		onInit:null,//初始化时的动作
		onCancel:null//删除掉某个文件后的回调函数，可传入参数file
	}
		
	var option = $.extend(defaults,opts);
	
	//将文件的单位由bytes转换为KB或MB，若第二个参数指定为true，则永远转换为KB
	var formatFileSize = function(size,byKB){
		if (size> 1024 * 1024&&!byKB){
			size = (Math.round(size * 100 / (1024 * 1024)) / 100).toString() + 'MB';
		}
		else{
			size = (Math.round(size * 100 / 1024) / 100).toString() + 'KB';
		}
		return size;
	}
	//根据文件序号获取文件
	var getFile = function(index,files){
		for(var i=0;i<files.length;i++){	   
		  if(files[i].index == index){
			  return files[i];
			}
		}
		return false;
	}
	
	//将输入的文件类型字符串转化为数组,原格式为*.jpg;*.png
	var getFileTypes = function(str){
		var result = [];
		var arr1 = str.split(";");
		for(var i=0,len=arr1.length;i<len;i++){
			result.push(arr1[i].split(".").pop());
		}
		return result;
	}
	
	this.each(function(){
		var _this = $(this);
		//先添加上file按钮和上传列表
		var instanceNumber = $('.uploadify').length+1;
		var inputStr = '<input id="select_btn_'+instanceNumber+'" class="selectbtn" style="display:none;" type="file" name="fileselect[]"';
		inputStr += option.multi ? ' multiple' : '';
		inputStr += ' accept="';
		inputStr += getFileTypes(option.fileTypeExts).join(",");
		inputStr += '"/>';
		inputStr += '<a id="file_upload_'+instanceNumber+'-button" href="javascript:void(0)" class="uploadify-button">';
		inputStr += option.buttonText;
		inputStr += '</a>';
		var uploadFileListStr = '<div id="file_upload_'+instanceNumber+'-queue" class="uploadify-queue"></div>';
		_this.append(inputStr+uploadFileListStr);	
		
		
		//创建文件对象
	  var fileObj = {
		  fileInput: _this.find('.selectbtn'),				//html file控件
		  uploadFileList : _this.find('.uploadify-queue'),
		  url: option.uploader,						//ajax地址
		  fileFilter: [],					//过滤后的文件数组
		  filter: function(files) {		//选择文件组的过滤方法
			  var arr = [];
			  var typeArray = getFileTypes(option.fileTypeExts);
			  if(typeArray.length>0){
				  for(var i=0,len=files.length;i<len;i++){
				  	var thisFile = files[i];
				  		if(parseInt(formatFileSize(thisFile.size,true))>option.fileSizeLimit){
				  			alert('文件'+thisFile.name+'大小超出限制！');
				  			continue;
				  		}
						if($.inArray(thisFile.name.split('.').pop(),typeArray)>=0 || $.inArray('*',typeArray)>=0){
							arr.push(thisFile);	
						}
						else{
							alert('文件'+thisFile.name+'类型不允许！');
						}  	
					}	
				}
			  return arr;  	
		  },
		  //文件选择后
		  onSelect: function(files){
				for(var i=0,len=files.length;i<len;i++){
					var file = files[i];
					//处理模板中使用的变量
					var $html = $(option.itemTemplate.replace(/\${fileID}/g,'fileupload_'+instanceNumber+'_'+file.index).replace(/\${fileName}/g,file.name).replace(/\${fileSize}/g,formatFileSize(file.size)).replace(/\${instanceID}/g,_this.attr('id')));
					//如果是自动上传，去掉上传按钮
					if(option.auto){
						$html.find('.uploadbtn').remove();
					}
					this.uploadFileList.append($html);
					
					//判断是否显示已上传文件大小
					if(option.showUploadedSize){
						var num = '<span class="progressnum"><span class="uploadedsize">0KB</span>/<span class="totalsize">${fileSize}</span></span>'.replace(/\${fileSize}/g,formatFileSize(file.size));
						$html.find('.uploadify-progress').after(num);
					}
					
					//判断是否显示上传百分比	
					if(option.showUploadedPercent){
						var percentText = '<span class="up_percent">0%</span>';
						$html.find('.uploadify-progress').after(percentText);
					}
					var theLast = false;
					if (i == len - 1) {
						theLast = true;
					}
					//判断是否是自动上传
					if(option.auto){
						this.funUploadFile(file, theLast);
					}
					else{
						//如果配置非自动上传，绑定上传事件
					 	$html.find('.uploadbtn').on('click',(function(file){
					 			return function(){fileObj.funUploadFile(file, theLast);}
					 		})(file));
					}
					//为删除文件按钮绑定删除文件事件
			 		$html.find('.delfilebtn').on('click',(function(file){
					 			return function(){fileObj.funDeleteFile(file.index);}
					 		})(file));
			 	}

			 
			},				
		  onProgress: function(file, loaded, total) {
				var eleProgress = _this.find('#fileupload_'+instanceNumber+'_'+file.index+' .uploadify-progress');
				var percent = (loaded / total * 100).toFixed(2) +'%';
				if(option.showUploadedSize){
					eleProgress.nextAll('.progressnum .uploadedsize').text(formatFileSize(loaded));
					eleProgress.nextAll('.progressnum .totalsize').text(formatFileSize(total));
				}
				if(option.showUploadedPercent){
					eleProgress.nextAll('.up_percent').text(percent);	
				}
				eleProgress.children('.uploadify-progress-bar').css('width',percent);
	  	},		//文件上传进度

		  /* 开发参数和内置方法分界线 */
		  
		  //获取选择文件，file控件
		  funGetFiles: function(e) {	  
			  // 获取文件列表对象
			  var files = e.target.files;
			  //继续添加文件
			  files = this.filter(files);
			  for(var i=0,len=files.length;i<len;i++){
			  	this.fileFilter.push(files[i]);	
			  }
			  this.funDealFiles(files);
			  return this;
		  },
		  
		  //选中文件的处理与回调
		  funDealFiles: function(files) {
			  var fileCount = _this.find('.uploadify-queue .uploadify-queue-item').length;//队列中已经有的文件个数
			  for(var i=0,len=files.length;i<len;i++){
				  files[i].index = ++fileCount;
				  files[i].id = files[i].index;
				  }
			  //执行选择回调
			  this.onSelect(files);
			  
			  return this;
		  },
		  
		  //删除对应的文件
		  funDeleteFile: function(index) {
			  for (var i = 0,len=this.fileFilter.length; i<len; i++) {
					  var file = this.fileFilter[i];
					  if (file.index == index) {
						  this.fileFilter.splice(i,1);
						  _this.find('#fileupload_'+instanceNumber+'_'+index).fadeOut();
						  option.onCancel&&option.onCancel(file);	
						  break;
					  }
			  }
			  return this;
		  },
		  
		  //文件上传
		  funUploadFile: function(file, theLast) {
			  var xhr = false;
			  try{
				 xhr=new XMLHttpRequest();//尝试创建 XMLHttpRequest 对象，除 IE 外的浏览器都支持这个方法。
			  }catch(e){	  
				xhr=ActiveXobject("Msxml12.XMLHTTP");//使用较新版本的 IE 创建 IE 兼容的对象（Msxml2.XMLHTTP）。
			  }
			  
			  if (xhr.upload) {
				  // 上传中
				  xhr.upload.addEventListener("progress", function(e) {
					  fileObj.onProgress(file, e.loaded, e.total);
				  }, false);
	  
				  // 文件上传成功或是失败
				  xhr.onreadystatechange = function(e) {
					  if (xhr.readyState == 4) {
						  if (xhr.status == 200) {
							  //校正进度条和上传比例的误差
							  var thisfile = _this.find('#fileupload_'+instanceNumber+'_'+file.index);
							  thisfile.find('.uploadify-progress-bar').css('width','100%');
								option.showUploadedSize&&thisfile.find('.uploadedsize').text(thisfile.find('.totalsize').text());
								option.showUploadedPercent&&thisfile.find('.up_percent').text('100%');

							  option.onUploadSuccess&&option.onUploadSuccess(file, xhr.responseText);
							  //在指定的间隔时间后删掉进度条
							  setTimeout(function(){
							  	_this.find('#fileupload_'+instanceNumber+'_'+file.index).fadeOut();
							  },option.removeTimeout);
						  } else {
							  option.onUploadError&&option.onUploadError(file, xhr.responseText);		
						  }
						  option.onUploadComplete&&option.onUploadComplete(file,xhr.responseText);
						  //清除文件选择框中的已有值
						  fileObj.fileInput.val('');

						  if (theLast) {
						  	  option.onUploadAllComplete&&option.onUploadAllComplete(file,xhr.responseText);
						  }
					  }
				  };
	  
	  			option.onUploadStart&&option.onUploadStart();	
				  // 开始上传
				  xhr.open(option.method, this.url, true);
				  xhr.setRequestHeader("X-Requested-With", "XMLHttpRequest");
				  var fd = new FormData();
				  fd.append(option.fileObjName,file);
				  if(option.formData){
				  	for(key in option.formData){
				  		fd.append(key,option.formData[key]);
				  	}
				  }
				  
				  xhr.send(fd);
			  }	
			  
				  
		  },
		  
		  init: function() {	  
			  //文件选择控件选择
			  if (this.fileInput.length>0) {
				  this.fileInput.change(function(e) { 
				  	fileObj.funGetFiles(e); 
				  });	
			  }
			  
			  //点击上传按钮时触发file的click事件
			  _this.find('.uploadify-button').on('click',function(){
				  _this.find('.selectbtn').trigger('click');
				});
			  
			  option.onInit&&option.onInit();
		  }
  	};

		//初始化文件对象
		fileObj.init();
	}); 
}	

})(jQuery);

(function(){!function(a){return"function"==typeof define&&define.amd?define(["jquery"],a):a(window.jQuery)}(function(a){var b,c,d,e,f,g,h,i=[].slice;c=function(){function b(b){this.current_flag=null,this.controllers={},this.alias_maps={},this.$inputor=a(b),this.setIframe(),this.listen()}return b.prototype.createContainer=function(b){return 0===(this.$el=a("#atwho-container",b)).length?a(b.body).append(this.$el=a("<div id='atwho-container'></div>")):void 0},b.prototype.setIframe=function(a,b){var c;return null==b&&(b=!1),a?(this.window=a.contentWindow,this.document=a.contentDocument||this.window.document,this.iframe=a):(this.document=document,this.window=window,this.iframe=null),(this.iframeStandalone=b)?(null!=(c=this.$el)&&c.remove(),this.createContainer(this.document)):this.createContainer(document)},b.prototype.controller=function(a){var b,c,d,e;if(this.alias_maps[a])c=this.controllers[this.alias_maps[a]];else{e=this.controllers;for(d in e)if(b=e[d],d===a){c=b;break}}return c?c:this.controllers[this.current_flag]},b.prototype.set_context_for=function(a){return this.current_flag=a,this},b.prototype.reg=function(a,b){var c,e;return c=(e=this.controllers)[a]||(e[a]=new d(this,a)),b.alias&&(this.alias_maps[b.alias]=a),c.init(b),this},b.prototype.listen=function(){return this.$inputor.on("keyup.atwhoInner",function(a){return function(b){return a.on_keyup(b)}}(this)).on("keydown.atwhoInner",function(a){return function(b){return a.on_keydown(b)}}(this)).on("scroll.atwhoInner",function(a){return function(b){var c;return null!=(c=a.controller())?c.view.hide(b):void 0}}(this)).on("blur.atwhoInner",function(a){return function(b){var c;return(c=a.controller())?c.view.hide(b,c.get_opt("display_timeout")):void 0}}(this)).on("click.atwhoInner",function(a){return function(b){var c;return null!=(c=a.controller())?c.view.hide(b):void 0}}(this))},b.prototype.shutdown=function(){var a,b,c;c=this.controllers;for(b in c)a=c[b],a.destroy(),delete this.controllers[b];return this.$inputor.off(".atwhoInner"),this.$el.remove()},b.prototype.dispatch=function(){return a.map(this.controllers,function(a){return function(b){var c;return(c=b.get_opt("delay"))?(clearTimeout(a.delayedCallback),a.delayedCallback=setTimeout(function(){return b.look_up()?a.set_context_for(b.at):void 0},c)):b.look_up()?a.set_context_for(b.at):void 0}}(this))},b.prototype.on_keyup=function(b){var c;switch(b.keyCode){case f.ESC:b.preventDefault(),null!=(c=this.controller())&&c.view.hide();break;case f.DOWN:case f.UP:case f.CTRL:a.noop();break;case f.P:case f.N:b.ctrlKey||this.dispatch();break;default:this.dispatch()}},b.prototype.on_keydown=function(b){var c,d;if(c=null!=(d=this.controller())?d.view:void 0,c&&c.visible())switch(b.keyCode){case f.ESC:b.preventDefault(),c.hide(b);break;case f.UP:b.preventDefault(),c.prev();break;case f.DOWN:b.preventDefault(),c.next();break;case f.P:if(!b.ctrlKey)return;b.preventDefault(),c.prev();break;case f.N:if(!b.ctrlKey)return;b.preventDefault(),c.next();break;case f.TAB:case f.ENTER:if(!c.visible())return;b.preventDefault(),c.choose(b);break;default:a.noop()}},b}(),d=function(){function b(b,c){this.app=b,this.at=c,this.$inputor=this.app.$inputor,this.id=this.$inputor[0].id||this.uid(),this.setting=null,this.query=null,this.pos=0,this.cur_rect=null,this.range=null,0===(this.$el=a("#atwho-ground-"+this.id,this.app.$el)).length&&this.app.$el.append(this.$el=a("<div id='atwho-ground-"+this.id+"'></div>")),this.model=new g(this),this.view=new h(this)}return b.prototype.uid=function(){return(Math.random().toString(16)+"000000000").substr(2,8)+(new Date).getTime()},b.prototype.init=function(b){return this.setting=a.extend({},this.setting||a.fn.atwho["default"],b),this.view.init(),this.model.reload(this.setting.data)},b.prototype.destroy=function(){return this.trigger("beforeDestroy"),this.model.destroy(),this.view.destroy(),this.$el.remove()},b.prototype.call_default=function(){var b,c,d;d=arguments[0],b=2<=arguments.length?i.call(arguments,1):[];try{return e[d].apply(this,b)}catch(f){return c=f,a.error(""+c+" Or maybe At.js doesn't have function "+d)}},b.prototype.trigger=function(a,b){var c,d;return null==b&&(b=[]),b.push(this),c=this.get_opt("alias"),d=c?""+a+"-"+c+".atwho":""+a+".atwho",this.$inputor.trigger(d,b)},b.prototype.callbacks=function(a){return this.get_opt("callbacks")[a]||e[a]},b.prototype.get_opt=function(a){var b;try{return this.setting[a]}catch(c){return b=c,null}},b.prototype.content=function(){return this.$inputor.is("textarea, input")?this.$inputor.val():this.$inputor.text()},b.prototype.catch_query=function(){var a,b,c,d,e,f;return b=this.content(),a=this.$inputor.caret("pos",{iframe:this.app.iframe}),f=b.slice(0,a),d=this.callbacks("matcher").call(this,this.at,f,this.get_opt("start_with_space")),"string"==typeof d&&d.length<=this.get_opt("max_len",20)?(e=a-d.length,c=e+d.length,this.pos=e,d={text:d,head_pos:e,end_pos:c},this.trigger("matched",[this.at,d.text])):(d=null,this.view.hide()),this.query=d},b.prototype.rect=function(){var b,c,d;if(b=this.$inputor.caret("offset",this.pos-1,{iframe:this.app.iframe}))return this.app.iframe&&!this.app.iframeStandalone&&(c=a(this.app.iframe).offset(),b.left+=c.left,b.top+=c.top),this.$inputor.is("[contentEditable]")&&(b=this.cur_rect||(this.cur_rect=b)),d=this.app.document.selection?0:2,{left:b.left,top:b.top,bottom:b.top+b.height+d}},b.prototype.reset_rect=function(){return this.$inputor.is("[contentEditable]")?this.cur_rect=null:void 0},b.prototype.mark_range=function(){var a;if(this.$inputor.is("[contentEditable]"))return this.app.window.getSelection&&(a=this.app.window.getSelection()).rangeCount>0?this.range=a.getRangeAt(0):this.app.document.selection?this.ie8_range=this.app.document.selection.createRange():void 0},b.prototype.insert_content_for=function(b){var c,d,e;return d=b.data("value"),e=this.get_opt("insert_tpl"),this.$inputor.is("textarea, input")||!e?d:(c=a.extend({},b.data("item-data"),{"atwho-data-value":d,"atwho-at":this.at}),this.callbacks("tpl_eval").call(this,e,c))},b.prototype.insert=function(b){var c,d,e,f,g,h,i,j,k;return c=this.$inputor,k=this.callbacks("inserting_wrapper").call(this,c,b,this.get_opt("suffix")),c.is("textarea, input")?(h=c.val(),i=h.slice(0,Math.max(this.query.head_pos-this.at.length,0)),j=""+i+k+h.slice(this.query.end_pos||0),c.val(j),c.caret("pos",i.length+k.length,{iframe:this.app.iframe})):(f=this.range)?(e=f.startOffset-(this.query.end_pos-this.query.head_pos)-this.at.length,f.setStart(f.endContainer,Math.max(e,0)),f.setEnd(f.endContainer,f.endOffset),f.deleteContents(),d=a(k,this.app.document)[0],f.insertNode(d),f.setEndAfter(d),f.collapse(!1),g=this.app.window.getSelection(),g.removeAllRanges(),g.addRange(f)):(f=this.ie8_range)&&(f.moveStart("character",this.query.end_pos-this.query.head_pos-this.at.length),f.pasteHTML(k),f.collapse(!1),f.select()),c.is(":focus")||c.focus(),c.change()},b.prototype.render_view=function(a){var b;return b=this.get_opt("search_key"),a=this.callbacks("sorter").call(this,this.query.text,a.slice(0,1001),b),this.view.render(a.slice(0,this.get_opt("limit")))},b.prototype.look_up=function(){var b,c;if(b=this.catch_query())return c=function(a){return a&&a.length>0?this.render_view(a):this.view.hide()},this.model.query(b.text,a.proxy(c,this)),b},b}(),g=function(){function b(a){this.context=a,this.at=this.context.at,this.storage=this.context.$inputor}return b.prototype.destroy=function(){return this.storage.data(this.at,null)},b.prototype.saved=function(){return this.fetch()>0},b.prototype.query=function(a,b){var c,d,e;return c=this.fetch(),d=this.context.get_opt("search_key"),c=this.context.callbacks("filter").call(this.context,a,c,d)||[],e=this.context.callbacks("remote_filter"),c.length>0||!e&&0===c.length?b(c):e.call(this.context,a,b)},b.prototype.fetch=function(){return this.storage.data(this.at)||[]},b.prototype.save=function(a){return this.storage.data(this.at,this.context.callbacks("before_save").call(this.context,a||[]))},b.prototype.load=function(a){return!this.saved()&&a?this._load(a):void 0},b.prototype.reload=function(a){return this._load(a)},b.prototype._load=function(b){return"string"==typeof b?a.ajax(b,{dataType:"json"}).done(function(a){return function(b){return a.save(b)}}(this)):this.save(b)},b}(),h=function(){function b(b){this.context=b,this.$el=a("<div class='atwho-view'><ul class='atwho-view-ul'></ul></div>"),this.timeout_id=null,this.context.$el.append(this.$el),this.bind_event()}return b.prototype.init=function(){var a;return a=this.context.get_opt("alias")||this.context.at.charCodeAt(0),this.$el.attr({id:"at-view-"+a})},b.prototype.destroy=function(){return this.$el.remove()},b.prototype.bind_event=function(){var b;return b=this.$el.find("ul"),b.on("mouseenter.atwho-view","li",function(c){return b.find(".cur").removeClass("cur"),a(c.currentTarget).addClass("cur")}).on("click",function(a){return function(b){return a.choose(b),b.preventDefault()}}(this))},b.prototype.visible=function(){return this.$el.is(":visible")},b.prototype.choose=function(a){var b,c;return(b=this.$el.find(".cur")).length&&(c=this.context.insert_content_for(b),this.context.insert(this.context.callbacks("before_insert").call(this.context,c,b),b),this.context.trigger("inserted",[b,a]),this.hide(a)),this.context.get_opt("hide_without_suffix")?this.stop_showing=!0:void 0},b.prototype.reposition=function(b){var c,d,e,f;return f=this.context.app.iframeStandalone?this.context.app.window:window,b.bottom+this.$el.height()-a(f).scrollTop()>a(f).height()&&(b.bottom=b.top-this.$el.height()),b.left>(d=a(f).width()-this.$el.width()-5)&&(b.left=d),c={left:b.left,top:b.bottom},null!=(e=this.context.callbacks("before_reposition"))&&e.call(this.context,c),this.$el.offset(c),this.context.trigger("reposition",[c])},b.prototype.next=function(){var a,b;return a=this.$el.find(".cur").removeClass("cur"),b=a.next(),b.length||(b=this.$el.find("li:first")),b.addClass("cur")},b.prototype.prev=function(){var a,b;return a=this.$el.find(".cur").removeClass("cur"),b=a.prev(),b.length||(b=this.$el.find("li:last")),b.addClass("cur")},b.prototype.show=function(){var a;return this.stop_showing?void(this.stop_showing=!1):(this.context.mark_range(),this.visible()||(this.$el.show(),this.context.trigger("shown")),(a=this.context.rect())?this.reposition(a):void 0)},b.prototype.hide=function(a,b){var c;if(this.visible())return isNaN(b)?(this.context.reset_rect(),this.$el.hide(),this.context.trigger("hidden",[a])):(c=function(a){return function(){return a.hide()}}(this),clearTimeout(this.timeout_id),this.timeout_id=setTimeout(c,b))},b.prototype.render=function(b){var c,d,e,f,g,h,i;if(!(a.isArray(b)&&b.length>0))return void this.hide();for(this.$el.find("ul").empty(),d=this.$el.find("ul"),g=this.context.get_opt("tpl"),h=0,i=b.length;i>h;h++)e=b[h],e=a.extend({},e,{"atwho-at":this.context.at}),f=this.context.callbacks("tpl_eval").call(this.context,g,e),c=a(this.context.callbacks("highlighter").call(this.context,f,this.context.query.text)),c.data("item-data",e),d.append(c);return this.show(),this.context.get_opt("highlight_first")?d.find("li:first").addClass("cur"):void 0},b}(),f={DOWN:40,UP:38,ESC:27,TAB:9,ENTER:13,CTRL:17,P:80,N:78},e={before_save:function(b){var c,d,e,f;if(!a.isArray(b))return b;for(f=[],d=0,e=b.length;e>d;d++)c=b[d],f.push(a.isPlainObject(c)?c:{name:c});return f},matcher:function(a,b,c){var d,e,f,g;return a=a.replace(/[\-\[\]\/\{\}\(\)\*\+\?\.\\\^\$\|]/g,"\\$&"),c&&(a="(?:^|\\s)"+a),f=decodeURI("%C3%80"),g=decodeURI("%C3%BF"),e=new RegExp(""+a+"([A-Za-z"+f+"-"+g+"0-9_+-]*)$|"+a+"([^\\x00-\\xff]*)$","gi"),d=e.exec(b),d?d[2]||d[1]:null},filter:function(a,b,c){var d,e,f,g;for(g=[],e=0,f=b.length;f>e;e++)d=b[e],~new String(d[c]).toLowerCase().indexOf(a.toLowerCase())&&g.push(d);return g},remote_filter:null,sorter:function(a,b,c){var d,e,f,g;if(!a)return b;for(g=[],e=0,f=b.length;f>e;e++)d=b[e],d.atwho_order=new String(d[c]).toLowerCase().indexOf(a.toLowerCase()),d.atwho_order>-1&&g.push(d);return g.sort(function(a,b){return a.atwho_order-b.atwho_order})},tpl_eval:function(a,b){var c;try{return a.replace(/\$\{([^\}]*)\}/g,function(a,c){return b[c]})}catch(d){return c=d,""}},highlighter:function(a,b){var c;return b?(c=new RegExp(">\\s*(\\w*?)("+b.replace("+","\\+")+")(\\w*)\\s*<","ig"),a.replace(c,function(a,b,c,d){return"> "+b+"<strong>"+c+"</strong>"+d+" <"})):a},before_insert:function(a){return a},inserting_wrapper:function(a,b,c){var d,e;return d=""===c?c:c||" ",a.is("textarea, input")?""+b+d:"true"===a.attr("contentEditable")?(d=""===c?c:c||"&nbsp;",/firefox/i.test(navigator.userAgent)?e="<span>"+b+d+"</span>":(c="<span contenteditable='false'>"+d+"<span>",e="<span contenteditable='false'>"+b+c+"</span>"),this.app.document.selection&&(e="<span contenteditable='true'>"+b+"</span>"),e):void 0}},b={load:function(a,b){var c;return(c=this.controller(a))?c.model.load(b):void 0},setIframe:function(a,b){return this.setIframe(a,b),null},run:function(){return this.dispatch()},destroy:function(){return this.shutdown(),this.$inputor.data("atwho",null)}},a.fn.atwho=function(d){var e,f;return f=arguments,e=null,this.filter('textarea, input, [contenteditable=""], [contenteditable=true]').each(function(){var g,h;return(h=(g=a(this)).data("atwho"))||g.data("atwho",h=new c(this)),"object"!=typeof d&&d?b[d]&&h?e=b[d].apply(h,Array.prototype.slice.call(f,1)):a.error("Method "+d+" does not exist on jQuery.caret"):h.reg(d.at,d)}),e||this},a.fn.atwho["default"]={at:void 0,alias:void 0,data:null,tpl:"<li data-value='${atwho-at}${name}'>${name}</li>",insert_tpl:"<span id='${id}'>${atwho-data-value}</span>",callbacks:e,search_key:"name",suffix:void 0,hide_without_suffix:!1,start_with_space:!0,highlight_first:!0,limit:5,max_len:20,display_timeout:300,delay:null}})}).call(this);
/*
 * ----------------------------------------------------------------------------
 * "THE BEER-WARE LICENSE" (Revision 42):
 * <jevin9@gmail.com> wrote this file. As long as you retain this notice you
 * can do whatever you want with this stuff. If we meet some day, and you think
 * this stuff is worth it, you can buy me a beer in return. Jevin O. Sewaruth
 * ----------------------------------------------------------------------------
 *
 * Autogrow Textarea Plugin Version v3.0
 * http://www.technoreply.com/autogrow-textarea-plugin-3-0
 * 
 * THIS PLUGIN IS DELIVERD ON A PAY WHAT YOU WHANT BASIS. IF THE PLUGIN WAS USEFUL TO YOU, PLEASE CONSIDER BUYING THE PLUGIN HERE :
 * https://sites.fastspring.com/technoreply/instant/autogrowtextareaplugin
 *
 * Date: October 15, 2012
 */

jQuery.fn.autoGrow=function(){return this.each(function(){var createMirror=function(textarea){jQuery(textarea).after('<div class="autogrow-textarea-mirror"></div>');return jQuery(textarea).next(".autogrow-textarea-mirror")[0]};var sendContentToMirror=function(textarea){mirror.innerHTML=String(textarea.value).replace(/&/g,"&amp;").replace(/"/g,"&quot;").replace(/'/g,"&#39;").replace(/</g,"&lt;").replace(/>/g,"&gt;").replace(/\n/g,"<br />")+".<br/>.";if(jQuery(textarea).height()!=jQuery(mirror).height())jQuery(textarea).height(jQuery(mirror).height())};
var growTextarea=function(){sendContentToMirror(this)};var mirror=createMirror(this);mirror.style.display="none";mirror.style.wordWrap="break-word";mirror.style.padding=jQuery(this).css("padding");mirror.style.width=jQuery(this).css("width");mirror.style.fontFamily=jQuery(this).css("font-family");mirror.style.fontSize=jQuery(this).css("font-size");mirror.style.lineHeight=jQuery(this).css("line-height");this.style.overflow="hidden";this.style.minHeight=this.rows+"em";this.onkeyup=growTextarea;sendContentToMirror(this)})};

(function(a){a.fn.cftoaster=function(b){var d=a.extend({},a.fn.cftoaster.options,b);return this.each(function(){d.element=a(this);if(!c(d)){a.cftoaster._addToQueue(d)}else{a.cftoaster._destroy(d)}});function c(e){var g="";for(var f=0;f<=a.cftoaster.DESTROY_COMMAND.length;f++){if(!e.hasOwnProperty(f)){break}g+=e[f]}return g==a.cftoaster.DESTROY_COMMAND}};a.fn.cftoaster.options={content:"This is a toast message eh",element:"body",animationTime:150,showTime:3000,maxWidth:250,backgroundColor:"#1a1a1a",fontColor:"#eaeaea",bottomMargin:75}})(jQuery);jQuery.extend({cftoaster:{NAMESPACE:"cf_toaster",DESTROY_COMMAND:"destroy",MAIN_CSS_CLASS:"cf_toaster",_queue:[],_addToQueue:function(a){this._queue.push(a);if(a.element&&!this._isShowingToastMessage(a.element)){this._showNextInQueue(a.element)}},_removeFromQueue:function(c){if(c){for(var b in this._queue){var a=this._queue[b];if($(a.element).is(c)){this._queue.splice(b,1)}}}else{this._queue=[]}},_destroy:function(a){var b=a&&a.element?a.element:undefined;if(b){$(b).find("."+this.MAIN_CSS_CLASS).remove()}else{$("."+this.MAIN_CSS_CLASS).remove()}this._removeFromQueue(b)},_isShowingToastMessage:function(b){var a=false;if(b){a=$(b).find("."+this.MAIN_CSS_CLASS).size()>0}return a},_showNextInQueue:function(e){var a;for(var d=0;d<this._queue.length;d++){var g=this._queue[d];if($(g.element).is(e)){a=g;this._queue.splice(d,1);break}}if(a){var c=$("<div/>").addClass("background").css("background",a.backgroundColor);var f=$("<div/>").addClass("content").html(a.content).css("width",a.maxWidth+"px").css("color",a.fontColor);var b=$("<div/>").addClass(this.MAIN_CSS_CLASS).hide().append(c).append(f);$(e).append(b);var h=-$(b).outerWidth()/2+"px";$(b).css("bottom",a.bottomMargin+"px").css("margin-left",h);$(b).stop().fadeIn(a.animationTime).delay(a.showTime).fadeOut(a.animationTime,function(){$(this).remove();$.cftoaster._showNextInQueue(e)})}},setDefaults:function(a){var b=$.extend({},$.fn.cftoaster.options,a);$.fn.cftoaster.options=b}}});
jQuery(document).ready(function(e){var t=0;var n="https://static.golangjob.cn/github_logo.gif";var r="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABQAAAAqCAMAAACEJ4viAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAyRpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw/eHBhY2tldCBiZWdpbj0i77u/IiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8+IDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuMC1jMDYxIDY0LjE0MDk0OSwgMjAxMC8xMi8wNy0xMDo1NzowMSAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvIiB4bWxuczp4bXBNTT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RSZWY9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZVJlZiMiIHhtcDpDcmVhdG9yVG9vbD0iQWRvYmUgUGhvdG9zaG9wIENTNS4xIE1hY2ludG9zaCIgeG1wTU06SW5zdGFuY2VJRD0ieG1wLmlpZDpEQjIyNkJEQkM0NjYxMUUxOEFDQzk3ODcxRDkzRjhCRSIgeG1wTU06RG9jdW1lbnRJRD0ieG1wLmRpZDpEQjIyNkJEQ0M0NjYxMUUxOEFDQzk3ODcxRDkzRjhCRSI+IDx4bXBNTTpEZXJpdmVkRnJvbSBzdFJlZjppbnN0YW5jZUlEPSJ4bXAuaWlkOkRCMjI2QkQ5QzQ2NjExRTE4QUNDOTc4NzFEOTNGOEJFIiBzdFJlZjpkb2N1bWVudElEPSJ4bXAuZGlkOkRCMjI2QkRBQzQ2NjExRTE4QUNDOTc4NzFEOTNGOEJFIi8+IDwvcmRmOkRlc2NyaXB0aW9uPiA8L3JkZjpSREY+IDwveDp4bXBtZXRhPiA8P3hwYWNrZXQgZW5kPSJyIj8+h1kA9gAAAK5QTFRF+fn5sbGx8fHx09PTmpqa2dnZ/f3919fX9PT00NDQ1dXVpKSk+vr6+/v7vb298vLyycnJ8/PztLS0zc3N6enp/v7+q6ur2NjY9/f3srKy/Pz8p6en7u7uoaGhnJyc4eHhtbW1pqam6Ojo9fX17e3toqKirKys1NTUzs7Ox8fHwcHBwMDA5eXlnZ2dpaWl0dHR9vb25ubm4uLi3d3dqqqqwsLCv7+/oKCgmZmZ////8yEsbwAAAMBJREFUeNrE0tcOgjAUBuDSliUoMhTEvfdef9//xUQjgaLX0Ium/ZLT/+SkRPxZpGykvuf5VMJogy5jY9yjDHcWFhqlcRuHc4o6B1QK0BDg+hcZgNDh3NWTwzItH/bRrhvT+g3zSxZkNGCZpoWGIbU0a3Y6zV5VA6keyeDxiw62P0gUqEW0FbDim4nVikFJbU2zZXybUEaxhCqOQqyh5/G0wpWICUwthyqwD4InOMuXJ7/gs7WkoPdVg1vykF8CDACEFanKO3aSYwAAAABJRU5ErkJggg==";e(".github-widget").each(function(){if(t==0)e("head").append('<style type="text/css">'+".github-box *{-webkit-box-sizing:content-box;-moz-box-sizing:content-box;box-sizing:content-box;}"+".github-box{font-family:helvetica,arial,sans-serif;font-size:13px;line-height:18px;background:#fafafa;border:1px solid #ddd;color:#666;border-radius:3px}"+".github-box a{color:#4183c4;border:0;text-decoration:none}"+".github-box .github-box-title{position:relative;border-bottom:1px solid #ddd;border-radius:3px 3px 0 0;background:#fcfcfc;background:-moz-linear-gradient(#fcfcfc,#ebebeb);background:-webkit-linear-gradient(#fcfcfc,#ebebeb);}"+".github-box .github-box-title h3{word-wrap:break-word;font-family:helvetica,arial,sans-serif;font-weight:normal;font-size:16px;color:gray;margin:0;padding:10px 10px 10px 80px;background:url("+n+") 7px center no-repeat; width: auto;}"+".github-box .github-box-title h3 .repo{font-weight:bold}"+".github-box .github-box-title .github-stats{float:right;position:absolute;top:8px;right:10px;font-size:11px;font-weight:bold;line-height:21px;height:auto;min-height:21px}"+".github-box .github-box-title .github-stats a{display:inline-block;height:21px;color:#666;border:1px solid #ddd;border-radius:3px;padding:0 5px 0 18px;background: white url("+r+") no-repeat}"+".github-box .github-box-title .github-stats .watchers{border-right:1px solid #ddd}"+".github-box .github-box-title .github-stats .forks{background-position:-4px -21px;padding-left:15px}"+".github-box .github-box-content{padding:10px;font-weight:300}"+".github-box .github-box-content p{margin:0}"+".github-box .github-box-content .link{font-weight:bold}"+".github-box .github-box-download{position:relative;border-top:1px solid #ddd;background:white;border-radius:0 0 3px 3px;padding:10px;height:auto;min-height:24px;}"+".github-box .github-box-download .updated{word-wrap:break-word;margin:0;font-size:11px;color:#666;line-height:24px;font-weight:300;width:auto}"+".github-box .github-box-download .updated strong{font-weight:bold;color:#000}"+".github-box .github-box-download .download{float:right;position:absolute;top:10px;right:10px;height:24px;line-height:24px;font-size:12px;color:#666;font-weight:bold;text-shadow:0 1px 0 rgba(255,255,255,0.9);padding:0 10px;border:1px solid #ddd;border-bottom-color:#bbb;border-radius:3px;background:#f5f5f5;background:-moz-linear-gradient(#f5f5f5,#e5e5e5);background:-webkit-linear-gradient(#f5f5f5,#e5e5e5);}"+".github-box .github-box-download .download:hover{color:#527894;border-color:#cfe3ed;border-bottom-color:#9fc7db;background:#f1f7fa;background:-moz-linear-gradient(#f1f7fa,#dbeaf1);background:-webkit-linear-gradient(#f1f7fa,#dbeaf1);}"+"@media (max-width: 767px) {"+".github-box .github-box-title{height:auto;min-height:60px}"+".github-box .github-box-title h3 .repo{display:block}"+".github-box .github-box-title .github-stats a{display:block;clear:right;float:right;}"+".github-box .github-box-download{height:auto;min-height:46px;}"+".github-box .github-box-download .download{top:32px;}"+"}"+"</style>");t++;var s=e(this),o,u=s.data("repo"),a=u.split("/")[0],f=u.split("/")[1],l="http://github.com/"+a,c="http://github.com/"+a+"/"+f;o=e('<div class="github-box repo">'+'<div class="github-box-title">'+"<h3>"+'<a class="owner" href="'+l+'" title="'+l+'">'+a+"</a>"+"/"+'<a class="repo" href="'+c+'" title="'+c+'">'+f+"</a>"+"</h3>"+'<div class="github-stats">'+'<a class="watchers" href="'+c+'/watchers" title="See watchers">?</a>'+'<a class="forks" href="'+c+'/network/members" title="See forkers">?</a>'+"</div>"+"</div>"+'<div class="github-box-content">'+'<p class="description"><span></span> — <a href="'+c+'#readme">Read More</a></p>'+'<p class="link"></p>'+"</div>"+'<div class="github-box-download">'+'<div class="updated"></div>'+'<a class="download" href="'+c+'/zipball/master" title="Get an archive of this repository">Download as zip</a>'+"</div>"+"</div>");o.appendTo(s);e.ajax({url:"https://api.github.com/repos/"+u,dataType:"jsonp",success:function(t){var n=t.data,r,i="unknown";if(n.pushed_at){r=new Date(n.pushed_at);i=r.getMonth()+1+"-"+r.getDate()+"-"+r.getFullYear()}o.find(".watchers").text(n.watchers);o.find(".forks").text(n.forks);o.find(".description span").text(n.description);o.find(".updated").html("Latest commit to the <strong>"+n.default_branch+"</strong> branch on "+i);if(n.homepage!=null)o.find(".link").append(e("<a />").attr("href",n.homepage).text(n.homepage))}})})});

/*
 * Metadata - jQuery plugin for parsing metadata from elements
 *
 * Copyright (c) 2006 John Resig, Yehuda Katz, J�örn Zaefferer, Paul McLanahan
 *
 * Dual licensed under the MIT and GPL licenses:
 *   http://www.opensource.org/licenses/mit-license.php
 *   http://www.gnu.org/licenses/gpl.html
 *
 * Revision: $Id: jquery.metadata.js 4187 2007-12-16 17:15:27Z joern.zaefferer $
 *
 */

/**
 * Sets the type of metadata to use. Metadata is encoded in JSON, and each property
 * in the JSON will become a property of the element itself.
 *
 * There are three supported types of metadata storage:
 *
 *   attr:  Inside an attribute. The name parameter indicates *which* attribute.
 *          
 *   class: Inside the class attribute, wrapped in curly braces: { }
 *   
 *   elem:  Inside a child element (e.g. a script tag). The
 *          name parameter indicates *which* element.
 *          
 * The metadata for an element is loaded the first time the element is accessed via jQuery.
 *
 * As a result, you can define the metadata type, use $(expr) to load the metadata into the elements
 * matched by expr, then redefine the metadata type and run another $(expr) for other elements.
 * 
 * @name $.metadata.setType
 *
 * @example <p id="one" class="some_class {item_id: 1, item_label: 'Label'}">This is a p</p>
 * @before $.metadata.setType("class")
 * @after $("#one").metadata().item_id == 1; $("#one").metadata().item_label == "Label"
 * @desc Reads metadata from the class attribute
 * 
 * @example <p id="one" class="some_class" data="{item_id: 1, item_label: 'Label'}">This is a p</p>
 * @before $.metadata.setType("attr", "data")
 * @after $("#one").metadata().item_id == 1; $("#one").metadata().item_label == "Label"
 * @desc Reads metadata from a "data" attribute
 * 
 * @example <p id="one" class="some_class"><script>{item_id: 1, item_label: 'Label'}</script>This is a p</p>
 * @before $.metadata.setType("elem", "script")
 * @after $("#one").metadata().item_id == 1; $("#one").metadata().item_label == "Label"
 * @desc Reads metadata from a nested script element
 * 
 * @param String type The encoding type
 * @param String name The name of the attribute to be used to get metadata (optional)
 * @cat Plugins/Metadata
 * @descr Sets the type of encoding to be used when loading metadata for the first time
 * @type undefined
 * @see metadata()
 */

(function($) {

$.extend({
	metadata : {
		defaults : {
			type: 'class',
			name: 'metadata',
			cre: /({.*})/,
			single: 'metadata'
		},
		setType: function( type, name ){
			this.defaults.type = type;
			this.defaults.name = name;
		},
		get: function( elem, opts ){
			var settings = $.extend({},this.defaults,opts);
			// check for empty string in single property
			if ( !settings.single.length ) settings.single = 'metadata';
			
			var data = $.data(elem, settings.single);
			// returned cached data if it already exists
			if ( data ) return data;
			
			data = "{}";
			
			if ( settings.type == "class" ) {
				var m = settings.cre.exec( elem.className );
				if ( m )
					data = m[1];
			} else if ( settings.type == "elem" ) {
				if( !elem.getElementsByTagName )
					return undefined;
				var e = elem.getElementsByTagName(settings.name);
				if ( e.length )
					data = $.trim(e[0].innerHTML);
			} else if ( elem.getAttribute != undefined ) {
				var attr = elem.getAttribute( settings.name );
				if ( attr )
					data = attr;
			}
			
			if ( data.indexOf( '{' ) <0 )
			data = "{" + data + "}";
			
			data = eval("(" + data + ")");
			
			$.data( elem, settings.single, data );
			return data;
		}
	}
});

/**
 * Returns the metadata object for the first member of the jQuery object.
 *
 * @name metadata
 * @descr Returns element's metadata object
 * @param Object opts An object contianing settings to override the defaults
 * @type jQuery
 * @cat Plugins/Metadata
 */
$.fn.metadata = function( opts ){
	return $.metadata.get( this[0], opts );
};

})(jQuery);
// Simplified Chinese
jQuery.timeago.settings.strings = {
  prefixAgo: null,
  prefixFromNow: "从现在开始",
  suffixAgo: "之前",
  suffixFromNow: null,
  seconds: "不到1分钟",
  minute: "大约1分钟",
  minutes: "%d分钟",
  hour: "大约1小时",
  hours: "大约%d小时",
  day: "1天",
  days: "%d天",
  month: "大约1个月",
  months: "%d月",
  year: "大约1年",
  years: "%d年",
  numbers: [],
  wordSeparator: ""
};

function md5cycle(x, k) {
	var a = x[0], b = x[1], c = x[2], d = x[3];

	a = ff(a, b, c, d, k[0], 7, -680876936);
	d = ff(d, a, b, c, k[1], 12, -389564586);
	c = ff(c, d, a, b, k[2], 17,  606105819);
	b = ff(b, c, d, a, k[3], 22, -1044525330);
	a = ff(a, b, c, d, k[4], 7, -176418897);
	d = ff(d, a, b, c, k[5], 12,  1200080426);
	c = ff(c, d, a, b, k[6], 17, -1473231341);
	b = ff(b, c, d, a, k[7], 22, -45705983);
	a = ff(a, b, c, d, k[8], 7,  1770035416);
	d = ff(d, a, b, c, k[9], 12, -1958414417);
	c = ff(c, d, a, b, k[10], 17, -42063);
	b = ff(b, c, d, a, k[11], 22, -1990404162);
	a = ff(a, b, c, d, k[12], 7,  1804603682);
	d = ff(d, a, b, c, k[13], 12, -40341101);
	c = ff(c, d, a, b, k[14], 17, -1502002290);
	b = ff(b, c, d, a, k[15], 22,  1236535329);

	a = gg(a, b, c, d, k[1], 5, -165796510);
	d = gg(d, a, b, c, k[6], 9, -1069501632);
	c = gg(c, d, a, b, k[11], 14,  643717713);
	b = gg(b, c, d, a, k[0], 20, -373897302);
	a = gg(a, b, c, d, k[5], 5, -701558691);
	d = gg(d, a, b, c, k[10], 9,  38016083);
	c = gg(c, d, a, b, k[15], 14, -660478335);
	b = gg(b, c, d, a, k[4], 20, -405537848);
	a = gg(a, b, c, d, k[9], 5,  568446438);
	d = gg(d, a, b, c, k[14], 9, -1019803690);
	c = gg(c, d, a, b, k[3], 14, -187363961);
	b = gg(b, c, d, a, k[8], 20,  1163531501);
	a = gg(a, b, c, d, k[13], 5, -1444681467);
	d = gg(d, a, b, c, k[2], 9, -51403784);
	c = gg(c, d, a, b, k[7], 14,  1735328473);
	b = gg(b, c, d, a, k[12], 20, -1926607734);

	a = hh(a, b, c, d, k[5], 4, -378558);
	d = hh(d, a, b, c, k[8], 11, -2022574463);
	c = hh(c, d, a, b, k[11], 16,  1839030562);
	b = hh(b, c, d, a, k[14], 23, -35309556);
	a = hh(a, b, c, d, k[1], 4, -1530992060);
	d = hh(d, a, b, c, k[4], 11,  1272893353);
	c = hh(c, d, a, b, k[7], 16, -155497632);
	b = hh(b, c, d, a, k[10], 23, -1094730640);
	a = hh(a, b, c, d, k[13], 4,  681279174);
	d = hh(d, a, b, c, k[0], 11, -358537222);
	c = hh(c, d, a, b, k[3], 16, -722521979);
	b = hh(b, c, d, a, k[6], 23,  76029189);
	a = hh(a, b, c, d, k[9], 4, -640364487);
	d = hh(d, a, b, c, k[12], 11, -421815835);
	c = hh(c, d, a, b, k[15], 16,  530742520);
	b = hh(b, c, d, a, k[2], 23, -995338651);

	a = ii(a, b, c, d, k[0], 6, -198630844);
	d = ii(d, a, b, c, k[7], 10,  1126891415);
	c = ii(c, d, a, b, k[14], 15, -1416354905);
	b = ii(b, c, d, a, k[5], 21, -57434055);
	a = ii(a, b, c, d, k[12], 6,  1700485571);
	d = ii(d, a, b, c, k[3], 10, -1894986606);
	c = ii(c, d, a, b, k[10], 15, -1051523);
	b = ii(b, c, d, a, k[1], 21, -2054922799);
	a = ii(a, b, c, d, k[8], 6,  1873313359);
	d = ii(d, a, b, c, k[15], 10, -30611744);
	c = ii(c, d, a, b, k[6], 15, -1560198380);
	b = ii(b, c, d, a, k[13], 21,  1309151649);
	a = ii(a, b, c, d, k[4], 6, -145523070);
	d = ii(d, a, b, c, k[11], 10, -1120210379);
	c = ii(c, d, a, b, k[2], 15,  718787259);
	b = ii(b, c, d, a, k[9], 21, -343485551);

	x[0] = add32(a, x[0]);
	x[1] = add32(b, x[1]);
	x[2] = add32(c, x[2]);
	x[3] = add32(d, x[3]);
}

function cmn(q, a, b, x, s, t) {
	a = add32(add32(a, q), add32(x, t));
	return add32((a << s) | (a >>> (32 - s)), b);
}

function ff(a, b, c, d, x, s, t) {
	return cmn((b & c) | ((~b) & d), a, b, x, s, t);
}

function gg(a, b, c, d, x, s, t) {
	return cmn((b & d) | (c & (~d)), a, b, x, s, t);
}

function hh(a, b, c, d, x, s, t) {
	return cmn(b ^ c ^ d, a, b, x, s, t);
}

function ii(a, b, c, d, x, s, t) {
	return cmn(c ^ (b | (~d)), a, b, x, s, t);
}

function md51(s) {
	txt = '';
	var n = s.length,
	state = [1732584193, -271733879, -1732584194, 271733878], i;
	for (i=64; i<=s.length; i+=64) {
		md5cycle(state, md5blk(s.substring(i-64, i)));
	}
	s = s.substring(i-64);
	var tail = [0,0,0,0, 0,0,0,0, 0,0,0,0, 0,0,0,0];
	for (i=0; i<s.length; i++)
	tail[i>>2] |= s.charCodeAt(i) << ((i%4) << 3);
	tail[i>>2] |= 0x80 << ((i%4) << 3);
	if (i > 55) {
		md5cycle(state, tail);
		for (i=0; i<16; i++) tail[i] = 0;
	}
	tail[14] = n*8;
	md5cycle(state, tail);
	return state;
}

/* there needs to be support for Unicode here,
 * unless we pretend that we can redefine the MD-5
 * algorithm for multi-byte characters (perhaps
 * by adding every four 16-bit characters and
 * shortening the sum to 32 bits). Otherwise
 * I suggest performing MD-5 as if every character
 * was two bytes--e.g., 0040 0025 = @%--but then
 * how will an ordinary MD-5 sum be matched?
 * There is no way to standardize text to something
 * like UTF-8 before transformation; speed cost is
 * utterly prohibitive. The JavaScript standard
 * itself needs to look at this: it should start
 * providing access to strings as preformed UTF-8
 * 8-bit unsigned value arrays.
 */
function md5blk(s) { /* I figured global was faster.   */
	var md5blks = [], i; /* Andy King said do it this way. */
	for (i=0; i<64; i+=4) {
		md5blks[i>>2] = s.charCodeAt(i)
		+ (s.charCodeAt(i+1) << 8)
		+ (s.charCodeAt(i+2) << 16)
		+ (s.charCodeAt(i+3) << 24);
	}
	return md5blks;
}

var hex_chr = '0123456789abcdef'.split('');

function rhex(n)
{
	var s='', j=0;
	for(; j<4; j++)
	s += hex_chr[(n >> (j * 8 + 4)) & 0x0F]
	+ hex_chr[(n >> (j * 8)) & 0x0F];
	return s;
}

function hex(x) {
	for (var i=0; i<x.length; i++)
	x[i] = rhex(x[i]);
	return x.join('');
}

function md5(s) {
	return hex(md51(s));
}

/* this function is much faster,
so if possible we use it. Some IEs
are the only ones I know of that
need the idiotic second function,
generated by an if clause.  */

function add32(a, b) {
	return (a + b) & 0xFFFFFFFF;
}

if (md5('hello') != '5d41402abc4b2a76b9719d911017c592') {
	function add32(x, y) {
		var lsw = (x & 0xFFFF) + (y & 0xFFFF),
		msw = (x >> 16) + (y >> 16) + (lsw >> 16);
		return (msw << 16) | (lsw & 0xFFFF);
	}
}

(function ($) {
    var $this;
    var $ajaxUrl = '';
    $.fn.pasteUploadImage = function (ajaxUrl) {
        $this = $(this);
        $ajaxUrl = ajaxUrl;
        $this.on('paste', function (event) {
            var filename, image, pasteEvent, text;
            pasteEvent = event.originalEvent;
            if (pasteEvent.clipboardData && pasteEvent.clipboardData.items) {
                image = isImage(pasteEvent);
                if (image) {
                    event.preventDefault();
                    filename = getFilename(pasteEvent) || "image.png";
                    text = "{{" + filename + "(uploading...)}}";
                    pasteText(text);
                    return uploadFile(image.getAsFile(), filename);
                }
            }
        });
        $this.on('drop', function (event) {
            var filename, image, pasteEvent, text;
            pasteEvent = event.originalEvent;
            if (pasteEvent.dataTransfer && pasteEvent.dataTransfer.files) {
                image = isImageForDrop(pasteEvent);
                if (image) {
                    event.preventDefault();
                    filename = pasteEvent.dataTransfer.files[0].name || "image.png";
                    text = "{{" + filename + "(uploading...)}}";
                    pasteText(text);
                    return uploadFile(image, filename);
                }
            }
        });
        return true
    };

    pasteText = function (text) {
        var afterSelection, beforeSelection, caretEnd, caretStart, textEnd;
        caretStart = $this[0].selectionStart;
        caretEnd = $this[0].selectionEnd;
        textEnd = $this.val().length;
        beforeSelection = $this.val().substring(0, caretStart);
        afterSelection = $this.val().substring(caretEnd, textEnd);
        $this.val(beforeSelection + text + afterSelection);
        $this.get(0).setSelectionRange(caretStart + text.length, caretEnd + text.length);
        return $this.trigger("input");
    };
    isImage = function (data) {
        var i, item;
        i = 0;
        while (i < data.clipboardData.items.length) {
            item = data.clipboardData.items[i];
            if (item.type.indexOf("image") !== -1) {
                return item;
            }
            i++;
        }
        return false;
    };
    isImageForDrop = function (data) {
        var i, item;
        i = 0;
        while (i < data.dataTransfer.files.length) {
            item = data.dataTransfer.files[i];
            if (item.type.indexOf("image") !== -1) {
                return item;
            }
            i++;
        }
        return false;
    };
    getFilename = function (e) {
        var value;
        if (window.clipboardData && window.clipboardData.getData) {
            value = window.clipboardData.getData("Text");
        } else if (e.clipboardData && e.clipboardData.getData) {
            value = e.clipboardData.getData("text/plain");
        }
        value = value.split("\r");
        return value[0];
    };
    getMimeType = function (file, filename) {
        var mimeType = file.type;
        var extendName = filename.substring(filename.lastIndexOf('.') + 1);
        if (mimeType != 'image/' + extendName) {
            return 'image/' + extendName;
        }
        return mimeType
    };
    uploadFile = function (file, filename) {
        var formData = new FormData();
        formData.append('imageFile', file);
        formData.append("mimeType", getMimeType(file, filename));

        $.ajax({
            url: $ajaxUrl,
            data: formData,
            type: 'post',
            processData: false,
            contentType: false,
            dataType: 'json',
            xhrFields: {
                withCredentials: true
            },
            success: function (data) {
                if (data.success) {
                    return insertToTextArea(filename, data.message);
                }
                return replaceLoadingTest(filename);
            },
            error: function (xOptions, textStatus) {
                replaceLoadingTest(filename);
                console.log(xOptions.responseText);
            }
        });
    };
    insertToTextArea = function (filename, url) {
        return $this.val(function (index, val) {
            return val.replace("{{" + filename + "(uploading...)}}", "![" + filename + "](" + url + ")" + "\n");
        });
    };
    replaceLoadingTest = function (filename) {
        return $this.val(function (index, val) {
            return val.replace("{{" + filename + "(uploading...)}}", filename + "\n");
        });
    };
})(jQuery);

var QRCode;!function(){function a(a){this.mode=c.MODE_8BIT_BYTE,this.data=a,this.parsedData=[];for(var b=[],d=0,e=this.data.length;e>d;d++){var f=this.data.charCodeAt(d);f>65536?(b[0]=240|(1835008&f)>>>18,b[1]=128|(258048&f)>>>12,b[2]=128|(4032&f)>>>6,b[3]=128|63&f):f>2048?(b[0]=224|(61440&f)>>>12,b[1]=128|(4032&f)>>>6,b[2]=128|63&f):f>128?(b[0]=192|(1984&f)>>>6,b[1]=128|63&f):b[0]=f,this.parsedData=this.parsedData.concat(b)}this.parsedData.length!=this.data.length&&(this.parsedData.unshift(191),this.parsedData.unshift(187),this.parsedData.unshift(239))}function b(a,b){this.typeNumber=a,this.errorCorrectLevel=b,this.modules=null,this.moduleCount=0,this.dataCache=null,this.dataList=[]}function i(a,b){if(void 0==a.length)throw new Error(a.length+"/"+b);for(var c=0;c<a.length&&0==a[c];)c++;this.num=new Array(a.length-c+b);for(var d=0;d<a.length-c;d++)this.num[d]=a[d+c]}function j(a,b){this.totalCount=a,this.dataCount=b}function k(){this.buffer=[],this.length=0}function m(){return"undefined"!=typeof CanvasRenderingContext2D}function n(){var a=!1,b=navigator.userAgent;return/android/i.test(b)&&(a=!0,aMat=b.toString().match(/android ([0-9]\.[0-9])/i),aMat&&aMat[1]&&(a=parseFloat(aMat[1]))),a}function r(a,b){for(var c=1,e=s(a),f=0,g=l.length;g>=f;f++){var h=0;switch(b){case d.L:h=l[f][0];break;case d.M:h=l[f][1];break;case d.Q:h=l[f][2];break;case d.H:h=l[f][3]}if(h>=e)break;c++}if(c>l.length)throw new Error("Too long data");return c}function s(a){var b=encodeURI(a).toString().replace(/\%[0-9a-fA-F]{2}/g,"a");return b.length+(b.length!=a?3:0)}a.prototype={getLength:function(){return this.parsedData.length},write:function(a){for(var b=0,c=this.parsedData.length;c>b;b++)a.put(this.parsedData[b],8)}},b.prototype={addData:function(b){var c=new a(b);this.dataList.push(c),this.dataCache=null},isDark:function(a,b){if(0>a||this.moduleCount<=a||0>b||this.moduleCount<=b)throw new Error(a+","+b);return this.modules[a][b]},getModuleCount:function(){return this.moduleCount},make:function(){this.makeImpl(!1,this.getBestMaskPattern())},makeImpl:function(a,c){this.moduleCount=4*this.typeNumber+17,this.modules=new Array(this.moduleCount);for(var d=0;d<this.moduleCount;d++){this.modules[d]=new Array(this.moduleCount);for(var e=0;e<this.moduleCount;e++)this.modules[d][e]=null}this.setupPositionProbePattern(0,0),this.setupPositionProbePattern(this.moduleCount-7,0),this.setupPositionProbePattern(0,this.moduleCount-7),this.setupPositionAdjustPattern(),this.setupTimingPattern(),this.setupTypeInfo(a,c),this.typeNumber>=7&&this.setupTypeNumber(a),null==this.dataCache&&(this.dataCache=b.createData(this.typeNumber,this.errorCorrectLevel,this.dataList)),this.mapData(this.dataCache,c)},setupPositionProbePattern:function(a,b){for(var c=-1;7>=c;c++)if(!(-1>=a+c||this.moduleCount<=a+c))for(var d=-1;7>=d;d++)-1>=b+d||this.moduleCount<=b+d||(this.modules[a+c][b+d]=c>=0&&6>=c&&(0==d||6==d)||d>=0&&6>=d&&(0==c||6==c)||c>=2&&4>=c&&d>=2&&4>=d?!0:!1)},getBestMaskPattern:function(){for(var a=0,b=0,c=0;8>c;c++){this.makeImpl(!0,c);var d=f.getLostPoint(this);(0==c||a>d)&&(a=d,b=c)}return b},createMovieClip:function(a,b,c){var d=a.createEmptyMovieClip(b,c),e=1;this.make();for(var f=0;f<this.modules.length;f++)for(var g=f*e,h=0;h<this.modules[f].length;h++){var i=h*e,j=this.modules[f][h];j&&(d.beginFill(0,100),d.moveTo(i,g),d.lineTo(i+e,g),d.lineTo(i+e,g+e),d.lineTo(i,g+e),d.endFill())}return d},setupTimingPattern:function(){for(var a=8;a<this.moduleCount-8;a++)null==this.modules[a][6]&&(this.modules[a][6]=0==a%2);for(var b=8;b<this.moduleCount-8;b++)null==this.modules[6][b]&&(this.modules[6][b]=0==b%2)},setupPositionAdjustPattern:function(){for(var a=f.getPatternPosition(this.typeNumber),b=0;b<a.length;b++)for(var c=0;c<a.length;c++){var d=a[b],e=a[c];if(null==this.modules[d][e])for(var g=-2;2>=g;g++)for(var h=-2;2>=h;h++)this.modules[d+g][e+h]=-2==g||2==g||-2==h||2==h||0==g&&0==h?!0:!1}},setupTypeNumber:function(a){for(var b=f.getBCHTypeNumber(this.typeNumber),c=0;18>c;c++){var d=!a&&1==(1&b>>c);this.modules[Math.floor(c/3)][c%3+this.moduleCount-8-3]=d}for(var c=0;18>c;c++){var d=!a&&1==(1&b>>c);this.modules[c%3+this.moduleCount-8-3][Math.floor(c/3)]=d}},setupTypeInfo:function(a,b){for(var c=this.errorCorrectLevel<<3|b,d=f.getBCHTypeInfo(c),e=0;15>e;e++){var g=!a&&1==(1&d>>e);6>e?this.modules[e][8]=g:8>e?this.modules[e+1][8]=g:this.modules[this.moduleCount-15+e][8]=g}for(var e=0;15>e;e++){var g=!a&&1==(1&d>>e);8>e?this.modules[8][this.moduleCount-e-1]=g:9>e?this.modules[8][15-e-1+1]=g:this.modules[8][15-e-1]=g}this.modules[this.moduleCount-8][8]=!a},mapData:function(a,b){for(var c=-1,d=this.moduleCount-1,e=7,g=0,h=this.moduleCount-1;h>0;h-=2)for(6==h&&h--;;){for(var i=0;2>i;i++)if(null==this.modules[d][h-i]){var j=!1;g<a.length&&(j=1==(1&a[g]>>>e));var k=f.getMask(b,d,h-i);k&&(j=!j),this.modules[d][h-i]=j,e--,-1==e&&(g++,e=7)}if(d+=c,0>d||this.moduleCount<=d){d-=c,c=-c;break}}}},b.PAD0=236,b.PAD1=17,b.createData=function(a,c,d){for(var e=j.getRSBlocks(a,c),g=new k,h=0;h<d.length;h++){var i=d[h];g.put(i.mode,4),g.put(i.getLength(),f.getLengthInBits(i.mode,a)),i.write(g)}for(var l=0,h=0;h<e.length;h++)l+=e[h].dataCount;if(g.getLengthInBits()>8*l)throw new Error("code length overflow. ("+g.getLengthInBits()+">"+8*l+")");for(g.getLengthInBits()+4<=8*l&&g.put(0,4);0!=g.getLengthInBits()%8;)g.putBit(!1);for(;;){if(g.getLengthInBits()>=8*l)break;if(g.put(b.PAD0,8),g.getLengthInBits()>=8*l)break;g.put(b.PAD1,8)}return b.createBytes(g,e)},b.createBytes=function(a,b){for(var c=0,d=0,e=0,g=new Array(b.length),h=new Array(b.length),j=0;j<b.length;j++){var k=b[j].dataCount,l=b[j].totalCount-k;d=Math.max(d,k),e=Math.max(e,l),g[j]=new Array(k);for(var m=0;m<g[j].length;m++)g[j][m]=255&a.buffer[m+c];c+=k;var n=f.getErrorCorrectPolynomial(l),o=new i(g[j],n.getLength()-1),p=o.mod(n);h[j]=new Array(n.getLength()-1);for(var m=0;m<h[j].length;m++){var q=m+p.getLength()-h[j].length;h[j][m]=q>=0?p.get(q):0}}for(var r=0,m=0;m<b.length;m++)r+=b[m].totalCount;for(var s=new Array(r),t=0,m=0;d>m;m++)for(var j=0;j<b.length;j++)m<g[j].length&&(s[t++]=g[j][m]);for(var m=0;e>m;m++)for(var j=0;j<b.length;j++)m<h[j].length&&(s[t++]=h[j][m]);return s};for(var c={MODE_NUMBER:1,MODE_ALPHA_NUM:2,MODE_8BIT_BYTE:4,MODE_KANJI:8},d={L:1,M:0,Q:3,H:2},e={PATTERN000:0,PATTERN001:1,PATTERN010:2,PATTERN011:3,PATTERN100:4,PATTERN101:5,PATTERN110:6,PATTERN111:7},f={PATTERN_POSITION_TABLE:[[],[6,18],[6,22],[6,26],[6,30],[6,34],[6,22,38],[6,24,42],[6,26,46],[6,28,50],[6,30,54],[6,32,58],[6,34,62],[6,26,46,66],[6,26,48,70],[6,26,50,74],[6,30,54,78],[6,30,56,82],[6,30,58,86],[6,34,62,90],[6,28,50,72,94],[6,26,50,74,98],[6,30,54,78,102],[6,28,54,80,106],[6,32,58,84,110],[6,30,58,86,114],[6,34,62,90,118],[6,26,50,74,98,122],[6,30,54,78,102,126],[6,26,52,78,104,130],[6,30,56,82,108,134],[6,34,60,86,112,138],[6,30,58,86,114,142],[6,34,62,90,118,146],[6,30,54,78,102,126,150],[6,24,50,76,102,128,154],[6,28,54,80,106,132,158],[6,32,58,84,110,136,162],[6,26,54,82,110,138,166],[6,30,58,86,114,142,170]],G15:1335,G18:7973,G15_MASK:21522,getBCHTypeInfo:function(a){for(var b=a<<10;f.getBCHDigit(b)-f.getBCHDigit(f.G15)>=0;)b^=f.G15<<f.getBCHDigit(b)-f.getBCHDigit(f.G15);return(a<<10|b)^f.G15_MASK},getBCHTypeNumber:function(a){for(var b=a<<12;f.getBCHDigit(b)-f.getBCHDigit(f.G18)>=0;)b^=f.G18<<f.getBCHDigit(b)-f.getBCHDigit(f.G18);return a<<12|b},getBCHDigit:function(a){for(var b=0;0!=a;)b++,a>>>=1;return b},getPatternPosition:function(a){return f.PATTERN_POSITION_TABLE[a-1]},getMask:function(a,b,c){switch(a){case e.PATTERN000:return 0==(b+c)%2;case e.PATTERN001:return 0==b%2;case e.PATTERN010:return 0==c%3;case e.PATTERN011:return 0==(b+c)%3;case e.PATTERN100:return 0==(Math.floor(b/2)+Math.floor(c/3))%2;case e.PATTERN101:return 0==b*c%2+b*c%3;case e.PATTERN110:return 0==(b*c%2+b*c%3)%2;case e.PATTERN111:return 0==(b*c%3+(b+c)%2)%2;default:throw new Error("bad maskPattern:"+a)}},getErrorCorrectPolynomial:function(a){for(var b=new i([1],0),c=0;a>c;c++)b=b.multiply(new i([1,g.gexp(c)],0));return b},getLengthInBits:function(a,b){if(b>=1&&10>b)switch(a){case c.MODE_NUMBER:return 10;case c.MODE_ALPHA_NUM:return 9;case c.MODE_8BIT_BYTE:return 8;case c.MODE_KANJI:return 8;default:throw new Error("mode:"+a)}else if(27>b)switch(a){case c.MODE_NUMBER:return 12;case c.MODE_ALPHA_NUM:return 11;case c.MODE_8BIT_BYTE:return 16;case c.MODE_KANJI:return 10;default:throw new Error("mode:"+a)}else{if(!(41>b))throw new Error("type:"+b);switch(a){case c.MODE_NUMBER:return 14;case c.MODE_ALPHA_NUM:return 13;case c.MODE_8BIT_BYTE:return 16;case c.MODE_KANJI:return 12;default:throw new Error("mode:"+a)}}},getLostPoint:function(a){for(var b=a.getModuleCount(),c=0,d=0;b>d;d++)for(var e=0;b>e;e++){for(var f=0,g=a.isDark(d,e),h=-1;1>=h;h++)if(!(0>d+h||d+h>=b))for(var i=-1;1>=i;i++)0>e+i||e+i>=b||(0!=h||0!=i)&&g==a.isDark(d+h,e+i)&&f++;f>5&&(c+=3+f-5)}for(var d=0;b-1>d;d++)for(var e=0;b-1>e;e++){var j=0;a.isDark(d,e)&&j++,a.isDark(d+1,e)&&j++,a.isDark(d,e+1)&&j++,a.isDark(d+1,e+1)&&j++,(0==j||4==j)&&(c+=3)}for(var d=0;b>d;d++)for(var e=0;b-6>e;e++)a.isDark(d,e)&&!a.isDark(d,e+1)&&a.isDark(d,e+2)&&a.isDark(d,e+3)&&a.isDark(d,e+4)&&!a.isDark(d,e+5)&&a.isDark(d,e+6)&&(c+=40);for(var e=0;b>e;e++)for(var d=0;b-6>d;d++)a.isDark(d,e)&&!a.isDark(d+1,e)&&a.isDark(d+2,e)&&a.isDark(d+3,e)&&a.isDark(d+4,e)&&!a.isDark(d+5,e)&&a.isDark(d+6,e)&&(c+=40);for(var k=0,e=0;b>e;e++)for(var d=0;b>d;d++)a.isDark(d,e)&&k++;var l=Math.abs(100*k/b/b-50)/5;return c+=10*l}},g={glog:function(a){if(1>a)throw new Error("glog("+a+")");return g.LOG_TABLE[a]},gexp:function(a){for(;0>a;)a+=255;for(;a>=256;)a-=255;return g.EXP_TABLE[a]},EXP_TABLE:new Array(256),LOG_TABLE:new Array(256)},h=0;8>h;h++)g.EXP_TABLE[h]=1<<h;for(var h=8;256>h;h++)g.EXP_TABLE[h]=g.EXP_TABLE[h-4]^g.EXP_TABLE[h-5]^g.EXP_TABLE[h-6]^g.EXP_TABLE[h-8];for(var h=0;255>h;h++)g.LOG_TABLE[g.EXP_TABLE[h]]=h;i.prototype={get:function(a){return this.num[a]},getLength:function(){return this.num.length},multiply:function(a){for(var b=new Array(this.getLength()+a.getLength()-1),c=0;c<this.getLength();c++)for(var d=0;d<a.getLength();d++)b[c+d]^=g.gexp(g.glog(this.get(c))+g.glog(a.get(d)));return new i(b,0)},mod:function(a){if(this.getLength()-a.getLength()<0)return this;for(var b=g.glog(this.get(0))-g.glog(a.get(0)),c=new Array(this.getLength()),d=0;d<this.getLength();d++)c[d]=this.get(d);for(var d=0;d<a.getLength();d++)c[d]^=g.gexp(g.glog(a.get(d))+b);return new i(c,0).mod(a)}},j.RS_BLOCK_TABLE=[[1,26,19],[1,26,16],[1,26,13],[1,26,9],[1,44,34],[1,44,28],[1,44,22],[1,44,16],[1,70,55],[1,70,44],[2,35,17],[2,35,13],[1,100,80],[2,50,32],[2,50,24],[4,25,9],[1,134,108],[2,67,43],[2,33,15,2,34,16],[2,33,11,2,34,12],[2,86,68],[4,43,27],[4,43,19],[4,43,15],[2,98,78],[4,49,31],[2,32,14,4,33,15],[4,39,13,1,40,14],[2,121,97],[2,60,38,2,61,39],[4,40,18,2,41,19],[4,40,14,2,41,15],[2,146,116],[3,58,36,2,59,37],[4,36,16,4,37,17],[4,36,12,4,37,13],[2,86,68,2,87,69],[4,69,43,1,70,44],[6,43,19,2,44,20],[6,43,15,2,44,16],[4,101,81],[1,80,50,4,81,51],[4,50,22,4,51,23],[3,36,12,8,37,13],[2,116,92,2,117,93],[6,58,36,2,59,37],[4,46,20,6,47,21],[7,42,14,4,43,15],[4,133,107],[8,59,37,1,60,38],[8,44,20,4,45,21],[12,33,11,4,34,12],[3,145,115,1,146,116],[4,64,40,5,65,41],[11,36,16,5,37,17],[11,36,12,5,37,13],[5,109,87,1,110,88],[5,65,41,5,66,42],[5,54,24,7,55,25],[11,36,12],[5,122,98,1,123,99],[7,73,45,3,74,46],[15,43,19,2,44,20],[3,45,15,13,46,16],[1,135,107,5,136,108],[10,74,46,1,75,47],[1,50,22,15,51,23],[2,42,14,17,43,15],[5,150,120,1,151,121],[9,69,43,4,70,44],[17,50,22,1,51,23],[2,42,14,19,43,15],[3,141,113,4,142,114],[3,70,44,11,71,45],[17,47,21,4,48,22],[9,39,13,16,40,14],[3,135,107,5,136,108],[3,67,41,13,68,42],[15,54,24,5,55,25],[15,43,15,10,44,16],[4,144,116,4,145,117],[17,68,42],[17,50,22,6,51,23],[19,46,16,6,47,17],[2,139,111,7,140,112],[17,74,46],[7,54,24,16,55,25],[34,37,13],[4,151,121,5,152,122],[4,75,47,14,76,48],[11,54,24,14,55,25],[16,45,15,14,46,16],[6,147,117,4,148,118],[6,73,45,14,74,46],[11,54,24,16,55,25],[30,46,16,2,47,17],[8,132,106,4,133,107],[8,75,47,13,76,48],[7,54,24,22,55,25],[22,45,15,13,46,16],[10,142,114,2,143,115],[19,74,46,4,75,47],[28,50,22,6,51,23],[33,46,16,4,47,17],[8,152,122,4,153,123],[22,73,45,3,74,46],[8,53,23,26,54,24],[12,45,15,28,46,16],[3,147,117,10,148,118],[3,73,45,23,74,46],[4,54,24,31,55,25],[11,45,15,31,46,16],[7,146,116,7,147,117],[21,73,45,7,74,46],[1,53,23,37,54,24],[19,45,15,26,46,16],[5,145,115,10,146,116],[19,75,47,10,76,48],[15,54,24,25,55,25],[23,45,15,25,46,16],[13,145,115,3,146,116],[2,74,46,29,75,47],[42,54,24,1,55,25],[23,45,15,28,46,16],[17,145,115],[10,74,46,23,75,47],[10,54,24,35,55,25],[19,45,15,35,46,16],[17,145,115,1,146,116],[14,74,46,21,75,47],[29,54,24,19,55,25],[11,45,15,46,46,16],[13,145,115,6,146,116],[14,74,46,23,75,47],[44,54,24,7,55,25],[59,46,16,1,47,17],[12,151,121,7,152,122],[12,75,47,26,76,48],[39,54,24,14,55,25],[22,45,15,41,46,16],[6,151,121,14,152,122],[6,75,47,34,76,48],[46,54,24,10,55,25],[2,45,15,64,46,16],[17,152,122,4,153,123],[29,74,46,14,75,47],[49,54,24,10,55,25],[24,45,15,46,46,16],[4,152,122,18,153,123],[13,74,46,32,75,47],[48,54,24,14,55,25],[42,45,15,32,46,16],[20,147,117,4,148,118],[40,75,47,7,76,48],[43,54,24,22,55,25],[10,45,15,67,46,16],[19,148,118,6,149,119],[18,75,47,31,76,48],[34,54,24,34,55,25],[20,45,15,61,46,16]],j.getRSBlocks=function(a,b){var c=j.getRsBlockTable(a,b);if(void 0==c)throw new Error("bad rs block @ typeNumber:"+a+"/errorCorrectLevel:"+b);for(var d=c.length/3,e=[],f=0;d>f;f++)for(var g=c[3*f+0],h=c[3*f+1],i=c[3*f+2],k=0;g>k;k++)e.push(new j(h,i));return e},j.getRsBlockTable=function(a,b){switch(b){case d.L:return j.RS_BLOCK_TABLE[4*(a-1)+0];case d.M:return j.RS_BLOCK_TABLE[4*(a-1)+1];case d.Q:return j.RS_BLOCK_TABLE[4*(a-1)+2];case d.H:return j.RS_BLOCK_TABLE[4*(a-1)+3];default:return void 0}},k.prototype={get:function(a){var b=Math.floor(a/8);return 1==(1&this.buffer[b]>>>7-a%8)},put:function(a,b){for(var c=0;b>c;c++)this.putBit(1==(1&a>>>b-c-1))},getLengthInBits:function(){return this.length},putBit:function(a){var b=Math.floor(this.length/8);this.buffer.length<=b&&this.buffer.push(0),a&&(this.buffer[b]|=128>>>this.length%8),this.length++}};var l=[[17,14,11,7],[32,26,20,14],[53,42,32,24],[78,62,46,34],[106,84,60,44],[134,106,74,58],[154,122,86,64],[192,152,108,84],[230,180,130,98],[271,213,151,119],[321,251,177,137],[367,287,203,155],[425,331,241,177],[458,362,258,194],[520,412,292,220],[586,450,322,250],[644,504,364,280],[718,560,394,310],[792,624,442,338],[858,666,482,382],[929,711,509,403],[1003,779,565,439],[1091,857,611,461],[1171,911,661,511],[1273,997,715,535],[1367,1059,751,593],[1465,1125,805,625],[1528,1190,868,658],[1628,1264,908,698],[1732,1370,982,742],[1840,1452,1030,790],[1952,1538,1112,842],[2068,1628,1168,898],[2188,1722,1228,958],[2303,1809,1283,983],[2431,1911,1351,1051],[2563,1989,1423,1093],[2699,2099,1499,1139],[2809,2213,1579,1219],[2953,2331,1663,1273]],o=function(){var a=function(a,b){this._el=a,this._htOption=b};return a.prototype.draw=function(a){function g(a,b){var c=document.createElementNS("http://www.w3.org/2000/svg",a);for(var d in b)b.hasOwnProperty(d)&&c.setAttribute(d,b[d]);return c}var b=this._htOption,c=this._el,d=a.getModuleCount();Math.floor(b.width/d),Math.floor(b.height/d),this.clear();var h=g("svg",{viewBox:"0 0 "+String(d)+" "+String(d),width:"100%",height:"100%",fill:b.colorLight});h.setAttributeNS("http://www.w3.org/2000/xmlns/","xmlns:xlink","http://www.w3.org/1999/xlink"),c.appendChild(h),h.appendChild(g("rect",{fill:b.colorDark,width:"1",height:"1",id:"template"}));for(var i=0;d>i;i++)for(var j=0;d>j;j++)if(a.isDark(i,j)){var k=g("use",{x:String(i),y:String(j)});k.setAttributeNS("http://www.w3.org/1999/xlink","href","#template"),h.appendChild(k)}},a.prototype.clear=function(){for(;this._el.hasChildNodes();)this._el.removeChild(this._el.lastChild)},a}(),p="svg"===document.documentElement.tagName.toLowerCase(),q=p?o:m()?function(){function a(){this._elImage.src=this._elCanvas.toDataURL("image/png"),this._elImage.style.display="block",this._elCanvas.style.display="none"}function d(a,b){var c=this;if(c._fFail=b,c._fSuccess=a,null===c._bSupportDataURI){var d=document.createElement("img"),e=function(){c._bSupportDataURI=!1,c._fFail&&_fFail.call(c)},f=function(){c._bSupportDataURI=!0,c._fSuccess&&c._fSuccess.call(c)};return d.onabort=e,d.onerror=e,d.onload=f,d.src="data:image/gif;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==",void 0}c._bSupportDataURI===!0&&c._fSuccess?c._fSuccess.call(c):c._bSupportDataURI===!1&&c._fFail&&c._fFail.call(c)}if(this._android&&this._android<=2.1){var b=1/window.devicePixelRatio,c=CanvasRenderingContext2D.prototype.drawImage;CanvasRenderingContext2D.prototype.drawImage=function(a,d,e,f,g,h,i,j){if("nodeName"in a&&/img/i.test(a.nodeName))for(var l=arguments.length-1;l>=1;l--)arguments[l]=arguments[l]*b;else"undefined"==typeof j&&(arguments[1]*=b,arguments[2]*=b,arguments[3]*=b,arguments[4]*=b);c.apply(this,arguments)}}var e=function(a,b){this._bIsPainted=!1,this._android=n(),this._htOption=b,this._elCanvas=document.createElement("canvas"),this._elCanvas.width=b.width,this._elCanvas.height=b.height,a.appendChild(this._elCanvas),this._el=a,this._oContext=this._elCanvas.getContext("2d"),this._bIsPainted=!1,this._elImage=document.createElement("img"),this._elImage.style.display="none",this._el.appendChild(this._elImage),this._bSupportDataURI=null};return e.prototype.draw=function(a){var b=this._elImage,c=this._oContext,d=this._htOption,e=a.getModuleCount(),f=d.width/e,g=d.height/e,h=Math.round(f),i=Math.round(g);b.style.display="none",this.clear();for(var j=0;e>j;j++)for(var k=0;e>k;k++){var l=a.isDark(j,k),m=k*f,n=j*g;c.strokeStyle=l?d.colorDark:d.colorLight,c.lineWidth=1,c.fillStyle=l?d.colorDark:d.colorLight,c.fillRect(m,n,f,g),c.strokeRect(Math.floor(m)+.5,Math.floor(n)+.5,h,i),c.strokeRect(Math.ceil(m)-.5,Math.ceil(n)-.5,h,i)}this._bIsPainted=!0},e.prototype.makeImage=function(){this._bIsPainted&&d.call(this,a)},e.prototype.isPainted=function(){return this._bIsPainted},e.prototype.clear=function(){this._oContext.clearRect(0,0,this._elCanvas.width,this._elCanvas.height),this._bIsPainted=!1},e.prototype.round=function(a){return a?Math.floor(1e3*a)/1e3:a},e}():function(){var a=function(a,b){this._el=a,this._htOption=b};return a.prototype.draw=function(a){for(var b=this._htOption,c=this._el,d=a.getModuleCount(),e=Math.floor(b.width/d),f=Math.floor(b.height/d),g=['<table style="border:0;border-collapse:collapse;">'],h=0;d>h;h++){g.push("<tr>");for(var i=0;d>i;i++)g.push('<td style="border:0;border-collapse:collapse;padding:0;margin:0;width:'+e+"px;height:"+f+"px;background-color:"+(a.isDark(h,i)?b.colorDark:b.colorLight)+';"></td>');g.push("</tr>")}g.push("</table>"),c.innerHTML=g.join("");var j=c.childNodes[0],k=(b.width-j.offsetWidth)/2,l=(b.height-j.offsetHeight)/2;k>0&&l>0&&(j.style.margin=l+"px "+k+"px")},a.prototype.clear=function(){this._el.innerHTML=""},a}();QRCode=function(a,b){if(this._htOption={width:256,height:256,typeNumber:4,colorDark:"#000000",colorLight:"#ffffff",correctLevel:d.H},"string"==typeof b&&(b={text:b}),b)for(var c in b)this._htOption[c]=b[c];"string"==typeof a&&(a=document.getElementById(a)),this._android=n(),this._el=a,this._oQRCode=null,this._oDrawing=new q(this._el,this._htOption),this._htOption.text&&this.makeCode(this._htOption.text)},QRCode.prototype.makeCode=function(a){this._oQRCode=new b(r(a,this._htOption.correctLevel),this._htOption.correctLevel),this._oQRCode.addData(a),this._oQRCode.make(),this._el.title=a,this._oDrawing.draw(this._oQRCode),this.makeImage()},QRCode.prototype.makeImage=function(){"function"==typeof this._oDrawing.makeImage&&(!this._android||this._android>=3)&&this._oDrawing.makeImage()},QRCode.prototype.clear=function(){this._oDrawing.clear()},QRCode.CorrectLevel=d}();