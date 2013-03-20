-- 角色表
INSERT INTO `role`(`name`) VALUES('站长');
INSERT INTO `role`(`name`) VALUES('副站长');
INSERT INTO `role`(`name`) VALUES('超级管理员');
INSERT INTO `role`(`name`) VALUES('社区管理员');
INSERT INTO `role`(`name`) VALUES('资源管理员');
INSERT INTO `role`(`name`) VALUES('酷站管理员');
INSERT INTO `role`(`name`) VALUES('高级会员');
INSERT INTO `role`(`name`) VALUES('中级会员');
INSERT INTO `role`(`name`) VALUES('初级会员');

-- 帖子节点表
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(0, 'Golang', 'Go语言基本问题探讨');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(1, 'Go基础', 'Go语言基础、语法、规范等');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(1, 'Go标准库', 'Go语言标准库使用、例子、源码等');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(1, 'Go源码', '深入Go语言内部实现，分享Go语言源码学习心得');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(1, 'Go Web开发', '用Go语言进行Web开发');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(1, 'Go问与答', '任何关于Go语言的问题都可以到这里提');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(1, 'Go动态', 'Go语言相关资讯和最新动态');
INSERT INTO `topics_node`(`parent`, `name`, `intro`) VALUES(1, 'Go开发工具', '交流Go开发工具的使用');

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