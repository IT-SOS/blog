/*
   Copyright (c) [2021] IT.SOS
   kn is licensed under Mulan PSL v2.
   You can use this software according to the terms and conditions of the Mulan PSL v2.
   You may obtain a copy of Mulan PSL v2 at:
            http://license.coscl.org.cn/MulanPSL2
   THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
   See the Mulan PSL v2 for more details.
*/

package main

import (
	"fmt"
	"gitee.com/itsos/golibs/config"
	"gitee.com/itsos/golibs/global/consts"
	C "gitee.com/itsos/studynotes/config"
	_ "gitee.com/itsos/studynotes/docs"
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12"
)

var c = C.C

func main() {
	config.Init()
	app := iris.New()
	var url string
	if c.GetActive() == consts.EnvProduct {
		return
	}
	if c.GetActive() == consts.EnvTest {
		url = c.GetSwaggerUrl()
	} else {
		url = fmt.Sprintf("%s://localhost:%s", c.GetScheme(), c.GetSwaggerPort())
	}
	swaggerUrl := swagger.URL(url + "/swagger/doc.json")
	app.Get("/swagger/{any:path}", swagger.WrapHandler(swaggerFiles.Handler, swaggerUrl))

	fmt.Printf("Swagger on: %s/swagger/index.html\n", c.GetSwaggerUrl())
	fmt.Printf("Swagger on: http://localhost:%s/swagger/index.html\n", c.GetSwaggerPort())
	app.Listen(":" + c.GetSwaggerPort())
}
