// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package g

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

var RedisConnPool *redis.Pool

func InitRedisConnPool() {
	auth, dsn := formatRedisAddr(viper.GetString("redis.dsn"))
	maxIdle := viper.GetInt("redis.max_idle")
	idleTimeout := 240 * time.Second

	connTimeout := time.Duration(viper.GetInt("redis.conn_timeout")) * time.Millisecond
	readTimeout := time.Duration(viper.GetInt("redis.read_timeout")) * time.Millisecond
	writeTimeout := time.Duration(viper.GetInt("redis.write_timeout")) * time.Millisecond

	RedisConnPool = &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", dsn, connTimeout, readTimeout, writeTimeout)
			if err != nil {
				return nil, err
			}
			if auth != "" {
				if _, err := c.Do("AUTH", auth); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: PingRedis,
	}
}

func formatRedisAddr(addrConfig string) (string, string) {
	if redisAddr := strings.Split(addrConfig, "@"); len(redisAddr) == 1 {
		return "", redisAddr[0]
	} else {
		return strings.Join(redisAddr[0:len(redisAddr)-1], "@"), redisAddr[len(redisAddr)-1]
	}
}

func PingRedis(c redis.Conn, t time.Time) error {
	_, err := c.Do("ping")
	if err != nil {
		log.Infoln("[ERROR] ping redis fail", err)
	}
	return err
}
