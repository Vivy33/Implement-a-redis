package main

import "fmt"

type sds struct {
    data string
}

func newSDS(data string) *sds {
    return &sds{data: data}
}

func (s *sds) append(data string) {
    s.data += data
}

func (s *sds) length() int {
    return len(s.data)
}

func (s *sds) get() string {
    return s.data
}