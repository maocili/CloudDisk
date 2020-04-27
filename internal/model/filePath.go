package model

type FolderPath struct {
	Id         int    `pk xorm:"int(11)"  `  // 自增长
	Uid        string `xorm:"varchar(256)" ` //用户id
	PrefixPath string `xorm:"varchar(256)  ` // 根文件夹的PathId
	PathId     string `xorm:"varchar(256)" ` // hash(PathName+PrefixPath+uid)
	FolderName string `xorm:"varchar(256)" ` //用户的路径名
}

type FolderList struct {
	PathId     string
	FolderName string
	Children   interface{}
}

type Zones struct {
	Uid   string
	Zones string `json:"zones"`
}

type TreeList struct {
	Name  string `json:"name"`
	Zones string `json:"zones"`
	Leaf  bool   `json:"leaf"`
}
