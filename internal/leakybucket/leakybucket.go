package leakybucket

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

var ErrBlocked = errors.New("blocked")

type LeakyBucket struct {
	maxLoadLogin    int64
	maxLoadPassword int64
	maxLoadIP       int64

	loginMutex   sync.Mutex
	loginBuckets map[string]*bucket
	loginCh      chan string

	passwordMutex   sync.Mutex
	passwordBuckets map[string]*bucket
	passwordCh      chan string

	ipAddrMutex   sync.Mutex
	ipAddrBuckets map[string]*bucket
	ipAddrCh      chan string

	ttl time.Duration
}

func NewLeakyBucket(ctx context.Context, maxLogin, maxPassword, maxIP int64, ttl time.Duration) *LeakyBucket {
	b := &LeakyBucket{
		maxLoadLogin:    maxLogin,
		maxLoadPassword: maxPassword,
		maxLoadIP:       maxIP,

		loginBuckets: map[string]*bucket{},
		loginCh:      make(chan string),

		passwordBuckets: map[string]*bucket{},
		passwordCh:      make(chan string),

		ipAddrBuckets: map[string]*bucket{},
		ipAddrCh:      make(chan string),

		ttl: ttl,
	}

	go func() {
		for {
			select {
			case login := <-b.loginCh:
				b.loginMutex.Lock()
				fmt.Println("delete from logins", login)
				delete(b.loginBuckets, login)
				b.loginMutex.Unlock()
			case password := <-b.passwordCh:
				b.passwordMutex.Lock()
				fmt.Println("delete from password", password)
				delete(b.passwordBuckets, password)
				b.passwordMutex.Unlock()
			case ip := <-b.ipAddrCh:
				b.ipAddrMutex.Lock()
				fmt.Println("delete from ip", ip)
				delete(b.ipAddrBuckets, ip)
				b.ipAddrMutex.Unlock()
			case <-ctx.Done():
				break
			}
		}
	}()

	return b
}

func (s *LeakyBucket) Try(login string, password string, ip string) error {
	var eg errgroup.Group

	eg.Go(func() error {
		s.loginMutex.Lock()
		defer s.loginMutex.Unlock()

		return checkBucket(s.loginBuckets, login, s.maxLoadLogin, s.loginCh, s.ttl)
	})

	eg.Go(func() error {
		s.passwordMutex.Lock()
		defer s.passwordMutex.Unlock()

		return checkBucket(s.passwordBuckets, password, s.maxLoadPassword, s.passwordCh, s.ttl)
	})

	eg.Go(func() error {
		s.ipAddrMutex.Lock()
		defer s.ipAddrMutex.Unlock()

		return checkBucket(s.ipAddrBuckets, ip, s.maxLoadIP, s.ipAddrCh, s.ttl)
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func checkBucket(
	buckets map[string]*bucket,
	bucketKey string,
	maxLoad int64,
	delCh chan<- string,
	ttl time.Duration,
) error {
	b, ok := buckets[bucketKey]
	if !ok {
		b = newBucket(bucketKey, maxLoad, delCh, ttl)
		buckets[bucketKey] = b
	}

	return b.try()
}
