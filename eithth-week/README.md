
## 1. 使用 redis benchmark 工具, 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能

```
# redis-benchmark -n 100000 -q -t get,set -d 10
SET: 65487.89 requests per second, p50=0.543 msec
GET: 68027.21 requests per second, p50=0.519 msec

# redis-cli flushall
OK

# redis-benchmark -n 100000 -q -t get,set -d 20
SET: 64599.48 requests per second, p50=0.543 msec
GET: 67159.17 requests per second, p50=0.527 msec

# redis-cli flushall
OK

# redis-benchmark -n 100000 -q -t get,set -d 50
SET: 60975.61 requests per second, p50=0.575 msec
GET: 62227.75 requests per second, p50=0.567 msec

# redis-cli flushall
OK

# redis-benchmark -n 100000 -q -t get,set -d 100
SET: 64516.13 requests per second, p50=0.551 msec
GET: 65146.58 requests per second, p50=0.543 msec

# redis-cli flushall
OK

# redis-benchmark -n 100000 -q -t get,set -d 200
SET: 67567.57 requests per second, p50=0.527 msec
GET: 68306.01 requests per second, p50=0.519 msec

# redis-cli flushall
OK

# redis-benchmark -n 100000 -q -t get,set -d 1024
SET: 60679.61 requests per second, p50=0.583 msec
GET: 65189.05 requests per second, p50=0.543 msec

# redis-cli flushall
OK

# redis-benchmark -n 100000 -q -t get,set -d 5120
SET: 53879.31 requests per second, p50=0.655 msec
GET: 54083.29 requests per second, p50=0.647 msec

# redis-cli flushall
OK

```
## 2. 写入一定量的 kv 数据, 根据数据大小 1w-50w 自己评估, 结合写入前后的 info memory 信息 ,分析上述不同 value 大小下，平均每个 key 的占用内存空间。

```
#  生成二进制执行文件
# CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build RedisInfoMemory.go

# 写入 10w 个key

## 执行结果
# ./RedisInfoMemory --size=10
Memory usage for writing 100000 keys with value size of 10.00 B (10) :
	used_memory: total usage 7.87 MB (8248576), each key usage 82.00 B (82)
	used_memory_rss: total usage 7.77 MB (8146944), each key usage 81.00 B (81)

# ./RedisInfoMemory --size=20
Memory usage for writing 100000 keys with value size of 20.00 B (20) :
	used_memory: total usage 8.63 MB (9048576), each key usage 90.00 B (90)
	used_memory_rss: total usage 7.54 MB (7905280), each key usage 79.00 B (79)
	
# ./RedisInfoMemory --size 50
Memory usage for writing 100000 keys with value size of 50.00 B (50) :
	used_memory: total usage 11.68 MB (12248576), each key usage 122.00 B (122)
	used_memory_rss: total usage 7.36 MB (7720960), each key usage 77.00 B (77)

# ./RedisInfoMemory --size 100
Memory usage for writing 100000 keys with value size of 100.00 B (100) :
	used_memory: total usage 17.02 MB (17848576), each key usage 178.00 B (178)
	used_memory_rss: total usage 12.56 MB (13168640), each key usage 131.00 B (131)

# ./RedisInfoMemory --size 200
Memory usage for writing 100000 keys with value size of 200.00 B (200) :
	used_memory: total usage 27.70 MB (29048576), each key usage 290.00 B (290)
	used_memory_rss: total usage 23.04 MB (24158208), each key usage 241.00 B (241)

# ./RedisInfoMemory --size 1024
Memory usage for writing 100000 keys with value size of 1.00 KB (1024) :
	used_memory_rss: total usage 120.28 MB (126124032), each key usage 1.23 KB (1261)
	used_memory: total usage 128.41 MB (134648576), each key usage 1.31 KB (1346)

# ./RedisInfoMemory --size 5120
Memory usage for writing 100000 keys with value size of 5.00 KB (5120) :
	used_memory: total usage 592.28 MB (621048576), each key usage 6.06 KB (6210)
	used_memory_rss: total usage 594.41 MB (623284224), each key usage 6.09 KB (6232)
```

## 运行环境
```
CentOS Linux release 8.5.2111, 2G 内存
Redis server 单机版， 版本v6.2.6

```