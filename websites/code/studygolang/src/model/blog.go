package model

import (
	"logger"
	"util"
)

// wordpress文章信息
type Article struct {
	Id          int    `json:"ID"`
	PostTitle   string `json:"post_title"`
	PostContent string `json:"post_content"`
	PostStatus  string `json:"post_status"` // 只查=publish的
	PostName    string `json:"post_name"`   // 链接后缀
	PostDate    string `json:"post_date"`

	// 链接
	PostUri string
	// 数据库访问对象
	*Dao
}

func NewArticle() *Article {
	return &Article{
		Dao: &Dao{tablename: "go_posts"},
	}
}

func (this *Article) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *Article) FindAll(selectCol ...string) ([]*Article, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	articleList := make([]*Article, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		article := NewArticle()
		err = this.Scan(rows, colNum, article.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("Article FindAll Scan Error:", err)
			continue
		}
		articleList = append(articleList, article)
	}
	return articleList, nil
}

func (this *Article) Where(condition string) *Article {
	this.Dao.Where(condition)
	return this
}

// 为了支持连写
func (this *Article) Limit(limit string) *Article {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *Article) Order(order string) *Article {
	this.Dao.Order(order)
	return this
}

func (this *Article) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"ID":           &this.Id,
		"post_title":   &this.PostTitle,
		"post_content": &this.PostContent,
		"post_status":  &this.PostStatus,
		"post_name":    &this.PostName,
		"post_date":    &this.PostDate,
	}
}
