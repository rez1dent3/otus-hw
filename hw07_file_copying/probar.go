package main

import (
	"fmt"
	"math"
	"strings"
	"sync/atomic"
)

type progressBar struct {
	cur, limit, max int64
}

func (b *progressBar) Write(p []byte) (n int, err error) {
	atomic.AddInt64(&b.cur, int64(len(p)))

	if b.max == 0 {
		b.max = math.MaxInt64
	}

	batchSize := b.limit
	if b.max < b.limit {
		batchSize = b.max
	} else if b.limit == 0 {
		batchSize = b.max
	}

	percent := int64((float64(b.cur) / float64(batchSize)) * 100)
	rate := int(percent / 2)
	if rate == 0 || rate > 50 {
		rate = 1
	}

	fmt.Printf("\r[%-50s]% 3d%% %d/%d", strings.Repeat("#", rate), percent, b.cur, batchSize)

	return len(p), nil
}
