package serve

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type fileInfo struct {
	FileName string
	FileSize string
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
