package main

import (
	"fmt"
	"strings"
	"time"
)

// 表示一个 Redis 命令的定义，包含命令名称和处理函数。
type redisCommand struct {
	name    string // 命令名称（如 SET、GET）
	handler func(client *redisClient, args []string) // 命令处理函数
}

// 定义了支持的 Redis 命令及其处理逻辑。
var commands = []*redisCommand{
	{
		name: "SET",
		handler: func(c *redisClient, args []string) {
			// SET 命令的处理逻辑
			if len(args) != 3 {
				c.writeResponse("ERR wrong number of arguments for 'SET' command")
				return
			}
			c.server.db.setKey(args[1], args[2]) // 调用数据库的 setKey 方法
			c.writeResponse("OK")
		},
	},
	{
		name: "GET",
		handler: func(c *redisClient, args []string) {
			// GET 命令的处理逻辑
			if len(args) != 2 {
				c.writeResponse("ERR wrong number of arguments for 'GET' command")
				return
			}
			value := c.server.db.getKey(args[1]) // 调用数据库的 getKey 方法
			if value == "" {
				c.writeResponse("(nil)") // 如果没有值，则返回 nil
			} else {
				c.writeResponse(value)
			}
		},
	},
	{
		name: "DEL",
		handler: func(c *redisClient, args []string) {
			// DEL 命令的处理逻辑
			if len(args) != 2 {
				c.writeResponse("ERR wrong number of arguments for 'DEL' command")
				return
			}
			c.server.db.deleteKey(args[1]) // 调用数据库的 deleteKey 方法
			c.writeResponse("OK")
		},
	},
	{
		name: "EXPIRE",
		handler: func(c *redisClient, args []string) {
			// EXPIRE 命令的处理逻辑
			if len(args) != 3 {
				c.writeResponse("ERR wrong number of arguments for 'EXPIRE' command")
				return
			}
			expireTime, err := time.ParseDuration(args[2] + "s") // 解析过期时间
			if err != nil {
				c.writeResponse("ERR invalid expire time")
				return
			}
			c.server.db.setExpire(args[1], expireTime) // 调用数据库的 setExpire 方法
			c.writeResponse("OK")
		},
	},
}