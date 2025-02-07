package main

import (
    "fmt"
    "log"
    "net"
    "time"
)

func main() {
	rdbFile := "dump.rdb" // RDB 持久化文件路径
	aofFile := "appendonly.aof" // AOF 持久化文件路径

	// 创建一个新的 Redis 服务器实例
	server := newRedisServer("localhost", 6379, rdbFile, aofFile)
	log.Println("Redis server started on :6379")

	// 启动服务器
	if err := server.start(); err != nil {
		log.Fatal(err)
	}
}