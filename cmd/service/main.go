package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	"github.com/make-it-git/otus-antibruteforce/internal/config"
	"github.com/make-it-git/otus-antibruteforce/internal/leakybucket"
	"github.com/make-it-git/otus-antibruteforce/internal/service"
	storage "github.com/make-it-git/otus-antibruteforce/internal/storage_redis"
	api "github.com/make-it-git/otus-antibruteforce/pkg/antibruteforce/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config errro: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisDSN,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	err = redisClient.Ping().Err()
	if err != nil {
		log.Fatalf("redis error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	redisStorage := storage.NewRedisStorage(redisClient)
	leakyBucket := leakybucket.NewLeakyBucket(
		ctx,
		cfg.LimitLogin,
		cfg.LimitPassword,
		cfg.LimitIP,
		time.Duration(cfg.TTL)*time.Second,
	)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	api.RegisterAntiBruteforceServer(grpcServer, service.NewService(redisStorage, leakyBucket))
	listen, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		log.Fatalf("bind error: %v at address %s", err, cfg.Listen)
	}

	shutdown := make(chan error, 1)
	go func() {
		err := grpcServer.Serve(listen)
		if err != nil {
			shutdown <- err
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case x := <-interrupt:
		cancel()
		log.Printf("Received interrupt `%v`.", x)
	case err := <-shutdown:
		cancel()
		log.Printf("Received shutdown message: %v", err)
	}

	grpcServer.GracefulStop()
}
