package main

import (
    "fmt"
    "log"
    "net"
    "time"
)

type RedisServer struct {
    addr     string
    port     int
    db       *redisDb
    handler  *RESPHandler
}

func newRedisServer(addr string, port int) *RedisServer {
    return &RedisServer{
        addr:     addr,
        port:     port,
        db:       newRedisDb(0),
        handler:  NewRESPHandler(),
    }
}

func (s *RedisServer) start() {
    listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.addr, s.port))
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()

    log.Printf("Redis server listening on %s:%d\n", s.addr, s.port)

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("Error accepting connection:", err)
            continue
        }

        go s.handleConnection(conn)
    }
}

func (s *RedisServer) handleConnection(conn net.Conn) {
    defer conn.Close()

    redisConn := NewConnection(conn)

    for {
        value, err := redisConn.Read()
        if err != nil {
            log.Println("Error reading from connection:", err)
            return
        }

        reply := s.handler.Handle(value, s.db)
        err = redisConn.Write(reply)
        if err != nil {
            log.Println("Error writing to connection:", err)
            return
        }
    }
}

func (s *RedisServer) cleanExpiredKeys() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        now := uint64(time.Now().Unix())
        s.db.expires.dictIterator(func(key, value interface{}) bool {
            if expire, ok := value.(*robj); ok {
                if expire.ptr.(uint64) <= now {
                    s.db.dict.dictDelete(key)
                    s.db.expires.dictDelete(key)
                }
            }
            return true
        })
    }
}

func main() {
    // 创建一个新的 Redis 服务端实例
    server := newRedisServer("localhost", 6379)

    // 启动清理过期键的协程
    go server.cleanExpiredKeys()

    // 启动 Redis 服务端，开始监听客户端连接
    server.start()
}