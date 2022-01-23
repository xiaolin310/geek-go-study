package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/thanhpk/randstr"
)

func getMemoryInfo(ctx context.Context, rdb *redis.Client, keyCount, valueSize int) (map[string]int64, error) {
	resultBeforeSet, err := getMemoryUsage(ctx, rdb)
	if err != nil {
		return nil, err
	}
	fmt.Println("Memory usage before set: ", resultBeforeSet)
	// 写入kv数据
	for i := 0; i < keyCount; i++ {
		key := fmt.Sprintf("key_%06d", i)
		rdb.Set(ctx, key, randstr.Bytes(valueSize), 0)
	}

	resultAfterSet, err := getMemoryUsage(ctx, rdb)
	if err != nil {
		return nil, err
	}
	fmt.Println("Memory usage after set: ", resultAfterSet)
	rdb.FlushAll(ctx)

	return calDiffSize(resultBeforeSet, resultAfterSet), nil

}

func calDiffSize(before, after map[string]int64) map[string]int64 {
	res := make(map[string]int64)
	for k, valBefore := range before {
		if valAfter, ok := after[k]; ok {
			res[k] = valAfter - valBefore
		}
	}
	return res
}

func getMemoryUsage(ctx context.Context, rdb *redis.Client) (map[string]int64, error) {
	value, err := rdb.Info(ctx, "memory").Result()
	if err != nil {
		return nil, err
	}
	res := make(map[string]int64)
	// value 格式，used_memory:893944
	scanner := bufio.NewScanner(strings.NewReader(value))
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, ":")
		// 获取 used_memory 和 used_memory_rss 指标
		switch split[0] {
			case "used_memory", "used_memory_rss" :
				val, err := strconv.ParseInt(split[1], 10, 64)
				if err != nil {
					return nil, err
				}
				res[split[0]] = val
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func humanFormatSize(size int) string {
	units := []string{"%.2f B", "%.2f KB", "%.2f MB", "%.2f GB"}
	val := float64(size)
	var i int
	for ; i < len(units); i++ {
		if val < 1024 {
			break
		}
		val = val / 1024
	}

	return fmt.Sprintf(units[i] + " (%d)", val, size)

}

func showStats(values map[string]int64, count int) string {
	builder := strings.Builder{}
	// 平均每个key的占用内存空间
	for k, v := range values {
		builder.WriteString(fmt.Sprintf("\t%s: total usage %s, each key usage %s\n",
			k, humanFormatSize(int(v)), humanFormatSize(int(v)/count)))
	}
	return builder.String()
}

func main() {
	redisAddr := flag.String("addr", "localhost:6379", "redis host")
	valueSize := flag.Int("size", 10, "value size in bytes")
	keyCount := flag.Int("count", 100000, "key value count")
	flag.Parse()

	rdb := redis.NewClient(&redis.Options{
		Addr: *redisAddr,
	})

	ctx := context.Background()
	result, err := getMemoryInfo(ctx, rdb, *keyCount, *valueSize)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Memory usage for writing %d keys with value size of %s :\n%s\n",
		*keyCount, humanFormatSize(*valueSize), showStats(result, *keyCount))
}
