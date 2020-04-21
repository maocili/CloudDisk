package rds

import (
	rPool "CloudDisk/tools/redis"

	"github.com/garyburd/redigo/redis"
)

func SaveToken(username, token string) {
	rPool := rPool.RedisPool().Get()
	defer rPool.Close()

	rPool.Do("MULTI")
	rPool.Do("HSET", "TOKEN_"+token, "username", username)
	rPool.Do("EXEC")
}

func QueryToken(token string) (string, error) {
	rPool := rPool.RedisPool().Get()
	defer rPool.Close()

	data, err := redis.String(rPool.Do("HGET", "TOKEN"+token, "username"))
	if err != nil {
		return "", nil
	}

	return data, nil
}
