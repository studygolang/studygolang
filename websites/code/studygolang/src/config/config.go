// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package config

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"process"
)

// 项目根目录
var ROOT string

var Config map[string]string

func init() {
	binDir, err := process.ExecutableDir()
	if err != nil {
		panic(err)
	}
	ROOT = path.Dir(binDir)

	// Load配置文件
	configFile := ROOT + "/conf/config.json"
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	Config = make(map[string]string)
	err = json.Unmarshal(content, &Config)
	if err != nil {
		panic(err)
	}
}
