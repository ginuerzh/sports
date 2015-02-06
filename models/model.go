// model
package models

import (
	"github.com/garyburd/redigo/redis"
	"labix.org/v2/mgo"
	"log"
	"time"
)

var (
	mgoSession *mgo.Session
	pool       *redis.Pool
)

func Init() {
	var err error
	mgoSession, err = mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				log.Println(err)
				return nil, err
			}
			/*
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			*/
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

type Session struct {
	mgoSess   *mgo.Session
	redisConn redis.Conn
}

func (s *Session) Clone() *Session {
	return &Session{
		mgoSess:   mgoSession.Clone(),
		redisConn: pool.Get(),
	}
}

func (m *Session) Close() error {
	m.mgoSess.Close()
	return m.redisConn.Close()
}
