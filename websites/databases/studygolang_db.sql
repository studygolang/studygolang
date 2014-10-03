/*---------------------------------------------------------------------------*
  NAME: topics
  用途：帖子内容表
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `topics`;
CREATE TABLE `topics` (
  `tid` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `content` text NOT NULL,
  `nid` int unsigned NOT NULL COMMENT '节点id',
  `uid` int unsigned NOT NULL COMMENT '帖子作者',
  `lastreplyuid` int unsigned NOT NULL DEFAULT 0 COMMENT '最后回复者',
  `lastreplytime` timestamp NOT NULL DEFAULT 0 COMMENT '最后回复时间',
  `flag` tinyint NOT NULL DEFAULT 0 COMMENT '审核标识,0-未审核;1-已审核;2-审核删除;3-用户自己删除',
  `editor_uid` int unsigned NOT NULL DEFAULT 0 COMMENT '最后编辑人',
  `ctime` timestamp NOT NULL DEFAULT 0,
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`tid`),
  KEY `uid` (`uid`),
  KEY `nid` (`nid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

alter table `studygolang`.`topics` add column `editor_uid` int UNSIGNED DEFAULT '0' NOT NULL COMMENT '最后编辑人' after `lastreplytime`

/*---------------------------------------------------------------------------*
  NAME: topics_ex
  用途：帖子扩展表（计数）
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `topics_ex`;
CREATE TABLE `topics_ex` (
  `tid` int unsigned NOT NULL,
  `view` int unsigned NOT NULL DEFAULT 0 COMMENT '浏览数',
  `reply` int unsigned NOT NULL DEFAULT 0 COMMENT '回复数',
  `like` int unsigned NOT NULL DEFAULT 0 COMMENT '喜欢数',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`tid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*---------------------------------------------------------------------------*
  NAME: topics_node
  用途：帖子节点表
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `topics_node`;
CREATE TABLE `topics_node` (
  `nid` int unsigned NOT NULL AUTO_INCREMENT,
  `parent` int unsigned NOT NULL DEFAULT 0 COMMENT '父节点id，无父节点为0',
  `name` varchar(20) NOT NULL COMMENT '节点名',
  `intro` varchar(50) NOT NULL DEFAULT '' COMMENT '节点简介',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`nid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


/*---------------------------------------------------------------------------*
  NAME: comments
  用途：评论表（帖子回复、博客文章评论等，统一处理）
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments` (
  `cid` int unsigned NOT NULL AUTO_INCREMENT,
  `objid` int unsigned NOT NULL COMMENT '对象id，属主（评论给谁）',
  `objtype` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '类型,0-帖子;1-博客;2-资源;3-酷站',
  `content` text NOT NULL,
  `uid` int unsigned NOT NULL COMMENT '回复者',
  `floor` int unsigned NOT NULL COMMENT '第几楼',
  `flag` tinyint NOT NULL DEFAULT 0 COMMENT '审核标识,0-未审核;1-已审核;2-审核删除;3-用户自己删除',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`cid`),
  UNIQUE KEY (`objid`,`objtype`,`floor`),
  KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*---------------------------------------------------------------------------*
  NAME: likes
  用途：喜欢表（帖子回复、博客文章评论等，统一处理）
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `likes`;
CREATE TABLE `likes` (
  `cid` int unsigned NOT NULL AUTO_INCREMENT,
  `objid` int unsigned NOT NULL COMMENT '对象id，属主（评论给谁）',
  `objtype` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '类型,0-帖子;1-博客;2-资源;3-酷站',
  `content` text NOT NULL,
  `uid` int unsigned NOT NULL COMMENT '回复者',
  `floor` int unsigned NOT NULL COMMENT '第几楼',
  `flag` tinyint NOT NULL DEFAULT 0 COMMENT '审核标识,0-未审核;1-已审核;2-审核删除;3-用户自己删除',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`cid`),
  UNIQUE KEY (`objid`,`objtype`,`floor`),
  KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


/*---------------------------------------------------------------------------*
  NAME: views
  用途：帖子用户最后阅读表（帖子回复、博客文章评论等，统一处理）
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `views`;
CREATE TABLE `views` (
  `cid` int unsigned NOT NULL AUTO_INCREMENT,
  `objid` int unsigned NOT NULL COMMENT '对象id，属主（评论给谁）',
  `objtype` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '类型,0-帖子;1-博客;2-资源;3-酷站',
  `content` text NOT NULL,
  `uid` int unsigned NOT NULL COMMENT '回复者',
  `floor` int unsigned NOT NULL COMMENT '第几楼',
  `flag` tinyint NOT NULL DEFAULT 0 COMMENT '审核标识,0-未审核;1-已审核;2-审核删除;3-用户自己删除',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`cid`),
  UNIQUE KEY (`objid`,`objtype`,`floor`),
  KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


/*---------------------------------------------------------------------------*
  NAME: user_login
  用途：用户登录表
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `user_login`;
CREATE TABLE `user_login` (
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
   
/*---------------------------------------------------------------------------*
  NAME: bind_user
  用途：第三方绑定表
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `bind_user`;
CREATE TABLE `bind_user` (
  `uid` int unsigned NOT NULL,
  `type` tinyint NOT NULL DEFAULT 0 COMMENT '绑定的第三方类型',
  `email` varchar(128) NOT NULL DEFAULT '',
  `tuid` int unsigned NOT NULL DEFAULT 0 COMMENT '第三方uid',
  `username` varchar(20) NOT NULL COMMENT '用户名',
  `token` varchar(50) NOT NULL COMMENT '第三方access_token',
  `refresh` varchar(50) NOT NULL COMMENT '第三方refresh_token',
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*---------------------------------------------------------------------------*
  NAME: user_info
  用途：用户信息表
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user_info` (
  `uid` int unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(128) NOT NULL DEFAULT '',
  `open` tinyint NOT NULL DEFAULT 1 COMMENT '邮箱是否公开，默认公开',
  `username` varchar(20) NOT NULL COMMENT '用户名',
  `name` varchar(20) NOT NULL DEFAULT '' COMMENT '姓名',
  `avatar` varchar(128) NOT NULL DEFAULT '' COMMENT '头像(暂时使用http://www.gravatar.com)',
  `city` varchar(10) NOT NULL DEFAULT '居住地',
  `company` varchar(64) NOT NULL DEFAULT '',
  `github` varchar(20) NOT NULL DEFAULT '',
  `weibo` varchar(20) NOT NULL DEFAULT '',
  `website` varchar(50) NOT NULL DEFAULT '' COMMENT '个人主页，博客',
  `monlog` varchar(140) NOT NULL DEFAULT '' COMMENT '个人状态，签名，独白',
  `introduce` text NOT NULL COMMENT '个人简介',
  `status` tinyint unsigned NOT NULL DEFAULT '' COMMENT '用户账号状态。0-默认；1-已审核；2-拒绝；3-冻结；4-停号',
  `ctime` timestamp NOT NULL DEFAULT 0,
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`uid`),
  UNIQUE KEY (`username`),
  UNIQUE KEY (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*---------------------------------------------------------------------------*
  NAME: user_active
  用途：活跃用户表
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `user_active`;
CREATE TABLE `user_active` (
  `uid` int unsigned NOT NULL,
  `email` varchar(128) NOT NULL,
  `username` varchar(20) NOT NULL COMMENT '用户名',
  `weight` smallint NOT NULL DEFAULT 1 COMMENT '活跃度，越大越活跃',
  `avatar` varchar(128) NOT NULL DEFAULT '' COMMENT '头像(暂时使用http://www.gravatar.com)',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`uid`),
  UNIQUE KEY (`username`),
  UNIQUE KEY (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*---------------------------------------------------------------------------*
  NAME: role
  用途：角色表，常驻内存
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
  `roleid` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '角色名',
  `op_user` varchar(20) NOT NULL DEFAULT '' COMMENT '操作人',
  `ctime` timestamp NOT NULL DEFAULT 0,
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`roleid`),
  UNIQUE KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*---------------------------------------------------------------------------*
  NAME: authority
  用途：权限表，常驻内存
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `authority`;
CREATE TABLE `authority` (
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*---------------------------------------------------------------------------*
  NAME: role_authority
  用途：角色拥有的权限表
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `role_authority`;
CREATE TABLE `role_authority` (
  `roleid` int unsigned NOT NULL,
  `aid` int unsigned NOT NULL,
  `op_user` varchar(20) NOT NULL COMMENT '操作人',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`roleid`, `aid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*---------------------------------------------------------------------------*
  NAME: user_role
  用途：用户角色表（用户是什么角色，可以多个角色）
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `user_role`;
CREATE TABLE `user_role` (
  `uid` int unsigned NOT NULL,
  `roleid` int unsigned NOT NULL,
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`uid`, `roleid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


/*---------------------------------------------------------------------------*
  NAME: message
  用途：短消息（私信）
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message` (
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
) COMMENT 'message 短消息（私信）' ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*---------------------------------------------------------------------------*
  NAME: system_message
  用途：系统消息表
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `system_message`;
CREATE TABLE `system_message` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `msgtype` tinyint NOT NULL DEFAULT 0 COMMENT '系统消息类型',
  `hasread` enum('未读','已读') NOT NULL DEFAULT '未读',
  `to` int unsigned NOT NULL COMMENT '发给谁',
  `ext` text NOT NULL COMMENT '额外信息',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY (`to`)
) COMMENT 'system_message 系统消息表' ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*---------------------------------------------------------------------------*
  NAME: wiki
  用途：wiki页（需要考虑审核问题？）
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `wiki`;
CREATE TABLE `wiki` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL COMMENT 'wiki标题',
  `content` longtext NOT NULL COMMENT 'wiki内容',
  `uri` varchar(50) NOT NULL COMMENT 'uri',
  `uid` int unsigned NOT NULL COMMENT '作者',
  `cuid` varchar(100) NOT NULL DEFAULT '' COMMENT '贡献者',
  `ctime` timestamp NOT NULL DEFAULT 0,
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`uri`)
) COMMENT 'wiki页' ENGINE=InnoDB DEFAULT CHARSET=utf8;


/*---------------------------------------------------------------------------*
  NAME: resource
  用途：资源表：包括Golang资源下载
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `resource`;
CREATE TABLE `resource` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL COMMENT '资源标题',
  `form` enum('只是链接','包括内容'),
  `content` longtext NOT NULL COMMENT '资源内容',
  `url` varchar(150) NOT NULL COMMENT '链接url',
  `uid` int unsigned NOT NULL COMMENT '作者',
  `catid` int unsigned NOT NULL COMMENT '所属类别',
  `ctime` timestamp NOT NULL DEFAULT 0,
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) COMMENT '资源' ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*---------------------------------------------------------------------------*
  NAME: resource_ex
  用途：资源扩展表（计数）
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `resource_ex`;
CREATE TABLE `resource_ex` (
  `id` int unsigned NOT NULL,
  `viewnum` int unsigned NOT NULL DEFAULT 0 COMMENT '浏览数',
  `cmtnum` int unsigned NOT NULL DEFAULT 0 COMMENT '回复数',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) COMMENT '资源扩展表' ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*---------------------------------------------------------------------------*
  NAME: resource_category
  用途：资源分类表
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `resource_category`;
CREATE TABLE `resource_category` (
  `catid` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(20) NOT NULL COMMENT '分类名',
  `intro` varchar(50) NOT NULL DEFAULT '' COMMENT '分类简介',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`catid`)
) COMMENT '资源分类表' ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*---------------------------------------------------------------------------*
  NAME: articles
  用途：网络文章聚合表
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `articles`;
CREATE TABLE `articles` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `domain` varchar(50) NOT NULL DEFAULT '' COMMENT '来源域名（不一定是顶级域名）',
  `name` varchar(30) NOT NULL DEFAULT '' COMMENT '来源名称',
  `title` varchar(127) NOT NULL DEFAULT '' COMMENT '文章标题',
  `cover` varchar(127) NOT NULL DEFAULT '' COMMENT '图片封面',
  `author` varchar(255) NOT NULL DEFAULT '' COMMENT '文章作者(可能带html)',
  `author_txt` varchar(30) NOT NULL DEFAULT '' COMMENT '文章作者(纯文本)',
  `lang` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '语言：0-中文；1-英文',
  `pub_date` varchar(20) NOT NULL DEFAULT '' COMMENT '发布时间',
  `url` varchar(127) NOT NULL DEFAULT '' COMMENT '文章原始链接',
  `content` text NOT NULL COMMENT '正文(带html)',
  `txt` text NOT NULL COMMENT '正文(纯文本)',
  `tags` varchar(50) NOT NULL DEFAULT '' COMMENT '文章tag，逗号分隔',
  `status` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '状态：0-初始抓取；1-已上线；2-下线(审核拒绝)',
  `op_user` varchar(20) NOT NULL DEFAULT '' COMMENT '操作人',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`url`),
  KEY (`domain`),
  KEY (`ctime`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '网络文章聚合表';

/*---------------------------------------------------------------------------*
  NAME: crawl_rule
  用途：网站抓取规则表
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `crawl_rule`;
CREATE TABLE `crawl_rule` (
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
  `op_user` varchar(20) NOT NULL DEFAULT '' COMMENT '操作人',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`domain`,`subpath`),
  KEY (`ctime`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '网站抓取规则表';

/*---------------------------------------------------------------------------*
  NAME: 动态表（go动态；本站动态等）
  用途：动态表
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `dynamic`;
CREATE TABLE `dynamic` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '动态内容',
  `dmtype` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '类型：0-Go动态；1-本站动态',
  `url` varchar(255) NOT NULL DEFAULT '' COMMENT '链接',
  `seq` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '顺序（越大越在前）',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY (`seq`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '动态表';