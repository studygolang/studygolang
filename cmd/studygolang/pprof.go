// Copyright 2018 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: meission	meission@aliyun.com

package main

import (
	"fmt"
	"net/http"
	"net/http/pprof"
)

// Pprof start http pprof.
func Pprof(addr string) {
	ps := http.NewServeMux()
	ps.HandleFunc("/debug/pprof/", pprof.Index)
	ps.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	ps.HandleFunc("/debug/pprof/profile", pprof.Profile)
	ps.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	go func() {
		if err := http.ListenAndServe(addr, ps); err != nil {
			fmt.Println("pprof exit:", err)
		}
	}()
}
