// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package logic

import (
	"fmt"
	"model"
	"unicode/utf8"
)

var (
	publishObservable Observable
	modifyObservable  Observable
	commentObservable Observable
	viewObservable    Observable
)

func init() {
	publishObservable = NewConcreteObservable(actionPublish)
	publishObservable.AddObserver(&UserWeightObserver{})
	publishObservable.AddObserver(&TodayActiveObserver{})
	publishObservable.AddObserver(&UserRichObserver{})

	modifyObservable = NewConcreteObservable(actionModify)
	modifyObservable.AddObserver(&UserWeightObserver{})
	modifyObservable.AddObserver(&TodayActiveObserver{})
	modifyObservable.AddObserver(&UserRichObserver{})

	commentObservable = NewConcreteObservable(actionComment)
	commentObservable.AddObserver(&UserWeightObserver{})
	commentObservable.AddObserver(&TodayActiveObserver{})
	commentObservable.AddObserver(&UserRichObserver{})

	viewObservable = NewConcreteObservable(actionView)
	viewObservable.AddObserver(&UserWeightObserver{})
	viewObservable.AddObserver(&TodayActiveObserver{})
}

type Observer interface {
	Update(action string, uid, objtype, objid int)
}

type Observable interface {
	// AddObserver 登记一个新的观察者
	AddObserver(o Observer)
	// Detach 删除一个已经登记过的观察者
	RemoveObserver(o Observer)
	// NotifyObservers 通知所有登记过的观察者
	NotifyObservers(uid, objtype, objid int)
}

const (
	actionPublish = "publish"
	actionModify  = "modify"
	actionComment = "comment"
	actionView    = "view"
)

type ConcreteObservable struct {
	observers []Observer
	action    string
}

func NewConcreteObservable(action string) *ConcreteObservable {
	return &ConcreteObservable{
		action:    action,
		observers: make([]Observer, 0, 8),
	}
}

func (this *ConcreteObservable) AddObserver(o Observer) {
	this.observers = append(this.observers, o)
}

func (this *ConcreteObservable) RemoveObserver(o Observer) {
	if len(this.observers) == 0 {
		return
	}

	var indexToRemove int

	for i, observer := range this.observers {
		if observer == o {
			indexToRemove = i
			break
		}
	}

	this.observers = append(this.observers[:indexToRemove], this.observers[indexToRemove+1:]...)
}

func (this *ConcreteObservable) NotifyObservers(uid, objtype, objid int) {
	for _, observer := range this.observers {
		observer.Update(this.action, uid, objtype, objid)
	}
}

/////////////////////////// 具体观察者 ////////////////////////////////////////

type UserWeightObserver struct{}

func (this *UserWeightObserver) Update(action string, uid, objtype, objid int) {
	if action == actionPublish {
		DefaultUser.IncrUserWeight("uid", uid, 20)
	} else if action == actionModify {
		DefaultUser.IncrUserWeight("uid", uid, 2)
	} else if action == actionComment {
		DefaultUser.IncrUserWeight("uid", uid, 5)
	} else if action == actionView {
		DefaultUser.IncrUserWeight("uid", uid, 1)
	}
}

type TodayActiveObserver struct{}

func (*TodayActiveObserver) Update(action string, uid, objtype, objid int) {
	if action == actionPublish {
		DefaultRank.GenDAURank(uid, 20)
	} else if action == actionModify {
		DefaultRank.GenDAURank(uid, 2)
	} else if action == actionComment {
		DefaultRank.GenDAURank(uid, 5)
	} else if action == actionView {
		DefaultRank.GenDAURank(uid, 1)
	}
}

type UserRichObserver struct{}

var objType2MissType = map[int]int{
	model.TypeTopic:    model.MissionTypeTopic,
	model.TypeArticle:  model.MissionTypeArticle,
	model.TypeResource: model.MissionTypeResource,
	model.TypeWiki:     model.MissionTypeWiki,
	model.TypeBook:     model.MissionTypeBook,
	model.TypeProject:  model.MissionTypeProject,
}

// Update 如果是回复，则 objid 是 cid
func (UserRichObserver) Update(action string, uid, objtype, objid int) {
	user := DefaultUser.FindOne(nil, "uid", uid)

	var (
		typ   int
		award int
		desc  string
	)

	if action == actionPublish || action == actionComment {
		var comment *model.Comment
		if action == actionComment {
			comment = DefaultComment.findById(objid)
			if comment.Cid != objid {
				return
			}

			objid = comment.Objid

			award = -5
			typ = model.MissionTypeReply
		} else {
			award = -20
			typ = objType2MissType[objtype]
		}

		switch objtype {
		case model.TypeTopic:
			topic := DefaultTopic.findByTid(objid)
			if topic.Tid != objid {
				return
			}
			if action == actionComment {
				desc = fmt.Sprintf(`创建了长度为 %d 个字符的回复 › <a href="/topics/%d">%s</a>`,
					utf8.RuneCountInString(comment.Content),
					objid,
					topic.Title)

				if uid != topic.Uid {
					// 主题发起人获得收益
					replyDesc := fmt.Sprintf(`收到 <a href="/user/%s">%s</a> 的回复 › <a href="/topics/%d">%s</a>`,
						user.Username,
						user.Username,
						objid,
						topic.Title)
					author := DefaultUser.FindOne(nil, "uid", topic.Uid)
					DefaultUserRich.IncrUserRich(author, model.MissionTypeReplied, 5, replyDesc)
				}
			} else {
				desc = fmt.Sprintf(`创建了长度为 %d 个字符的主题 › <a href="/topics/%d">%s</a>`,
					utf8.RuneCountInString(topic.Content),
					objid,
					topic.Title)
			}

		case model.TypeArticle:
			article, err := DefaultArticle.FindById(nil, objid)
			if err != nil {
				return
			}
			if action == actionComment {
				desc = fmt.Sprintf(`创建了长度为 %d 个字符的回复 › <a href="/articles/%d">%s</a>`,
					utf8.RuneCountInString(comment.Content),
					objid,
					article.Title)
				if article.Domain == WebsiteSetting.Domain && user.Username != article.Author {
					// 文章发起人获得收益
					replyDesc := fmt.Sprintf(`收到 <a href="/user/%s">%s</a> 的回复 › <a href="/articles/%d">%s</a>`,
						user.Username,
						user.Username,
						objid,
						article.Title)
					author := DefaultUser.FindOne(nil, "username", article.Author)
					DefaultUserRich.IncrUserRich(author, model.MissionTypeReplied, 5, replyDesc)
				}
			} else {
				desc = fmt.Sprintf(`发表了长度为 %d 个字符的文章 › <a href="/articles/%d">%s</a>`,
					utf8.RuneCountInString(article.Txt),
					objid,
					article.Title)
			}
		case model.TypeResource:
			resource := DefaultResource.findById(objid)
			if resource.Id != objid {
				return
			}
			if action == actionComment {
				desc = fmt.Sprintf(`创建了长度为 %d 个字符的回复 › <a href="/resources/%d">%s</a>`,
					utf8.RuneCountInString(comment.Content),
					objid,
					resource.Title)

				if uid != resource.Uid {
					// 资源发起人获得收益
					replyDesc := fmt.Sprintf(`收到 <a href="/user/%s">%s</a> 的回复 › <a href="/resources/%d">%s</a>`,
						user.Username,
						user.Username,
						objid,
						resource.Title)
					author := DefaultUser.FindOne(nil, "uid", resource.Uid)
					DefaultUserRich.IncrUserRich(author, model.MissionTypeReplied, 5, replyDesc)
				}
			} else {

				desc = fmt.Sprintf(`分享了一个资源 › <a href="/resources/%d">%s</a>`,
					objid,
					resource.Title)
			}
		case model.TypeProject:
			project := DefaultProject.FindOne(nil, objid)
			if project == nil || project.Id != objid {
				return
			}
			if action == actionComment {
				desc = fmt.Sprintf(`创建了长度为 %d 个字符的回复 › <a href="/p/%d">%s</a>`,
					utf8.RuneCountInString(comment.Content),
					objid,
					project.Category+project.Name)

				if user.Username != project.Username {
					// 项目发起人获得收益
					replyDesc := fmt.Sprintf(`收到 <a href="/user/%s">%s</a> 的回复 › <a href="/p/%d">%s</a>`,
						user.Username,
						user.Username,
						objid,
						project.Category+project.Name)
					author := DefaultUser.FindOne(nil, "username", project.Username)
					DefaultUserRich.IncrUserRich(author, model.MissionTypeReplied, 5, replyDesc)
				}
			} else {
				desc = fmt.Sprintf(`发布了一个开源项目 › <a href="/p/%d">%s</a>`,
					objid,
					project.Category+project.Name)
			}
		case model.TypeWiki:
			wiki := DefaultWiki.FindById(nil, objid)
			if wiki == nil || wiki.Id != objid {
				return
			}
			if action == actionComment {
				desc = fmt.Sprintf(`创建了长度为 %d 个字符的回复 › <a href="/wiki/%s">%s</a>`,
					utf8.RuneCountInString(comment.Content),
					wiki.Uri,
					wiki.Title)

				if uid != wiki.Uid {
					// WIKI发起人获得收益
					replyDesc := fmt.Sprintf(`收到 <a href="/user/%s">%s</a> 的回复 › <a href="/wiki/%d">%s</a>`,
						user.Username,
						user.Username,
						objid,
						wiki.Title)
					author := DefaultUser.FindOne(nil, "uid", wiki.Uid)
					DefaultUserRich.IncrUserRich(author, model.MissionTypeReplied, 5, replyDesc)
				}
			} else {
				desc = fmt.Sprintf(`创建了长度为 %d 个字符的WIKI › <a href="/wiki/%s">%s</a>`,
					utf8.RuneCountInString(wiki.Content),
					wiki.Uri,
					wiki.Title)
			}
		case model.TypeBook:
			book, err := DefaultGoBook.FindById(nil, objid)
			if err != nil || book.Id != objid {
				return
			}
			if action == actionComment {
				desc = fmt.Sprintf(`创建了长度为 %d 个字符的回复 › <a href="/book/%d">%s</a>`,
					utf8.RuneCountInString(comment.Content),
					book.Id,
					book.Name)
			}
		}
	} else if action == actionModify {
		// TODO：修改暂时不消耗铜币
		// DefaultUserRich.IncrUserRich(uid, model.MissionTypeModify, -2, desc)
	} else if action == actionView {
		return
	}

	DefaultUserRich.IncrUserRich(user, typ, award, desc)
}
