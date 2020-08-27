// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package consistenthash

import (
	"hash/crc32"
	"math"
	"strconv"
	"testing"

	"github.com/q191201771/naza/pkg/assert"
	"github.com/q191201771/naza/pkg/nazalog"
)

func TestConsistentHash(t *testing.T) {
	ch := New(1024)
	_, err := ch.Get("aaa")
	assert.Equal(t, ErrIsEmpty, err)

	ch.Add("127.0.0.1")
	ch.Add("0.0.0.0", "8.8.8.8")
	ch.Del("127.0.0.1", "8.8.8.8")
	ch.Add("114.114.114.114", "255.255.255.255", "1.1.1.1", "2.2.2.2", "3.3.3.3")
	// 经过上面模拟的节点扩容和缩容(宕机)，剩下的节点应该如下：
	exptectedNodes := []string{
		"0.0.0.0",
		"114.114.114.114",
		"255.255.255.255",
		"1.1.1.1",
		"2.2.2.2",
		"3.3.3.3",
	}
	actualNodes := ch.Nodes()
	assert.Equal(t, len(exptectedNodes), len(actualNodes))
	for _, en := range exptectedNodes {
		_, ok := actualNodes[en]
		assert.Equal(t, true, ok)
	}
	// consistenthash: map[0.0.0.0:627880983 1.1.1.1:829300198 114.114.114.114:564406193 2.2.2.2:875872728 255.255.255.255:636425375 3.3.3.3:761081819]
	// 每个节点的总区间范围，相加结果为：4,294,967,296 (即2^32)
	nazalog.Debugf("consistenthash: %+v", actualNodes)

	counts := make(map[string]int)
	// 从一致性哈希里获取真实节点，看是否大致均衡，每个节点获取的次数相加，应该和总次数相等
	// 测试情况：map[0.0.0.0:2397 1.1.1.1:2907 114.114.114.114:2062 2.2.2.2:3412 255.255.255.255:2320 3.3.3.3:3286]
	// 总次数相加即为16384
	for i := 0; i < 16384; i++ {
		node, err := ch.Get(strconv.Itoa(i))
		assert.Equal(t, nil, err)
		counts[node]++
	}
	nazalog.Debugf("%+v", counts)
}

func TestConsistentHash_Nodes(t *testing.T) {
	nodesGolden := []string{
		"0.0.0.0",
		"114.114.114.114",
		"255.255.255.255",
		"1.1.1.1",
		"2.2.2.2",
		"3.3.3.3",
	}
	// 每个真实节点对应虚拟节点数，对应虚拟节点越多，分布越均衡
	j := 5  // 1，起始测试
	k := 10 // 16384，2^14，结束
	// 下面循环每次*2步进，也可以以其他步长
	for i := j; i <= k; i = i << 1 {
		nazalog.Debugf("-----%d-----", i)
		// 每个真实节点对应多少个虚拟节点
		ch := New(i)
		// 节点(大致)均衡分布到哈希环上
		ch.Add(nodesGolden...)
		nodes := ch.Nodes()
		// 所有节点的范围之和，多次验证是否为2^32
		count := uint64(0)
		for k, v := range nodes {
			count += uint64(v)
			// 该真实节点总范围占的比例
			nazalog.Debugf("%s: %+v", k, float32(v)/float32(math.MaxUint32+1))
		}
		// 断言，保证总和是2^32
		assert.Equal(t, uint64(math.MaxUint32+1), count)
	}
}

func TestCorner(t *testing.T) {
	// 简单测试
	ch := New(1, func(option *Option) {
		option.hfn = crc32.ChecksumIEEE
	})

	ch = New(1)
	nodes := ch.Nodes()
	assert.Equal(t, nil, nodes)

	ch = New(1)
	ch.Add("127.0.0.1")
	nodes = ch.Nodes()
	exptectedNodes := map[string]uint64{
		"127.0.0.1": math.MaxUint32 + 1,
	}
	assert.Equal(t, exptectedNodes, nodes)
}
