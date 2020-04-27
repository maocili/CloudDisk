package model

type FileData struct {
	Id       int    `xorm:"int(11)" form:"id" json:"id"`
	FileSha1 string `xorm:"pk char(40)" form:"file_sha1" json:"file_sha1"`
	FileSize string `xorm :"varchar(254)" form:"file_size" json:"file_size"`
	FileAddr string `xorm:"varchar(1024)" form:"file_addr" json:"file_addr"`
}

type UserFile struct {
	Id       int    `xorm:"int(11)" form :"id" json:"id"`
	Uid      string `xorm:"varchar(254)" form:"uid" json:"uid"`
	FileSha1 string `xorm:"char(40)" form:"file_sha1" json:"file_sha1"`
	FileName string `xorm:"varchar(254)" form:"file_name" json:"file_name"`
	// UserPath string `xorm:"char(40)" form:"user_path" json:"user_path"`
	PathId string `xorm:"char(254)" form:"path_id" json:"path_id"`
}
