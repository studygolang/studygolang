// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

/*
sets version information for the binary where it is imported.
The version can be retrieved either from the -version command line argument.

To include in a project simply import the package.

The version and compile date is stored in App variables and
are supposed to be set during compile time. Typically this is done by the
install(bash/bat). Or date is binary modify time.

To set these manually use -ldflags together with -X, like in this example:

	go install -ldflags "-X global/Build xxxxx"

*/

package global

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/studygolang/studygolang/internal/model"

	"github.com/polaris1119/config"
)

const (
	DefaultCDNHttp  = "http://test.static.golangjob.cn/"
	DefaultCDNHttps = "https://static.golangjob.cn/"
)

var Build string

type app struct {
	Name    string
	Build   string
	Version string
	Date    time.Time

	Copyright string

	// 启动时间
	LaunchTime time.Time
	Uptime     time.Duration

	Env string

	Host string
	Port string

	BaseURL string

	// CDN 资源域名
	CDNHttp  string
	CDNHttps string

	Domain string

	locker sync.Mutex
}

var App = &app{}

var showVersion = flag.Bool("version", false, "Print version of this binary")

const (
	DEV  = "dev"
	TEST = "test"
	PRO  = "pro"
)

func init() {
	App.Name = os.Args[0]
	App.Version = "V4.0.0"
	App.Build = Build
	App.LaunchTime = time.Now()

	fileInfo, err := os.Stat(os.Args[0])
	if err != nil {
		panic(err)
	}

	App.Date = fileInfo.ModTime()

	App.Env = config.ConfigFile.MustValue("global", "env")

	App.CDNHttp = config.ConfigFile.MustValue("qiniu", "http_domain", DefaultCDNHttp)
	App.CDNHttps = config.ConfigFile.MustValue("qiniu", "https_domain", DefaultCDNHttps)
}

func (this *app) Init(domain string) {
	this.Domain = config.ConfigFile.MustValue("global", "domain", domain)
}

func (this *app) SetUptime() {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.Uptime = time.Now().Sub(this.LaunchTime)
}

func (this *app) SetCopyright() {
	curYear := time.Now().Year()
	this.locker.Lock()
	defer this.locker.Unlock()
	if curYear == model.WebsiteSetting.StartYear {
		this.Copyright = fmt.Sprintf("%d %s", model.WebsiteSetting.StartYear, model.WebsiteSetting.Domain)
	} else {
		this.Copyright = fmt.Sprintf("%d-%d %s", model.WebsiteSetting.StartYear, curYear, model.WebsiteSetting.Domain)
	}
}

func (this *app) CanonicalCDN(isHTTPS bool) string {
	cdnDomain := this.CDNHttp
	if isHTTPS {
		cdnDomain = this.CDNHttps
	}
	if !strings.HasSuffix(cdnDomain, "/") {
		cdnDomain += "/"
	}

	return cdnDomain
}

func PrintVersion(w io.Writer) {
	if !flag.Parsed() {
		flag.Parse()
	}

	if showVersion == nil || !*showVersion {
		return
	}

	fmt.Fprintf(w, "Binary: %s\n", App.Name)
	fmt.Fprintf(w, "Version: %s\n", App.Version)
	fmt.Fprintf(w, "Build: %s\n", App.Build)
	fmt.Fprintf(w, "Compile date: %s\n", App.Date.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "Env: %s\n", App.Env)
	os.Exit(0)
}

func OnlineEnv() bool {
	return App.Env == PRO
}
