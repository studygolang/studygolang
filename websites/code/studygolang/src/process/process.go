// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package process

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// 保存pid
func SavePidTo(pidFile string) error {
	pidPath := filepath.Dir(pidFile)
	if err := os.MkdirAll(pidPath, 0777); err != nil {
		return err
	}
	return ioutil.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0777)
}

// 获得可执行程序所在目录
func ExecutableDir() (string, error) {
	pathAbs, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Dir(pathAbs), nil
}
