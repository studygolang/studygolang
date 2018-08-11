CREATE TABLE IF NOT EXISTS `website_setting` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(63) NOT NULL DEFAULT '' COMMENT '网站名称',
  `domain` varchar(63) NOT NULL DEFAULT '' COMMENT '网站域名',
  `only_https` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否只支持HTTPS',
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
  `index_nav` varchar(2044) NOT NULL DEFAULT '' COMMENT '首页顶部导航，json 格式',
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
  `lastreplytime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '最后回复时间',
  `flag` tinyint NOT NULL DEFAULT 0 COMMENT '审核标识,0-未审核;1-已审核;2-审核删除;3-用户自己删除',
  `editor_uid` int unsigned NOT NULL DEFAULT 0 COMMENT '最后编辑人',
  `top` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '置顶，0否，1置顶',
  `top_time` int unsigned NOT NULL DEFAULT 0 COMMENT '置顶时间',
  `tags` varchar(63) NOT NULL DEFAULT '' COMMENT 'tag，逗号分隔',
  `permission` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '访问权限：0-公开；1-登录用户可见；2-关注的人可见',
  `ctime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`tid`),
  KEY `uid` (`uid`),
  KEY `nid` (`nid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '主题内容表';

CREATE TABLE IF NOT EXISTS `topics_ex` (
  `tid` int unsigned NOT NULL,
  `view` int unsigned NOT NULL DEFAULT 0 COMMENT '浏览数',
  `reply` int unsigned NOT NULL DEFAULT 0 COMMENT '回复数',
  `like` int unsigned NOT NULL DEFAULT 0 COMMENT '喜欢数',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`tid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '主题扩展表（计数）';

CREATE TABLE IF NOT EXISTS `topic_append` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `tid` int unsigned NOT NULL DEFAULT 0 COMMENT '主题 TID',
  `content` text NOT NULL COMMENT '附言内容',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `tid` (`tid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '主题附言表';

CREATE TABLE IF NOT EXISTS `topics_node` (
  `nid` int unsigned NOT NULL AUTO_INCREMENT,
  `parent` int unsigned NOT NULL DEFAULT 0 COMMENT '父节点id，无父节点为0',
  `logo` varchar(127) NOT NULL DEFAULT '' COMMENT '节点logo',
  `name` varchar(20) NOT NULL DEFAULT '' COMMENT '节点名',
  `ename` varchar(15) NOT NULL DEFAULT '' COMMENT '节点英文名，用于导航',
  `intro` varchar(127) NOT NULL DEFAULT '' COMMENT '节点简介',
  `seq` smallint unsigned NOT NULL DEFAULT 0 COMMENT '节点排序，小的在前',
  `show_index` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '主题是否在首页显示',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`nid`),
  INDEX `idx_ename` (`ename`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '帖子节点表';

CREATE TABLE IF NOT EXISTS `recommend_node` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(20) NOT NULL DEFAULT '' COMMENT '虚拟节点名',
  `parent` int unsigned NOT NULL DEFAULT 0 COMMENT '父节点id，无父节点为0',
  `nid` int unsigned NOT NULL COMMENT 'topics_node nid，虚拟节点为0',
  `seq` smallint(6) NOT NULL DEFAULT '0' COMMENT '节点排序，小的在前',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '导航推荐节点';

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
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` int unsigned NOT NULL DEFAULT 0 COMMENT '本站用户UID',
  `type` tinyint NOT NULL DEFAULT 0 COMMENT '绑定的第三方类型,0-github',
  `email` varchar(128) NOT NULL DEFAULT '' COMMENT '第三方邮箱',
  `tuid` int unsigned NOT NULL DEFAULT 0 COMMENT '第三方uid',
  `username` varchar(20) NOT NULL DEFAULT '' COMMENT '第三方用户名',
  `name` varchar(31) NOT NULL DEFAULT '' COMMENT '姓名',
  `access_token` varchar(50) NOT NULL DEFAULT ''  COMMENT '第三方access_token',
  `refresh_token` varchar(50) NOT NULL DEFAULT '' COMMENT '第三方refresh_token',
  `expire` int unsigned NOT NULL DEFAULT 0 COMMENT '过期时间',
  `avatar` varchar(127) NOT NULL DEFAULT '' COMMENT '第三方头像',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_user_type` (`username`,`type`),
  INDEX idx_uid (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '第三方绑定表';

CREATE TABLE IF NOT EXISTS `user_info` (
  `uid` int unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(128) NOT NULL DEFAULT '',
  `open` tinyint NOT NULL DEFAULT 0 COMMENT '邮箱是否公开，默认不公开',
  `username` varchar(20) NOT NULL COMMENT '用户名',
  `name` varchar(20) NOT NULL DEFAULT '' COMMENT '姓名',
  `avatar` varchar(128) NOT NULL DEFAULT '' COMMENT '头像(如果为空，则使用http://www.gravatar.com)',
  `city` varchar(10) NOT NULL DEFAULT '' COMMENT '居住地',
  `company` varchar(63) NOT NULL DEFAULT '' COMMENT '公司',
  `github` varchar(31) NOT NULL DEFAULT '' COMMENT 'Github昵称',
  `weibo` varchar(31) NOT NULL DEFAULT '' COMMENT '微博昵称',
  `website` varchar(63) NOT NULL DEFAULT '' COMMENT '个人主页，博客',
  `monlog` varchar(140) NOT NULL DEFAULT '' COMMENT '个人状态，签名，独白',
  `introduce` varchar(2022) NOT NULL COMMENT '个人简介',
  `unsubscribe` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否退订本站邮件，0-否；1-是',
  `is_third` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否通过第三方账号注册',
  `balance` int unsigned NOT NULL DEFAULT 0 COMMENT '财富余额（铜币）',
  `dau_auth` int unsigned NOT NULL DEFAULT 0 COMMENT '控制用户权限，如能否发文章等',
  `status` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '用户账号状态。0-默认；1-已审核；2-拒绝；3-冻结；4-停号',
  `is_root` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否超级用户，不受权限控制：1-是',
  `ctime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
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
  `ctime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
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
  `ctime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
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
  `tags` varchar(63) NOT NULL DEFAULT '' COMMENT 'tag，逗号分隔',
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
  `lastreplyuid` int unsigned NOT NULL DEFAULT 0 COMMENT '最后回复者',
  `lastreplytime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '最后回复时间',
  `tags` varchar(63) NOT NULL DEFAULT '' COMMENT 'tag，逗号分隔',
  `ctime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
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
  `tags` varchar(63) NOT NULL DEFAULT '' COMMENT '文章tag，逗号分隔',
  `css` varchar(255) NOT NULL DEFAULT '' COMMENT '需要额外引入的css样式',
  `viewnum` int unsigned NOT NULL DEFAULT 0 COMMENT '浏览数',
  `cmtnum` int unsigned NOT NULL DEFAULT 0 COMMENT '评论数',
  `likenum` int unsigned NOT NULL DEFAULT 0 COMMENT '赞数',
  `lastreplyuid` int unsigned NOT NULL DEFAULT 0 COMMENT '最后回复者',
  `lastreplytime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '最后回复时间',
  `top` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '置顶，0否，1置顶',
  `markdown` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否是markwon格式：0-否，1-是',
  `gctt` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否是 gctt 翻译：0-否则；1-是',
  `status` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '状态：0-初始抓取；1-已上线；2-下线(审核拒绝)',
  `op_user` varchar(20) NOT NULL DEFAULT '' COMMENT '操作人',
  `ctime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`url`),
  KEY (`top`),
  KEY (`author_txt`),
  KEY (`domain`),
  KEY (`mtime`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '网络文章聚合表';

CREATE TABLE IF NOT EXISTS `article_gctt` (
  `article_id` int unsigned NOT NULL COMMENT '文章ID',
  `author` varchar(31) NOT NULL DEFAULT '' COMMENT '原文作者',
  `author_url` varchar(127) NOT NULL DEFAULT '' COMMENT '原文作者的主页',
  `translator` varchar(31) NOT NULl DEFAULT '' COMMENT 'gctt 译者',
  `checker` varchar(31) NOT NULl DEFAULT '' COMMENT 'gctt 校对',
  `url` varchar(255) NOT NULL DEFAULT '' COMMENT '原文链接',
  PRIMARY KEY (`article_id`),
  UNIQUE KEY (`url`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'gctt 翻译文章信息表';

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
  `lastreplyuid` int unsigned NOT NULL DEFAULT 0 COMMENT '最后回复者',
  `lastreplytime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '最后回复时间',
  `status` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '状态：0-新建；1-已上线；2-下线(审核拒绝)',
  `ctime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '加入时间',
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

CREATE TABLE IF NOT EXISTS `book` (
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
  `tags` varchar(63) NOT NULL DEFAULT '' COMMENT 'tag，逗号分隔',
  `is_free` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否免费；0-否；1-是',
  `online_url` varchar(127) NOT NULL DEFAULT '' COMMENT '在线阅读url',
  `download_url` varchar(127) NOT NULL DEFAULT '' COMMENT '下载url',
  `buy_url` varchar(127) NOT NULL DEFAULT '' COMMENT '购买url',
  `price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '参考价格',
  `lastreplyuid` int unsigned NOT NULL DEFAULT 0 COMMENT '最后回复者',
  `lastreplytime` timestamp NOT NULL DEFAULT '2010-01-01 00:00:00' COMMENT '最后回复时间',
  `viewnum` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '浏览数',
  `cmtnum` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '评论数',
  `likenum` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '赞数（推荐数）',
  `uid` int unsigned NOT NULL DEFAULT 0 COMMENT '分享人UID',
  `created_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  KEY `name` (`name`),
  KEY `created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='图书表';

CREATE TABLE IF NOT EXISTS `advertisement` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(31) NOT NULL DEFAULT '' COMMENT '广告名称',
  `ad_type` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '广告类型：0-直接在html显示；1-js 操作 html',
  `code` varchar(1022) NOT NULL DEFAULT '' COMMENT '广告内容代码(html、js等)',
  `source` varchar(31) NOT NULL DEFAULT '' COMMENT '广告来源，如 百度联盟，阿里云',
  `is_online` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否在线：0-下线；1-在线',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '广告表';

CREATE TABLE IF NOT EXISTS  `page_ad` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `path` varchar(31) NOT NULL DEFAULT '' COMMENT '页面路径',
  `ad_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '广告ID',
  `position` varchar(15) NOT NULL DEFAULT '' COMMENT '广告在页面的位置，英文字符串',
  `is_online` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否在线：0-下线；1-在线',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_path` (`path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='页面广告管理表';

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

CREATE TABLE IF NOT EXISTS `user_setting` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(31) NOT NULL DEFAULT '' COMMENT '配置项名称',
  `value` int NOT NULL DEFAULT 0 COMMENT '配置项值',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '配置项说明',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE `uniq_key`(`key`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户行为信息设置';

CREATE TABLE IF NOT EXISTS `user_balance_detail` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` int unsigned NOT NULL DEFAULT 0 COMMENT '用户UID',
  `type` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '类型',
  `num` int NOT NULL DEFAULT 0 COMMENT '数额，负数表示减少，正数表示增加',
  `balance` int unsigned NOT NULL DEFAULT 0 COMMENT '余额（铜币）',
  `desc` varchar(1022) NOT NULL DEFAULT '' COMMENT '具体原因，支持html格式',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_uid`(`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户余额明细';

CREATE TABLE IF NOT EXISTS `mission` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(31) NOT NULL DEFAULT '' COMMENT '任务名，如每日登录奖励',
  `type` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '类型 1-每日登录奖励，2-初始资本，3-分享获得',
  `fixed` int unsigned NOT NULl DEFAULT 0 COMMENT '固定奖励多少铜币',
  `min` int unsigned NOT NULL DEFAULT 0 COMMENT '奖励最少铜币，连续型任务',
  `max` int unsigned NOT NULL DEFAULT 0 COMMENT '奖励最多铜币，连续型任务',
  `incr` int unsigned NOT NULL DEFAULT 0 COMMENT '连续登录增量，连续型任务',
  `state` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '状态: 0-正常，未完成；1-已过期；2-已下线',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '任务表';

CREATE TABLE IF NOT EXISTS `user_login_mission` (
  `uid` int unsigned NOT NULL DEFAULT 0 COMMENT '用户UID',
  `date` int unsigned NOT NULL DEFAULT 0 COMMENT '最新领取日期',
  `award` int unsigned NOT NULL DEFAULT 0 COMMENT '最新领取的奖励（铜币）',
  `days` int unsigned NOT NULL DEFAULT 0 COMMENT '连续登录领取天数',
  `total_days` int unsigned NOT NULL DEFAULT 0 COMMENT '总登录领取天数',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户登录任务';

CREATE TABLE IF NOT EXISTS `user_recharge` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` int unsigned NOT NULL DEFAULT 0 COMMENT '用户UID',
  `amount` int unsigned NOT NULL DEFAULT 0 COMMENT '充值金额',
  `channel` varchar(15) NOT NULL DEFAULT '' COMMENT '充值渠道：alipay或wechatpay',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '充值备注',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '充值时间',
  PRIMARY KEY (`id`),
  INDEX `idx_uid`(`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户充值记录表';

CREATE TABLE IF NOT EXISTS `feed` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(127) NOT NULL DEFAULT '' COMMENT '标题',
  `objid` int unsigned NOT NULL DEFAULT 0 COMMENT '对象id',
  `objtype` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '类型,0-主题;1-文章;2-资源;3-wiki;4-project',
  `uid` int unsigned NOT NULL DEFAULT 0 COMMENT '发布人UID',
  `author` varchar(31) NOT NULL DEFAULT '' COMMENT '外站作者',
  `nid` int unsigned NOT NULL DEFAULT 0 COMMENT '主题的nid或资源的catid',
  `lastreplyuid` int unsigned NOT NULL DEFAULT 0 COMMENT '最后回复者',
  `lastreplytime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '最后回复时间',
  `tags` varchar(63) NOT NULL DEFAULT '' COMMENT 'tag，逗号分隔',
  `cmtnum` int unsigned NOT NULL DEFAULT 0 COMMENT '评论数',
  `top` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '置顶，0否，1置顶',
  `state` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '状态：0-正常；1-下线',
  `created_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_objid_type` (`objid`, `objtype`),
  INDEX `idx_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='网站关键资源动态表';

CREATE TABLE `view_record` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `objid` int(10) unsigned NOT NULL COMMENT '对象id，属主',
  `objtype` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '类型,0-帖子;1-博客;2-资源;3-wiki;4-项目;5-图书',
  `uid` int unsigned NOT NULL DEFAULT 0 COMMENT '浏览人UID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_obj_uid` (`objid`,`objtype`,`uid`),
  INDEX `idx_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户浏览记录表';

CREATE TABLE `view_source` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `objid` int(10) unsigned NOT NULL COMMENT '对象id，属主',
  `objtype` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '类型,0-帖子;1-博客;2-资源;3-wiki;4-项目;5-图书',
  `google` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '来源谷歌数量',
  `baidu` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '来源百度数量',
  `bing` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '来源必应数量',
  `sogou` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '来源搜狗数量',
  `so` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '来源360数量',
  `other` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '其他来源数量',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_obj` (`objid`,`objtype`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='浏览来源表';

CREATE TABLE `gift` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(63) NOT NULL DEFAULT '' COMMENT '物品名称',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '详细描述',
  `price` int unsigned NOT NULL DEFAULT 0 COMMENT '价格（铜币数）',
  `total_num` int unsigned NOT NULL DEFAULT 0 COMMENT '总数量',
  `remain_num` int unsigned NOT NULL DEFAULT 0 COMMENT '剩余数量',
  `expire_time` int unsigned NOT NULL DEFAULT 0 COMMENT '有效期',
  `supplier` varchar(31) NOT NULL DEFAULT '' COMMENT '合作供应商',
  `buy_limit` int unsigned NOT NULL DEFAULT 0 COMMENT '兑换数量限制',
  `typ` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '类型：0-兑换码；1-折扣',
  `state` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '状态,0-未上线;1-已上线;2-已下线;3-过期',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '物品表';

CREATE TABLE `gift_redeem` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `gift_id` int unsigned NOT NULL DEFAULT 0 COMMENT '物品ID',
  `code` varchar(15) NOT NULL DEFAULT '' COMMENT '兑换码',
  `exchange` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否已兑换：0-否；1-是',
  `uid` int unsigned NOT NULL DEFAULT 0 COMMENT '兑换者UID',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '物品兑换码';

CREATE TABLE `user_exchange_record` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `gift_id` int unsigned NOT NULL DEFAULT 0 COMMENT '物品ID',
  `uid` int unsigned NOT NULL DEFAULT 0 COMMENT '兑换者UID',
  `remark` varchar(63) NOT NULL DEFAULT '' COMMENT '物品说明',
  `expire_time` int unsigned NOT NULL DEFAULT 0 COMMENT '过期时间',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_gid` (`gift_id`),
  INDEX `idx_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '物品用户兑换记录';

CREATE TABLE IF NOT EXISTS `gctt_user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(31) NOT NULL DEFAULT '' COMMENT 'Github 用户名',
  `avatar` varchar(127) NOT NULL DEFAULT '' COMMENT 'github 头像',
  `uid` int unsigned NOT NULL DEFAULT 0 COMMENT '本站 uid',
  `joined_at` int unsigned NOT NULL DEFAULT 0 COMMENT '加入GCTT时间，第一个 pr 时间',
  `last_at` int unsigned NOT NULL DEFAULT 0 COMMENT '最后一个 pr 时间',
  `num` int unsigned NOT NULl DEFAULT 0 COMMENT '翻译的文章数',
  `words` int unsigned NOT NULl DEFAULT 0 COMMENT '翻译的字数',
  `avg_time` int unsigned NOT NULl DEFAULT 0 COMMENT '平均每篇用时（秒）',
  `role` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '角色，如 组长、选题、校对等。0-译者；1-组长；2-选题；3-校对；4-核心成员',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`username`),
  INDEX idx_uid (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'GCTT 用户表';

CREATE TABLE IF NOT EXISTS `gctt_git` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(31) NOT NULL DEFAULT '' COMMENT 'Github 用户名',
  `pr` int unsigned NOT NULL DEFAULT 0 COMMENT '完成翻译时的 PR 编号',
  `title` varchar(127) NOT NULL DEFAULT '' COMMENT 'github 上文章名（也是文件名）',
  `md5` char(32) NOT NULL DEFAULT '' COMMENT '标题 md5',
  `translating_at` int unsigned NOT NULL DEFAULT 0 COMMENT '开始翻译时间',
  `translated_at` int unsigned NOT NULL DEFAULT 0 COMMENT '完成翻译时间',
  `words` int unsigned NOT NULL DEFAULT 0 COMMENT '字数',
  `article_id` int unsigned NOT NULL DEFAULT 0 COMMENT '本站 article id',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`md5`),
  INDEX (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'GCTT github 文章翻译信息表';


CREATE TABLE IF NOT EXISTS `gctt_timeline` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `content` varchar(1022) NOT NULL DEFAULT '' COMMENT '内容',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'GCTT 大事记';


CREATE TABLE IF NOT EXISTS `gctt_issue` (
  `id` int unsigned NOT NULL DEFAULT 0 COMMENT '选题的 issue 编号',
  `translator` varchar(31) NOT NULL DEFAULT '' COMMENT '译者 Github 用户名',
  `email` varchar(63) NOT NULL DEFAULT '' COMMENT '译者邮箱',
  `title` varchar(127) NOT NULL DEFAULT '' COMMENT 'issue 标题',
  `translating_at` int unsigned NOT NULL DEFAULT 0 COMMENT '开始翻译时间（认领时间）',
  `translated_at` int unsigned NOT NULL DEFAULT 0 COMMENT '完成翻译时间（close 时间）',
  `label` varchar(31) NOT NULL DEFAULT '' COMMENT '标签，如：已认领，待认领',
  `state` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '0-opened；1-closed',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX (`label`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'GCTT github 选题 issue 列表';


CREATE TABLE IF NOT EXISTS `subject` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '专栏ID',
  `name` varchar(31) NOT NULL DEFAULT '' COMMENT '专栏名',
  `cover` varchar(127) NOT NULL DEFAULT '' COMMENT '专栏封面',
  `description` varchar(1023) NOT NULL DEFAULT '' COMMENT '专栏描述（公告）',
  `uid` int unsigned NOT NULL DEFAULT 0 COMMENT '创建者UID',
  `contribute` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '是否允许投稿, 0-不允许；1-允许',
  `audit` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '投稿是否需要审核, 0-不需要；1-需要',
  `article_num` int unsigned NOT NULL DEFAULT 0 COMMENT '收录的文章数',
  `created_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '专栏';


CREATE TABLE IF NOT EXISTS `subject_admin` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `sid` int unsigned NOT NULL DEFAULT 0 COMMENT '专栏ID',
  `uid` int unsigned NOT NULL DEFAULT 0 COMMENT '管理员UID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`sid`,`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '专栏管理员（不包括创建者）';


CREATE TABLE IF NOT EXISTS `subject_article` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `sid` int unsigned NOT NULL DEFAULT 0 COMMENT '专栏ID',
  `article_id` int unsigned NOT NULL DEFAULT 0 COMMENT '文章ID',
  `state` tinyint unsigned NOT NULl DEFAULT 0 COMMENT '状态：0-新投稿（待审核）；1-上线；2-下线（删除）',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`sid`,`article_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '专栏文章列表';


CREATE TABLE IF NOT EXISTS `subject_follower` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `sid` int unsigned NOT NULL DEFAULT 0 COMMENT '专栏ID',
  `uid` int unsigned NOT NULL DEFAULT 0 COMMENT '关注者UID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`sid`,`uid`),
  INDEX (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '专栏关注者';

CREATE TABLE IF NOT EXISTS `download` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '自增',
  `version` varchar(31) NOT NULL DEFAULT '' COMMENT '版本号',
  `filename` varchar(63) NOT NULL DEFAULT '' COMMENT '文件名',
  `kind` varchar(31) NOT NULL DEFAULT '' COMMENT '类型',
  `os` varchar(31) NOT NULL DEFAULT '' COMMENT '操作系统',
  `arch` varchar(31) NOT NULL DEFAULT '' COMMENT '架构',
  `size` int unsigned NOT NULL DEFAULT 0 COMMENT '大小，单位 MB',
  `checksum` varchar(64) NOT NULL DEFAULT '' COMMENT 'SHA1/256 校验和',
  `is_recommend` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否推荐（推荐的高亮显示）',
  `category` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '0-Archived versions; 1-Stable versions; 2-Unstable versions;',
  `seq` int unsigned NOT NULL DEFAULT 0 COMMENT '排序，越大越靠前',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '下载信息表';

CREATE TABLE `wechat_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `openid` varchar(127) NOT NULL DEFAULT '' COMMENT '用户的标识，对当前公众号/小程序唯一',
  `nickname` varchar(127) NOT NULL DEFAULT '' COMMENT '用户的昵称',
  `session_key` varchar(127) NOT NULL DEFAULT '' COMMENT '小程序返回的 session_key',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '用户微信头像',
  `open_info` varchar(1024) NOT NULL DEFAULT '' COMMENT '用户微信的其他信息，json格式',
  `uid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '用户UID',
  `created_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `openid` (`openid`),
  KEY `uid` (`uid`),
  KEY `updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='微信用户绑定表';
