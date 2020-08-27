// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package lru

// 标准库中的双向链表，实现为了一个循环链表，最后一个节点的Back()和第一个节点的Front()都是 &l.root (root是该链表中保存的哨兵节点)
import "container/list"

type LRU struct {
	c int                           // capacity
	m map[interface{}]*list.Element // mapping key -> index  用来快捷查看是否存在指定的key，存在则获取
	l *list.List                    // value
}

type pair struct {
	k interface{}
	v interface{}
}

func New(capacity int) *LRU {
	return &LRU{
		c: capacity,
		m: make(map[interface{}]*list.Element),
		l: list.New(),
	}
}

// 注意：
// 1. 无论插入前，元素是否已经存在，插入后，元素都会存在于lru容器中
// 2. 插入元素时，也会更新热度（不管插入前元素是否已经存在）
// @return 插入前元素已经存在则返回false
func (lru *LRU) Put(k interface{}, v interface{}) bool {
	var (
		exist bool
		e     *list.Element // 标准库中的双向循环链表，定义一个链表节点
	)
	e, exist = lru.m[k]
	if exist {
		// 成员map里已存在该key则删除，并从链表里删除节点，后面再写入 操作都是O(1)
		lru.l.Remove(e)
		delete(lru.m, k)
	}

	// 头部更热
	// 插入链表头部并添加到map，此处插入链表是送的value，里面会构造链表节点
	e = lru.l.PushFront(pair{k, v})
	lru.m[k] = e

	// 超出lru容量(链表长度上限)则淘汰(删除)最后一个链表节点，并从map删除该链表中成员对应的key()
	if lru.l.Len() > lru.c {
		// 类型断言，指定具体类型为 pair struct，然后从其中取成员
		k = lru.l.Back().Value.(pair).k

		lru.l.Remove(lru.l.Back())
		delete(lru.m, k)
	}

	// 原来map里已经存在该key则返回false，无论返回什么，lru都会更新热度
	return !exist
}

func (lru *LRU) Get(k interface{}) (v interface{}, exist bool) {
	e, exist := lru.m[k]
	if !exist {
		return nil, false
	}
	pair := e.Value.(pair)
	// 存在则还要更新该成员到链表头部
	lru.l.MoveToFront(e)
	return pair.v, true
}

func (lru *LRU) Size() int {
	return lru.l.Len()
}
