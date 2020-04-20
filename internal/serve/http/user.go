package serve

import (
	"CloudDisk/internal/dao"
	"CloudDisk/internal/model"
	"CloudDisk/tools"
	"net/http"

	"github.com/gin-gonic/gin"
)

type siginInfo struct {
	Username string `json:"username"`
	PassWord string `json:"password"`
}

type registerInfo struct {
	Username   string `json:"username"`
	PassWord   string `json:"password"`
	RePassWord string `json:"repassword"`
}

func Sigin(c *gin.Context) {
	var sigin siginInfo
	var response model.ResponseBody

	defer c.Abort()

	if err := c.ShouldBindJSON(&sigin); err != nil {
		response.Code = 20001
		response.Msg = err.Error()
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response.Code = 20000
	response.Msg = sigin.Username + ":" + sigin.PassWord
	c.JSON(http.StatusOK, response)
}

func Register(c *gin.Context) {
	var register registerInfo
	var response model.ResponseBody

	defer c.Abort()

	if err := c.ShouldBindJSON(&register); err != nil {
		response.Code = 20001
		response.Msg = err.Error()
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	if register.PassWord != register.RePassWord {
		response.Code = 20001
		response.Msg = "两次密码不正确"
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}
	var userinfo model.UserInfo

	userinfo.Uid = tools.GenerateUID(register.Username)
	userinfo.Username = register.Username
	userinfo.Password = tools.MakePasswd(register.PassWord, userinfo.Uid)

	if err := dao.Signup(userinfo); err != nil {
		response.Code = -1
		response.Msg = err.Error()
		c.JSONP(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response.Code = 20000
	response.Msg = "注册成功"
	c.JSON(http.StatusOK, response)
}

func Isexist(c *gin.Context) {
	var response model.ResponseBody
	var username siginInfo

	if err := c.ShouldBindJSON(&username); err != nil {
		response.Code = 20001
		response.Msg = err.Error()
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	isUsername, err := dao.IsexistUsername(username.Username)
	if err != nil {
		response.Code = 20002
		response.Msg = err.Error()
		c.JSON(http.StatusOK, response)
		c.Abort()
		return
	}
	if isUsername {
		response.Msg = "true"
	} else {
		response.Msg = "false"
	}
	response.Code = 20000
	c.JSONP(http.StatusOK, response)
	return

}
