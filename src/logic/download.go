// Copyright 2018 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"net/http"
	"strings"

	"model"

	"golang.org/x/net/context"

	. "db"

	"github.com/PuerkitoBio/goquery"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
)

type DownloadLogic struct{}

var DefaultDownload = DownloadLogic{}

func (DownloadLogic) FindAll(ctx context.Context) []*model.Download {
	downloads := make([]*model.Download, 0)
	err := MasterDB.Desc("seq").Find(&downloads)
	if err != nil {
		logger.Errorln("DownloadLogic FindAll Error:", err)
	}

	return downloads
}

func (DownloadLogic) RecordDLTimes(ctx context.Context, filename string) error {
	MasterDB.Where("filename=?", filename).Incr("times", 1).Update(new(model.Download))

	return nil
}

func (DownloadLogic) AddNewDownload(ctx context.Context, version, selector string) error {
	objLog := GetLogger(ctx)

	resp, err := http.Get("https://golang.google.cn/dl/")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	doc.Find(selector).Each(func(i int, versionSel *goquery.Selection) {
		idVal, exists := versionSel.Attr("id")
		if !exists {
			return
		}

		if idVal != version {
			return
		}

		versionSel.Find("table tbody tr").Each(func(j int, dlSel *goquery.Selection) {
			download := &model.Download{
				Version: version,
			}

			if dlSel.HasClass("highlight") {
				download.IsRecommend = true
			}

			dlSel.Find("td").Each(func(k int, fieldSel *goquery.Selection) {
				val := fieldSel.Text()
				switch k {
				case 0:
					download.Filename = val
				case 1:
					download.Kind = val
				case 2:
					download.OS = val
				case 3:
					download.Arch = val
				case 4:
					download.Size = goutils.MustInt(strings.TrimRight(val, "MB"))
				case 5:
					download.Checksum = val
				}
			})

			has, err := MasterDB.Where("filename=?", download.Filename).Exist(new(model.Download))
			if err != nil || has {
				return
			}

			_, err = MasterDB.Insert(download)
			if err != nil {
				objLog.Errorln("insert download error:", err, "version:", version)
			}
		})

		MasterDB.Exec("UPDATE download SET seq=id WHERE seq=0")
	})

	return nil
}
