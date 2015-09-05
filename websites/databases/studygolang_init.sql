-- 角色表
INSERT INTO `role`(`name`) VALUES('站长');
INSERT INTO `role`(`name`) VALUES('副站长');
INSERT INTO `role`(`name`) VALUES('超级管理员');
INSERT INTO `role`(`name`) VALUES('社区管理员');
INSERT INTO `role`(`name`) VALUES('资源管理员');
INSERT INTO `role`(`name`) VALUES('文章管理员');
INSERT INTO `role`(`name`) VALUES('晨读管理员');
INSERT INTO `role`(`name`) VALUES('高级会员');
INSERT INTO `role`(`name`) VALUES('中级会员');
INSERT INTO `role`(`name`) VALUES('初级会员');

INSERT INTO `authority` (`aid`, `name`, `menu1`, `menu2`, `route`, `op_user`, `ctime`, `mtime`)
VALUES
	(1, '用户管理', 0, 0, '', 'polaris', '2014-08-17 11:06:50', '2014-08-16 10:51:13'),
	(2, '权限管理', 1, 0, '/admin/user/auth/list', 'polaris', '2014-08-17 11:06:50', '2014-08-16 12:05:55'),
	(3, '权限查询', 1, 2, '/admin/user/auth/query.html', 'polaris', '2014-08-17 11:06:50', '2014-08-17 11:06:50'),
	(4, '新建权限', 1, 2, '/admin/user/auth/new', 'polaris', '2014-08-17 11:06:50', '2014-08-17 21:18:56'),
	(5, '修改权限', 1, 2, '/admin/user/auth/modify', 'polaris', '2014-08-18 23:56:11', '2014-10-04 23:16:31'),
	(6, '删除权限', 1, 2, '/admin/user/auth/del', 'polaris', '2014-08-17 22:12:30', '2014-08-17 22:12:30'),
	(7, '角色管理', 1, 0, '/admin/user/role/list', 'polaris', '2014-08-17 23:03:40', '2014-08-17 23:03:40'),
	(8, '角色查询', 1, 7, '/admin/user/role/query.html', 'polaris', '2014-08-18 23:54:44', '2014-08-18 23:54:44'),
	(9, '新建角色', 1, 7, '/admin/user/role/new', 'polaris', '2014-08-18 23:55:10', '2014-08-18 23:55:10'),
	(10, '修改角色', 1, 7, '/admin/user/role/modify', 'polaris', '2014-08-18 23:55:24', '2014-08-18 23:55:24'),
	(11, '删除角色', 1, 7, '/admin/user/role/del', 'polaris', '2014-08-18 23:57:20', '2014-08-18 23:57:44'),
	(12, '用户管理', 1, 0, '/admin/user/user/list', 'polaris', '2014-08-19 00:20:41', '2014-08-19 00:20:41'),
	(13, '用户查询', 1, 12, '/admin/user/user/query.html', 'polaris', '2014-08-19 08:43:55', '2014-08-19 08:43:55'),
	(14, '用户详情', 1, 12, '/admin/user/user/detail', 'polaris', '2014-08-19 21:14:07', '2014-08-19 21:14:07');

INSERT INTO `role_authority` (`roleid`, `aid`, `op_user`, `ctime`)
    VALUES
    	(0, 19, 'polaris', '2015-07-14 13:49:03'),
    	(0, 20, 'polaris', '2015-07-14 13:49:03'),
    	(0, 21, 'polaris', '2015-07-14 13:49:03'),
    	(0, 22, 'polaris', '2015-07-14 13:49:03'),
    	(0, 23, 'polaris', '2015-07-14 13:49:03'),
    	(0, 24, 'polaris', '2015-07-14 13:49:03'),
    	(0, 25, 'polaris', '2015-07-14 13:49:03'),
    	(0, 26, 'polaris', '2015-07-14 13:49:03'),
    	(0, 27, 'polaris', '2015-07-14 13:49:03'),
    	(0, 28, 'polaris', '2014-10-30 22:56:21'),
    	(0, 29, 'polaris', '2014-10-30 22:56:21'),
    	(0, 30, 'polaris', '2014-10-30 22:56:21'),
    	(0, 31, 'polaris', '2014-10-30 22:56:21'),
    	(0, 37, 'polaris', '2015-07-14 13:49:03'),
    	(1, 1, 'polaris', '2015-07-14 13:54:30'),
    	(1, 2, 'polaris', '2015-07-14 13:54:30'),
    	(1, 3, 'polaris', '2015-07-14 13:54:30'),
    	(1, 4, 'polaris', '2015-07-14 13:54:30'),
    	(1, 5, 'polaris', '2015-07-14 13:54:30'),
    	(1, 6, 'polaris', '2015-07-14 13:54:30'),
    	(1, 7, 'polaris', '2015-07-14 13:54:30'),
    	(1, 8, 'polaris', '2015-07-14 13:54:30'),
    	(1, 9, 'polaris', '2015-07-14 13:54:30'),
    	(1, 10, 'polaris', '2015-07-14 13:54:30'),
    	(1, 11, 'polaris', '2015-07-14 13:54:30'),
    	(1, 12, 'polaris', '2015-07-14 13:54:30'),
    	(1, 13, 'polaris', '2015-07-14 13:54:30'),
    	(1, 14, 'polaris', '2015-07-14 13:54:30'),
    	(1, 15, 'polaris', '2015-07-14 13:54:30'),
    	(1, 16, 'polaris', '2015-07-14 13:54:30'),
    	(1, 17, 'polaris', '2015-07-14 13:54:30'),
    	(1, 18, 'polaris', '2015-07-14 13:54:30'),
    	(1, 19, 'polaris', '2015-07-14 13:54:30'),
    	(1, 20, 'polaris', '2015-07-14 13:54:30'),
    	(1, 21, 'polaris', '2015-07-14 13:54:30'),
    	(1, 22, 'polaris', '2015-07-14 13:54:30'),
    	(1, 23, 'polaris', '2015-07-14 13:54:30'),
    	(1, 24, 'polaris', '2015-07-14 13:54:30'),
    	(1, 25, 'polaris', '2015-07-14 13:54:30'),
    	(1, 26, 'polaris', '2015-07-14 13:54:30'),
    	(1, 27, 'polaris', '2015-07-14 13:54:30'),
    	(1, 28, 'polaris', '2015-07-14 13:54:30'),
    	(1, 29, 'polaris', '2015-07-14 13:54:30'),
    	(1, 30, 'polaris', '2015-07-14 13:54:30'),
    	(1, 31, 'polaris', '2015-07-14 13:54:30'),
    	(1, 32, 'polaris', '2015-07-14 13:54:30'),
    	(1, 33, 'polaris', '2015-07-14 13:54:30'),
    	(1, 34, 'polaris', '2015-07-14 13:54:30'),
    	(1, 35, 'polaris', '2015-07-14 13:54:30'),
    	(1, 36, 'polaris', '2015-07-14 13:54:30'),
    	(1, 37, 'polaris', '2015-07-14 13:54:30'),
    	(1, 38, 'polaris', '2015-07-14 13:54:30'),
    	(6, 19, 'polaris', '2015-07-14 13:51:31'),
    	(6, 20, 'polaris', '2015-07-14 13:51:31'),
    	(6, 21, 'polaris', '2015-07-14 13:51:31'),
    	(6, 22, 'polaris', '2015-07-14 13:51:31'),
    	(6, 23, 'polaris', '2015-07-14 13:51:31'),
    	(6, 24, 'polaris', '2015-07-14 13:51:31'),
    	(6, 25, 'polaris', '2015-07-14 13:51:31'),
    	(6, 26, 'polaris', '2015-07-14 13:51:31'),
    	(6, 27, 'polaris', '2015-07-14 13:51:31'),
    	(6, 37, 'polaris', '2015-07-14 13:51:31'),
    	(7, 28, 'polaris', '2014-10-31 10:12:19'),
    	(7, 29, 'polaris', '2014-10-31 10:12:19'),
    	(7, 30, 'polaris', '2014-10-31 10:12:19'),
    	(7, 31, 'polaris', '2014-10-31 10:12:19');

INSERT INTO `user_role` (`uid`, `roleid`, `ctime`)
        VALUES
        	(1, 1, '2013-03-15 14:38:09');


-- 帖子节点表
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(0, 'Golang', 'Go语言基本问题探讨');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(1, 'Go基础', 'Go语言基础、语法、规范等');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(1, 'Go标准库', 'Go语言标准库使用、例子、源码等');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(1, 'Go源码', '深入Go语言内部实现，分享Go语言源码学习心得');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(1, 'Go Web开发', '用Go语言进行Web开发');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(1, 'Go问与答', '任何关于Go语言的问题都可以到这里提');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(1, 'Go动态', 'Go语言相关资讯和最新动态');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(1, 'Go开发工具', '交流Go开发工具的使用');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(1, 'Go实战', 'Go语言实际使用经验交流');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(1, 'Go工具链', '(go tool chain)Go提供的各种工具学习、使用');

INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(0, '开源控', 'Golang本身开源，自然gopher们都是开源控');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(9, 'Go Web框架', '开源的Go Web框架，你知道多少？使用了吗？');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(9, 'Go第三方库', '第三方Go库，你都用了哪些？');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(9, 'Go代码分享', '和gopher们分享您自己的Go代码吧');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(9, 'Go资料', '好多Go语言资料啊，学习学习……');

INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(0, 'Study Golang', 'Golang China，Go语言学习园地，中文社区');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(14, '社区公告', '社区最新动态、公共以及其他信息');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(14, '反馈', '使用过程中遇到了问题，可以在这里提交');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(14, '社区开发', '对开发该社区感兴趣的可以一起加入进来哦');

INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(0, '分享', '分享生活、学习、工作等方方面面');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(18, 'Markdown', '当下Markdown是相当火，本站就使用Markdown发帖，有必要聊聊它的使用');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(18, '招聘', '发布Go语言招聘信息');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(18, '杂谈', 'Go相关或不太相关的杂谈');


INSERT INTO `resource_category`(`name`, `intro`) VALUES('精彩文章', '分享来自互联网关于Go语言的精彩文章');
INSERT INTO `resource_category`(`name`, `intro`) VALUES('开源项目', '收集优秀的开源项目、第三方库');
INSERT INTO `resource_category`(`name`, `intro`) VALUES('Go语言资料', 'Go语言书籍、资料，提供下载地址或在线链接');
INSERT INTO `resource_category`(`name`, `intro`) VALUES('其他资源', '分享跟Go相关或其他有用的资源');
