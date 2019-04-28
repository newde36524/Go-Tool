package redistool

import (
	"crypto/tls"
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
	c      redis.Conn
	addr   string
	option *RedisClientOption
}

//NewRedisClient 实例化Redis客户端新实例
//@option Redis客户端配置
func NewRedisClient(option *RedisClientOption) *RedisClient {
	return &RedisClient{
		option: option,
	}
}

//Connect 连接Redis服务器 "ip:port"
func (redisClient *RedisClient) Connect(addr string) error {
	redisClient.Close() //可关闭旧连接，连接新地址
	redisClient.addr = addr
	conn, err := redisClient.Clone()
	if err != nil {
		logs.Error(err)
	}
	redisClient.c = conn.c
	return err
}

//Clone 相同的配置创建一个新的RedisClient实例
func (redisClient *RedisClient) Clone() (*RedisClient, error) {
	conn, err := redisClient.create()
	if err != nil {
		return nil, err
	}
	redisClient.c = conn
	return redisClient, err
}

//create 用相同的配置创建一个新的redis连接
func (redisClient *RedisClient) create() (redis.Conn, error) {
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
	if redisClient.c == nil {
		return nil
	}
	return redisClient.c.Close()
}

//Set 设置一个键值对
func (redisClient *RedisClient) Set(key, value string) (string, error) {
	return redis.String(redisClient.c.Do("SET", key, value))
}

//Get 获取一个键值
func (redisClient *RedisClient) Get(key string) (string, error) {
	return redis.String(redisClient.c.Do("GET", key))
}

//HSet 设置一个哈希映射的键值对
func (redisClient *RedisClient) HSet(hash, key, value string) (int64, error) {
	return redis.Int64(redisClient.c.Do("HSET", hash, key, value))
}

//HGet 获取一个哈希映射的键值
func (redisClient *RedisClient) HGet(hash, key string) (string, error) {
	return redis.String(redisClient.c.Do("HGET", hash, key))
}

//LPop 自1.0.0起可用.
//时间复杂度:O(1)
//删除并返回存储在列表中的第一个元素key
func (redisClient *RedisClient) LPop(listName string) (string, error) {
	return redis.String(redisClient.c.Do("LPOP", listName))
}

//LPush .
func (redisClient *RedisClient) LPush(listName, value string) (int, error) {
	return redis.Int(redisClient.c.Do("LPUSH", listName, value))
}

//LRange .
func (redisClient *RedisClient) LRange(listName string, startIndex, endIndex int) ([]string, error) {
	return redis.Strings(redisClient.c.Do("LRANGE", listName, startIndex, endIndex))
}

//RPop .
func (redisClient *RedisClient) RPop(listName string) (string, error) {
	return redis.String(redisClient.c.Do("RPOP", listName))
}

//RPush .
func (redisClient *RedisClient) RPush(listName, value string) (int, error) {
	return redis.Int(redisClient.c.Do("RPUSH", listName, value))
}

//Exists 是否存在
func (redisClient *RedisClient) Exists(key string) (bool, error) {
	return redis.Bool(redisClient.c.Do("EXISTS", key))
}

//Publish 发布
func (redisClient *RedisClient) Publish(channelName, msg string) (int, error) {
	reply, err := redisClient.c.Do("PUBLISH", channelName, msg)
	if err != nil {
		return redis.Int(reply, err)
	}
	err = redisClient.c.Flush()
	return redis.Int(reply, err)
}

//Subscript 订阅
func (redisClient *RedisClient) Subscript(onMessage func(string), channel string) error {
	return redisClient.Subscript2(onMessage, channel, nil)
}

//Subscript2 订阅
func (redisClient *RedisClient) Subscript2(onMessage func(string), channel string, onError func(error)) error {
	conn, err := redisClient.create()
	if err != nil {
		return err
	}
	psc := redis.PubSubConn{Conn: conn}
	err = psc.Subscribe(channel)
	if err == nil {
		go func(psc *redis.PubSubConn) {
			for {
				switch v := psc.Receive().(type) {
				case redis.Message:
					onMessage(string(v.Data))
				case redis.Subscription:
					logs.Infof("%s: %s %d\n", v.Channel, v.Kind, v.Count)
				case error:
					if onError != nil {
						onError(v)
					}
					psc.Close()
					return
				}
			}
		}(&psc)
	}
	return err
}
