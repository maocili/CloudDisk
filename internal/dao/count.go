package dao

import (
	"CloudDisk/internal/model"
	"CloudDisk/tools/mysql"
)

func FileCount(uid string) int64 {
	userFile := model.UserFile{
		Uid: uid,
	}
	fileCount, _ := mysql.DB.Count(&userFile)
	return fileCount
}

func FolderCount(uid string) int64 {
	userFolder := model.FolderPath{
		Uid: uid,
	}
	folderCount, _ := mysql.DB.Count(&userFolder)
	return folderCount
}

type storageCount struct {
	Storage int64
}

func StorageCount(uid string) int64 {

	scount := storageCount{}
	a, _ := mysql.DB.SQL("SELECT sum(file_data.file_size) as Storage FROM user_file JOIN file_data  ON user_file.file_sha1 = file_data.file_sha1 WHERE user_file.uid = ?", uid).Cols("Storage").Get(&scount)
	if a {
		return scount.Storage / (1024 * 1024)
	} else {
		return 0
	}

}
