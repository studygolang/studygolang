package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	echo "github.com/labstack/echo/v4"
	. "github.com/studygolang/studygolang/db"
	"xorm.io/xorm"
	"xorm.io/xorm/core"
)

func TestWikiListPreviousReturnsAdjacentPage(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	engine, err := xorm.NewEngineWithDB("mysql", "", core.FromDB(sqlDB))
	if err != nil {
		t.Fatal(err)
	}
	oldDB := MasterDB
	MasterDB = engine
	t.Cleanup(func() {
		MasterDB = oldDB
		_ = engine.Close()
	})

	wikiRows := sqlmock.NewRows([]string{"id", "uid", "title", "uri"})
	for id := 41; id <= 61; id++ {
		wikiRows.AddRow(id, 1, fmt.Sprintf("wiki-%d", id), fmt.Sprintf("wiki-%d", id))
	}
	mock.ExpectQuery("SELECT .*wiki.*").WithArgs(40).WillReturnRows(wikiRows)
	mock.ExpectQuery("SELECT .*user_info.*").WithArgs(1).WillReturnRows(
		sqlmock.NewRows([]string{"uid", "username"}).AddRow(1, "polaris"),
	)

	recorder := httptest.NewRecorder()
	ctx := echo.New().NewContext(
		httptest.NewRequest(http.MethodGet, "/api/v1/wiki?lastid=40&direction=prev", nil),
		recorder,
	)
	if err := (WikiController{}).List(ctx); err != nil {
		t.Fatalf("List() error = %v", err)
	}

	var payload struct {
		Code int `json:"code"`
		Data struct {
			Wikis []struct {
				ID int `json:"id"`
			} `json:"wikis"`
			Page struct {
				HasPrev bool `json:"has_prev"`
				PrevID  int  `json:"prev_id"`
				HasNext bool `json:"has_next"`
				NextID  int  `json:"next_id"`
			} `json:"page"`
		} `json:"data"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.Code != 0 {
		t.Fatalf("code = %d, body = %s", payload.Code, recorder.Body.String())
	}
	if len(payload.Data.Wikis) != 20 {
		t.Fatalf("返回 %d 条，want 20", len(payload.Data.Wikis))
	}
	if got := payload.Data.Wikis[0].ID; got != 60 {
		t.Errorf("首条 ID = %d, want 60", got)
	}
	if got := payload.Data.Wikis[len(payload.Data.Wikis)-1].ID; got != 41 {
		t.Errorf("末条 ID = %d, want 41", got)
	}
	if !payload.Data.Page.HasPrev || payload.Data.Page.PrevID != 60 {
		t.Errorf("上一页游标 = %+v, want has_prev=true prev_id=60", payload.Data.Page)
	}
	if !payload.Data.Page.HasNext || payload.Data.Page.NextID != 41 {
		t.Errorf("下一页游标 = %+v, want has_next=true next_id=41", payload.Data.Page)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
