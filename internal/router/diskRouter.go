package router

import (
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

	file := r.Group("/file")
	{
		file.POST("/upload", serve.Upload)
	}

	// sms := r.Group("/sms")
	// {
	// 	sms.POST("/sendcode", server.SendCode)
	// }
	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("phonenumber", vali.ValidatorPhoneNumber)
	// }
}
