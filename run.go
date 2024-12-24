func main() {
	// 创建一个新的 Redis 服务端实例
	server := newRedisServer("localhost", 6379)

	// 启动 Redis 服务端，开始监听客户端连接
	server.start()
}