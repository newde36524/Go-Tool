package redistool

import (
	"crypto/tls"
	"time"

	"github.com/gomodule/redigo/redis"
)

//RedisClientOption Redis客户端配置项
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

//RedisClient Redis客户端
type RedisClient struct {
	c      redis.Conn
	addr   string
	option *RedisClientOption
}

//NewRedisClient 实例化Redis客户端新实例
//@option Redis客户端配置
func NewRedisClient(option *RedisClientOption) *RedisClient {
	var result *RedisClient = &RedisClient{
		option: option,
	}
	return result
}

//Connect 连接Redis服务器 "ip:port"
func (this *RedisClient) Connect(addr string) error {
	this.Close()
	this.addr = addr
	option := this.option
	c, err := redis.Dial("tcp", this.addr, []redis.DialOption{
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
	this.c = c
	return err
}

//Close 关闭Redis客户端
func (this *RedisClient) Close() error {
	if this.c == nil {
		return nil
	}
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
