# studygolang

[![Build Status](https://travis-ci.org/studygolang/studygolang.svg?branch=master)](https://travis-ci.org/studygolang/studygolang)

[Go语言中文网 - Golang中文社区](https://studygolang.com "Go语言中文网 - Golang中文社区") 源码

网站上线时间：2013-03-15 14:38:09

目前在线运行的分支是 Master。欢迎有兴趣的 gopher 们参与进来，一起构建一个完善的 Go 语言中文网，Go 语言爱好者的学习家园，参与方式请参考：https://studygolang.com/topics/4092

## 本地搭建一个 Go语言中文网

要求 Go 1.16+

1、下载源码到本地某个目录

```shell
git clone https://github.com/studygolang/studygolang
```

2、编译

进入 studygolang 项目目录，执行如下命令：

```shell
// unix
make build
// windows
install.bat
```

这样便编译好了 studygolang

3、在 studygolang 源码中的 bin 目录下应该有了 studygolang 可执行文件。

接下来启动 studygolang。

```shell
// unix
make start
// windows
start.bat
```

或者

```shell
// unix
bin/studygolang
// windows
bin\studygolang.exe
```

一切顺利的话，studygolang 应该就启动了。

4、验证

在浏览器中输入：http://127.0.0.1:8088

应该就能看到了。

接下来你会看到图形化安装界面，一步步照做吧。

* 如果之后有出现页面空白，请查看 error.log 是否有错误

## 参与我们

fork + PR。如果有修改 js 和 css，请执行 gulp （需要先安装 gulp）。注意，Node 版本为：v10.16.2

## 使用该项目搭建的网站

- [Go语言中文网](https://studygolang.com)
