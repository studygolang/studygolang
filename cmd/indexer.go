// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package cmd

import (
	"math/rand"
	"time"

	"github.com/studygolang/studygolang/modules/server"

	"github.com/polaris1119/config"
	"github.com/polaris1119/keyword"
	"github.com/polaris1119/logger"
)

func Indexer() {
	logger.Init(config.ROOT+"/log", config.ConfigFile.MustValue("global", "log_level", "DEBUG"))
	go keyword.Extractor.Init(keyword.DefaultProps, true, config.ROOT+"/data/programming.txt,"+config.ROOT+"/data/dictionary.txt")

	server.IndexingServer()

	select {}
}
