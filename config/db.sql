CREATE TABLE `website_setting` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(63) NOT NULL DEFAULT '' COMMENT '网站名称',
  `domain` varchar(63) NOT NULL DEFAULT '' COMMENT '网站域名',
  `title_suffix` varchar(63) NOT NULL DEFAULT '' COMMENT '标题后缀',
  `favicon` varchar(127) NOT NULL DEFAULT '' COMMENT '自定义favicon',
  `logo` varchar(127) NOT NULL DEFAULT '' COMMENT '自定义logo',
  `start_year` int unsigned NOT NULL DEFAULT 0 COMMENT '网站运营开始年份',
  `blog_url` varchar(127) NOT NULL DEFAULT '' COMMENT '独立博客url，没有则留空',
  `slogan` varchar(127) NOT NULL DEFAULT '' COMMENT '网站slogan，在页脚最后',
  `beian` varchar(63) NOT NULL DEFAULT '' COMMENT '网站备案信息',
  `reading_menu` varchar(127) NOT NULL DEFAULT '' COMMENT '技术晨读菜单名，留空则用默认的',
  `docs_menu` varchar(255) NOT NULL DEFAULT '' COMMENT '官方文档菜单，json格式，留空则用默认',
  `footer_nav` varchar(1022) NOT NULL DEFAULT '' COMMENT '底部导航，json格式',
  `friends_logo` varchar(1022) NOT NULL DEFAULT '' COMMENT '底部友情logo，json格式',
  `project_df_logo` varchar(255) NOT NULL DEFAULT '' COMMENT '开源项目默认logo',
  `seo_keywords` varchar(63) NOT NULL DEFAULT '' COMMENT '页面 seo 通用keywords',
  `seo_description` varchar(255) NOT NULL DEFAULT '' COMMENT '页面 seo 通用description',
  `created_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='网站设置信息';

CREATE TABLE IF NOT EXISTS `topics` (
  `tid` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `content` text NOT NULL,
  `nid` int unsigned NOT NULL COMMENT '节点id',
  `uid` int unsigned NOT NULL COMMENT '帖子作者',
  `lastreplyuid` int unsigned NOT NULL DEFAULT 0 COMMENT '最后回复者',
  `lastreplytime` timestamp NOT NULL DEFAULT 0 COMMENT '最后回复时间',
  `flag` tinyint NOT NULL DEFAULT 0 COMMENT '审核标识,0-未审核;1-已审核;2-审核删除;3-用户自己删除',
  `editor_uid` int unsigned NOT NULL DEFAULT 0 COMMENT '最后编辑人',
  `top` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '置顶，0否，1置顶',
  `ctime` timestamp NOT NULL DEFAULT 0,
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`tid`),
  KEY `uid` (`uid`),
  KEY `nid` (`nid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '帖子内容表';

CREATE TABLE IF NOT EXISTS `topics_ex` (
  `tid` int unsigned NOT NULL,
  `view` int unsigned NOT NULL DEFAULT 0 COMMENT '浏览数',
  `reply` int unsigned NOT NULL DEFAULT 0 COMMENT '回复数',
  `like` int unsigned NOT NULL DEFAULT 0 COMMENT '喜欢数',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`tid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '帖子扩展表（计数）';

CREATE TABLE IF NOT EXISTS `topics_node` (
  `nid` int unsigned NOT NULL AUTO_INCREMENT,
  `parent` int unsigned NOT NULL DEFAULT 0 COMMENT '父节点id，无父节点为0',
  `name` varchar(20) NOT NULL COMMENT '节点名',
  `intro` varchar(50) NOT NULL DEFAULT '' COMMENT '节点简介',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`nid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '帖子节点表';

CREATE TABLE IF NOT EXISTS `comments` (
  `cid` int unsigned NOT NULL AUTO_INCREMENT,
  `objid` int unsigned NOT NULL COMMENT '对象id，属主（评论给谁）',
  `objtype` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '类型,0-帖子;1-博文;2-资源;3-wiki;4-project',
  `content` text NOT NULL,
  `uid` int unsigned NOT NULL COMMENT '回复者',
  `floor` int unsigned NOT NULL COMMENT '第几楼',
  `flag` tinyint NOT NULL DEFAULT 0 COMMENT '审核标识,0-未审核;1-已审核;2-审核删除;3-用户自己删除',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`cid`),
  UNIQUE KEY (`objid`,`objtype`,`floor`),
  KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '评论表（帖子回复、博客文章评论等，统一处理）';

CREATE TABLE IF NOT EXISTS `likes` (
  `uid` int unsigned NOT NULL DEFAULT 0 COMMENT '喜欢人的uid',
  `objtype` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '类型,0-帖子;1-博文;2-资源;3-wiki',
  `objid` int unsigned NOT NULL DEFAULT 0 COMMENT '对象id，属主',
  `flag` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '1-喜欢；2-不喜欢（暂时不支持）',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`uid`,`objtype`,`objid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '喜欢表';

CREATE TABLE IF NOT EXISTS `user_login` (
  `uid` int unsigned NOT NULL,
  `email` varchar(128) NOT NULL DEFAULT '',
  `username` varchar(20) NOT NULL COMMENT '用户名',
  `passcode` char(12) NOT NULL DEFAULT '' COMMENT '加密随机数',
  `passwd` char(32) NOT NULL DEFAULT '' COMMENT 'md5密码',
  `login_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后一次登录时间（主动登录或cookie登录）',
  PRIMARY KEY (`uid`),
  UNIQUE KEY (`username`),
  UNIQUE KEY (`email`),
  KEY `logintime` (`login_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户登录表';
   
CREATE TABLE IF NOT EXISTS `bind_user` (
  `uid` int unsigned NOT NULL,
  `type` tinyint NOT NULL DEFAULT 0 COMMENT '绑定的第三方类型',
  `email` varchar(128) NOT NULL DEFAULT '',
  `tuid` int unsigned NOT NULL DEFAULT 0 COMMENT '第三方uid',
  `username` varchar(20) NOT NULL COMMENT '用户名',
  `token` varchar(50) NOT NULL COMMENT '第三方access_token',
  `refresh` varchar(50) NOT NULL COMMENT '第三方refresh_token',
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '第三方绑定表';

CREATE TABLE IF NOT EXISTS `user_info` (
  `uid` int unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(128) NOT NULL DEFAULT '',
  `open` tinyint NOT NULL DEFAULT 1 COMMENT '邮箱是否公开，默认公开',
  `username` varchar(20) NOT NULL COMMENT '用户名',
  `name` varchar(20) NOT NULL DEFAULT '' COMMENT '姓名',
  `avatar` varchar(128) NOT NULL DEFAULT '' COMMENT '头像(如果为空，则使用http://www.gravatar.com)',
  `city` varchar(10) NOT NULL DEFAULT '' COMMENT '居住地',
  `company` varchar(64) NOT NULL DEFAULT '',
  `github` varchar(20) NOT NULL DEFAULT '',
  `weibo` varchar(20) NOT NULL DEFAULT '',
  `website` varchar(50) NOT NULL DEFAULT '' COMMENT '个人主页，博客',
  `monlog` varchar(140) NOT NULL DEFAULT '' COMMENT '个人状态，签名，独白',
  `introduce` text NOT NULL COMMENT '个人简介',
  `unsubscribe` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否退订本站邮件，0-否；1-是',
  `status` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '用户账号状态。0-默认；1-已审核；2-拒绝；3-冻结；4-停号',
  `is_root` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否超级用户，不受权限控制：1-是',
  `ctime` timestamp NOT NULL DEFAULT 0,
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`uid`),
  UNIQUE KEY (`username`),
  UNIQUE KEY (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户信息表';

CREATE TABLE IF NOT EXISTS `user_active` (
  `uid` int unsigned NOT NULL,
  `email` varchar(128) NOT NULL,
  `username` varchar(20) NOT NULL COMMENT '用户名',
  `weight` smallint NOT NULL DEFAULT 1 COMMENT '活跃度，越大越活跃',
  `avatar` varchar(128) NOT NULL DEFAULT '' COMMENT '头像(如果为空，则使用http://www.gravatar.com)',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`uid`),
  UNIQUE KEY (`username`),
  UNIQUE KEY (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '活跃用户表';

CREATE TABLE IF NOT EXISTS `role` (
  `roleid` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '角色名',
  `op_user` varchar(20) NOT NULL DEFAULT '' COMMENT '操作人',
  `ctime` timestamp NOT NULL DEFAULT 0,
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`roleid`),
  UNIQUE KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '角色表，常驻内存';

CREATE TABLE IF NOT EXISTS `authority` (
  `aid` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '权限名',
  `menu1` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '所属一级菜单，本身为一级菜单，则为0',
  `menu2` int unsigned NOT NULL DEFAULT 0 COMMENT '所属二级菜单，本身为二级菜单，则为0',
  `route` varchar(128) NOT NULL DEFAULT '' COMMENT '路由（权限）',
  `op_user` varchar(20) NOT NULL COMMENT '操作人',
  `ctime` timestamp NOT NULL DEFAULT 0,
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`aid`),
  KEY (`route`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '权限表，常驻内存';

CREATE TABLE IF NOT EXISTS `role_authority` (
  `roleid` int unsigned NOT NULL,
  `aid` int unsigned NOT NULL,
  `op_user` varchar(20) NOT NULL DEFAULT '' COMMENT '操作人',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`roleid`, `aid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '角色拥有的权限表';

CREATE TABLE IF NOT EXISTS `user_role` (
  `uid` int unsigned NOT NULL,
  `roleid` int unsigned NOT NULL,
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`uid`, `roleid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户角色表（用户是什么角色，可以多个角色）';


CREATE TABLE IF NOT EXISTS `message` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `content` text NOT NULL COMMENT '消息内容',
  `hasread` enum('未读','已读') NOT NULL DEFAULT '未读',
  `from` int unsigned NOT NULL DEFAULT 0 COMMENT '来自谁',
  `fdel` enum('未删','已删') NOT NULL DEFAULT '未删' COMMENT '发送方删除标识',
  `to` int unsigned NOT NULL COMMENT '发给谁',
  `tdel` enum('未删','已删') NOT NULL DEFAULT '未删' COMMENT '接收方删除标识',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY (`to`),
  KEY (`from`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'message 短消息（私信）';

CREATE TABLE IF NOT EXISTS `system_message` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `msgtype` tinyint NOT NULL DEFAULT 0 COMMENT '系统消息类型',
  `hasread` enum('未读','已读') NOT NULL DEFAULT '未读',
  `to` int unsigned NOT NULL COMMENT '发给谁',
  `ext` text NOT NULL COMMENT '额外信息',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY (`to`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'system_message 系统消息表';

CREATE TABLE IF NOT EXISTS `wiki` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL COMMENT 'wiki标题',
  `content` longtext NOT NULL COMMENT 'wiki内容',
  `uri` varchar(50) NOT NULL COMMENT 'uri',
  `uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '作者',
  `cuid` varchar(100) NOT NULL DEFAULT '' COMMENT '贡献者uid,多个逗号分隔',
  `viewnum` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '浏览数',
  `ctime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uri` (`uri`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'wiki页';

CREATE TABLE IF NOT EXISTS `resource` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL COMMENT '资源标题',
  `form` enum('只是链接','包括内容'),
  `content` longtext NOT NULL COMMENT '资源内容',
  `url` varchar(150) NOT NULL COMMENT '链接url',
  `uid` int unsigned NOT NULL COMMENT '作者',
  `catid` int unsigned NOT NULL COMMENT '所属类别',
  `ctime` timestamp NOT NULL DEFAULT 0,
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY (`url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '资源';

CREATE TABLE IF NOT EXISTS `resource_ex` (
  `id` int unsigned NOT NULL,
  `viewnum` int unsigned NOT NULL DEFAULT 0 COMMENT '浏览数',
  `cmtnum` int unsigned NOT NULL DEFAULT 0 COMMENT '回复数',
  `likenum` int unsigned NOT NULL DEFAULT 0 COMMENT '喜欢数',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '资源扩展表';

CREATE TABLE IF NOT EXISTS `resource_category` (
  `catid` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(20) NOT NULL COMMENT '分类名',
  `intro` varchar(50) NOT NULL DEFAULT '' COMMENT '分类简介',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`catid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '资源分类表';

CREATE TABLE IF NOT EXISTS `articles` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `domain` varchar(50) NOT NULL DEFAULT '' COMMENT '来源域名（不一定是顶级域名）',
  `name` varchar(30) NOT NULL DEFAULT '' COMMENT '来源名称',
  `title` varchar(127) NOT NULL DEFAULT '' COMMENT '文章标题',
  `cover` varchar(127) NOT NULL DEFAULT '' COMMENT '图片封面',
  `author` varchar(1024) NOT NULL DEFAULT '' COMMENT '文章作者(可能带html)',
  `author_txt` varchar(127) NOT NULL DEFAULT '' COMMENT '文章作者(纯文本)',
  `lang` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '语言：0-中文；1-英文',
  `pub_date` varchar(20) NOT NULL DEFAULT '' COMMENT '发布时间',
  `url` varchar(255) NOT NULL DEFAULT '' COMMENT '文章原始链接',
  `content` longtext NOT NULL COMMENT '正文(带html)',
  `txt` text NOT NULL COMMENT '正文(纯文本)',
  `tags` varchar(50) NOT NULL DEFAULT '' COMMENT '文章tag，逗号分隔',
  `css` varchar(255) NOT NULL DEFAULT '' COMMENT '需要额外引入的css样式',
  `viewnum` int unsigned NOT NULL DEFAULT 0 COMMENT '浏览数',
  `cmtnum` int unsigned NOT NULL DEFAULT 0 COMMENT '评论数',
  `likenum` int unsigned NOT NULL DEFAULT 0 COMMENT '赞数',
  `top` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '置顶，0否，1置顶',
  `status` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '状态：0-初始抓取；1-已上线；2-下线(审核拒绝)',
  `op_user` varchar(20) NOT NULL DEFAULT '' COMMENT '操作人',
  `ctime` timestamp NOT NULL DEFAULT 0,
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`url`),
  KEY (`top`),
  KEY (`domain`),
  KEY (`ctime`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '网络文章聚合表';

CREATE TABLE IF NOT EXISTS `crawl_rule` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `domain` varchar(50) NOT NULL DEFAULT '' COMMENT '来源域名（不一定是顶级域名）',
  `subpath` varchar(20) NOT NULL DEFAULT '' COMMENT '域名下面紧接着的path（区别同一网站多个路径不同抓取规则）',
  `name` varchar(30) NOT NULL DEFAULT '' COMMENT '来源名称',
  `lang` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '主要语言：0-中文;1-英文',
  `title` varchar(127) NOT NULL DEFAULT '' COMMENT '文章标题规则',
  `in_url` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '作者信息是否在url中；0-否;1-是(是的时候，author表示在url中的位置)',
  `author` varchar(127) NOT NULL DEFAULT '' COMMENT '文章作者规则',
  `pub_date` varchar(127) NOT NULL DEFAULT '' COMMENT '发布时间规则',
  `content` varchar(127) NOT NULL DEFAULT '' COMMENT '正文规则',
  `ext` varchar(1022) NOT NULL DEFAULT '' COMMENT '扩展，比如附加 css 等，json格式',
  `op_user` varchar(20) NOT NULL DEFAULT '' COMMENT '操作人',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`domain`,`subpath`),
  KEY (`ctime`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '网站抓取规则表';

CREATE TABLE IF NOT EXISTS `auto_crawl_rule` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `website` varchar(20) NOT NULL DEFAULT '' COMMENT '网站标识，有可能是网站某个子模块',
  `all_url` varchar(127) NOT NULL DEFAULT '' COMMENT '全量url，关键词占位符使用%s',
  `incr_url` varchar(127) NOT NULL DEFAULT '' COMMENT '增量url，关键词占位符使用%s',
  `keywords` varchar(63) NOT NULL DEFAULT '' COMMENT '搜索关键词，多个逗号分隔',
  `list_selector` varchar(31) NOT NULL DEFAULT '' COMMENT '列表选择器',
  `result_selector` varchar(31) NOT NULL DEFAULT '' COMMENT '结果选择器，获取具体文章的 url',
  `page_field` varchar(20) NOT NULL DEFAULT '' COMMENT '分页字段名',
  `max_page` int unsigned NOT NULL DEFAULT 0 COMMENT '全量最多抓取多少页',
  `ext` varchar(1023) NOT NULL DEFAULT '' COMMENT '扩展信息，某些网站的特殊配置，json格式',
  `status` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '状态：0-自动抓取；1-停止抓取',
  `op_user` varchar(20) NOT NULL DEFAULT '' COMMENT '操作人',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `website` (`website`),
  KEY `mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='网站自动抓取规则表';

CREATE TABLE IF NOT EXISTS `dynamic` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '动态内容',
  `dmtype` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '类型：0-Go动态；1-本站动态',
  `url` varchar(255) NOT NULL DEFAULT '' COMMENT '链接',
  `seq` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '顺序（越大越在前）',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY (`seq`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '动态表';

CREATE TABLE IF NOT EXISTS `search_stat` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `keyword` varchar(127) NOT NULL DEFAULT '' COMMENT '搜索词',
  `times` int unsigned NOT NULL DEFAULT 0 COMMENT '次数',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`keyword`),
  KEY (`times`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '搜索词统计';

CREATE TABLE IF NOT EXISTS `favorites` (
  `uid` int unsigned NOT NULL DEFAULT 0 COMMENT '用户uid',
  `objtype` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '类型,0-帖子;1-博文;2-资源;3-wiki',
  `objid` int unsigned NOT NULL DEFAULT 0 COMMENT '对象id，属主',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`uid`,`objtype`,`objid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户收藏';

CREATE TABLE IF NOT EXISTS `open_project` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '项目id',
  `name` varchar(127) NOT NULL DEFAULT '' COMMENT '项目名(软件名)，如 Docker',
  `category` varchar(127) NOT NULL DEFAULT '' COMMENT '项目类别，如 Linux 容器引擎',
  `uri` varchar(127) NOT NULL DEFAULT '' COMMENT '项目uri，访问使用（如/p/docker 中的 docker)',
  `home` varchar(127) NOT NULL DEFAULT '' COMMENT '项目首页',
  `doc` varchar(127) NOT NULL DEFAULT '' COMMENT '项目文档地址',
  `download` varchar(127) NOT NULL DEFAULT '' COMMENT '项目下载地址',
  `src` varchar(127) NOT NULL DEFAULT '' COMMENT '源码地址',
  `logo` varchar(127) NOT NULL DEFAULT '' COMMENT '项目logo',
  `desc` text NOT NULL COMMENT '项目描述',
  `repo` varchar(127) NOT NULL DEFAULT '' COMMENT '源码uri部分，方便repo widget插件使用',
  `author` varchar(127) NOT NULL DEFAULT '' COMMENT '作者',
  `licence` varchar(127) NOT NULL DEFAULT '' COMMENT '授权协议',
  `lang` varchar(127) NOT NULL DEFAULT '' COMMENT '开发语言',
  `os` varchar(127) NOT NULL DEFAULT '' COMMENT '操作系统（多个逗号分隔）',
  `tags` varchar(127) NOT NULL DEFAULT '' COMMENT 'tag，逗号分隔',
  `username` varchar(127) NOT NULL DEFAULT '' COMMENT '收录人',
  `viewnum` int unsigned NOT NULL DEFAULT 0 COMMENT '浏览数',
  `cmtnum` int unsigned NOT NULL DEFAULT 0 COMMENT '评论数',
  `likenum` int unsigned NOT NULL DEFAULT 0 COMMENT '赞数',
  `status` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '状态：0-新建；1-已上线；2-下线(审核拒绝)',
  `ctime` timestamp NOT NULL DEFAULT 0 COMMENT '加入时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY (`uri`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '开源项目';

CREATE TABLE IF NOT EXISTS `morning_reading` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '晨读内容',
  `rtype` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '晨读类别：0-Go技术晨读;1-综合技术晨读',
  `inner` int unsigned NOT NULL DEFAULT 0 COMMENT '本站文章id，如果外站文章，则为0',
  `url` varchar(255) NOT NULL DEFAULT '' COMMENT '文章链接，本站文章时为空',
  `moreurls` varchar(1024) NOT NULL DEFAULT '' COMMENT '可能顺带推荐多篇文章；url逗号分隔',
  `clicknum` int unsigned NOT NULL DEFAULT 0 COMMENT '点击数',
  `username` varchar(20) NOT NULL DEFAULT '' COMMENT '发布人',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '技术晨读表';

CREATE TABLE IF NOT EXISTS `image` (
  `pid` int unsigned NOT NULL AUTO_INCREMENT,
  `md5` char(32) NOT NULL DEFAULT '' COMMENT '图片md5',
  `path` varchar(127) NOT NULL DEFAULT '' COMMENT '图片路径（不包括域名）',
  `size` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '图片的大小（字节）',
  `width` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '图片宽度',
  `height` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '图片高度',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`pid`),
  UNIQUE KEY `md5` (`md5`),
  KEY `created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='图片表';

CREATE TABLE `book` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(63) NOT NULL DEFAULT '' COMMENT '书名',
  `ename` varchar(63) NOT NULL DEFAULT '' COMMENT '英文书名',
  `lang` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '语言：0-中文；1-英文',
  `author` varchar(63) NOT NULL DEFAULT '' COMMENT '作者',
  `translator` varchar(31) NOT NULL DEFAULT '' COMMENT '译者',
  `cover` varchar(127) NOT NULL DEFAULT '' COMMENT '封面',
  `pub_date` varchar(10) NOT NULL DEFAULT '' COMMENT '出版日期，2017-04-01',
  `desc` varchar(2046) NOT NULL DEFAULT '' COMMENT '简介',
  `catalogue` text NOT NULL COMMENT '目录',
  `is_free` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否免费；0-否；1-是',
  `online_url` varchar(127) NOT NULL DEFAULT '' COMMENT '在线阅读url',
  `download_url` varchar(127) NOT NULL DEFAULT '' COMMENT '下载url',
  `buy_url` varchar(127) NOT NULL DEFAULT '' COMMENT '购买url',
  `price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '参考价格',
  `viewnum` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '浏览数',
  `cmtnum` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '评论数',
  `likenum` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '赞数（推荐数）',
  `created_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  KEY `name` (`name`),
  KEY `created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='图书表';

CREATE TABLE IF NOT EXISTS `advertisement` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(20) NOT NULL DEFAULT '' COMMENT '广告名称',
  `code` varchar(255) NOT NULL DEFAULT '' COMMENT '广告内容代码(html、js等)',
  `source` varchar(20) NOT NULL DEFAULT '' COMMENT '广告来源，如 baidu_union，shiyanlou',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '广告表';

CREATE TABLE IF NOT EXISTS `page_ad` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `ad_id` varchar(20) NOT NULL DEFAULT '' COMMENT '广告名称',
  `code` varchar(255) NOT NULL DEFAULT '' COMMENT '广告内容代码(html、js等)',
  `source` varchar(20) NOT NULL DEFAULT '' COMMENT '广告来源，如 baidu_union，shiyanlou',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '页面广告管理表';

CREATE TABLE IF NOT EXISTS `friend_link` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(15) NOT NULL DEFAULT '' COMMENT '链接名',
  `url` varchar(255) NOT NULL DEFAULT '' COMMENT '链接URL',
  `seq` smallint unsigned NOT NULL DEFAULT 100 COMMENT '排序',
  `logo` varchar(63) NOT NULL DEFAULT '' COMMENT 'LOGO url',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '友情链接';

CREATE TABLE IF NOT EXISTS `learning_material` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(31) NOT NULL DEFAULT '' COMMENT '标题',
  `url` varchar(63) NOT NULL DEFAULT '' COMMENT '资料URL',
  `type` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '类型，0-文本；1-视频',
  `seq` smallint unsigned NOT NULL DEFAULT 100 COMMENT '排序',
  `first_url` varchar(63) NOT NULL DEFAULT '' COMMENT '开始学习的url',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '成体系的学习资料';

CREATE TABLE IF NOT EXISTS `default_avatar` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `filename` varchar(31) NOT NULL DEFAULT '' COMMENT '图像文件名',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '默认头像';

