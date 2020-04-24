package rds

import (
	rPool "CloudDisk/tools/redis"

	"github.com/garyburd/redigo/redis"
)

func QueryUploadIdHash(uploadId string) (string, error) {
	rPool := rPool.RedisPool().Get()
	defer rPool.Close()

	filehash, err := redis.String(rPool.Do("HGET", "MP_"+uploadId, "filehash"))
	if err != nil {
		return "", nil
	}

	return filehash, nil
}

func QueryUploadIdSize(uploadId string) (string, error) {
	rPool := rPool.RedisPool().Get()
	defer rPool.Close()

	filesize, err := redis.String(rPool.Do("HGET", "MP_"+uploadId, "filesize"))
	if err != nil {
		return "", nil
	}

	return filesize, nil
}

func DelUpload(UploadId string) {
	rPool := rPool.RedisPool().Get()
	defer rPool.Close()

	rPool.Do("MULTI")
	rPool.Do("DEL", "MP_"+UploadId)
	rPool.Do("EXEC")
}
