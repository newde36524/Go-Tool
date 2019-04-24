package redistool

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/issue9/logs"

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
	C      redis.Conn
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
func (redisClient *RedisClient) Connect(addr string) error {
	redisClient.Close()
	redisClient.addr = addr
	conn, err := redisClient.Clone()
	if err != nil {
		logs.Error(err)
	}
	redisClient.C = conn
	return err
}
func (redisClient *RedisClient) Clone() (redis.Conn, error) {
	if redisClient.addr == "" {
		panic("服务端地址不能为空")
	}
	option := redisClient.option
	return redis.Dial("tcp", redisClient.addr, []redis.DialOption{
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
}

//Close 关闭Redis客户端
func (redisClient *RedisClient) Close() error {
	if redisClient.C == nil {
		return nil
	}
	return redisClient.C.Close()
}
func (redisClient *RedisClient) Set(key, value string) (string, error) {
	return redis.String(redisClient.C.Do("SET", key, value))
}
func (redisClient *RedisClient) Get(key string) (string, error) {
	return redis.String(redisClient.C.Do("GET", key))
}
func (redisClient *RedisClient) HSet(hash, key, value string) (int64, error) {
	return redis.Int64(redisClient.C.Do("HSET", hash, key, value))
}
func (redisClient *RedisClient) HGet(hash, key string) (string, error) {
	return redis.String(redisClient.C.Do("HGET", hash, key))
}

//自1.0.0起可用.
//时间复杂度:O(1)
//删除并返回存储在列表中的第一个元素key
func (redisClient *RedisClient) LPop(listName string) (string, error) {
	return redis.String(redisClient.C.Do("LPOP", listName))
}
func (redisClient *RedisClient) LPush(listName, value string) (int, error) {
	return redis.Int(redisClient.C.Do("LPUSH", listName, value))
}
func (redisClient *RedisClient) LRange(listName string, startIndex, endIndex int) ([]string, error) {
	return redis.Strings(redisClient.C.Do("LRANGE", listName, startIndex, endIndex))
}
func (redisClient *RedisClient) RPop(listName string) (string, error) {
	return redis.String(redisClient.C.Do("RPOP", listName))
}
func (redisClient *RedisClient) RPush(listName, value string) (int, error) {
	return redis.Int(redisClient.C.Do("RPUSH", listName, value))
}

func (redisClient *RedisClient) Exists(key string) (bool, error) {
	return redis.Bool(redisClient.C.Do("EXISTS", key))
}

//发布
func (redisClient *RedisClient) Publish(channelName, msg string) {
	redisClient.C.Do("PUBLISH", channelName, msg)
	redisClient.C.Flush()
}

//订阅
func (redisClient *RedisClient) Subscript(onMessage func(string), channel string) {
	conn, err := redisClient.Clone()
	if err != nil {
		logs.Error(err)
	}
	psc := redis.PubSubConn{Conn: conn}
	psc.Subscribe(channel)
	go func(psc *redis.PubSubConn) {
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				onMessage(string(v.Data))
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				fmt.Println(v)
				psc.Close()
				break
			}
		}
	}(&psc)
}
