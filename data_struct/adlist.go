package main

import (
    "container/list"
    "fmt"
)

type adlist struct {
    list *list.List
}

func newAdlist() *adlist {
    return &adlist{list: list.New()}
}

func (a *adlist) pushFront(value interface{}) {
    a.list.PushFront(value)
}

func (a *adlist) pushBack(value interface{}) {
    a.list.PushBack(value)
}

func (a *adlist) popFront() interface{} {
    front := a.list.Front()
    if front != nil {
        return a.list.Remove(front)
    }
    return nil
}

func (a *adlist) popBack() interface{} {
    back := a.list.Back()
    if back != nil {
        return a.list.Remove(back)
    }
    return nil
}

func (a *adlist) iterate() {
    for e := a.list.Front(); e != nil; e = e.Next() {
        fmt.Println(e.Value)
    }