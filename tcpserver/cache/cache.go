package cache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"go-entry-task/tcpserver/config"
	"go-entry-task/tcpserver/util"
	"log"
	"time"
)

var rdb *redis.Client

func init() {
	conf := config.DefaultRedisConfig()
	rdb = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		PoolSize: conf.PoolConn,
		//ReadTimeout: time.Millisecond * time.Duration(600),
		//WriteTimeout: time.Millisecond * time.Duration(600),
		IdleTimeout: time.Second * time.Duration(conf.IdleTimeout),
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		log.Println("Failed to connect to db entry_task, err: ", err.Error())
	}
}

//从redis中查询用户信息
func QueryUser(username string) (string, string, string, error) {
	result, err := rdb.Get(username).Result()
	hash := make(map[string]string)
	if err == nil {
		err = json.Unmarshal([]byte(result), &hash)
		nick := hash["nick"]
		passwd := hash["pass"]
		url := hash["url"]
		return nick, passwd, url, nil
	}
	return "", "", "", util.ErrRedisNotExist
}

//查询redis中是否有token
func QueryToken(token string, user string) bool {
	result, err := rdb.Get(token).Result()
	if err != nil || result != user {
		return false
	}
	return true
}

//把用户信息插入redis
func Set(key string, val interface{}, expire int) error {
	err := rdb.Set(key, val, time.Second*time.Duration(expire)).Err()
	return err
}

func Del(key string) error {
	err := rdb.Del(key).Err()
	return err
}
