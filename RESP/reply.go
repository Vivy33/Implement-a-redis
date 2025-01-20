package main

import (
    "fmt"
    "io"
)

// Reply 接口定义了所有 RESP 回复类型必须实现的方法
type Reply interface {
    WriteTo(w io.Writer) (int64, error)
}

// 表示简单字符串回复
type SimpleStringReply struct {
    Value string
}

// 将简单字符串回复写入 io.Writer
func (r *SimpleStringReply) WriteTo(w io.Writer) (int64, error) {
    n, err := fmt.Fprintf(w, "+%s\r\n", r.Value)
    return int64(n), err
}

// 表示错误回复
type ErrorReply struct {
    Value string
}

// 将错误回复写入 io.Writer
func (r *ErrorReply) WriteTo(w io.Writer) (int64, error) {
    n, err := fmt.Fprintf(w, "-%s\r\n", r.Value)
    return int64(n), err
}

// 表示整数回复
type IntegerReply struct {
    Value int64
}

// 将整数回复写入 io.Writer
func (r *IntegerReply) WriteTo(w io.Writer) (int64, error) {
    n, err := fmt.Fprintf(w, ":%d\r\n", r.Value)
    return int64(n), err
}

// 表示批量字符串回复
type BulkStringReply struct {
    Value string
}

// 将批量字符串回复写入 io.Writer
func (r *BulkStringReply) WriteTo(w io.Writer) (int64, error) {
    // 处理空字符串的特殊情况
    if r.Value == "" {
        n, err := w.Write([]byte("$-1\r\n"))
        return int64(n), err
    }
    // 写入字符串长度和内容
    n, err := fmt.Fprintf(w, "$%d\r\n%s\r\n", len(r.Value), r.Value)
    return int64(n), err
}

// 表示数组回复
type ArrayReply struct {
    Value []Reply
}

// 将数组回复写入 io.Writer
func (r *ArrayReply) WriteTo(w io.Writer) (int64, error) {
    // 写入数组长度
    n, err := fmt.Fprintf(w, "*%d\r\n", len(r.Value))
    if err != nil {
        return int64(n), err
    }
    // 遍历并写入数组中的每个元素
    for _, reply := range r.Value {
        m, err := reply.WriteTo(w)
        n += int(m)
        if err != nil {
            return int64(n), err
        }
    }
    return int64(n), nil
}
