package cache

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
)

// Clients 保存了所有的客户端访问原始数据
var Clients map[string]*UserClient

func init() {
	Clients = make(map[string]*UserClient, 20)
}

// UserClient 用户端
type UserClient struct {
	Host      string
	URL       string
	Token     string
	Timestamp int64
	MaxAge    int64
}

// NewUserClient new client
func NewUserClient(url, host string) *UserClient {
	var c UserClient
	c.Token = uuid.New().String()
	c.Timestamp = time.Now().UnixNano()
	c.URL = url
	c.Host = host
	Clients[c.Token] = &c
	return &c
}

// Refresh 刷新驻留时间
func (c *UserClient) Refresh() {
	c.MaxAge = time.Now().UnixNano() - c.Timestamp
}

// Release 释放
func (c *UserClient) Release() {
	rd := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       2,  // use default DB
	})
	var totalAge, views int
	totalAgeT, err := rd.HGet("watcher.total_age."+c.Host, c.URL).Result()
	if err == nil {
		totalAge, _ = strconv.Atoi(totalAgeT)
	}
	viewsT, err := rd.HGet("watcher.views."+c.Host, c.URL).Result()
	if err == nil {
		views, _ = strconv.Atoi(viewsT)
	}
	viewData, err := rd.HGet("watcher.view_data."+c.Host, c.URL).Result()
	rd.HSet("watcher.total_age."+c.Host, c.URL, int64(totalAge)+c.MaxAge)
	rd.HSet("watcher.views."+c.Host, c.URL, views+1)
	rd.HSet("watcher.view_data."+c.Host, c.URL, viewData+","+strconv.Itoa(int(c.MaxAge)))
	rd.HSet("watcher.avr_age."+c.Host, c.URL, math.Floor(float64(int64(totalAge)+c.MaxAge)/(float64(views+1)*1000000)))
	delete(Clients, c.Token)
	c = nil
}

// Watch 监听是否过期
func (c *UserClient) Watch() {
	t := time.NewTicker(time.Second * 5)
	defer t.Stop()
	for range t.C {
		if time.Now().UnixNano()-c.Timestamp > 10000000000+c.MaxAge {
			fmt.Println("release-", c.URL)
			c.Release()
			break
		}
	}
}

// GetUserClientByToken get client
func GetUserClientByToken(token string) *UserClient {
	if token == "" || Clients[token] == nil {
		return nil
	}
	return Clients[token]
}
