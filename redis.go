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
	rdb *redis.Client
)

// helps us use the database connection
func GetRedis() *redis.Client {
	return rdb
}

func ConnectToRedis() error {

	// parse the connection string in the .env file under the key "REDIS_URI"
	Options, err := redis.ParseURL(os.Getenv("REDIS_URI"))
	if err != nil {
		return err
	}

	// assign the redis client to the global variable
	rdb = redis.NewClient(Options)

	// send a ping command to redis
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	// no errors occurred return nil
	return nil
}
