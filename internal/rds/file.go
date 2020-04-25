package rds

import (
	rPool "CloudDisk/tools/redis"
	"strconv"
	"strings"

	"github.com/gomodule/redigo/redis"
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

// 查询上传状态
func QueryUploadStatus(uploadId string) (bool, string, string, error) {
	var chunkCount, totalCount int
	var filehash, filesize string
	rPool := rPool.RedisPool().Get()
	defer rPool.Close()

	data, err := redis.Values(rPool.Do("HGETALL", "MP_"+uploadId))
	if err != nil {
		return false, filehash, filesize, err
	}

	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
		if k == "filehash" {
			filehash = v
		}
		if k == "filesize" {
			filesize = v
		}

	}
	return totalCount == chunkCount, filehash, filesize, nil

}

func DelUpload(UploadId string) {
	rPool := rPool.RedisPool().Get()
	defer rPool.Close()

	rPool.Do("MULTI")
	rPool.Do("DEL", "MP_"+UploadId)
	rPool.Do("EXEC")
}
