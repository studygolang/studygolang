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
	(1, '用户管理', 0, 0, '', '', '2014-08-17 11:06:50', '2014-08-16 10:51:13'),
	(2, '权限管理', 1, 0, '/admin/user/auth/list', '', '2014-08-17 11:06:50', '2014-08-16 12:05:55'),
	(3, '权限查询', 1, 2, '/admin/user/auth/query.html', '', '2014-08-17 11:06:50', '2014-08-17 11:06:50'),
	(4, '新建权限', 1, 2, '/admin/user/auth/new', '', '2014-08-17 11:06:50', '2014-08-17 21:18:56'),
	(5, '修改权限', 1, 2, '/admin/user/auth/modify', '', '2014-08-18 23:56:11', '2014-10-04 23:16:31'),
	(6, '删除权限', 1, 2, '/admin/user/auth/del', '', '2014-08-17 22:12:30', '2014-08-17 22:12:30'),
	(7, '角色管理', 1, 0, '/admin/user/role/list', '', '2014-08-17 23:03:40', '2014-08-17 23:03:40'),
	(8, '角色查询', 1, 7, '/admin/user/role/query.html', '', '2014-08-18 23:54:44', '2014-08-18 23:54:44'),
	(9, '新建角色', 1, 7, '/admin/user/role/new', '', '2014-08-18 23:55:10', '2014-08-18 23:55:10'),
	(10, '修改角色', 1, 7, '/admin/user/role/modify', '', '2014-08-18 23:55:24', '2014-08-18 23:55:24'),
	(11, '删除角色', 1, 7, '/admin/user/role/del', '', '2014-08-18 23:57:20', '2014-08-18 23:57:44'),
	(12, '用户管理', 1, 0, '/admin/user/user/list', '', '2014-08-19 00:20:41', '2014-08-19 00:20:41'),
	(13, '用户查询', 1, 12, '/admin/user/user/query.html', '', '2014-08-19 08:43:55', '2014-08-19 08:43:55'),
	(14, '用户详情', 1, 12, '/admin/user/user/detail', '', '2014-08-19 21:14:07', '2014-08-19 21:14:07'),
	(15, '社区管理', 0, 0, '', '', '2014-09-03 21:18:10', '2014-09-03 21:18:10'),
	(16, '帖子管理', 15, 0, '/admin/community/topic/list', '', '2014-09-03 22:23:08', '2014-09-03 22:23:08'),
	(17, '帖子查询', 15, 16, '/admin/community/topic/query.html', '', '2014-09-04 07:28:33', '2014-09-04 07:28:33'),
	(18, '帖子修改', 15, 16, '/admin/community/topic/modify', '', '2014-09-04 22:42:47', '2014-09-04 22:42:47'),
	(19, '抓取管理', 0, 0, '', '', '2014-09-08 11:24:08', '2014-09-08 11:24:08'),
	(20, '文章管理', 19, 0, '/admin/crawl/article/list', '', '2014-09-08 11:24:39', '2014-09-08 11:24:39'),
	(21, '规则管理', 19, 0, '/admin/crawl/rule/list', '', '2014-09-08 11:50:19', '2014-09-08 11:50:19'),
	(22, '文章修改', 19, 20, '/admin/crawl/article/modify', '', '2014-09-30 09:07:44', '2014-09-30 09:07:44'),
	(23, '文章删除', 19, 20, '/admin/crawl/article/del', '', '2014-09-30 09:09:01', '2014-09-30 09:09:01'),
	(24, '新建规则', 19, 21, '/admin/crawl/rule/new', '', '2014-10-04 01:19:42', '2014-10-04 01:19:41'),
	(25, '抓取文章', 19, 20, '/admin/crawl/article/new', '', '2014-10-05 00:02:41', '2014-10-05 00:02:40'),
	(26, '文章查询', 19, 20, '/admin/crawl/article/query.html', '', '2014-10-05 08:08:48', '2014-10-05 08:08:47'),
	(27, '修改规则', 19, 21, '/admin/crawl/rule/modify', '', '2014-10-20 23:14:37', '2014-10-20 23:14:33'),
	(28, '晨读管理', 0, 0, '', '', '2014-10-30 22:53:30', '2014-10-30 22:53:24'),
	(29, '晨读列表', 28, 0, '/admin/reading/list', '', '2014-10-30 22:53:53', '2014-10-30 22:53:47'),
	(30, '晨读查询', 28, 29, '/admin/reading/query.html', '', '2014-10-30 22:54:11', '2014-10-30 22:54:05'),
	(31, '发布晨读', 28, 29, '/admin/reading/publish', '', '2014-10-30 22:54:28', '2014-10-30 22:54:22'),
	(32, '管理工具', 0, 0, '', '', '2014-11-02 11:43:19', '2014-11-02 11:43:13'),
	(33, '生成Sitemap', 32, 0, '/admin/tool/sitemap', '', '2014-11-02 11:43:41', '2014-11-02 11:43:35'),
	(34, '项目管理', 15, 0, '/admin/community/project/list', '', '2014-11-03 21:54:48', '2014-11-03 21:54:42'),
	(35, '项目查询', 15, 34, '/admin/community/project/query.html', '', '2014-11-03 21:55:07', '2014-11-03 21:55:01'),
	(36, '项目上下线', 15, 34, '/admin/community/project/update_status', '', '2014-11-03 21:55:26', '2014-11-03 21:55:19'),
	(37, '规则查询', 19, 21, '/admin/crawl/rule/query.html', '', '2014-11-08 12:38:20', '2014-11-08 12:38:13'),
	(38, '用户修改', 1, 12, '/admin/user/user/modify', '', '2015-07-14 13:53:53', '2015-07-14 13:53:53'),
	(39, '设置', 0, 0, '', '', '2017-05-21 16:03:00', '2017-05-21 16:03:59'),
	(40, '常规', 39, 0, '/admin/setting/genneral/modify', '', '2017-05-21 16:05:00', '2017-05-21 16:05:46'),
	(41, '导航', 39, 0, '/admin/setting/nav/modify', '', '2017-05-21 18:01:00', '2017-05-21 18:01:16'),
	(42, '节点管理', 15, 0, '/admin/community/node/list', 'polaris', '2017-09-01 22:23:08', '2017-09-01 23:10:38'),
	(43, '编辑/新增节点', 15, 42, '/admin/community/node/modify', 'polaris', '2017-09-01 22:23:08', '2017-09-01 23:11:09');


INSERT INTO `website_setting` (`id`, `name`, `domain`, `title_suffix`, `favicon`, `logo`, `start_year`, `blog_url`, `reading_menu`, `docs_menu`, `slogan`, `beian`, `friends_logo`, `footer_nav`, `project_df_logo`, `index_nav`, `created_at`, `updated_at`)
VALUES
	(1, 'Go语言中文网', 'studygolang.com', '- Go语言中文网 - Golang中文社区', '/static/img/go.ico', '/static/img/logo1.png', 2013, 'http://blog.studygolang.com', '', '', 'Go语言中文网，中国 Golang 社区，致力于构建完善的 Golang 中文社区，Go语言爱好者的学习家园。', '京ICP备14030343号-1', '[{\"image\":\"http://qiniutek.com/images/logo-2.png\",\"url\":\"https://portal.qiniu.com/signup?code=3lfz4at7pxfma\",\"name\":\"\",\"width\":\"290px\",\"height\":\"45px\"}]', '[{\"name\":\"关于\",\"url\":\"/wiki/about\",\"outer_site\":false},{\"name\":\"贡献者\",\"url\":\"/wiki/contributors\",\"outer_site\":false},{\"name\":\"帮助推广\",\"url\":\"/wiki\",\"outer_site\":false},{\"name\":\"反馈\",\"url\":\"/topics/node/16\",\"outer_site\":false},{\"name\":\"Github\",\"url\":\"https://github.com/studygolang\",\"outer_site\":true},{\"name\":\"新浪微博\",\"url\":\"http://weibo.com/studygolang\",\"outer_site\":true},{\"name\":\"内嵌Wide\",\"url\":\"/wide/playground\",\"outer_site\":false},{\"name\":\"免责声明\",\"url\":\"/wiki/duty\",\"outer_site\":false}]', '', '[{"tab":"all"}]', '2017-05-21 10:22:00', '2017-05-21 21:30:56');

INSERT INTO `friend_link` (`id`, `name`, `url`, `seq`, `logo`, `created_at`)
VALUES
	(1, 'Go语言中文网', 'http://studygolang.com', 0, '', '2017-05-21 14:52:07');

INSERT INTO `user_setting` (`id`, `key`, `value`, `remark`, `created_at`)
VALUES
	(1, 'new_user_wait', 0, '新用户注册多久能发布帖子，单位秒，0表示没限制', '2017-05-30 18:11:31'),
	(2, 'can_edit_time', 300, '发布后多久内能够编辑，单位秒', '2017-05-30 18:12:53');

INSERT INTO `mission` (`id`, `name`, `type`, `fixed`, `min`, `max`, `incr`, `state`, `created_at`)
VALUES
	(1, '初始资本', 2, 2000, 0, 0, 0, 0, '2017-06-03 22:44:59'),
	(2, '每日登录任务', 1, 0, 25, 50, 5, 0, '2017-06-05 13:35:16');

