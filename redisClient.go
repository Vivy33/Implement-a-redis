package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// redisClient 用于处理每个客户端连接的请求
type redisClient struct {
	conn   net.Conn      // 客户端的网络连接
	server *redisServer  // 服务端引用
	db     *redisDb      // 客户端选择的数据库（Redis 默认选择数据库 0）
}

// 创建一个新的 redisClient 实例
func newRedisClient(conn net.Conn, server *redisServer) *redisClient {
	return &redisClient{
		conn:   conn,
		server: server,
		db:     server.dbs[0], // 默认选择数据库 0
	}
}

// 处理客户端请求，读取并解析命令
func (c *redisClient) handleRequest() {
	defer c.conn.Close() // 确保连接在处理完请求后关闭

	// 使用 bufio 包来高效读取客户端请求
	reader := bufio.NewReader(c.conn)
	for {
		// 读取客户端请求行，直到遇到换行符
		request, err := reader.ReadString('\n')
		if err != nil {
			// 如果读取失败，打印错误并返回
			fmt.Println("Error reading from client:", err)
			return
		}

		// 清除请求的换行符（\n）
		request = strings.TrimSpace(request)

		// 处理客户端的命令
		c.processCommand(request)
	}
}

// 处理客户端发送的命令
func (c *redisClient) processCommand(command string) {
	// 拆分命令和参数
	args := strings.Fields(command)
	if len(args) == 0 {
		return // 如果没有命令，则直接返回
	}

	// 获取命令名称（如 SET、GET 等）
	cmd := args[0]
	switch cmd {
	case "SET":
		// 如果命令是 SET，处理 SET 命令
		if len(args) != 3 {
			c.writeResponse("ERR wrong number of arguments for 'SET' command")
			return
		}
		c.db.setKey(args[1], args[2])
		c.writeResponse("OK")

	case "GET":
		// 如果命令是 GET，处理 GET 命令
		if len(args) != 2 {
			c.writeResponse("ERR wrong number of arguments for 'GET' command")
			return
		}
		value := c.db.getKey(args[1])
		if value == "" {
			c.writeResponse("(nil)") // 如果值不存在，返回 nil
		} else {
			c.writeResponse(value)
		}

	default:
		// 对于未识别的命令，返回错误响应
		c.writeResponse("ERR unknown command '" + cmd + "'")
	}
}

// 向客户端发送响应
func (c *redisClient) writeResponse(response string) {
	// 将响应写入客户端
	c.conn.Write([]byte(response + "\r\n"))
}