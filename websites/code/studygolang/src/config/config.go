// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"

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
	ROOT = path.Dir(binDir) + "/"

	if !strings.Contains(ROOT, "studygolang") {
		ROOT, _ = os.Getwd()
		ROOT = ROOT[:strings.Index(ROOT, "src")]
	}

	// Load配置文件
	configFile := ROOT + "conf/config.json"
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

const Gt = ">"

type Conf map[string]interface{}

func ParseConfig(filename string, store interface{}) (Conf, error) {
	content, err := ioutil.ReadFile(ROOT + filename)
	if err != nil {
		return nil, err
	}

	var conf Conf
	if store == nil {
		store = &conf
	} else {
		storeType := reflect.TypeOf(store)
		if storeType.Kind() != reflect.Ptr {
			return nil, errors.New("store must be pointer or nil")
		}
	}

	err = json.Unmarshal(content, store)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (self Conf) String(xpath string) string {
	var value string

	if !strings.Contains(xpath, Gt) {
		value, _ = self[xpath].(string)
		return value
	}

	ret := self.Convert(xpath)
	if ret != nil {
		value, _ = ret.(string)
	}

	return value
}

func (self Conf) Int(xpath string) int {
	var value int

	if !strings.Contains(xpath, Gt) {
		value, _ = strconv.Atoi(self[xpath].(string))
		return value
	}

	ret := self.Convert(xpath)
	if ret != nil {
		value, _ = strconv.Atoi(ret.(string))
	}

	return value
}

func (self Conf) Bool(xpath string) bool {
	var value bool

	if !strings.Contains(xpath, Gt) {
		value, _ = self[xpath].(bool)
		return value
	}

	ret := self.Convert(xpath)
	if ret != nil {
		value, _ = ret.(bool)
	}

	return value
}

func (self Conf) Float64(xpath string) float64 {
	var value float64

	if !strings.Contains(xpath, Gt) {
		value, _ = self[xpath].(float64)
		return value
	}

	ret := self.Convert(xpath)
	if ret != nil {
		value, _ = ret.(float64)
	}

	return value
}

// xpath 使用 ">" 方式分隔json嵌套层级，如 a>b>c
func (self Conf) Convert(xpath string) interface{} {
	paths := strings.Split(xpath, Gt)

	middleValue := self[paths[0]]
	paths = paths[1:]

	for _, p := range paths {
		switch v := middleValue.(type) {
		case map[string]interface{}:
			if vi, ok := v[p]; !ok {
				return nil
			} else {
				middleValue = vi
			}
		case []interface{}:
			pInt, err := strconv.Atoi(p)
			if err != nil {
				return nil
			}

			if len(v) <= pInt {
				return nil
			}
			middleValue = v[pInt]
		}
	}

	return middleValue
}
