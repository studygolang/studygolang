// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"errors"
	"fmt"
	"model"
	"os"
	"regexp"
	"strconv"
	"time"
	"util"

	"github.com/gorilla/schema"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/nosql"
	"github.com/polaris1119/times"
	"golang.org/x/net/context"
)

var schemaDecoder = schema.NewDecoder()

func init() {
	schemaDecoder.SetAliasTag("json")
	schemaDecoder.IgnoreUnknownKeys(true)
}

var (
	NotModifyAuthorityErr = errors.New("没有修改权限")
	NotFoundErr           = errors.New("Not Found")
)

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

// parseAtUser 解析 @某人
func parseAtUser(ctx context.Context, content string) string {
	reg := regexp.MustCompile(`@([^\s@]{4,20})`)
	return reg.ReplaceAllStringFunc(content, func(matched string) string {
		username := matched[1:]

		// 校验 username 是否存在
		user := DefaultUser.FindOne(ctx, "username", username)
		if user.Username != username {
			return matched
		}
		return fmt.Sprintf(`<a href="/user/%s" title="%s">%s</a>`, username, matched, matched)
	})
}

// CanEdit 判断能否编辑
func CanEdit(me *model.Me, curModel interface{}) bool {
	if me == nil {
		return false
	}

	canEditTime := time.Duration(UserSetting[model.KeyCanEditTime]) * time.Second
	switch entity := curModel.(type) {
	case *model.Topic:
		if me.Uid != entity.Uid && me.IsAdmin && roleCanEdit(model.TopicAdmin, me) {
			return true
		}

		if time.Now().Sub(time.Time(entity.Ctime)) > canEditTime {
			return false
		}

		if me.Uid == entity.Uid {
			return true
		}
	case *model.Article:
		if me.IsAdmin && roleCanEdit(model.ArticleAdmin, me) {
			return true
		}

		// 文章的能编辑时间是15天
		if time.Now().Sub(time.Time(entity.Ctime)) > 15*86400*time.Second {
			return false
		}

		if me.Username == entity.Author {
			return true
		}
	case *model.Resource:
		if time.Now().Sub(time.Time(entity.Ctime)) > canEditTime {
			return false
		}

		if me.Uid == entity.Uid {
			return true
		}
	case *model.OpenProject:
		if me.IsAdmin && roleCanEdit(model.Administrator, me) {
			return true
		}

		// 开源项目的能编辑时间是30天
		if time.Now().Sub(time.Time(entity.Ctime)) > 30*86400*time.Second {
			return false
		}

		if me.Username == entity.Username {
			return true
		}
	case *model.Wiki:
		if me.IsAdmin && roleCanEdit(model.Administrator, me) {
			return true
		}
		if time.Now().Sub(time.Time(entity.Ctime)) > canEditTime {
			return false
		}

		if me.Uid == entity.Uid {
			return true
		}
	case *model.Book:
		if me.IsAdmin && roleCanEdit(model.Administrator, me) {
			return true
		}
		if time.Now().Sub(time.Time(entity.CreatedAt)) > canEditTime {
			return false
		}

		if me.Uid == entity.Uid {
			return true
		}
	case *model.Comment:
		if me.IsAdmin && roleCanEdit(model.Administrator, me) {
			return true
		}
		if time.Now().Sub(time.Time(entity.Ctime)) > canEditTime {
			return false
		}

		if me.Uid == entity.Uid {
			return true
		}
	case map[string]interface{}:
		if adminCanEdit(entity, me) {
			return true
		}
		if ctime, ok := entity["ctime"]; ok {
			if time.Now().Sub(time.Time(ctime.(model.OftenTime))) > canEditTime {
				return false
			}
		}

		if createdAt, ok := entity["created_at"]; ok {
			if time.Now().Sub(time.Time(createdAt.(model.OftenTime))) > canEditTime {
				return false
			}
		}

		if uid, ok := entity["uid"]; ok {
			if me.Uid == uid.(int) {
				return true
			}
		}

		if username, ok := entity["username"]; ok {
			if me.Username == username.(string) {
				return true
			}
		}
	}

	return false
}

func CanPublish(dauAuth, objtype int) bool {
	if dauAuth == 0 {
		return true
	}

	switch objtype {
	case model.TypeTopic:
		return (dauAuth & model.DauAuthTopic) == model.DauAuthTopic
	case model.TypeArticle:
		return (dauAuth & model.DauAuthArticle) == model.DauAuthArticle
	case model.TypeResource:
		return (dauAuth & model.DauAuthResource) == model.DauAuthResource
	case model.TypeProject:
		return (dauAuth & model.DauAuthProject) == model.DauAuthProject
	case model.TypeWiki:
		return (dauAuth & model.DauAuthWiki) == model.DauAuthWiki
	case model.TypeBook:
		return (dauAuth & model.DauAuthBook) == model.DauAuthBook
	case model.TypeComment:
		return (dauAuth & model.DauAuthComment) == model.DauAuthComment
	case model.TypeTop:
		return (dauAuth & model.DauAuthTop) == model.DauAuthTop
	default:
		return true
	}
}

// NeedCaptcha 是否需要验证码：
//  - 新客注册后一段时间内需要
//  - 发布内容太频繁（一天次数太多、间隔太快）
func NeedCaptcha(user *model.Me) bool {
	// 注册后 30 分钟内发布需要验证码
	if user.CreatedAt.Add(30 * time.Minute).After(time.Now()) {
		return true
	}

	// 发布内容是否太频繁
	redis := nosql.NewRedisFromPool()
	defer redis.Close()

	publishTimes := redis.GET(getPublishTimesKey(user.Uid))
	if goutils.MustInt(publishTimes) > UserSetting[model.KeyPublishTimes] {
		return true
	}

	lastTimestampStr := redis.GET(getLastPublishTimeKey(user.Uid))
	lastTimestamp := goutils.MustInt64(lastTimestampStr)
	if time.Now().Unix()-lastTimestamp < int64(UserSetting[model.KeyPublishInterval]) {
		return true
	}

	return false
}

// incrPublishTimes 增加用户发布次数
func incrPublishTimes(uid int) {
	redis := nosql.NewRedisFromPool()
	defer redis.Close()

	key := getPublishTimesKey(uid)
	redis.INCR(key)
	redis.EXPIRE(key, 86401)
}

// recordLastPubishTime 记录用户上次发布时间
func recordLastPubishTime(uid int) {
	redis := nosql.NewRedisFromPool()
	defer redis.Close()

	key := getLastPublishTimeKey(uid)
	redis.SET(key, time.Now().Unix(), 86400)
}

func getPublishTimesKey(uid int) string {
	return "publish:times:user:" + strconv.Itoa(uid) + ":date:" + times.Format("Ymd")
}

func getLastPublishTimeKey(uid int) string {
	return "last:publish:time:user:" + strconv.Itoa(uid)
}

func website() string {
	host := "http://"
	if WebsiteSetting.OnlyHttps {
		host = "https://"
	}
	return host + WebsiteSetting.Domain
}

func adminCanEdit(entity map[string]interface{}, me *model.Me) bool {
	if uid, ok := entity["uid"]; ok {
		if me.Uid != uid.(int) && me.IsAdmin && roleCanEdit(model.Administrator, me) {
			return true
		}
		return false
	}

	if username, ok := entity["username"]; ok {
		if me.Username != username.(string) && me.IsAdmin && roleCanEdit(model.Administrator, me) {
			return true
		}
		return false
	}

	return false
}

func roleCanEdit(typRoleID int, me *model.Me) bool {
	if me.IsRoot {
		return true
	}

	if util.InSlice(typRoleID, me.RoleIds) {
		return true
	}

	for _, roleID := range me.RoleIds {
		if roleID <= model.Administrator {
			return true
		}
	}

	return false
}
