package Storage

import (
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

// RedisInst 实例
var RedisInst *redisStorage

// func init() {
// 	RedisInst = &redisStorage{}
// 	RedisInst.OpenRedis("127.0.0.1:6382", 0)
// }

//RedisStorage 连接Redis
type redisStorage struct {
	c redis.Conn
}

// OpenRedis opens redis as entity storage
func (me *redisStorage) OpenRedis(url string, dbindex int) error {
	c, err := redis.DialURL(url)
	if err != nil {
		return errors.Wrap(err, "redis dail failed")
	}

	if dbindex >= 0 {
		if _, err := c.Do("SELECT", dbindex); err != nil {
			return errors.Wrap(err, "redis select db failed")
		}
	}
	me.c = c
	return nil
}

// Write 存入redis
func (me *redisStorage) Write(typeName string, key string, personalID int, data []byte) error {
	if personalID <= 0 {
		personalID = 100000
	}
	_, err := me.c.Do("HSET", typeName, personalID, data)
	return err
}
