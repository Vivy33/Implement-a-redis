package main

import (
	"fmt"
	"net"
	"time"
)

// 表示一个 Redis 服务器实例，包含主机地址、端口、数据库和客户端列表。
type redisServer struct {
	host          string // 服务器主机地址
	port          int // 服务器端口
	db            *redisDb // 数据库实例
	activeClients []*redisClient // 当前活跃的客户端列表
}

// 创建一个新的 Redis 服务器实例，并加载 RDB 和 AOF 文件。
func newRedisServer(host string, port int, rdbFile, aofFile string) *redisServer {
	db := newRedisDb(rdbFile, aofFile)
	// 加载 RDB 和 AOF 文件
	db.loadRDB()
	db.loadAOF()
	return &redisServer{
		host: host,
		port: port,
		db:   db,
	}
}

// 启动一个协程定期清理过期键。
func (s *redisServer) cleanExpiredKeys() {
	go s.db.cleanExpiredKeys()
}

// start 启动 Redis 服务器，监听客户端连接并处理请求。
func (s *redisServer) start() {
	address := fmt.Sprintf("%s:%d", s.host, s.port)
	ln, err := net.Listen("tcp", address) // 监听 TCP 端口
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Server started on", address)

	s.cleanExpiredKeys() // 启动过期键清理协程

	for {
		conn, err := ln.Accept() // 接受客户端连接
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		client := newRedisClient(conn, s) // 创建新的客户端实例
		s.activeClients = append(s.activeClients, client)
		go client.handleRequest() // 使用 Goroutine 处理客户端请求
	}
}