package main

import (
	"gitee.com/itsos/golibs/v2/db/mysql"
	"gitee.com/itsos/studynotes/datamodels"
	"gitee.com/itsos/studynotes/services"
	"github.com/kataras/golog"
	"sync"
)

var wg sync.WaitGroup

func main() {
	db := mysql.NewMysql()
	article := make([]datamodels.Article, 0)
	err := db.Where("is_del=?", 1).Desc("utime").Find(&article)
	if err != nil {
		panic(err)
	}

	content := datamodels.ArticleContent{}

	esData := services.EsData{}
	for _, v := range article {
		db.Where("aid=?", v.Id).Find(&content)
		esData.Id = v.Id
		esData.Title = v.Title
		esData.Intro = v.Intro
		esData.IsState = v.IsState
		esData.Utime = v.Utime
		esData.Ctime = v.Ctime
		esData.IsDel = v.IsDel
		esData.Data = content.Data
		wg.Add(1)
		go services.EsSync(esData, &wg)
	}
	wg.Wait()
	golog.Info("Done")
}