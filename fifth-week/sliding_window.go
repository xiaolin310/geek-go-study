package fifth_week

import (
	"github.com/pkg/errors"
	"sync"
	"time"
)


type SlidingWindow struct {
	bucketDuration time.Duration // 桶的时间间隔，默认1s
	winDuration    time.Duration // sliding window的长度
	maxReq         int // 大窗口时间内允许的最大请求数
	windows        []*bucket // bucket 合集
	mu             sync.Mutex
}

type bucket struct {
	begin time.Time // 时间起点
	count int       // 落在此bucket内的请求数
}

func NewSlidingWindow(bucketDuration, winDuration time.Duration, maxReq int) *SlidingWindow {
	return &SlidingWindow{
		bucketDuration: bucketDuration,
		winDuration: winDuration,
		maxReq: maxReq,
	}
}

func (w *SlidingWindow) Allow() bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	timeoutOffset := -1
	// 将过期的移出窗口
	for i, bu := range w.windows {
		if bu.begin.Add(w.winDuration).After(now) {
			break
		}
		timeoutOffset = i
	}
	if timeoutOffset > -1 {
		w.windows = w.windows[timeoutOffset+1:]
	}

	// 判断请求是否超限
	var ok bool
	if w.countReq() < w.maxReq {
		ok = true
	}

	// 记录本次请求数
	var lastBucket *bucket
	if len(w.windows) > 0 {
		lastBucket = w.windows[len(w.windows)-1]
		if lastBucket.begin.Add(w.bucketDuration).Before(now) {
			lastBucket = &bucket{begin: now, count: 1}
			w.windows = append(w.windows, lastBucket)
		} else {
			lastBucket.count++
		}
	} else {
		lastBucket = &bucket{begin: now, count: 1}
		w.windows = append(w.windows, lastBucket)
	}
	return ok

}

func (w *SlidingWindow) countReq() int {
	var count int
	for _, bu := range w.windows {
		count += bu.count
	}
	return count
}

func (w *SlidingWindow) Run(foo func()) error {
	if !w.Allow() {
		return errors.New("Requests Nums Exceed")
	}
	foo()
	return nil
}