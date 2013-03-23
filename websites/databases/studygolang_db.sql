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
  `ctime` timestamp NOT NULL DEFAULT 0,
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`tid`),
  KEY `uid` (`uid`),
  KEY `nid` (`nid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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
  PRIMARY KEY (`uid`),
  UNIQUE KEY (`username`),
  UNIQUE KEY (`email`)
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
  `city` varchar(10) NOT NULL DEFAULT '',
  `company` varchar(64) NOT NULL DEFAULT '',
  `github` varchar(20) NOT NULL DEFAULT '',
  `weibo` varchar(20) NOT NULL DEFAULT '',
  `website` varchar(50) NOT NULL DEFAULT '' COMMENT '个人主页，博客',
  `status` varchar(140) NOT NULL DEFAULT '' COMMENT '个人状态，签名',
  `introduce` text NOT NULL COMMENT '个人简介',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`aid`),
  UNIQUE KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*---------------------------------------------------------------------------*
  NAME: role_authority
  用途：角色拥有的权限表
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `role_authority`;
CREATE TABLE `role_authority` (
  `roleid` int unsigned NOT NULL,
  `aid` int unsigned NOT NULL,
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
  NAME: notification
  用途：通知表
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `notification`;
CREATE TABLE `notification` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `content` text NOT NULL,
  `read` tinyint NOT NULL DEFAULT 0 COMMENT '是否已读',
  `uid` int unsigned NOT NULL,
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*---------------------------------------------------------------------------*
  NAME: message
  用途：短消息（暂不实现）
*---------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `content` text NOT NULL,
  `read` tinyint NOT NULL DEFAULT 0 COMMENT '是否已读',
  `uid` int unsigned NOT NULL,
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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
