/*
   Copyright (c) [2021] IT.SOS
   kn is licensed under Mulan PSL v2.
   You can use this software according to the terms and conditions of the Mulan PSL v2.
   You may obtain a copy of Mulan PSL v2 at:
            http://license.coscl.org.cn/MulanPSL2
   THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
   See the Mulan PSL v2 for more details.
*/

package controllers

import (
	"github.com/it-sos/blog/services"
	"github.com/kataras/iris/v12"
	"strconv"
	"time"
)

type IndexController struct {
	Ctx       iris.Context
	StartTime time.Time
}

// GetArticleList
// @Tags 博客前台接口
// @Summary 首页文章列表
// @Description 分页接口通过page递增获取更多
// @Accept json
// @Produce json
// @Param keyword query string false "关键词"
// @Param page query integer true "页码"
// @Param size query integer true "每页条数"
// @Success 200 {object} []vo.ArticleVO "列表数据"
// @Failure 400 {object} cerrors.Errors "error"
// @Router /article/list [get]
func (c *IndexController) GetArticleList() interface{} {
	page, _ := strconv.Atoi(c.Ctx.FormValue("page"))
	size, _ := strconv.Atoi(c.Ctx.FormValue("size"))
	keyword := c.Ctx.FormValue("keyword")
	isLogin := services.SAuthService.IsLogin(c.Ctx.GetHeader("token"))
	return services.SArticle.GetListPage(isLogin, page, size, keyword)
}

// GetArticleRank
// @Tags 博客前台接口
// @Summary 文章访问TOP50
// @Description 根据访问量排序
// @Accept json
// @Produce json
// @Success 200 {object} []vo.ArticleAccessTimesVO "列表数据"
// @Failure 400 {object} cerrors.Errors "error"
// @Router /article/rank [get]
func (c *IndexController) GetArticleRank() interface{} {
	isLogin := services.SAuthService.IsLogin(c.Ctx.GetHeader("token"))
	return services.SArticle.GetRank(isLogin)
}

// GetArticleContent
// @Tags 博客前台接口
// @Summary 文章内容详情
// @Description 文章内容详情
// @Accept json
// @Produce json
// @Param title query string true "文章标题"
// @Success 200 {object} []vo.ArticleContentVO "详情数据"
// @Failure 400 {object} cerrors.Errors "error"
// @Router /article/content [get]
func (c *IndexController) GetArticleContent() interface{} {
	title := c.Ctx.FormValue("title")
	isLogin := services.SAuthService.IsLogin(c.Ctx.GetHeader("token"))
	return services.SArticle.GetContent(isLogin, title)
}
