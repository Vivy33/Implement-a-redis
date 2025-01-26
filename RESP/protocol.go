package resp

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "strconv"
)

// RESP 类型常量
const (
    SimpleString = '+'
    Error        = '-'
    Integer      = ':'
    BulkString   = '$'
    Array        = '*'
)

// 常见错误
var (
    ErrInvalidSyntax = errors.New("invalid RESP syntax")
    ErrUnexpectedEOF = errors.New("unexpected end of input")
)

// RESPValue 表示一个 RESP 值
type RESPValue struct {
    Type  byte
    Str   string
    Num   int64
    Array []RESPValue
}

// ParseRESP 从 reader 中解析一个 RESP 值
func ParseRESP(reader *bufio.Reader) (RESPValue, error) {
    typ, err := reader.ReadByte()
    if err != nil {
        return RESPValue{}, fmt.Errorf("read type byte: %w", err)
    }

    switch typ {
    case SimpleString, Error:
        return parseLineResp(typ, reader)
    case Integer:
        return parseIntegerResp(reader)
    case BulkString:
        return parseBulkStringResp(reader)
    case Array:
        return parseArrayResp(reader)
    default:
        return RESPValue{}, fmt.Errorf("%w: unknown type %c", ErrInvalidSyntax, typ)
    }
}

func parseLineResp(typ byte, reader *bufio.Reader) (RESPValue, error) {
    line, err := reader.ReadString('\n')
    if err != nil {
        return RESPValue{}, fmt.Errorf("read line: %w", err)
    }
    return RESPValue{Type: typ, Str: line[:len(line)-2]}, nil
}

func parseIntegerResp(reader *bufio.Reader) (RESPValue, error) {
    line, err := reader.ReadString('\n')
    if err != nil {
        return RESPValue{}, fmt.Errorf("read integer line: %w", err)
    }
    num, err := strconv.ParseInt(line[:len(line)-2], 10, 64)
    if err != nil {
        return RESPValue{}, fmt.Errorf("parse integer: %w", err)
    }
    return RESPValue{Type: Integer, Num: num}, nil
}

func parseBulkStringResp(reader *bufio.Reader) (RESPValue, error) {
    line, err := reader.ReadString('\n')
    if err != nil {
        return RESPValue{}, fmt.Errorf("read bulk string length: %w", err)
    }
    length, err := strconv.Atoi(line[:len(line)-2])
    if err != nil {
        return RESPValue{}, fmt.Errorf("parse bulk string length: %w", err)
    }
    if length == -1 {
        return RESPValue{Type: BulkString, Str: ""}, nil
    }
    buf := make([]byte, length+2) // +2 for CRLF
    _, err = io.ReadFull(reader, buf)
    if err != nil {
        return RESPValue{}, fmt.Errorf("read bulk string content: %w", err)
    }
    return RESPValue{Type: BulkString, Str: string(buf[:length])}, nil
}

func parseArrayResp(reader *bufio.Reader) (RESPValue, error) {
    line, err := reader.ReadString('\n')
    if err != nil {
        return RESPValue{}, fmt.Errorf("read array length: %w", err)
    }
    length, err := strconv.Atoi(line[:len(line)-2])
    if err != nil {
        return RESPValue{}, fmt.Errorf("parse array length: %w", err)
    }
    if length == -1 {
        return RESPValue{Type: Array, Array: nil}, nil
    }
    array := make([]RESPValue, length)
    for i := 0; i < length; i++ {
        value, err := ParseRESP(reader)
        if err != nil {
            return RESPValue{}, fmt.Errorf("parse array element %d: %w", i, err)
        }
        array[i] = value
    }
    return RESPValue{Type: Array, Array: array}, nil
}

// WriteTo 将 RESPValue 写入 io.Writer
func (v RESPValue) WriteTo(w io.Writer) (int64, error) {
    var total int64
    switch v.Type {
    case SimpleString, Error:
        n, err := fmt.Fprintf(w, "%c%s\r\n", v.Type, v.Str)
        return int64(n), err
    case Integer:
        n, err := fmt.Fprintf(w, "%c%d\r\n", v.Type, v.Num)
        return int64(n), err
    case BulkString:
        n, err := fmt.Fprintf(w, "%c%d\r\n%s\r\n", v.Type, len(v.Str), v.Str)
        return int64(n), err
    case Array:
        n, err := fmt.Fprintf(w, "%c%d\r\n", v.Type, len(v.Array))
        total += int64(n)
        if err != nil {
            return total, err
        }
        for _, item := range v.Array {
            n, err := item.WriteTo(w)
            total += n
            if err != nil {
                return total, err
            }
        }
        return total, nil
    default:
        return 0, fmt.Errorf("%w: unknown type %c", ErrInvalidSyntax, v.Type)
    }
}
