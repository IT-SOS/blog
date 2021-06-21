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
	_ "gitee.com/itsos/golibs/tests/testsdb"
	"testing"
)

func TestTimesCache(t *testing.T) {
	NAuthTimes.Key("peng.yu").Incr()
	NAuthTimes.Key("peng.yu").Incr()
	NAuthTimes.Key("peng.yu").Incr()
	t.Log(NAuthTimes.Key("peng.yu").Incr())
	t.Log(NAuthTimes.Key("peng.yu").Get())
	NAuthTimes.Key("peng.yu").Clear()
	t.Log(NAuthTimes.Key("peng.yu").Decr())
}
