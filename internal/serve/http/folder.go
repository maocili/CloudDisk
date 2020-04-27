package serve

import (
	"CloudDisk/internal/dao"
	"CloudDisk/internal/folder"
	"CloudDisk/internal/model"
	"CloudDisk/internal/rds"
	"CloudDisk/tools"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetFolderList(c *gin.Context) {
	var zones model.Zones
	var response model.ResponseBody

	err := c.ShouldBindJSON(&zones)
	if err != nil {
		response.Code = 20001
		response.Msg = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}
	token, err := c.Cookie("token")
	zones.Uid, _ = rds.QueryTokenUid(token)

	treeData, err := folder.GetTreeList(zones)
	if err != nil {
		response.Code = 20001
		response.Msg = err.Error()
		response.Data = treeData
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response.Code = 20000
	response.Data = treeData
	c.JSON(http.StatusOK, response)
}

func DeleteFile(c *gin.Context) {
	var zones model.Zones
	var response model.ResponseBody

	err := c.ShouldBindJSON(&zones)
	if err != nil {
		response.Code = 20001
		response.Msg = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}
	token, err := c.Cookie("token")
	zones.Uid, _ = rds.QueryTokenUid(token)

	//删除文件夹或文件
	err = folder.DeleteZones(zones)
	if err != nil {
		response.Code = 20001
		response.Msg = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response.Code = 20000
	response.Msg = "删除成功"
	c.JSON(http.StatusOK, response)
}

type addFolder struct {
	Uid        string
	Zones      string `json:"zones"`
	FolderName string `json:"folderName"`
}

func AddFolder(c *gin.Context) {
	var newFolder addFolder
	var response model.ResponseBody

	if err := c.ShouldBindJSON(&newFolder); err != nil {
		response.Code = 20001
		response.Msg = err.Error()
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}
	token, _ := c.Cookie("token")
	newFolder.Uid, _ = rds.QueryTokenUid(token)

	folder := model.FolderPath{
		Uid:        newFolder.Uid,
		FolderName: newFolder.FolderName,
		PrefixPath: newFolder.Zones,
	}
	folder.PathId = tools.Sha1(folder.FolderName + folder.PrefixPath + folder.Uid)

	if err := dao.AddFolder(folder); err != nil {
		response.Code = 20001
		response.Msg = err.Error()
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	response.Code = 20000
	response.Msg = folder.FolderName + "新建成功"
	c.JSON(http.StatusOK, response)
	return
}

func Download(c *gin.Context) {
	filehash := c.Query("filehash")
	token, _ := c.Cookie("token")
	uid, _ := rds.QueryTokenUid(token)

	//获取用户的文件名
	fileName, err := dao.GetUserFileName(uid, filehash)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	// 获取文件信息
	fileData, err := dao.GetFileData(filehash)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	// 获得全部的chunk
	byteFile, err := tools.MergeChunk(fileName, fileData)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Writer.Write([]byte(err.Error()))
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/text/plain")
	c.Header("Accept-Length", fmt.Sprintf("%d", fileData.FileSize))
	c.Writer.Write(byteFile)

}
