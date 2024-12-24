package main

// redisCommand 用于存储命令的信息
type redisCommand struct {
	name    string              // 命令名称（如 SET、GET）
	handler func(client *redisClient, args []string) // 命令处理函数
}

// 创建一些基本的命令
var commands = []*redisCommand{
	{
		name: "SET", // SET 命令
		handler: func(c *redisClient, args []string) {
			// SET 命令的处理逻辑
			if len(args) != 3 {
				c.writeResponse("ERR wrong number of arguments for 'SET' command")
				return
			}
			c.db.setKey(args[1], args[2])
			c.writeResponse("OK")
		},
	},
	{
		name: "GET", // GET 命令
		handler: func(c *redisClient, args []string) {
			// GET 命令的处理逻辑
			if len(args) != 2 {
				c.writeResponse("ERR wrong number of arguments for 'GET' command")
				return
			}
			value := c.db.getKey(args[1])
			if value == "" {
				c.writeResponse("(nil)") // 如果没有值，则返回 nil
			} else {
				c.writeResponse(value)
			}
		},
	},
}