package dao

import (
	"CloudDisk/internal/model"
	"CloudDisk/tools/mysql"
	"errors"
	"fmt"

	"github.com/go-xorm/xorm"
)

func NewFile(filedata model.FileData) error {
	if _, err := mysql.DB.Insert(&filedata); err != nil {
		return err
	}
	return nil
}

func AddUserFile(userfile model.UserFile) error {

	if _, err := mysql.DB.Insert(&userfile); err != nil {
		return err
	}
	return nil
}

func DeleteFile(zones model.Zones) error {

	_, err := mysql.DB.Transaction(func(session *xorm.Session) (interface{}, error) {

		//检查uid + 路径是否存在
		has, err := session.Exist(&model.UserFile{
			FileSha1: zones.Zones,
			Uid:      zones.Uid,
		})
		switch {
		case has == false:
			return nil, errors.New("文件不存在")
		case err != nil:
			return nil, err
		}

		affected, err := session.Where("uid = ? and file_sha1 = ?", zones.Uid, zones.Zones).Delete(&model.UserFile{})
		return affected, err
	})
	return err
}

func IsexistFileHash(filehash string) (bool, error) {
	has, err := mysql.DB.Exist(&model.FileData{
		FileSha1: filehash,
	})

	return has, err
}

func IsexistUserFile(zones model.Zones) (bool, error) {
	has, err := mysql.DB.Exist(&model.UserFile{
		Uid:      zones.Uid,
		FileSha1: zones.Zones,
	})

	return has, err
}

func IsexistUserFolder(zones model.Zones) (bool, error) {
	has, err := mysql.DB.Exist(&model.FolderPath{
		Uid:    zones.Uid,
		PathId: zones.Zones,
	})

	return has, err
}

func GetFileList(uid string, zones string) ([]model.UserFile, error) {
	var userFile []model.UserFile

	err := mysql.DB.Table("user_file").Select("*").Where("uid = ? and path_id=?", uid, zones).Find(&userFile)

	return userFile, err
}

func GetUserFileName(uid, filehash string) (string, error) {

	var userFile model.UserFile
	var fileName string
	has, err := mysql.DB.Table(&userFile).Where("uid = ? and file_sha1 = ?", uid, filehash).Cols("file_name").Get(&fileName)
	fmt.Println(has, fileName)
	return fileName, err
}

func GetFileData(filehash string) (model.FileData, error) {
	var fileData []model.FileData

	err := mysql.DB.Table("file_data").Select("*").Where("file_sha1=?", filehash).Find(&fileData)

	return fileData[0], err
}
