package serve

import (
	"CloudDisk/internal/dao"
	"CloudDisk/internal/model"
	"CloudDisk/internal/rds"
	rPool "CloudDisk/tools/redis"
	"fmt"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

type fileInfo struct {
	FileName string
	FileSize string
}

type initInfo struct {
	Token    string
	FileSize int64  `json:"fileSize"`
	FilePath string `json:"filePath"`
	FileHash string `json:"fileHash"`
	FileName string `json:"fileName"`
}

type multiPart struct {
	UploadId   string
	ChunkSize  int
	ChunkCount int
}

type chunkInfo struct {
	UploadId   string
	Token      string
	ChunkIndex string
	ChunkFile  *multipart.FileHeader
}

type completeUploadInfo struct {
	Token    string
	UploadId string `json:"uploadId"`
	FileName string `json:"fileName`
	FilePath string `json:"filePath`
	FileHash string
}

func Upload(c *gin.Context) {
	// 单文件
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("FormFile", err.Error())
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", gin.H{
			"msg": err.Error(),
		}))
		return
	}

	// 上传文件至指定目录
	c.SaveUploadedFile(file, "./file/"+file.Filename)

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", gin.H{
		"msg": file.Filename,
	}))

	return
}

func UploadInit(c *gin.Context) {

	var initinfo initInfo
	var response model.ResponseBody

	if err := c.ShouldBindJSON(&initinfo); err != nil {
		response.Code = 20001
		response.Msg = err.Error()
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}
	initinfo.Token, _ = c.Cookie("token")

	//秒传
	has, err := dao.IsexistFileHash(initinfo.FileHash)
	if err != nil {
		response.Code = 20001
		response.Msg = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//TODO: 封装秒传
	if has {
		uid, err := rds.QueryTokenUid(initinfo.Token)
		if err != nil {
			response.Code = 20001
			response.Msg = "重新登录尝试"
			response.Data = initinfo
			c.JSON(http.StatusBadRequest, response)
			c.Abort()
			return
		}
		userFile := model.UserFile{
			Uid:      uid,
			FileSha1: initinfo.FileHash,
			FileName: initinfo.FileName,
			UserPath: initinfo.FilePath,
		}

		if err := dao.AddUserFile(userFile); err != nil {
			response.Code = 20001
			response.Msg = err.Error()
			response.Data = userFile
			c.JSON(http.StatusBadRequest, response)
			c.Abort()
			return
		}

		response.Code = 20000
		response.Msg = "秒传成功"
		c.JSON(http.StatusOK, response)
		return
	}

	// redis 初始化分块上传信息
	rPool := rPool.RedisPool().Get()
	defer rPool.Close()

	multipart := multiPart{
		UploadId:   initinfo.Token + fmt.Sprint(time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int(math.Ceil(float64(float64(initinfo.FileSize) / float64(5*1024*1024)))),
	}

	rPool.Do("HSET", "MP_"+multipart.UploadId, "chunkcount", multipart.ChunkCount)
	rPool.Do("HSET", "MP_"+multipart.UploadId, "filehash", initinfo.FileHash)
	rPool.Do("HSET", "MP_"+multipart.UploadId, "filesize", initinfo.FileSize)

	response.Code = 20000
	response.Data = multipart
	c.JSON(http.StatusOK, response)
}

func ChunkUpload(c *gin.Context) {

	var response model.ResponseBody
	// 分块文件
	file, err := c.FormFile("file")
	uploadId, _ := c.GetPostForm("uploadid")
	chunkindex, _ := c.GetPostForm("chunkindex")

	var chunk = chunkInfo{
		UploadId:   uploadId,
		ChunkIndex: chunkindex,
		ChunkFile:  file,
	}

	if err != nil {
		response.Code = 20000
		response.Msg = err.Error()
		c.JSON(http.StatusOK, response)
		return
	}
	go func() {
		fpath := "file/" + chunk.UploadId + "/" + chunk.ChunkIndex
		os.MkdirAll(path.Dir(fpath), 0744)
		fd, err := os.Create(fpath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer fd.Close()
		// 上传文件至指定目录
		c.SaveUploadedFile(chunk.ChunkFile, fpath)
	}()

	rPool := rPool.RedisPool().Get()
	defer rPool.Close()
	rPool.Do("HSET", "MP_"+chunk.UploadId, "chkidx_"+chunk.ChunkIndex, 1)

	response.Code = 20000
	response.Msg = chunk.UploadId + ":" + chunk.ChunkIndex + "is ok"
	c.JSON(http.StatusOK, response)

	return
}

func FinishUploadHandler(c *gin.Context) {

	var completeInfo completeUploadInfo
	var response model.ResponseBody

	if err := c.ShouldBindJSON(&completeInfo); err != nil {
		response.Code = 20001
		response.Msg = err.Error()
		response.Data = completeInfo
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}
	completeInfo.Token, _ = c.Cookie("token")

	filehash, err := rds.QueryUploadIdHash(completeInfo.UploadId)
	if err != nil || filehash == "" {
		response.Code = 20001
		response.Msg = "Hash:" + "Uploadid不存在"
		response.Data = completeInfo
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	filesize, err := rds.QueryUploadIdSize(completeInfo.UploadId)
	if err != nil {
		response.Code = 20001
		response.Msg = "Size:" + err.Error()
		response.Data = completeInfo
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}
	filedata := model.FileData{
		FileSha1: filehash,
		FileSize: filesize,
		FileAddr: "./file/" + completeInfo.UploadId, // 系统存储路径
	}
	//更新文件
	if err := dao.NewFile(filedata); err != nil {
		response.Code = 20001
		response.Msg = err.Error()
		response.Data = filehash
	}

	// 更新用户文件
	// username, err := rds.QueryToken(completeInfo.Token)
	uid, err := rds.QueryTokenUid(completeInfo.Token)
	if err != nil {
		response.Code = 20001
		response.Msg = "重新登录尝试"
		response.Data = completeInfo
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}
	userFile := model.UserFile{
		Uid:      uid,
		FileSha1: filehash,
		FileName: completeInfo.FileName,
		UserPath: completeInfo.FilePath,
	}

	if err := dao.AddUserFile(userFile); err != nil {
		response.Code = 20001
		response.Msg = err.Error()
		response.Data = userFile
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	rds.DelUpload(completeInfo.UploadId)

	response.Code = 20000
	response.Msg = completeInfo.FileName + "上传完成"
	response.Data = completeInfo
	c.JSON(http.StatusOK, response)

}
