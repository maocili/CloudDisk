package dao

import (
	"CloudDisk/internal/model"
	"CloudDisk/tools/mysql"
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
