package main

import (
    "bufio"
    "io"
)

// 提供了一个流式接口来解析 RESP 数据
type ParseStream struct {
    reader *bufio.Reader
}

func NewParseStream(r io.Reader) *ParseStream {
    return &ParseStream{reader: bufio.NewReader(r)}
}

// 从流中读取并解析下一个 RESP 值
func (ps *ParseStream) Next() (RESPValue, error) {
    return ParseRESP(ps.reader)
}