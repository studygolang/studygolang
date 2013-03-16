// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

// +build debug

package logger

import (
	"os"
	"path/filepath"
	"runtime"
)

func Debugf(format string, args ...interface{}) {
	file, err := os.OpenFile(debug_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Printf(format, args...)
}

func Debugln(args ...interface{}) {
	file, err := os.OpenFile(debug_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer file.Close()
	// 加上文件调用和行号
	_, callerFile, line, ok := runtime.Caller(1)
	if ok {
		args = append([]interface{}{"文件：", filepath.Base(callerFile), "行号:", line}, args...)
	}
	New(file).Println(args...)
}
