package rds

import (
	rPool "CloudDisk/tools/redis"
)

func DeleteKey(keys string) {
	rPool := rPool.RedisPool().Get()
	defer rPool.Close()

	rPool.Do("DEL", keys)
}
