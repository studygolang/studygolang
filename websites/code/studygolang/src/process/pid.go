// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package process

import (
	"io/ioutil"
	"os"
	"strconv"
)

func SavePidTo(pidFile string) error {
	return ioutil.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0777)
}
