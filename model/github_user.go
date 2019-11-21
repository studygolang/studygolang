// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import (
	"code.gitea.io/sdk/gitea"
)

type GithubUser struct {
	Id        int    `json:"id"`
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Company   string `json:"company"`
	Blog      string `json:"blog"`
	Location  string `json:"location"`
}

type GiteaUser = gitea.User

func DisplayName(g *GiteaUser) string {
	if g.FullName == "" {
		return g.UserName
	}
	return g.FullName
}
