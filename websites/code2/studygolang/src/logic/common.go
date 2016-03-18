// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package logic

import (
	"os"

	"github.com/gorilla/schema"
	"github.com/polaris1119/logger"
	"golang.org/x/net/context"
)

var schemaDecoder = schema.NewDecoder()

func init() {
	schemaDecoder.SetAliasTag("json")
}

func GetLogger(ctx context.Context) *logger.Logger {
	if ctx == nil {
		return logger.New(os.Stdout)
	}

	_logger, ok := ctx.Value("logger").(*logger.Logger)
	if ok {
		return _logger
	}

	return logger.New(os.Stdout)
}
