# 简介

目前版本基于 [echo](https://github.com/labstack/echo) 框架开发，mysql 数据库操作使用 [xorm](http://books.studygolang.com/xorm)。

## 项目目录结构

```console
├── LICENSE
├── README.md
├── bin
│   ├── crawler 抓取程序
│   ├── indexer 索引程序
│   └── studygolang 网站主程序
├── config 配置文件目录
│   ├── db.sql 建表语句
│   ├── env.ini 配置文件
│   ├── env.sample.ini 配置文件样本，图形安装过程据此生成 env.ini
│   ├── init.sql 初始化脚本
│   └── solr_schema.xml solr 配置
├── data
│   ├── dictionary.txt github.com/huichen/sego 分词使用的词典
│   ├── max_online_num 记录在线历史最高人数
│   └── programming.txt 自定义词典
├── docs 文档
│   ├── README.md
│   └── intro.md
├── getpkg.bat
├── getpkg.sh 下载依赖
├── install.bat
├── install.sh 编译
├── log 日志目录
├── reload.bat
├── reload.sh 重启
├── robots.txt
├── sitemap 存放搜素引擎 sitemap
├── src 源码
│   ├── db
│   ├── global
│   ├── http
│   ├── logic
│   ├── model
│   ├── server
│   ├── util
│   └── vendor
├── start.bat
├── start.sh 启动网站
├── static 静态文件
│   ├── ckeditor
│   ├── css
│   ├── fonts
│   ├── img
│   ├── js
│   └── upload
├── stop.bat
├── stop.sh 停止网站
└── template 模板
    ├── 403.html
    ├── 404.html
    ├── 500.html
    ├── admin
    ├── articles
    ├── atom.html
    ├── books
    ├── common
    ├── download
    ├── email.html
    ├── favorite.html
    ├── feed
    ├── gift
    ├── index.html
    ├── install
    ├── link.html
    ├── login.html
    ├── markdown.html
    ├── messages
    ├── mission
    ├── notfound.html
    ├── pkgdoc.html
    ├── projects
    ├── readings
    ├── register.html
    ├── resources
    ├── rich
    ├── search.html
    ├── sidebar
    ├── sitemap.xml
    ├── sitemapindex.xml
    ├── top
    ├── topics
    ├── user
    ├── wide
    ├── wiki
    └── wr.html
```

### 源码的组织结构

```console
├── src
│   ├── db
│   ├── global
│   ├── http
│   ├── logic
│   ├── model
│   ├── server
│   ├── util
│   └── vendor
```

- db 包  
    负责初始化 xorm，构造数据库对象，需要数据库操作的地方，只需要如下方式使用即可：
    
    ```go
    import . "db"

    MasterDB.Where()....
    ```
    
- global 包  
    全局的一些对象，比如 App（包含网站的一些配置信息）、全局的一些 channel 等。
- http 包  
    包含 controller 和 middleware；其中的 http.go 文件封装了模板处理的一些通用逻辑。
- logic 包  
    所有业务逻辑
- model 包  
    数据库实体，ORM
- server 包  
    存放 main
- util 包  
    一些辅助函数
