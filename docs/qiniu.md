Go 中文网静态资源托管在七牛云上，对于 js 和 css，如果有变化，需要更新对应的文件，通过 [qshell](https://developer.qiniu.com/kodo/tools/1302/qshel) 工具可以做到

按文档安装完后，需要设置 account，之后执行类似如下命令来替换 js 或 css：

qshell fput studygolang static/dist/js/sg_base.min.js dist/js/sg_base.min.js -w

即：qshell fput <Bucket> <Key> <LocalFile> [<Overwrite>]
