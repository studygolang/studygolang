// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: javasgl	songganglin@gmail.com
package main

import (
	"server"

	"github.com/polaris1119/config"
	"github.com/polaris1119/logger"
)

func init() {

}

func main() {

	logger.Init(config.ROOT+"/log", config.ConfigFile.MustValue("global", "log_level", "DEBUG"))
	server.MigratorServer()

}
