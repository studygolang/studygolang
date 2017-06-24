// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	. "db"
	"model"
	"net/url"
	"time"

	"github.com/polaris1119/logger"
	"golang.org/x/net/context"
)

type GoBookLogic struct{}

var DefaultGoBook = GoBookLogic{}

func (self GoBookLogic) Publish(ctx context.Context, user *model.Me, form url.Values) (err error) {
	objLog := GetLogger(ctx)

	id := form.Get("id")
	isModify := id != ""

	book := &model.Book{}

	if isModify {
		_, err = MasterDB.Id(id).Get(book)
		if err != nil {
			objLog.Errorln("Publish Book find error:", err)
			return
		}

		if !CanEdit(user, book) {
			err = NotModifyAuthorityErr
			return
		}

		err = schemaDecoder.Decode(book, form)
		if err != nil {
			objLog.Errorln("Publish Book schema decode error:", err)
			return
		}
	} else {
		err = schemaDecoder.Decode(book, form)
		if err != nil {
			objLog.Errorln("Publish Book schema decode error:", err)
			return
		}

		book.Lastreplytime = model.NewOftenTime()
		book.Uid = user.Uid
	}

	var affected int64
	if !isModify {
		affected, err = MasterDB.Insert(book)
	} else {
		affected, err = MasterDB.Update(book)
	}

	if err != nil {
		objLog.Errorln("Publish Book error:", err)
		return
	}

	if affected == 0 {
		return
	}

	if isModify {
		go modifyObservable.NotifyObservers(user.Uid, model.TypeBook, book.Id)
	} else {
		go publishObservable.NotifyObservers(user.Uid, model.TypeBook, book.Id)
	}

	return
}

// FindBy 获取图书列表（分页）
func (GoBookLogic) FindBy(ctx context.Context, limit int, lastIds ...int) []*model.Book {
	objLog := GetLogger(ctx)

	dbSession := MasterDB.OrderBy("id DESC")

	if len(lastIds) > 0 && lastIds[0] > 0 {
		dbSession.And("id<?", lastIds[0])
	}

	books := make([]*model.Book, 0)
	err := dbSession.OrderBy("id DESC").Limit(limit).Find(&books)
	if err != nil {
		objLog.Errorln("GoBookLogic FindBy Error:", err)
		return nil
	}

	return books
}

// FindAll 支持多页翻看
func (GoBookLogic) FindAll(ctx context.Context, paginator *Paginator, orderBy string) []*model.Book {
	objLog := GetLogger(ctx)

	bookList := make([]*model.Book, 0)
	err := MasterDB.OrderBy(orderBy).Limit(paginator.PerPage(), paginator.Offset()).Find(&bookList)
	if err != nil {
		objLog.Errorln("GoBookLogic FindAll error:", err)
		return nil
	}

	return bookList
}

func (GoBookLogic) Count(ctx context.Context) int64 {
	objLog := GetLogger(ctx)

	var (
		total int64
		err   error
	)
	total, err = MasterDB.Count(new(model.Book))

	if err != nil {
		objLog.Errorln("GoBookLogic Count error:", err)
	}

	return total
}

// FindByIds 获取多个图书详细信息
func (GoBookLogic) FindByIds(ids []int) []*model.Book {
	if len(ids) == 0 {
		return nil
	}
	books := make([]*model.Book, 0)
	err := MasterDB.In("id", ids).Find(&books)
	if err != nil {
		logger.Errorln("GoBookLogic FindByIds error:", err)
		return nil
	}
	return books
}

// findByIds 获取多个图书详细信息 包内使用
func (GoBookLogic) findByIds(ids []int) map[int]*model.Book {
	if len(ids) == 0 {
		return nil
	}

	books := make(map[int]*model.Book)
	err := MasterDB.In("id", ids).Find(&books)
	if err != nil {
		logger.Errorln("GoBookLogic findByIds error:", err)
		return nil
	}
	return books
}

// FindById 获取一本图书信息
func (GoBookLogic) FindById(ctx context.Context, id interface{}) (*model.Book, error) {
	book := &model.Book{}
	_, err := MasterDB.Id(id).Get(book)
	if err != nil {
		logger.Errorln("book logic FindById Error:", err)
	}

	return book, err
}

// Total 图书总数
func (GoBookLogic) Total() int64 {
	total, err := MasterDB.Count(new(model.Book))
	if err != nil {
		logger.Errorln("GoBookLogic Total error:", err)
	}
	return total
}

// 图书评论
type BookComment struct{}

// UpdateComment 更新该图书的评论信息
// cid：评论id；objid：被评论对象id；uid：评论者；cmttime：评论时间
func (self BookComment) UpdateComment(cid, objid, uid int, cmttime time.Time) {
	// 更新评论数（TODO：暂时每次都更新表）
	_, err := MasterDB.Table(new(model.Book)).Id(objid).Incr("cmtnum", 1).Update(map[string]interface{}{
		"lastreplyuid":  uid,
		"lastreplytime": cmttime,
	})
	if err != nil {
		logger.Errorln("更新图书评论数失败：", err)
	}
}

func (self BookComment) String() string {
	return "book"
}

// SetObjinfo 实现 CommentObjecter 接口
func (self BookComment) SetObjinfo(ids []int, commentMap map[int][]*model.Comment) {
	books := DefaultGoBook.FindByIds(ids)
	if len(books) == 0 {
		return
	}

	for _, book := range books {
		objinfo := make(map[string]interface{})
		objinfo["title"] = book.Name
		objinfo["uri"] = model.PathUrlMap[model.TypeBook]
		objinfo["type_name"] = model.TypeNameMap[model.TypeBook]

		for _, comment := range commentMap[book.Id] {
			comment.Objinfo = objinfo
		}
	}
}

// 图书推荐
type BookLike struct{}

// 更新该图书的推荐数
// objid：被喜欢对象id；num: 喜欢数(负数表示取消喜欢)
func (self BookLike) UpdateLike(objid, num int) {
	// 更新喜欢数（TODO：暂时每次都更新表）
	_, err := MasterDB.Where("id=?", objid).Incr("likenum", num).Update(new(model.Book))
	if err != nil {
		logger.Errorln("更新图书喜欢数失败：", err)
	}
}

func (self BookLike) String() string {
	return "book"
}
