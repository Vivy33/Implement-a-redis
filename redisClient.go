package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// 表示一个 Redis 客户端连接，包含网络连接和服务器引用。
type redisClient struct {
	conn   net.Conn // 客户端的网络连接
	server *redisServer // 服务器引用
}

// 创建一个新的 Redis 客户端实例。
func newRedisClient(conn net.Conn, server *redisServer) *redisClient {
	return &redisClient{
		conn:   conn,
		server: server,
	}
}

// 处理客户端请求，读取并解析命令。
func (c *redisClient) handleRequest() {
	defer c.conn.Close() // 确保连接在处理完请求后关闭

	reader := bufio.NewReader(c.conn) // 使用 bufio 包来高效读取客户端请求
	for {
		request, err := reader.ReadString('\n') // 读取客户端请求行，直到遇到换行符
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}
		request = strings.TrimSpace(request) // 清除请求的换行符

		c.processCommand(request) // 处理客户端的命令
	}
}

// 处理客户端发送的命令。
func (c *redisClient) processCommand(command string) {
	args := strings.Fields(command) // 拆分命令和参数
	if len(args) == 0 {
		return // 如果没有命令，则直接返回
	}
	cmd := strings.ToUpper(args[0]) // 获取命令名称（如 SET、GET 等）

	// 遍历命令定义，找到匹配的命令并执行
	for _, cmdDef := range commands {
		if cmd == cmdDef.name {
			cmdDef.handler(c, args) // 执行命令处理函数
			return
		}
	}
	// 对于未识别的命令，返回错误响应
	c.writeResponse("ERR unknown command '" + cmd + "'")
}

// 向客户端发送响应。
func (c *redisClient) writeResponse(response string) {
	// 将响应写入客户端连接
	c.conn.Write([]byte(response + "\r\n"))
}