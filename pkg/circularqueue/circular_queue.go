// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package circularqueue

import "errors"

// 底层基于切片实现的固定容量大小的FIFO的环形队列

var ErrCircularQueue = errors.New("circular queue: fxxk")

type CircularQueue struct {
	capacity int
	core     []interface{}
	first    int // 队列头的索引
	last     int // 此处last是最后一个元素的后面，而不是最后一个元素的索引，即队列尾不存元素
}

func New(capacity int) *CircularQueue {
	return &CircularQueue{
		capacity: capacity + 1, // 队列尾不存数据，所以多留一个空间
		core:     make([]interface{}, capacity+1, capacity+1),
		first:    0,
		last:     0,
	}
}

// @return 如果队列满了，则返回错误
func (c *CircularQueue) PushBack(v interface{}) error {
	if c.Full() {
		return ErrCircularQueue
	}

	c.core[c.last] = v
	c.last = (c.last + 1) % c.capacity
	return nil
}

// @return 如果队列为空，则返回错误
func (c *CircularQueue) PopFront() (interface{}, error) {
	if c.Empty() {
		return nil, ErrCircularQueue
	}

	v := c.core[c.first]
	// 从队列头出队，所以队列头后移一个位置
	c.first = (c.first + 1) % c.capacity
	return v, nil
}

// @return 如果队列为空，则返回错误
func (c *CircularQueue) Front() (interface{}, error) {
	if c.Empty() {
		return nil, ErrCircularQueue
	}

	return c.core[c.first], nil
}

// @return 如果队列为空，则返回错误
func (c *CircularQueue) Back() (interface{}, error) {
	if c.Empty() {
		return nil, ErrCircularQueue
	}

	// last-1是最后一个元素的索引，由于是循环队列所以+capacity后取模
	return c.core[(c.last+c.capacity-1)%c.capacity], nil
}

// 获取第i个元素
func (c *CircularQueue) At(i int) (interface{}, error) {
	if i > c.Size()-1 {
		return nil, ErrCircularQueue
	}

	return c.core[(c.first+i)%c.capacity], nil
}

func (c *CircularQueue) Size() int {
	// last-1即最后一个元素，有可能跑到first前面去了(last-1-first+1 + cap)，所以加capacity后取模
	return (c.last + c.capacity - c.first) % c.capacity
}

func (c *CircularQueue) Full() bool {
	// 环形队列，所以取模
	return (c.last+1)%c.capacity == c.first
}

func (c *CircularQueue) Empty() bool {
	return c.first == c.last
}
