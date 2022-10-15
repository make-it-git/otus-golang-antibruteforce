package leakybucket

import (
	"sync"
	"time"
)

type bucket struct {
	m           sync.Mutex
	key         string
	load        int64
	maxItems    int64
	unusedSince time.Time
}

func newBucket(key string, maxLoad int64, delCh chan<- string, ttl time.Duration) *bucket {
	b := &bucket{
		key:      key,
		maxItems: maxLoad,
	}

	go func() {
		ticker := time.NewTicker(time.Duration(int64(time.Minute) / b.maxItems))
		defer ticker.Stop()
		for {
			<-ticker.C

			b.m.Lock()
			if b.load > 0 {
				b.load--
				b.m.Unlock()
				continue
			}
			b.m.Unlock()

			if b.unusedSince.IsZero() {
				b.unusedSince = time.Now()
				continue
			}

			if ttl < time.Since(b.unusedSince) {
				delCh <- b.key
				break
			}
		}
	}()

	return b
}

func (b *bucket) try() error {
	b.m.Lock()
	defer b.m.Unlock()

	if b.load < b.maxItems {
		b.load++
		return nil
	}

	return ErrBlocked
}
