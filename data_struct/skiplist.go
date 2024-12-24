package main

import (
    "fmt"
    "math/rand"
    "time"
)

const MAX_LEVEL = 16

type SkipListNode struct {
    value interface{}
    score int
    forward []*SkipListNode
}

type SkipList struct {
    level int
    header *SkipListNode
}

func NewSkipList() *SkipList {
    header := &SkipListNode{forward: make([]*SkipListNode, MAX_LEVEL)}
    return &SkipList{level: 1, header: header}
}

func (sl *SkipList) insert(score int, value interface{}) {
    update := make([]*SkipListNode, MAX_LEVEL)
    current := sl.header
    for i := sl.level - 1; i >= 0; i-- {
        for current.forward[i] != nil && current.forward[i].score < score {
            current = current.forward[i]
        }
        update[i] = current
    }

    level := sl.randomLevel()
    if level > sl.level {
        for i := sl.level; i < level; i++ {
            update[i] = sl.header
        }
        sl.level = level
    }

    newNode := &SkipListNode{
        score:  score,
        value:  value,
        forward: make([]*SkipListNode, level),
    }

    for i := 0; i < level; i++ {
        newNode.forward[i] = update[i].forward[i]
        update[i].forward[i] = newNode
    }
}

func (sl *SkipList) randomLevel() int {
    level := 1
    for rand.Float32() < 0.5 && level < MAX_LEVEL {
        level++
    }
    return level
}

func (sl *SkipList) printList() {
    for i := 0; i < sl.level; i++ {
        current := sl.header.forward[i]
        fmt.Printf("Level %d: ", i)
        for current != nil {
            fmt.Printf("%d->", current.score)
            current = current.forward[i]
        }
        fmt.Println("nil")
    }
}
