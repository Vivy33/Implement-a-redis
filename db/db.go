package main

import (
    "fmt"
	"time"
)

type redisDb struct {
	dict    *dict // 存储键值对
	expires *dict // 存储过期时间
	id      int    // 数据库 ID
}

func newRedisDb(id int) *redisDb {
	return &redisDb{
		dict:    newDict(),
		expires: newDict(),
		id:      id,
	}
}

func (r *redisDb) setKey(key *robj, val *robj) {
	if r.lookupKey(key) == nil {
		r.dbAdd(key, val)
	} else {
		r.dbOverwrite(key, val)
	}
	val.refcount++
}

func (r *redisDb) dbAdd(key *robj, val *robj) {
	r.dict.dictAdd(key.ptr, val)
}

func (r *redisDb) dbOverwrite(key *robj, val *robj) {
	r.dict.dictAdd(key.ptr, val)
}

func (r *redisDb) lookupKey(key *robj) *robj {
	return r.doLookupKey(key)
}

func (r *redisDb) doLookupKey(key *robj) *robj {
	return r.dict.dictFind(key.ptr)
}

func (r *redisDb) setExpire(key *robj, expire uint64) {
	kde := r.dict.dictFind(key.ptr)
	if kde != nil {
		// 为 key 设置过期时间
		r.expires.dictAdd(key.ptr, newRobj(0, expire))
	}
}
