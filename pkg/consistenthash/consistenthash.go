// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package consistenthash

import (
	"errors"
	"hash/crc32"
	"math"
	"sort"
	"strconv"
)

var ErrIsEmpty = errors.New("naza.consistenthash: is empty")

type ConsistentHash interface {
	Add(nodes ...string)
	Del(nodes ...string)
	Get(key string) (node string, err error)

	// @return: 返回的 map 的
	//          key 为添加到内部的 node，
	//          value 为该 node 在环上所占的 point 个数。
	//          我们可以通过各个 node 对应的 point 个数是否接近，来判断各 node 在环上的分布是否均衡。
	//          map 的所有 value 加起来应该等于 (math.MaxUint32 + 1) (math.MaxUint32为：1<<32 - 1，所以此处总和为：2^32)
	Nodes() map[string]uint64
}

type HashFunc func([]byte) uint32

type Option struct {
	hfn HashFunc
}

var defaultOption = Option{
	// 默认哈希算法指定为：循环冗余校验
	hfn: crc32.ChecksumIEEE,
}

type ModOption func(option *Option)

// @param dups: 每个实际的 node 转变成多少个环上的节点(每个node对应虚拟节点的数量)，必须大于等于1
// @param modOptions: 可修改内部的哈希函数，比如替换成murmur32的开源实现，可以这样：
//   import "github.com/spaolacci/murmur3"
//   import "github.com/q191201771/naza/pkg/consistenthash"
//
//   ch := consistenthash.New(1000, func(option *Option) {
//     option.hfn = func(bytes []byte) uint32 {
//       h := murmur3.New32()
//       h.Write(bytes)
//       return h.Sum32()
//     }
//   })
// ConsistentHash 是interface，面向接口编程，返回的 type consistentHash struct 对于包外不可见，其实现了ConsistentHash interface
func New(dups int, modOptions ...ModOption) ConsistentHash {
	option := defaultOption
	for _, fn := range modOptions {
		fn(&option)
	}

	return &consistentHash{
		point2node: make(map[uint32]string),
		dups:       dups,
		option:     option,
	}
}

type consistentHash struct {
	point2node map[uint32]string // 虚拟节点(哈希值)和对应原始节点名
	points     []uint32          // 所有虚拟节点
	dups       int               // 哈希环上每个真实节点对应的虚拟节点个数
	option     Option            // 可选，可设置哈希函数，默认CRC32
}

// 可以添加多个节点(最后一个参数，在类型前加`...`，表示可变参数可能为0个或多个)
func (ch *consistentHash) Add(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < ch.dups; i++ {
			// hash2point(struct绑定的method) 通过哈希函数hfn把string计算为一个uint32数值，hfn默认是CRC32算法
			// virtualKey 把两个入参拼接成一个string，如 node0
			point := ch.hash2point(virtualKey(node, i))
			// uint32作为key(若同样的节点名Add进来后覆盖原来位置？ 覆盖后下面的slice还是会新增虚拟节点成员，导致超过dups)
			// 所以Add的时候不能添加以前有的node？
			// 加了重复哈希也没什么影响，算区间的时候用本虚拟节点减前一个节点，区间为0；且删除时也会重新生成一下全量虚拟节点slice
			ch.point2node[point] = node
			// 算好的uint32即作为一个虚拟节点，保存到总的slice里(可能会有相等的uint32)
			ch.points = append(ch.points, point)
		}
	}
	// 哈希后计算出的uint32从小到大排序
	sortSlice(ch.points)
}

func (ch *consistentHash) Del(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < ch.dups; i++ {
			// 用Add时同样的算法算出哈希值(CRC32)
			point := ch.hash2point(virtualKey(node, i))
			// 根据key删除map成员，总的 points slice在下面会更新
			delete(ch.point2node, point)
		}
	}

	// 每次都根据map的所有key重新生成 points slice，即便slice里有相同的哈希值也一并去掉了
	ch.points = nil
	for k := range ch.point2node {
		ch.points = append(ch.points, k)
	}
	// 从小到大排序
	sortSlice(ch.points)
}

// 先根据传入的key算哈希从全部虚拟节点里拿一个就近的虚拟节点(均衡)，返回其对应的真实节点
func (ch *consistentHash) Get(key string) (node string, err error) {
	if len(ch.points) == 0 {
		return "", ErrIsEmpty
	}

	// 根据传入的key(string类型)计算哈希，此处没有加数字后缀(Add节点时自动添加了数字后缀算哈希，如node1/node2等形式)
	point := ch.hash2point(key)
	// 从数组中找出满足 point 值 >= key 所对应 point 值的最小的元素
	// sort.Search 函数利用二分查找，传入一个函数作为参数，若数据为升序则入参函数送>=目标值时为true，可找到第一个>=目标值的索引下标
	// (降序则<=为true，目标值是否存在则==目标值为true)
	index := sort.Search(len(ch.points), func(i int) bool {
		// 第一个>=虚拟节点的位置，到前一个节点的区间，为真实节点的一部分区间
		return ch.points[i] >= point
	})

	// 找不到key对应节点则取第一个(这种情况只会出现在最后一个虚拟节点和第一个之间)
	if index == len(ch.points) {
		index = 0
	}

	// 从map里取节点名，找不到对应key则取slice里第一个虚拟节点对应的节点名
	return ch.point2node[ch.points[index]], nil
}

// 返回真实节点各自的总区间范围，瓜分2^32
func (ch *consistentHash) Nodes() map[string]uint64 {
	if len(ch.points) == 0 {
		return nil
	}
	ret := make(map[string]uint64)
	prev := uint64(0)
	// 从小到大遍历虚拟节点，看节点之间的范围归属到哪个实际节点(0-第一个虚拟 和 最后一个到0(2^32)，都归到第一个节点)
	for _, point := range ch.points {
		// 虚拟节点对应的原始节点名
		node := ch.point2node[point]
		// 本次虚拟节点(哈希值)到前一个虚拟节点之间的一段，都算做本虚拟节点的范围，数字总大小即表示这个真实节点在哈希环上的范围
		ret[node] = ret[node] + uint64(point) - prev
		prev = uint64(point)
	}

	// 最后一个 node 到终点位置的 point 都归入第一个 node
	point := ch.points[len(ch.points)-1]
	node := ch.point2node[point]
	ret[node] = ret[node] + uint64(math.MaxUint32-point+1)
	return ret
}

func (ch *consistentHash) hash2point(key string) uint32 {
	return ch.option.hfn([]byte(key))
}

func virtualKey(node string, index int) string {
	return node + strconv.Itoa(index)
}

func sortSlice(a []uint32) {
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
}
