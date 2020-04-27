package router

import (
	"CloudDisk/internal/middleware"
	serve "CloudDisk/internal/serve/http"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DiskRouter(r *gin.Engine) {
	r.GET("/ping", func(context *gin.Context) {
		context.JSONP(http.StatusOK, gin.H{
			"code": 20000,
			"msg":  "pong",
		})
	})

	user := r.Group("/user")
	{
		user.POST("/login", serve.Sigin)
		user.POST("/register", serve.Register)
		user.POST("/isexist", serve.Isexist)
		user.GET("/info", serve.GetInfo)
		user.GET("/avatar", serve.Avatar)

	}

	file := r.Group("/file").Use(middleware.VailToken())
	{
		file.POST("/upload", serve.Upload)
		file.POST("/delete", serve.DeleteFile)
		file.POST("/chunk/init", serve.UploadInit)
		file.POST("/chunk/upload", serve.ChunkUpload)
		file.POST("/chunk/finish", serve.FinishUploadHandler)
	}
	r.GET("/download", serve.Download)

	folder := r.Group("folder").Use(middleware.VailToken())
	{
		folder.POST("/list", serve.GetFolderList)
	}

}
