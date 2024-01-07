/*
 * redis.go
 *
 * This file contains the core component to connect to the redis database, might add more stuff later
 */

package main

import (
	"context"
	"os"
	
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func ConnectToRedis() error {
	
	// parse the connection string in the .env file under the key "redisUrl"
	Options, err := redis.ParseURL(os.Getenv("redisUrl"))
	if err != nil {
		return err
	}
	
	// assign the redis client to the global variable
	rdb = redis.NewClient(Options)
	
	// send a ping command to redis
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	
	// no errors occurred return nil
	return nil
}
