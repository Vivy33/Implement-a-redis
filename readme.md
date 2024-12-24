项目结构
整个Redis项目主要可以分为以下几大块：

数据结构
主要是数据结构的实现，包括：
sds       动态字符串的实现
adlist    双端链表的实现
dict  字典的实现
zskiplist 跳跃表的实现
ziplist zipmap 压缩表、压缩字段的实现
hyperloglog hyperloglogr的实现

数据类型
object Redis 的对象（类型）系统实现。
t_string  字符串键的实现。
t_list  列表键的实现。
t_hash 散列键的实现。
t_set 集合键的实现。
t_zset 有序集合键的实现。
hyperloglog HyperLogLog 键的实现

数据库
db Redis的数据库实现
notify 数据库通知功能实现
rdb  rdb持久化实现
aof aof持久化实现

客户端和服务端相关
ae* Redis的事件处理器实现，redis自己实现了一套基于reactor模式的事件处理器
networking 网络连接库，负责网络相关的操作，比如发送，接收命令，协议解析，创建/销毁客户端
redis 单机Redis服务器实现

集群相关
replication  复制功能的实现
sentinel   sentinel的实现
cluster    集群的实现
