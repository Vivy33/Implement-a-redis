package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"sync"
	"time"
)

// 表示一个 Redis 数据库，包含键值对存储、过期时间存储以及持久化文件路径。
type redisDb struct {
	data     map[string]string // 存储键值对
	expires  map[string]time.Time // 存储键的过期时间
	mu       sync.RWMutex // 互斥锁，用于并发控制
	rdbFile  string // RDB 持久化文件路径
	aofFile  string // AOF 持久化文件路径
}

// 创建一个新的 Redis 数据库实例。
func newRedisDb(rdbFile, aofFile string) *redisDb {
	return &redisDb{
		data:    make(map[string]string),
		expires: make(map[string]time.Time),
		rdbFile: rdbFile,
		aofFile: aofFile,
	}
}

// 设置一个键值对，并将操作记录到 AOF 文件。
func (db *redisDb) setKey(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
	db.saveAOF(fmt.Sprintf("SET %s %s", key, value)) // 记录到 AOF 文件
}

// 获取一个键对应的值。
func (db *redisDb) getKey(key string) string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.data[key]
}

// 删除一个键，并将操作记录到 AOF 文件。
func (db *redisDb) deleteKey(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, key)
	delete(db.expires, key)
	db.saveAOF(fmt.Sprintf("DEL %s", key)) // 记录到 AOF 文件
}

// 为一个键设置过期时间，并将操作记录到 AOF 文件。
func (db *redisDb) setExpire(key string, expireTime time.Duration) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.expires[key] = time.Now().Add(expireTime)
	db.saveAOF(fmt.Sprintf("EXPIRE %s %d", key, int(expireTime.Seconds()))) // 记录到 AOF 文件
}

// 将当前数据库状态保存到 RDB 文件。
func (db *redisDb) saveRDB() error {
	db.mu.RLock()
	defer db.mu.RUnlock()
	file, err := os.Create(db.rdbFile) // 创建或覆盖 RDB 文件
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file) // 使用 gob 编码器
	return encoder.Encode(db.data) // 将数据编码并写入文件
}

// loadRDB 从 RDB 文件加载数据到内存。
func (db *redisDb) loadRDB() error {
	file, err := os.Open(db.rdbFile) // 打开 RDB 文件
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file) // 使用 gob 解码器
	return decoder.Decode(&db.data) // 将文件内容解码到内存
}

// 将命令追加到 AOF 文件。
func (db *redisDb) saveAOF(command string) {
	file, err := os.OpenFile(db.aofFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // 打开 AOF 文件
	if err != nil {
		fmt.Println("Error opening AOF file:", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(command + "\n") // 将命令写入文件
	if err != nil {
		fmt.Println("Error writing to AOF file:", err)
	}
}

// loadAOF 从 AOF 文件加载命令并执行。
func (db *redisDb) loadAOF() error {
	file, err := os.Open(db.aofFile) // 打开 AOF 文件
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file) // 创建扫描器
	for scanner.Scan() {
		line := scanner.Text()
		args := strings.Fields(line)
		cmd := strings.ToUpper(args[0])
		for _, cmdDef := range commands {
			if cmd == cmdDef.name {
				cmdDef.handler(nil, args) // 执行命令
				break
			}
		}
	}
	return scanner.Err()
}

// 定期清理过期的键。
func (db *redisDb) cleanExpiredKeys() {
	ticker := time.NewTicker(1 * time.Second) // 每秒触发一次
	defer ticker.Stop()
	for range ticker.C {
		db.mu.Lock()
		now := time.Now()
		for key, expireTime := range db.expires {
			if expireTime.Before(now) { // 如果键已过期
				delete(db.data, key)
				delete(db.expires, key)
			}
		}
		db.mu.Unlock()
	}
}