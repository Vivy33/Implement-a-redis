package main

import (
    "fmt"
    "strings"
)

// 处理 RESP 协议的命令
type RESPHandler struct {
    data map[string]string // 用于存储键值对的内存数据库
}

// 创建一个新的 RESPHandler 实例
func NewRESPHandler() *RESPHandler {
    return &RESPHandler{
        data: make(map[string]string), // 初始化空的数据存储
    }
}

// 处理接收到的 RESP 命令并返回响应
func (h *RESPHandler) Handle(value RESPValue) Reply {
    // 检查输入是否为数组类型（RESP 中的命令格式）
    if value.Type != Array {
        return &ErrorReply{Value: "ERR invalid request"}
    }

    // 获取命令名称并转换为大写
    command := strings.ToUpper(value.Array[0].Str)

    // 根据不同的命令执行相应的操作
    switch command {
    case "PING":
        // PING 命令：返回 PONG
        return &SimpleStringReply{Value: "PONG"}
    case "SET":
        // SET 命令：设置键值对
        if len(value.Array) != 3 {
            // 检查参数数量是否正确
            return &ErrorReply{Value: "ERR wrong number of arguments for 'set' command"}
        }
        // 存储键值对
        h.data[value.Array[1].Str] = value.Array[2].Str
        return &SimpleStringReply{Value: "OK"}
    case "GET":
        // GET 命令：获取键对应的值
        if len(value.Array) != 2 {
            // 检查参数数量是否正确
            return &ErrorReply{Value: "ERR wrong number of arguments for 'get' command"}
        }
        // 查找并返回键对应的值，如果不存在则返回空字符串
        if val, ok := h.data[value.Array[1].Str]; ok {
            return &BulkStringReply{Value: val}
        }
        return &BulkStringReply{Value: ""}
    default:
        // 未知命令：返回错误信息
        return &ErrorReply{Value: fmt.Sprintf("ERR unknown command '%s'", command)}
    }
}
