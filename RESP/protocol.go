package main

import (
    "bufio"
    "errors"
    "io"
    "strconv"
)

type RESPType byte

const (
    SimpleString RESPType = '+' // 简单字符串
    Error        RESPType = '-' // 错误
    Integer      RESPType = ':' // 整数
    BulkString   RESPType = '$' // 批量字符串
    Array        RESPType = '*' // 数组
)

// 存储解析后的 RESP 数据
type RESPValue struct {
    Type  RESPType    // 数据类型
    Str   string      // 用于存储字符串值
    Int   int64       // 用于存储整数值
    Array []RESPValue // 用于存储数组值
}

// 解析 RESP 协议数据
func ParseRESP(reader *bufio.Reader) (RESPValue, error) {
    // 读取类型标识符
    type, err := reader.ReadByte()
    if err != nil {
        return RESPValue{}, err
    }

    // 根据类型调用相应的解析函数
    switch RESPType(type) {
    case SimpleString, Error:
        return parseSimpleString(reader, RESPType(type))
    case Integer:
        return parseInteger(reader)
    case BulkString:
        return parseBulkString(reader)
    case Array:
        return parseArray(reader)
    default:
        return RESPValue{}, errors.New("Unknown RESP type")
    }
}

// 解析简单字符串或错误
func parseSimpleString(reader *bufio.Reader, type RESPType) (RESPValue, error) {
    line, err := reader.ReadString('\n')
    if err != nil {
        return RESPValue{}, err
    }
    // 去除末尾的 \r\n
    return RESPValue{Type: type, Str: line[:len(line)-2]}, nil
}

// 解析整数
func parseInteger(reader *bufio.Reader) (RESPValue, error) {
    line, err := reader.ReadString('\n')
    if err != nil {
        return RESPValue{}, err
    }
    // 解析整数值
    i, err := strconv.ParseInt(line[:len(line)-2], 10, 64)
    if err != nil {
        return RESPValue{}, err
    }
    return RESPValue{Type: Integer, Int: i}, nil
}

// 解析批量字符串
func parseBulkString(reader *bufio.Reader) (RESPValue, error) {
    // 读取字符串长度
    line, err := reader.ReadString('\n')
    if err != nil {
        return RESPValue{}, err
    }
    length, err := strconv.Atoi(line[:len(line)-2])
    if err != nil {
        return RESPValue{}, err
    }
    // 处理空字符串（长度为-1）
    if length == -1 {
        return RESPValue{Type: BulkString, Str: ""}, nil
    }
    // 读取指定长度的字符串
    buf := make([]byte, length+2) // +2 for \r\n
    _, err = io.ReadFull(reader, buf)
    if err != nil {
        return RESPValue{}, err
    }
    return RESPValue{Type: BulkString, Str: string(buf[:length])}, nil
}

// 解析数组
func parseArray(reader *bufio.Reader) (RESPValue, error) {
    // 读取数组长度
    line, err := reader.ReadString('\n')
    if err != nil {
        return RESPValue{}, err
    }
    count, err := strconv.Atoi(line[:len(line)-2])
    if err != nil {
        return RESPValue{}, err
    }
    // 解析数组中的每个元素
    array := make([]RESPValue, count)
    for i := 0; i < count; i++ {
        value, err := ParseRESP(reader)
        if err != nil {
            return RESPValue{}, err
        }
        array[i] = value
    }
    return RESPValue{Type: Array, Array: array}, nil
}
