package storage

import (
	"net"

	"github.com/go-redis/redis"
	"github.com/make-it-git/otus-antibruteforce/internal/contract"
)

type RedisStorage struct {
	client *redis.Client
}

const blackType = "B"

const whiteType = "W"

func NewRedisStorage(redisClient *redis.Client) *RedisStorage {
	return &RedisStorage{
		client: redisClient,
	}
}

func (s *RedisStorage) BlackListAdd(netAddr string) error {
	err := s.client.SAdd(blackType, netAddr).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisStorage) BlackListRemove(netAddr string) error {
	err := s.client.SRem(blackType, netAddr).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisStorage) WhiteListAdd(netAddr string) error {
	err := s.client.SAdd(whiteType, netAddr).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisStorage) WhiteListRemove(netAddr string) error {
	err := s.client.SRem(whiteType, netAddr).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisStorage) ClearLists() error {
	err := s.client.Del(whiteType).Err()
	if err != nil {
		return err
	}

	err = s.client.Del(blackType).Err()

	return err
}

func (s *RedisStorage) GetStatus(ip net.IP) (contract.NetAddrStatus, error) {
	blackList, err := s.client.SMembers(blackType).Result()
	if err != nil {
		return contract.Unknown, err
	}
	for _, blackNet := range blackList {
		_, ipNet, err := net.ParseCIDR(blackNet)
		if err != nil {
			return contract.Unknown, err
		}

		if ipNet.Contains(ip) {
			return contract.Blacklisted, nil
		}
	}

	whiteList, err := s.client.SMembers(whiteType).Result()
	if err != nil {
		return contract.Unknown, err
	}
	for _, whiteNet := range whiteList {
		_, ipNet, err := net.ParseCIDR(whiteNet)
		if err != nil {
			return contract.Unknown, err
		}

		if ipNet.Contains(ip) {
			return contract.Whitelisted, nil
		}
	}

	return contract.Unknown, nil
}
