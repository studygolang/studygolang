// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package logic_test

import (
	"testing"

	"github.com/studygolang/studygolang/internal/logic"

	"github.com/polaris1119/config"
	"github.com/polaris1119/logger"
)

func TestPullPR(t *testing.T) {
	logger.Init(config.ROOT+"/log", config.ConfigFile.MustValue("global", "log_level", "DEBUG"))

	err := logic.DefaultGithub.PullPR("studygolang/GCTT", true)
	if err != nil {
		t.Error("pull request error:", err)
	}
}

func TestSyncIssues(t *testing.T) {
	logger.Init(config.ROOT+"/log", config.ConfigFile.MustValue("global", "log_level", "DEBUG"))

	err := logic.DefaultGithub.SyncIssues("studygolang/GCTT", true)
	if err != nil {
		t.Error("SyncIssues error:", err)
	}
}

func TestIssueEvent(t *testing.T) {
	logger.Init(config.ROOT+"/log", config.ConfigFile.MustValue("global", "log_level", "DEBUG"))

	body := []byte(`{
		"action": "closed",
		"issue": {
		  "url": "https://api.github.com/repos/studygolang/GCTT/issues/110",
		  "repository_url": "https://api.github.com/repos/studygolang/GCTT",
		  "labels_url": "https://api.github.com/repos/studygolang/GCTT/issues/110/labels{/name}",
		  "comments_url": "https://api.github.com/repos/studygolang/GCTT/issues/110/comments",
		  "events_url": "https://api.github.com/repos/studygolang/GCTT/issues/110/events",
		  "html_url": "https://github.com/studygolang/GCTT/issues/110",
		  "id": 279211537,
		  "number": 110,
		  "title": "20171205 What’s the most common identifier in Go’s stdlib?",
		  "user": {
			"login": "polaris1119",
			"id": 899673,
			"avatar_url": "https://avatars1.githubusercontent.com/u/899673?v=4",
			"gravatar_id": "",
			"url": "https://api.github.com/users/polaris1119",
			"html_url": "https://github.com/polaris1119",
			"followers_url": "https://api.github.com/users/polaris1119/followers",
			"following_url": "https://api.github.com/users/polaris1119/following{/other_user}",
			"gists_url": "https://api.github.com/users/polaris1119/gists{/gist_id}",
			"starred_url": "https://api.github.com/users/polaris1119/starred{/owner}{/repo}",
			"subscriptions_url": "https://api.github.com/users/polaris1119/subscriptions",
			"organizations_url": "https://api.github.com/users/polaris1119/orgs",
			"repos_url": "https://api.github.com/users/polaris1119/repos",
			"events_url": "https://api.github.com/users/polaris1119/events{/privacy}",
			"received_events_url": "https://api.github.com/users/polaris1119/received_events",
			"type": "User",
			"site_admin": false
		  },
		  "labels": [
			{
			  "id": 768962805,
			  "url": "https://api.github.com/repos/studygolang/GCTT/labels/%E5%B7%B2%E8%AE%A4%E9%A2%86",
			  "name": "已认领",
			  "color": "5edb81",
			  "default": false
			}
		  ],
		  "state": "closed",
		  "locked": false,
		  "assignee": null,
		  "assignees": [
	  
		  ],
		  "milestone": null,
		  "comments": 1,
		  "created_at": "2017-12-05T01:22:18Z",
		  "updated_at": "2018-01-18T06:35:08Z",
		  "closed_at": "2018-01-18T06:35:08Z",
		  "author_association": "CONTRIBUTOR",
		  "body": "标题：What’s the most common identifier"
		},
		"repository": {
		  "id": 110936509,
		  "name": "GCTT",
		  "full_name": "studygolang/GCTT",
		  "owner": {
			"login": "studygolang",
			"id": 3772217,
			"avatar_url": "https://avatars3.githubusercontent.com/u/3772217?v=4",
			"gravatar_id": "",
			"url": "https://api.github.com/users/studygolang",
			"html_url": "https://github.com/studygolang",
			"followers_url": "https://api.github.com/users/studygolang/followers",
			"following_url": "https://api.github.com/users/studygolang/following{/other_user}",
			"gists_url": "https://api.github.com/users/studygolang/gists{/gist_id}",
			"starred_url": "https://api.github.com/users/studygolang/starred{/owner}{/repo}",
			"subscriptions_url": "https://api.github.com/users/studygolang/subscriptions",
			"organizations_url": "https://api.github.com/users/studygolang/orgs",
			"repos_url": "https://api.github.com/users/studygolang/repos",
			"events_url": "https://api.github.com/users/studygolang/events{/privacy}",
			"received_events_url": "https://api.github.com/users/studygolang/received_events",
			"type": "Organization",
			"site_admin": false
		  },
		  "private": false,
		  "html_url": "https://github.com/studygolang/GCTT",
		  "description": "GCTT Go中文网翻译组。",
		  "fork": false,
		  "url": "https://api.github.com/repos/studygolang/GCTT",
		  "forks_url": "https://api.github.com/repos/studygolang/GCTT/forks",
		  "keys_url": "https://api.github.com/repos/studygolang/GCTT/keys{/key_id}",
		  "collaborators_url": "https://api.github.com/repos/studygolang/GCTT/collaborators{/collaborator}",
		  "teams_url": "https://api.github.com/repos/studygolang/GCTT/teams",
		  "hooks_url": "https://api.github.com/repos/studygolang/GCTT/hooks",
		  "issue_events_url": "https://api.github.com/repos/studygolang/GCTT/issues/events{/number}",
		  "events_url": "https://api.github.com/repos/studygolang/GCTT/events",
		  "assignees_url": "https://api.github.com/repos/studygolang/GCTT/assignees{/user}",
		  "branches_url": "https://api.github.com/repos/studygolang/GCTT/branches{/branch}",
		  "tags_url": "https://api.github.com/repos/studygolang/GCTT/tags",
		  "blobs_url": "https://api.github.com/repos/studygolang/GCTT/git/blobs{/sha}",
		  "git_tags_url": "https://api.github.com/repos/studygolang/GCTT/git/tags{/sha}",
		  "git_refs_url": "https://api.github.com/repos/studygolang/GCTT/git/refs{/sha}",
		  "trees_url": "https://api.github.com/repos/studygolang/GCTT/git/trees{/sha}",
		  "statuses_url": "https://api.github.com/repos/studygolang/GCTT/statuses/{sha}",
		  "languages_url": "https://api.github.com/repos/studygolang/GCTT/languages",
		  "stargazers_url": "https://api.github.com/repos/studygolang/GCTT/stargazers",
		  "contributors_url": "https://api.github.com/repos/studygolang/GCTT/contributors",
		  "subscribers_url": "https://api.github.com/repos/studygolang/GCTT/subscribers",
		  "subscription_url": "https://api.github.com/repos/studygolang/GCTT/subscription",
		  "commits_url": "https://api.github.com/repos/studygolang/GCTT/commits{/sha}",
		  "git_commits_url": "https://api.github.com/repos/studygolang/GCTT/git/commits{/sha}",
		  "comments_url": "https://api.github.com/repos/studygolang/GCTT/comments{/number}",
		  "issue_comment_url": "https://api.github.com/repos/studygolang/GCTT/issues/comments{/number}",
		  "contents_url": "https://api.github.com/repos/studygolang/GCTT/contents/{+path}",
		  "compare_url": "https://api.github.com/repos/studygolang/GCTT/compare/{base}...{head}",
		  "merges_url": "https://api.github.com/repos/studygolang/GCTT/merges",
		  "archive_url": "https://api.github.com/repos/studygolang/GCTT/{archive_format}{/ref}",
		  "downloads_url": "https://api.github.com/repos/studygolang/GCTT/downloads",
		  "issues_url": "https://api.github.com/repos/studygolang/GCTT/issues{/number}",
		  "pulls_url": "https://api.github.com/repos/studygolang/GCTT/pulls{/number}",
		  "milestones_url": "https://api.github.com/repos/studygolang/GCTT/milestones{/number}",
		  "notifications_url": "https://api.github.com/repos/studygolang/GCTT/notifications{?since,all,participating}",
		  "labels_url": "https://api.github.com/repos/studygolang/GCTT/labels{/name}",
		  "releases_url": "https://api.github.com/repos/studygolang/GCTT/releases{/id}",
		  "deployments_url": "https://api.github.com/repos/studygolang/GCTT/deployments",
		  "created_at": "2017-11-16T07:10:44Z",
		  "updated_at": "2018-01-18T06:16:04Z",
		  "pushed_at": "2018-01-17T15:46:12Z",
		  "git_url": "git://github.com/studygolang/GCTT.git",
		  "ssh_url": "git@github.com:studygolang/GCTT.git",
		  "clone_url": "https://github.com/studygolang/GCTT.git",
		  "svn_url": "https://github.com/studygolang/GCTT",
		  "homepage": "https://studygolang.com/gctt",
		  "size": 4554,
		  "stargazers_count": 255,
		  "watchers_count": 255,
		  "language": "Shell",
		  "has_issues": true,
		  "has_projects": true,
		  "has_downloads": true,
		  "has_wiki": true,
		  "has_pages": false,
		  "forks_count": 105,
		  "mirror_url": null,
		  "archived": false,
		  "open_issues_count": 38,
		  "license": {
			"key": "apache-2.0",
			"name": "Apache License 2.0",
			"spdx_id": "Apache-2.0",
			"url": "https://api.github.com/licenses/apache-2.0"
		  },
		  "forks": 105,
		  "open_issues": 38,
		  "watchers": 255,
		  "default_branch": "master"
		},
		"organization": {
		  "login": "studygolang",
		  "id": 3772217,
		  "url": "https://api.github.com/orgs/studygolang",
		  "repos_url": "https://api.github.com/orgs/studygolang/repos",
		  "events_url": "https://api.github.com/orgs/studygolang/events",
		  "hooks_url": "https://api.github.com/orgs/studygolang/hooks",
		  "issues_url": "https://api.github.com/orgs/studygolang/issues",
		  "members_url": "https://api.github.com/orgs/studygolang/members{/member}",
		  "public_members_url": "https://api.github.com/orgs/studygolang/public_members{/member}",
		  "avatar_url": "https://avatars3.githubusercontent.com/u/3772217?v=4",
		  "description": ""
		},
		"sender": {
		  "login": "polaris1119",
		  "id": 899673,
		  "avatar_url": "https://avatars1.githubusercontent.com/u/899673?v=4",
		  "gravatar_id": "",
		  "url": "https://api.github.com/users/polaris1119",
		  "html_url": "https://github.com/polaris1119",
		  "followers_url": "https://api.github.com/users/polaris1119/followers",
		  "following_url": "https://api.github.com/users/polaris1119/following{/other_user}",
		  "gists_url": "https://api.github.com/users/polaris1119/gists{/gist_id}",
		  "starred_url": "https://api.github.com/users/polaris1119/starred{/owner}{/repo}",
		  "subscriptions_url": "https://api.github.com/users/polaris1119/subscriptions",
		  "organizations_url": "https://api.github.com/users/polaris1119/orgs",
		  "repos_url": "https://api.github.com/users/polaris1119/repos",
		  "events_url": "https://api.github.com/users/polaris1119/events{/privacy}",
		  "received_events_url": "https://api.github.com/users/polaris1119/received_events",
		  "type": "User",
		  "site_admin": false
		}
	  }`)
	err := logic.DefaultGithub.IssueEvent(nil, body)
	if err != nil {
		t.Error("SyncIssues error:", err)
	}
}
