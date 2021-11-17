package planner

import (
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
)

func init() {
	redisHost, okRS := os.LookupEnv("REDIS_SERVER")
	redisPass, okRP := os.LookupEnv("REDIS_PASSWORD")

	if okRS {
		if okRP {
			rdb = redis.NewClient(&redis.Options{
				Addr:     redisHost,
				Password: redisPass,
				DB:       0, // use default DB
			})
		} else {
			rdb = redis.NewClient(&redis.Options{
				Addr:     redisHost,
				Password: "", // no password set
				DB:       0,  // use default DB
			})
		}
	} else {
		if okRP {
			rdb = redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
				Password: redisPass,
				DB:       0, // use default DB
			})
		} else {
			rdb = redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
				Password: "", // no password set
				DB:       0,  // use default DB
			})
		}
	}

}
