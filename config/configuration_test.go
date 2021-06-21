/*
   Copyright (c) [2021] itsos
   kn is licensed under Mulan PSL v2.
   You can use this software according to the terms and conditions of the Mulan PSL v2.
   You may obtain a copy of Mulan PSL v2 at:
               http://license.coscl.org.cn/MulanPSL2
   THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
   See the Mulan PSL v2 for more details.
*/

package config

import (
	_ "gitee.com/itsos/golibs/tests"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Log(C.GetUrl())
	t.Log(C.GetSwaggerUrl())
	t.Log(C.GetQBuffer())
	t.Log(C.GetQHeight())
	t.Log(C.GetQWidth())
	t.Log(C.GetQUrl())
}
