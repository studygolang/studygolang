// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package global

var AuthorityChan = make(chan struct{}, 1)
var RoleChan = make(chan struct{}, 1)
var RoleAuthChan = make(chan struct{}, 1)
var UserSettingChan = make(chan struct{}, 1)
var TopicNodeChan = make(chan struct{}, 1)
