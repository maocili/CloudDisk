package rds

import (
	rPool "CloudDisk/tools/redis"

	"github.com/gomodule/redigo/redis"
)

func SaveToken(uid, username, token string) {
	rPool := rPool.RedisPool().Get()
	defer rPool.Close()

	rPool.Do("MULTI")
	rPool.Do("HSET", "TOKEN_"+token, "username", username)
	rPool.Do("HSET", "TOKEN_"+token, "uid", uid)
	rPool.Do("EXEC")
}

func QueryToken(token string) (string, error) {
	rPool := rPool.RedisPool().Get()
	defer rPool.Close()

	data, err := redis.String(rPool.Do("HGET", "TOKEN_"+token, "username"))
	if err != nil {
		return "", nil
	}

	return data, nil
}

func QueryTokenUid(token string) (string, error) {
	rPool := rPool.RedisPool().Get()
	defer rPool.Close()

	data, err := redis.String(rPool.Do("HGET", "TOKEN_"+token, "uid"))
	if err != nil {
		return "", nil
	}

	return data, nil
}
