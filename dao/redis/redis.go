package redis

import (
	"bluebellproject/setting"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var redisConn *redis.Pool

func Init(RedisConf *setting.RedisConfig) error {
	redisConn = &redis.Pool{
		MaxIdle:     RedisConf.MaxIdleConn,
		MaxActive:   RedisConf.MaxOpenConn,
		IdleTimeout: RedisConf.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d",
				RedisConf.Host,
				RedisConf.Port))
			if err != nil {
				return nil, err
			}
			if RedisConf.Password != "" {
				if _, err := c.Do("AUTH", RedisConf.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

func Close() {
	redisConn.Close()
}

func Set(key string, data interface{}, time int) error {
	conn := redisConn.Get()
	defer conn.Close()
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}
	return nil
}

func Exist(key string) bool {
	conn := redisConn.Get()
	defer conn.Close()
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

func Get(key string) ([]byte, error) {
	conn := redisConn.Get()
	defer conn.Close()
	value, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return value, err
}
func Delete(key string) (bool, error) {
	conn := redisConn.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(key string) error {
	conn := redisConn.Get()
	defer conn.Close()
	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}
	for _, key := range keys {
		_, err := Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}
