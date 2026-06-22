// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package logic

import (
	"fmt"
	"unicode/utf8"

	"github.com/studygolang/studygolang/internal/model"

	"github.com/polaris1119/logger"
)

var (
	publishObservable Observable
	modifyObservable  Observable
	commentObservable Observable
	ViewObservable    Observable
	appendObservable  Observable
	topObservable     Observable
	likeObservable    Observable
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

	ViewObservable = NewConcreteObservable(actionView)
	ViewObservable.AddObserver(&UserWeightObserver{})
	ViewObservable.AddObserver(&TodayActiveObserver{})
	ViewObservable.AddObserver(&FeedSeqObserver{})

	appendObservable = NewConcreteObservable(actionAppend)
	appendObservable.AddObserver(&UserWeightObserver{})
	appendObservable.AddObserver(&TodayActiveObserver{})
	appendObservable.AddObserver(&UserRichObserver{})

	topObservable = NewConcreteObservable(actionTop)
	topObservable.AddObserver(&UserWeightObserver{})
	topObservable.AddObserver(&TodayActiveObserver{})
	topObservable.AddObserver(&UserRichObserver{})

	likeObservable = NewConcreteObservable(actionLike)
	likeObservable.AddObserver(&UserWeightObserver{})
	likeObservable.AddObserver(&TodayActiveObserver{})
	likeObservable.AddObserver(&UserRichObserver{})
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
	actionAppend  = "append"
	actionTop     = "top"  // 置顶
	actionLike    = "like" // 喜欢（赞）
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
		// 每个 observer 独立 recover：任一 observer panic 不能中断其他 observer，
		// 也不能让 panic 冒泡到调用方（多处 go NotifyObservers 会让进程崩溃）。
		func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Errorln("NotifyObservers observer panic:", observer, "action:", this.action, "recover:", r)
				}
			}()
			observer.Update(this.action, uid, objtype, objid)
		}()
	}
}

// ///////////////////////// 具体观察者 ////////////////////////////////////////

type UserWeightObserver struct{}

func (this *UserWeightObserver) Update(action string, uid, objtype, objid int) {
	if uid == 0 {
		return
	}

	var weight int
	switch action {
	case actionPublish:
		weight = 20
	case actionModify:
		weight = 2
	case actionComment:
		weight = 5
	case actionView:
		weight = 1
	case actionAppend:
		weight = 15
	case actionTop:
		weight = 5
	case actionLike:
		weight = 3
	}

	DefaultUser.IncrUserWeight("uid", uid, weight)
}

type TodayActiveObserver struct{}

func (*TodayActiveObserver) Update(action string, uid, objtype, objid int) {
	if uid == 0 {
		return
	}

	var weight int

	switch action {
	case actionPublish:
		weight = 20
		// 记录当天发布次数和上次发布时时间
		incrPublishTimes(uid)
		recordLastPubishTime(uid)
	case actionModify:
		weight = 2
	case actionComment:
		weight = 5
	case actionView:
		weight = 1
	case actionAppend:
		weight = 15
	case actionTop:
		weight = 5
	case actionLike:
		weight = 5
	}

	DefaultRank.GenDAURank(uid, weight)
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
	if uid == 0 {
		return
	}

	user := DefaultUser.FindOne(nil, "uid", uid)
	// actor 可能在 action 之后被删除，FindOne 返回 &User{} 但 Uid=0。
	// 若不拦截，下面所有 user.Username 引用会拼出空字符串描述，
	// 最后 IncrUserRich(user, ...) 还会用 uid=0 写入 bogus 余额记录。
	if user.Uid == 0 {
		return
	}

	var (
		typ   int
		award int
		desc  string
	)

	if action == actionPublish || action == actionComment {
		var comment *model.Comment
		if action == actionComment {
			comment, _ = DefaultComment.FindById(objid)
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
					// 主题作者可能已注销：FindOne 返回 &User{} 但 Uid=0，
					// 直接 IncrUserRich 会写入 uid=0 的 bogus 余额记录。
					if author.Uid != 0 {
						DefaultUserRich.IncrUserRich(author, model.MissionTypeReplied, 5, replyDesc)
					}
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
					// 文章作者账号可能已注销（外部站文章 author 字段为昵称，也可能匹配不到本站用户），
					// FindOne 返回 &User{} 但 Uid=0，IncrUserRich 会写入 uid=0 的 bogus 记录。
					if author.Uid != 0 {
						DefaultUserRich.IncrUserRich(author, model.MissionTypeReplied, 5, replyDesc)
					}
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
					// 资源作者可能已注销：FindOne 返回 &User{} 但 Uid=0，
					// IncrUserRich 会写入 uid=0 的 bogus 余额记录。
					if author.Uid != 0 {
						DefaultUserRich.IncrUserRich(author, model.MissionTypeReplied, 5, replyDesc)
					}
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
					// 项目作者可能已注销：FindOne 返回 &User{} 但 Uid=0，
					// IncrUserRich 会写入 uid=0 的 bogus 余额记录。
					if author.Uid != 0 {
						DefaultUserRich.IncrUserRich(author, model.MissionTypeReplied, 5, replyDesc)
					}
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
					// 注意 URL 用 wiki.Uri（slug 字符串），不能用 objid（数字 id），
					// 否则生成的链接 /wiki/123 是 404。上方自己的 desc 也是 %s+Uri。
					replyDesc := fmt.Sprintf(`收到 <a href="/user/%s">%s</a> 的回复 › <a href="/wiki/%s">%s</a>`,
						user.Username,
						user.Username,
						wiki.Uri,
						wiki.Title)
					author := DefaultUser.FindOne(nil, "uid", wiki.Uid)
					// WIKI 作者可能已注销：FindOne 返回 &User{} 但 Uid=0，
					// IncrUserRich 会写入 uid=0 的 bogus 余额记录。
					if author.Uid != 0 {
						DefaultUserRich.IncrUserRich(author, model.MissionTypeReplied, 5, replyDesc)
					}
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
			} else {
				desc = fmt.Sprintf(`分享一本图书 › <a href="/book/%d">%s</a>`,
					book.Id,
					book.Name)
			}
		}
	} else if action == actionModify {
		// TODO：修改暂时不消耗铜币
		// DefaultUserRich.IncrUserRich(uid, model.MissionTypeModify, -2, desc)
		return
	} else if action == actionView {
		return
	} else if action == actionAppend {
		typ = model.MissionTypeAppend
		award = -15
		topic := DefaultTopic.findByTid(objid)
		desc = fmt.Sprintf(`为主题 › <a href="/topics/%d">%s</a> 增加附言`,
			topic.Tid,
			topic.Title)
	} else if action == actionTop {
		typ = model.MissionTypeTop
		award = -30000

		switch objtype {
		case model.TypeTopic:
			topic := DefaultTopic.findByTid(objid)
			desc = fmt.Sprintf(`将主题 › <a href="/topics/%d">%s</a> 置顶`,
				topic.Tid,
				topic.Title)
		case model.TypeArticle:
			article, _ := DefaultArticle.FindById(nil, objid)
			desc = fmt.Sprintf(`将文章 › <a href="/articles/%d">%s</a> 置顶`,
				article.Id,
				article.Title)
		}
	} else if action == actionLike {
		// TODO: 暂时不处理
		return
	}

	DefaultUserRich.IncrUserRich(user, typ, award, desc)
}

type FeedSeqObserver struct{}

func (this *FeedSeqObserver) Update(action string, uid, objtype, objid int) {
	if objid == 0 {
		return
	}

	if action == actionView {
		DefaultFeed.updateSeq(objid, objtype, 0, 0, 1)
	}
}
