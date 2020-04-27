package dao

import (
	"CloudDisk/internal/model"
	"CloudDisk/tools/mysql"
	"fmt"
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

func IsexistFileHash(filehash string) (bool, error) {
	has, err := mysql.DB.Exist(&model.FileData{
		FileSha1: filehash,
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
