// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

//go:build !windows && !plan9
// +build !windows,!plan9

package controller

import (
	"os"
	"syscall"
)

func (InstallController) reload() {
	syscall.Kill(os.Getpid(), syscall.SIGUSR2)
}
