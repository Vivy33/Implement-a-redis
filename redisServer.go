package main

import (
	"fmt"
	"net"
	"time"
)

// redisServer 保存 Redis 服务端配置信息、数据库等
type redisServer struct {
	host          string           // 服务端主机地址
	port          int              // 服务端端口
	dbs           []*redisDb      // 支持多个数据库（Redis 默认 16 个数据库）
	activeClients []*redisClient  // 当前连接的客户端列表
}

// 创建一个新的 redisServer 实例
func newRedisServer(host string, port int) *redisServer {
	// 初始化服务端，默认支持 16 个数据库
	server := &redisServer{
		host: host,
		port: port,
		dbs:  make([]*redisDb, 16), // Redis 默认支持 16 个数据库
	}

	// 初始化每个数据库
	for i := 0; i < len(server.dbs); i++ {
		server.dbs[i] = newRedisDb(i)
	}
	return server
}

// 定期清理过期键
func (s *redisServer) cleanExpiredKeys() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for range ticker.C {
        for _, db := range s.dbs {
            db.mu.Lock()
            now := time.Now()
            for key, expireTime := range db.expires {
                if expireTime.(time.Time).Before(now) {
                    delete(db.data, key)
                    delete(db.expires, key)
                }
            }
            db.mu.Unlock()
        }
    }
}

// 启动 Redis 服务端，监听客户端连接并处理请求
func (s *redisServer) start() {
	// 设置监听地址
	address := fmt.Sprintf("%s:%d", s.host, s.port)

	// 开始监听 TCP 端口
	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Server started on", address)

	// 持续接受客户端连接
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// 创建新的客户端实例并处理请求
		client := newRedisClient(conn, s)
		s.activeClients = append(s.activeClients, client)
		go client.handleRequest() // 使用 goroutine 处理客户端请求
	}
}