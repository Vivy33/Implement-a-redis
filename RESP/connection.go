package main

import (
    "bufio"
    "net"
)

// 表示一个 Redis 连接
type Connection struct {
    conn   net.Conn
    reader *bufio.Reader
    writer *bufio.Writer
}

func NewConnection(conn net.Conn) *Connection {
    return &Connection{
        conn:   conn,
        reader: bufio.NewReader(conn),
        writer: bufio.NewWriter(conn),
    }
}

// 从连接中读取一个 RESP 值
func (c *Connection) Read() (RESPValue, error) {
    return ParseRESP(c.reader)
}

// 将一个 RESP 值写入连接
func (c *Connection) Write(reply Reply) error {
    _, err := reply.WriteTo(c.writer)
    if err != nil {
        return err
    }
    return c.writer.Flush()
}

// 关闭连接
func (c *Connection) Close() error {
    return c.conn.Close()
}