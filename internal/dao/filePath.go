package dao

import (
	"CloudDisk/internal/model"
	"CloudDisk/tools/mysql"
	"errors"

	"github.com/go-xorm/xorm"
)

func GetFolderList(uid string, prefixPath string) ([]model.FolderPath, error) {
	// var folderList model.FolderList
	var folderPath []model.FolderPath

	err := mysql.DB.Table("folder_path").Select("*").Where("uid = ? and prefix_path=?", uid, prefixPath).Find(&folderPath)

	return folderPath, err
}

//新建文件夹
func AddFolder(folder model.FolderPath) error {

	// folder.PathId = tools.Sha1(folder.FolderName + folder.PrefixPath + folder.Uid)

	_, err := mysql.DB.Transaction(func(session *xorm.Session) (interface{}, error) {

		//检查账户是否已存在
		hasFolder, err := session.Table(&folder).Where("path_id = ?", folder.PathId).Exist()
		switch {
		case hasFolder == true:
			return nil, errors.New("路径已存在")
		case err != nil:
			return nil, err
		}

		if _, err := session.Insert(&folder); err != nil {
			println(err.Error())
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func DeleteFolder(zones model.Zones) error {

	if zones.Zones == "root" {
		return errors.New("无法删除root")
	}

	_, err := mysql.DB.Transaction(func(session *xorm.Session) (interface{}, error) {

		//检查uid + 路径是否存在
		has, err := session.Exist(&model.FolderPath{
			Uid:    zones.Uid,
			PathId: zones.Zones,
		})
		switch {
		case has == false:
			return nil, errors.New("文件夹不存在")
		case err != nil:
			return nil, err
		}

		affected, err := session.Where("uid = ? and path_id = ?", zones.Uid, zones.Zones).Delete(&model.FolderPath{})
		return affected, err
	})
	return err
}
