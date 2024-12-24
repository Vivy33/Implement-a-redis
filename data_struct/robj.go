package main

import (
	"fmt"
	"time"
)

// robj 结构体定义
type robj struct {
	rtype    uint8       
	encoding uint8       // 数据编码方式
	lru      uint32      // 最近最少使用时钟值（LRU）
	refcount int         // 引用计数
	ptr      interface{} // 指向实际数据（如字符串、整数等）
}

// 创建一个新的 robj 对象
func createObject(t uint8, ptr interface{}) *robj {
	return &robj{
		rtype:    t,
		encoding: 0,    // 编码方式默认值为 0
		refcount: 1,    
		ptr:      ptr,  // 实际数据指针
		lru:      lruClock(), // 设置 LRU 时钟
	}
}

// 获取当前的 LRU 时钟（模拟）
func lruClock() uint32 {
	// 当前时间戳（精确到秒），可以根据需要自定义
	return uint32(time.Now().Unix())
}
