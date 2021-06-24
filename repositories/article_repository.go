/*
   Copyright (c) [2021] IT.SOS
   kn is licensed under Mulan PSL v2.
   You can use this software according to the terms and conditions of the Mulan PSL v2.
   You may obtain a copy of Mulan PSL v2 at:
            http://license.coscl.org.cn/MulanPSL2
   THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
   See the Mulan PSL v2 for more details.
*/

package repositories

import (
	"gitee.com/itsos/golibs/db"
	"gitee.com/itsos/studynotes/datamodels"
)

const (
	IS_STATE_PRIVATE = 1
	IS_STATE_PUBLIC  = 2
)

type ArticleRepository interface {
	// Insert 新增
	Insert(p *datamodels.Article) (id uint)
	// Update 更新
	Update(p *datamodels.Article) (id uint)
	// Select 查询文章详细
	Select(p *datamodels.Article) (datamodels.Article, bool)
	// Content 文章内容
	Content(id uint) datamodels.ArticleContent
	// SelectMany 查询文章列表
	SelectMany(state []int8, offset int, limit int) (results []datamodels.Article)
	SelectManyByIds(ids []int) []datamodels.Article
}

type articleRepository struct {
}

var err error

func NewArticleRepository() ArticleRepository {
	return &articleRepository{}
}

func (ur *articleRepository) Content(id uint) datamodels.ArticleContent {
	content := &datamodels.ArticleContent{Aid: id}
	_, err = db.Conn.Get(content)
	if err != nil {
		panic(err)
	}
	return *content
}

func (ur *articleRepository) Select(p *datamodels.Article) (datamodels.Article, bool) {
	has, err := db.Conn.Get(p)
	if err != nil {
		panic(err)
	}
	return *p, has
}

func (ur *articleRepository) SelectMany(state []int8, offset int, limit int) (results []datamodels.Article) {
	article := make([]datamodels.Article, 0)
	err = db.Conn.In("is_state", state).Desc("utime").Limit(limit, offset).Find(&article)
	if err != nil {
		panic(err)
	}
	return article
}

func (ur *articleRepository) SelectManyByIds(ids []int) []datamodels.Article {
	article := make([]datamodels.Article, 0)
	err := db.Conn.In("id", ids).Find(&article)
	if err != nil {
		panic(err)
	}
	return article
}

func (ur *articleRepository) Insert(p *datamodels.Article) (id uint) {
	_, err = db.Conn.Insert(p)
	if err != nil {
		panic(err)
	}
	return p.Id
}

func (ur *articleRepository) Update(p *datamodels.Article) (id uint) {
	_, err = db.Conn.ID(p.Id).Update(p)
	if err != nil {
		panic(err)
	}
	return p.Id
}
