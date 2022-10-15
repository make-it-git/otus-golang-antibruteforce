package main

import (
	"flag"
	"log"
	"net"

	"github.com/go-redis/redis"
	"github.com/make-it-git/otus-antibruteforce/internal/config"
	storage "github.com/make-it-git/otus-antibruteforce/internal/storage_redis"
)

func main() {
	op := flag.String("operation", "", "whitelistadd, whitelistremove, blacklistadd, blacklistremove, clearlists")
	ip := flag.String("ip", "", "ip subnet for add/remove operation")
	_ = ip
	flag.Parse()

	if op == nil || *op == "" {
		log.Fatal("no operation provided")
	}
	if *op != "" {
		if ip == nil || *ip == "" {
			log.Fatalf("no ip address provided for %s", *op)
		}
		_, _, err := net.ParseCIDR(*ip)
		if err != nil {
			log.Fatal(err)
		}
	}

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

	redisStorage := storage.NewRedisStorage(redisClient)

	switch *op {
	case "whitelistadd":
		err := redisStorage.WhiteListAdd(*ip)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s added to whitelist", *ip)
	case "whitelistremove":
		err := redisStorage.WhiteListRemove(*ip)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s removed from whitelist", *ip)
	case "blacklistadd":
		err := redisStorage.BlackListAdd(*ip)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s added to blacklist", *ip)
	case "blacklistremove":
		err := redisStorage.BlackListRemove(*ip)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s removed from blacklist", *ip)
	case "clearlists":
		err := redisStorage.ClearLists()
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("unknown operation")
	}
}
