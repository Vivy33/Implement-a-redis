package main

import (
    "fmt"
)

type dict map[interface{}]interface{}

func (d dict) dictAdd(key interface{}, value interface{}) {
    d[key] = value
}

func (d dict) dictFind(key interface{}) interface{} {
    return d[key]
}

func (d dict) dictDelete(key interface{}) {
    delete(d, key)
}