studygolang
===========
[Golang中文社区 | Go语言学习园地](http://studygolang.com "Golang中文社区 | Go语言学习园地") 源码

网站上线时间：2013-03-15 14:38:09

目前只开发了主要功能，还有不少功能在不断开发中。欢迎有兴趣的gopher们参与进来，一起构建一个完善的 Golang 中文社区，Go语言爱好者的学习家园。

#目前在正开发或需要开发的功能
1. 热门节点
2. 小贴士
3. 喜欢
4. 收藏
5. 活跃用户（已完成）
6. 关注
7. 用第三方账号登录
8. 绑定github后显示其代码
9. 同步到微博？
10. wiki（已完成，有些细节待完善）
11. 资源
12. 酷站
13. 后台管理

# 本地搭建一个 Golang 社区 #

1、下载 studygolang 代码
	
	git clone https://github.com/studygolang/studygolang

2、下载安装依赖库（如果依赖库下载不下来可以联系我）

	cd studygolang/websites/code/thirdparty
	// windows 下执行
	getpkg.bat
	// linux/mac 下执行
	sh getpkg

3、编译并运行 studygolang

先编译

	// 接着上一步
	cd ../studygolang/
	// windows 下执行
	install.bat
	// linux/mac 下执行
	sh install
	
这样便编译好了 studygolang，下面运行 studygolang。（运行前可以根据需要修改 config/config.json 配置）

	// windows 下执行
	start.bat
	// linux/mac 下执行
	sh start

一切顺利的话，studygolang 应该就启动了。

4、浏览器中查看

在浏览器中输入：http://127.0.0.1:8080

应该就能看到了。

5、建立数据库

运行起来了，但没有建数据库。源码中有一个 databases 文件夹，里面有建表和初始化的sql语句。之前这些sql之前，在mysql数据库中建立一个数据库：studygolang，之后执行这些sql语句。

根据你的数据库设置，修改上面提到的 `config/config.json` 对应的配置，重新启动 studygolang.（通过restart脚本重新启动）