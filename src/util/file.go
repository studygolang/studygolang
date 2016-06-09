// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package util

import (
	"os"
)

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// 列出指定路径中的文件和目录
// 如果目录不存在，则返回空slice
func ScanDir(directory string) []string {
	file, err := os.Open(directory)
	if err != nil {
		return []string{}
	}
	names, err := file.Readdirnames(-1)
	if err != nil {
		return []string{}
	}
	return names
}

// 判断给定文件名是否是一个目录
// 如果文件名存在并且为目录则返回 true。如果 filename 是一个相对路径，则按照当前工作目录检查其相对路径。
func IsDir(filename string) bool {
	return isFileOrDir(filename, true)
}

// 判断给定文件名是否为一个正常的文件
// 如果文件存在且为正常的文件则返回 true
func IsFile(filename string) bool {
	return isFileOrDir(filename, false)
}

// 判断是文件还是目录，根据decideDir为true表示判断是否为目录；否则判断是否为文件
func isFileOrDir(filename string, decideDir bool) bool {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return false
	}
	isDir := fileInfo.IsDir()
	if decideDir {
		return isDir
	}
	return !isDir
}
