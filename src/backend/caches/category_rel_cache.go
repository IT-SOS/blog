/*
   Copyright (c) [2021] IT.SOS
   kn is licensed under Mulan PSL v2.
   You can use this software according to the terms and conditions of the Mulan PSL v2.
   You may obtain a copy of Mulan PSL v2 at:
            http://license.coscl.org.cn/MulanPSL2
   THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
   See the Mulan PSL v2 for more details.
*/

package caches

import (
	"fmt"
	"gitee.com/itsos/golibs/v2/db/redis"
	"golang.org/x/net/context"
)

type CategoryRel interface {
	Id(aid uint, types uint8) CategoryRelCmd
}

type categoryRel struct{}

type CategoryRelCmd interface {
	Get() []string
	Add(v uint)
	Exists() bool
	Remove(v uint) bool
}

type categoryRelCmd struct {
	aidType string
	db      redis.GoLibRedis
}

func (a *categoryRelCmd) Remove(v uint) bool {
	return a.db.SRem(context.Background(), a.aidType, v).Val() > 0
}

func (a *categoryRelCmd) Exists() bool {
	return a.db.Exists(context.Background(), a.aidType).Val() > 0
}

func (a *categoryRelCmd) Get() []string {
	return a.db.SMembers(context.Background(), a.aidType).Val()
}

func (a *categoryRelCmd) Add(v uint) {
	_, err := a.db.SAdd(context.Background(), a.aidType, v).Result()
	if err != nil {
		panic(err)
	}
}

const categoryRelPrefix = "category:%d_%d"

func (a *categoryRel) Id(aid uint, types uint8) CategoryRelCmd {
	return &categoryRelCmd{aidType: fmt.Sprintf(categoryRelPrefix, aid, types), db: redis.NewRedis()}
}

// CCategoryRel cache文章的标题和标签
var CCategoryRel CategoryRel = &categoryRel{}
