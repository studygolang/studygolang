// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package logger

import (
	"config"
	"io"
	"log"
	"os"
	"time"
)

var (
	// 日志文件
	info_file  = config.ROOT + "/log/info.log"
	debug_file = config.ROOT + "/log/debug.log"
	error_file = config.ROOT + "/log/error.log"
)

func init() {
	os.Mkdir(config.ROOT+"/log/", 0777)
}

type logger struct {
	*log.Logger
}

func Infof(format string, args ...interface{}) {
	file, err := openFile(info_file)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Printf(format, args...)
}

func Infoln(args ...interface{}) {
	file, err := openFile(info_file)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Println(args...)
}

func Errorf(format string, args ...interface{}) {
	file, err := openFile(error_file)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Printf(format, args...)
}

func Errorln(args ...interface{}) {
	file, err := openFile(error_file)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Println(args...)
}

func New(out io.Writer) *logger {
	return &logger{
		Logger: log.New(out, "", log.Ltime),
	}
}

func openFile(filename string) (*os.File, error) {
	filename += "-" + time.Now().Format("060102")

	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
}
