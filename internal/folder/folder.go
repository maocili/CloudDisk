package folder

import (
	"CloudDisk/internal/dao"
	"CloudDisk/internal/model"
)

func GetTreeList(zones model.Zones) ([]model.TreeList, error) {
	var treeData []model.TreeList

	folderList, err := dao.GetFolderList(zones.Uid, zones.Zones)
	if err != nil {
		return treeData, err
	}

	fileList, err := dao.GetFileList(zones.Uid, zones.Zones)
	if err != nil {
		return treeData, err
	}

	for i := 0; i < len(folderList); i++ {
		tempList := model.TreeList{
			Name:  folderList[i].FolderName,
			Zones: folderList[i].PathId,
		}
		treeData = append(treeData, tempList)
	}

	for i := 0; i < len(fileList); i++ {
		tempList := model.TreeList{
			Name:  fileList[i].FileName,
			Zones: fileList[i].FileSha1,
			Leaf:  true,
		}
		treeData = append(treeData, tempList)
	}

	return treeData, nil
}
