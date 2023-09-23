// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/studygolang/studygolang/internal/model"
	"github.com/studygolang/studygolang/util"

	. "github.com/studygolang/studygolang/db"

	"github.com/tidwall/gjson"

	"golang.org/x/net/context"

	"github.com/polaris1119/config"
	"github.com/polaris1119/nosql"
)

type WechatLogic struct{}

var DefaultWechat = WechatLogic{}

var jscodeRUL = "https://api.weixin.qq.com/sns/jscode2session"

// CheckSession 微信小程序登录凭证校验
func (self WechatLogic) CheckSession(ctx context.Context, code string) (*model.WechatUser, error) {
	objLog := GetLogger(ctx)

	appid := config.ConfigFile.MustValue("wechat.xcx", "appid")
	appsecret := config.ConfigFile.MustValue("wechat.xcx", "appsecret")

	checkLoginURL := fmt.Sprintf("%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		jscodeRUL, appid, appsecret, code)

	body, err := util.DoGet(checkLoginURL)
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(body)

	openidResult := result.Get("openid")
	if !openidResult.Exists() {
		objLog.Errorln("WechatLogic WxLogin error:", result.Raw)
		return nil, errors.New(result.Get("errmsg").String())
	}

	openid := openidResult.String()
	wechatUser := &model.WechatUser{}
	_, err = MasterDB.Where("openid=?", openid).Get(wechatUser)
	if err != nil {
		objLog.Errorln("WechatLogic WxLogin find wechat user error:", err)
		return nil, err
	}

	if wechatUser.Id == 0 {
		wechatUser.Openid = openid
		wechatUser.SessionKey = result.Get("session_key").String()
		_, err = MasterDB.Insert(wechatUser)
		if err != nil {
			objLog.Errorln("WechatLogic WxLogin insert wechat user error:", err)
			return nil, err
		}
	}

	return wechatUser, nil
}

func (self WechatLogic) Bind(ctx context.Context, id, uid int, userInfo string) (*model.WechatUser, error) {
	objLog := GetLogger(ctx)

	result := gjson.Parse(userInfo)

	wechatUser := &model.WechatUser{
		Uid:      uid,
		Nickname: result.Get("nickName").String(),
		Avatar:   result.Get("avatarUrl").String(),
		OpenInfo: userInfo,
	}
	_, err := MasterDB.ID(id).Update(wechatUser)
	if err != nil {
		objLog.Errorln("WechatLogic Bind update error:", err)
		return nil, err
	}

	return wechatUser, nil
}

func (self WechatLogic) FetchOrUpdateToken() (string, error) {
	var result = struct {
		AccessToken string
		ExpiresTime time.Time
	}{}

	filename := config.ROOT + "/data/wechat-token.json"
	if util.Exist(filename) {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return "", err
		}

		err = json.Unmarshal(b, &result)
		if err != nil {
			return "", err
		}

		if result.ExpiresTime.After(time.Now()) {
			return result.AccessToken, nil
		}
	}

	appid := config.ConfigFile.MustValue("wechat", "appid")
	appsecret := config.ConfigFile.MustValue("wechat", "appsecret")
	strURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appid, appsecret)

	b, err := util.DoGet(strURL)
	if err != nil {
		return "", err
	}
	gresult := gjson.ParseBytes(b)
	if gresult.Get("errmsg").Exists() {
		return "", errors.New(gresult.Get("errmsg").String())
	}

	result.AccessToken = gresult.Get("access_token").String()
	result.ExpiresTime = time.Now().Add(time.Duration(gresult.Get("expires_in").Int()-5) * time.Second)

	b, err = json.Marshal(result)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(filename, b, 0755)
	if err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

func (self WechatLogic) AutoReply(ctx context.Context, reqData []byte) (*model.WechatReply, error) {
	objLog := GetLogger(ctx)

	wechatMsg := &model.WechatMsg{}
	err := xml.Unmarshal(reqData, wechatMsg)
	if err != nil {
		objLog.Errorln("wechat autoreply xml unmarshal error:", err)
		return nil, err
	}

	switch wechatMsg.MsgType {
	case model.WeMsgTypeText:
		if strings.Contains(wechatMsg.Content, "晨读") {
			return self.readingContent(ctx, wechatMsg)
		} else if strings.Contains(wechatMsg.Content, "主题") || strings.Contains(wechatMsg.Content, "帖子") {
			return self.topicContent(ctx, wechatMsg)
		} else if strings.Contains(wechatMsg.Content, "文章") {
			return self.articleContent(ctx, wechatMsg)
		} else if strings.Contains(wechatMsg.Content, "资源") {
			return self.resourceContent(ctx, wechatMsg)
		} else if strings.Contains(wechatMsg.Content, "项目") {
			return self.projectContent(ctx, wechatMsg)
		} else if strings.Contains(wechatMsg.Content, "图书") {
			return self.bookContent(ctx, wechatMsg)
		} else {
			// 用户获取验证码用
			user := DefaultUser.FindOne(ctx, "username", wechatMsg.Content)
			if user.Uid > 0 {
				var content string
				// 获取微信用户信息
				if err = self.checkAndSave(ctx, wechatMsg); err != nil {
					content = err.Error()
				} else {
					content = self.genCaptcha(user.Username, wechatMsg.FromUserName)
				}
				return self.wechatResponse(ctx, content, wechatMsg)
			}

			// 关键词回复
			autoReply := &model.WechatAutoReply{}
			MasterDB.Where("word LIKE ?", "%"+wechatMsg.Content+"%").Get(autoReply)
			if autoReply.Id != 0 {
				wechatMsg.MsgType = autoReply.MsgType
				return self.wechatResponse(ctx, autoReply.Content, wechatMsg)
			}

			return self.searchContent(ctx, wechatMsg)
		}
	case model.WeMsgTypeEvent:
		switch wechatMsg.Event {
		case model.WeEventSubscribe:
			wechatMsg.MsgType = model.WeMsgTypeText
			welcomeText := strings.ReplaceAll(config.ConfigFile.MustValue("wechat", "subscribe"), "\\n", "\n")

			autoReply := &model.WechatAutoReply{}
			_, err = MasterDB.Where("typ=?", model.AutoReplyTypSubscribe).Get(autoReply)
			if err == nil {
				welcomeText = autoReply.Content
			}

			return self.wechatResponse(ctx, welcomeText, wechatMsg)
		}
	}

	return self.wechatResponse(ctx, "success", wechatMsg)
}

func (self WechatLogic) genCaptcha(username, openid string) string {
	num := rand.Intn(9000) + 1000
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()

	captcha := strconv.Itoa(num)
	redisClient.SET("wechat:captcha:$username:"+username, captcha+openid, 600)

	return captcha
}

func (self WechatLogic) CheckCaptchaAndActivate(ctx context.Context, me *model.Me, captcha string) error {
	openid, err := self.checkCaptchaAndFetch(ctx, me, captcha)
	if err != nil {
		return err
	}

	session := MasterDB.NewSession()
	defer session.Close()

	session.Begin()
	_, err = session.Table(new(model.WechatUser)).Where("openid=?", openid).Update(map[string]interface{}{
		"uid": me.Uid,
	})
	if err != nil {
		session.Rollback()
		return err
	}

	_, err = session.Table(new(model.User)).ID(me.Uid).Update(map[string]interface{}{
		"status": model.UserStatusAudit,
		"ctime":  time.Now().Add(-5 * time.Hour),
	})
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

func (self WechatLogic) CheckCaptchaAndBind(ctx context.Context, me *model.Me, captcha string) error {
	openid, err := self.checkCaptchaAndFetch(ctx, me, captcha)
	if err != nil {
		return err
	}

	session := MasterDB.NewSession()
	defer session.Close()

	session.Begin()
	_, err = session.Table(new(model.WechatUser)).Where("openid=?", openid).Update(map[string]interface{}{
		"uid": me.Uid,
	})
	if err != nil {
		session.Rollback()
		return err
	}

	_, err = session.Table(new(model.User)).ID(me.Uid).Update(map[string]interface{}{
		"ctime": time.Now().Add(-5 * time.Hour),
	})
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

func (self WechatLogic) checkCaptchaAndFetch(ctx context.Context, me *model.Me, captcha string) (string, error) {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()

	key := "wechat:captcha:$username:" + me.Username
	store := redisClient.GET(key)
	if store[:4] != captcha {
		return "", errors.New("验证码错误")
	}

	redisClient.DEL(key)

	return store[4:], nil
}

func (self WechatLogic) checkAndSave(ctx context.Context, wechatMsg *model.WechatMsg) error {
	accessToken, err := self.FetchOrUpdateToken()
	if err != nil {
		return err
	}

	wechatUser := &model.WechatUser{}
	_, err = MasterDB.Where("openid=?", wechatMsg.FromUserName).Get(wechatUser)
	if err != nil {
		return err
	}

	strURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN", accessToken, wechatMsg.FromUserName)
	b, err := util.DoGet(strURL)
	if err != nil {
		return err
	}

	result := gjson.ParseBytes(b)
	if result.Get("errmsg").Exists() {
		return errors.New(result.Get("errmsg").String())
	}

	// 已经存在
	if wechatUser.Openid != "" {
		wechatUser.Nickname = result.Get("nickname").String()
		wechatUser.Avatar = result.Get("headimgurl").String()
		wechatUser.OpenInfo = result.Raw

		_, err = MasterDB.ID(wechatUser.Id).Update(wechatUser)
	} else {
		wechatUser = &model.WechatUser{
			Openid:   result.Get("openid").String(),
			Nickname: result.Get("nickname").String(),
			Avatar:   result.Get("headimgurl").String(),
			OpenInfo: result.Raw,
		}
		_, err = MasterDB.InsertOne(wechatUser)
	}

	if wechatUser.Uid > 0 {
		return errors.New("该微信绑定过其他账号")
	}

	return err
}

func (self WechatLogic) topicContent(ctx context.Context, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {

	topics := DefaultTopic.FindRecent(5)

	respContentSlice := make([]string, len(topics))
	for i, topic := range topics {
		respContentSlice[i] = fmt.Sprintf("%d.《%s》 %s/topics/%d", i+1, topic.Title, website(), topic.Tid)
	}

	return self.wechatResponse(ctx, strings.Join(respContentSlice, "\n"), wechatMsg)
}

func (self WechatLogic) articleContent(ctx context.Context, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {

	articles := DefaultArticle.FindBy(ctx, 5)

	respContentSlice := make([]string, len(articles))
	for i, article := range articles {
		respContentSlice[i] = fmt.Sprintf("%d.《%s》 %s/articles/%d", i+1, article.Title, website(), article.Id)
	}

	return self.wechatResponse(ctx, strings.Join(respContentSlice, "\n"), wechatMsg)
}

func (self WechatLogic) resourceContent(ctx context.Context, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {

	resources := DefaultResource.FindBy(ctx, 5)

	respContentSlice := make([]string, len(resources))
	for i, resource := range resources {
		respContentSlice[i] = fmt.Sprintf("%d.《%s》 %s/resources/%d", i+1, resource.Title, website(), resource.Id)
	}

	return self.wechatResponse(ctx, strings.Join(respContentSlice, "\n"), wechatMsg)
}

func (self WechatLogic) projectContent(ctx context.Context, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {

	projects := DefaultProject.FindBy(ctx, 5)

	respContentSlice := make([]string, len(projects))
	for i, project := range projects {
		respContentSlice[i] = fmt.Sprintf("%d.《%s%s》 %s/p/%d", i+1, project.Category, project.Name, website(), project.Id)
	}

	return self.wechatResponse(ctx, strings.Join(respContentSlice, "\n"), wechatMsg)
}

func (self WechatLogic) bookContent(ctx context.Context, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {

	books := DefaultGoBook.FindBy(ctx, 5)

	respContentSlice := make([]string, len(books))
	for i, book := range books {
		respContentSlice[i] = fmt.Sprintf("%d.《%s》 %s/book/%d", i+1, book.Name, website(), book.Id)
	}

	return self.wechatResponse(ctx, strings.Join(respContentSlice, "\n"), wechatMsg)
}

func (self WechatLogic) readingContent(ctx context.Context, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {

	var formatContent = func(reading *model.MorningReading) string {
		if reading.Inner == 0 {
			return fmt.Sprintf("%s\n%s", reading.Content, reading.Url)
		}

		return fmt.Sprintf("%s\n%s/articles/%d", reading.Content, website(), reading.Inner)
	}

	var readings []*model.MorningReading
	if wechatMsg.Content == "最新晨读" {
		readings = DefaultReading.FindBy(ctx, 1, model.RtypeGo)
		if len(readings) == 0 {
			return self.wechatResponse(ctx, config.ConfigFile.MustValue("wechat", "not_found"), wechatMsg)
		}

		return self.wechatResponse(ctx, formatContent(readings[0]), wechatMsg)
	}

	readings = DefaultReading.FindBy(ctx, 3, model.RtypeGo)

	respContentSlice := make([]string, len(readings))
	for i, reading := range readings {
		respContentSlice[i] = fmt.Sprintf("%d. %s", i+1, formatContent(reading))
	}

	return self.wechatResponse(ctx, strings.Join(respContentSlice, "\n\n"), wechatMsg)
}

func (self WechatLogic) searchContent(ctx context.Context, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {
	objLog := GetLogger(ctx)

	respBody, err := DefaultSearcher.SearchByField("title", wechatMsg.Content, 0, 5)
	if err != nil {
		objLog.Errorln("wechat search by field error:", err)
		return nil, err
	}

	if respBody.NumFound == 0 {
		return self.wechatResponse(ctx, config.ConfigFile.MustValue("wechat", "not_found"), wechatMsg)
	}

	host := WebsiteSetting.Domain
	if WebsiteSetting.OnlyHttps {
		host = "https://" + host
	} else {
		host = "http://" + host
	}

	respContentSlice := make([]string, len(respBody.Docs))
	for i, doc := range respBody.Docs {
		url := ""

		switch doc.Objtype {
		case model.TypeTopic:
			url = fmt.Sprintf("%s/topics/%d", host, doc.Objid)
		case model.TypeArticle:
			url = fmt.Sprintf("%s/articles/%d", host, doc.Objid)
		case model.TypeResource:
			url = fmt.Sprintf("%s/resources/%d", host, doc.Objid)
		case model.TypeProject:
			url = fmt.Sprintf("%s/p/%d", host, doc.Objid)
		case model.TypeWiki:
			url = fmt.Sprintf("%s/wiki/%d", host, doc.Objid)
		case model.TypeBook:
			url = fmt.Sprintf("%s/book/%d", host, doc.Objid)
		}
		respContentSlice[i] = fmt.Sprintf("%d.《%s》 %s", i+1, doc.Title, url)
	}

	return self.wechatResponse(ctx, strings.Join(respContentSlice, "\n"), wechatMsg)
}

func (self WechatLogic) wechatResponse(ctx context.Context, respContent string, wechatMsg *model.WechatMsg) (*model.WechatReply, error) {
	wechatReply := &model.WechatReply{
		ToUserName:   &model.CData{Val: wechatMsg.FromUserName},
		FromUserName: &model.CData{Val: wechatMsg.ToUserName},
		MsgType:      &model.CData{Val: wechatMsg.MsgType},
		CreateTime:   time.Now().Unix(),
	}
	switch wechatMsg.MsgType {
	case model.WeMsgTypeText:
		wechatReply.Content = &model.CData{Val: respContent}
	case model.WeMsgTypeImage:
		wechatReply.Image = &model.WechatImage{
			MediaId: &model.CData{Val: respContent},
		}
	default:
		wechatReply.Content = &model.CData{Val: config.ConfigFile.MustValue("wechat", "not_found")}
	}

	return wechatReply, nil
}
