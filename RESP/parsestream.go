package main

import (
    "bufio"
    "io"
)

type ParseStream struct {
    reader *bufio.Reader
}

func NewParseStream(r io.Reader) *ParseStream {
    return &ParseStream{reader: bufio.NewReader(r)}
}

func (ps *ParseStream) Next() (RESPValue, error) {
    return ParseRESP(ps.reader)
}