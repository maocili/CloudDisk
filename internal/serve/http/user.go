package serve

import (
	"CloudDisk/internal/dao"
	"CloudDisk/internal/model"
	"CloudDisk/internal/rds"
	"CloudDisk/tools"
	"net/http"
	"time"

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
	userinfo := model.UserInfo{
		Uid:      tools.GenerateUID(sigin.Username),
		Username: sigin.Username,
	}
	userinfo.Password = tools.MakePasswd(sigin.PassWord, userinfo.Uid)

	isLogin, err := dao.Login(userinfo)
	if err != nil {
		response.Code = 20001
		response.Msg = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if isLogin {
		token := tools.Sha1(time.Now().String() + userinfo.Uid)
		rds.SaveToken(userinfo.Username, token)
		c.SetCookie("token", token, 3600, "/", "localhost", false, true)

		response.Code = 20000
		response.Msg = "成功"
		response.Data = map[string]string{"SET_TOKEN": token}

		c.JSON(http.StatusOK, response)
		return
	} else {
		response.Code = 20002
		response.Msg = "失败"
		c.JSON(http.StatusForbidden, response)
		return
	}

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

func GetInfo(c *gin.Context) {

	token := c.Query("token")

	infoData := map[string]string{
		"name":   token,
		"avatar": "123",
	}
	response := model.ResponseBody{
		Code: 20000,
		Data: infoData,
	}
	c.JSON(http.StatusOK, response)
}

func Avatar(c *gin.Context) {
	// imageName := c.Query("imageName")
	c.File("./others/avatar.gif")
	return
}
