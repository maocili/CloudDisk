package dao

import (
	"CloudDisk/internal/model"
	"CloudDisk/tools"
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

func AddFolder(folder model.FolderPath) error {

	folder.PathId = tools.Sha1(folder.FolderName + folder.PrefixPath + folder.Uid)

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
