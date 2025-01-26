package main

import (
    "bufio"
    "net"
)

// 表示一个 Redis 连接
// 用于管理与客户端的网络连接
// 它封装了底层的网络连接和用于高效读写的缓冲器
type Connection struct {
    conn   net.Conn
    reader *bufio.Reader
    writer *bufio.Writer
}

// 它接受一个 net.Conn 作为参数，并初始化相应的读写缓冲器
func NewConnection(conn net.Conn) *Connection {
    return &Connection{
        conn:   conn,
        reader: bufio.NewReader(conn),
        writer: bufio.NewWriter(conn),
    }
}

// 从连接中读取一个 RESP 值
// 如果读取过程中发生错误，它会返回该错误
func (c *Connection) Read() (RESPValue, error) {
    return ParseRESP(c.reader)
}

// 将一个 RESP 值写入连接
// 如果读取过程中发生错误，它会返回该错误
func (c *Connection) Write(reply Reply) error {
    _, err := reply.WriteTo(c.writer)
    if err != nil {
        return err
    }
    return c.writer.Flush()
}

// 关闭连接
// 它返回关闭过程中可能发生的任何错误
func (c *Connection) Close() error {
    return c.conn.Close()
}