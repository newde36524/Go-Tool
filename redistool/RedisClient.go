package redistool

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisClientOption struct {
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	ConnectTimeout time.Duration
	KeepAlive      time.Duration
	Database       int
	Password       string
	ClientName     string
	TLSConfig      *tls.Config
	TLSSkipVerify  bool
	UseTLS         bool
}

type RedisClient struct {
	c redis.Conn
}

func (this *RedisClient) Login(addr string, option *RedisClientOption) {
	c, err := redis.Dial("tcp", addr, []redis.DialOption{
		redis.DialReadTimeout(option.ReadTimeout),
		redis.DialWriteTimeout(option.WriteTimeout),
		redis.DialConnectTimeout(option.ConnectTimeout),
		redis.DialKeepAlive(option.KeepAlive),
		redis.DialDatabase(option.Database),
		redis.DialPassword(option.Password),
		redis.DialClientName(option.ClientName),
		redis.DialTLSConfig(option.TLSConfig),
		redis.DialTLSSkipVerify(option.TLSSkipVerify),
		redis.DialUseTLS(option.UseTLS),
	}...)
	if err != nil {
		fmt.Println("Connect to redis error:", err)
		return
	}
	this.c = c
}

func (this *RedisClient) Close() error {
	return this.c.Close()
}
func (this *RedisClient) Set(key, value string) (string, error) {
	return redis.String(this.c.Do("SET", key, value))
}
func (this *RedisClient) Get(key string) (string, error) {
	return redis.String(this.c.Do("GET", key))
}
func (this *RedisClient) HSet(hash, key, value string) (int64, error) {
	return redis.Int64(this.c.Do("HSET", hash, key, value))
}
func (this *RedisClient) HGet(hash, key string) (string, error) {
	return redis.String(this.c.Do("HGET", hash, key))
}

//自1.0.0起可用.
//时间复杂度:O(1)
//删除并返回存储在列表中的第一个元素key
func (this *RedisClient) LPop(listName string) (string, error) {
	return redis.String(this.c.Do("LPOP", listName))
}
func (this *RedisClient) LPush(listName, value string) (int, error) {
	return redis.Int(this.c.Do("LPUSH", listName, value))
}
func (this *RedisClient) LRange(listName string, startIndex, endIndex int) ([]string, error) {
	return redis.Strings(this.c.Do("LRANGE", listName, startIndex, endIndex))
}
func (this *RedisClient) RPop(listName string) (string, error) {
	return redis.String(this.c.Do("RPOP", listName))
}
func (this *RedisClient) RPush(listName, value string) (int, error) {
	return redis.Int(this.c.Do("RPUSH", listName, value))
}
