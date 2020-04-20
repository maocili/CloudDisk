package ini

import (
	"fmt"

	"gopkg.in/ini.v1"
)

var (
	//配置信息
	iniFile *ini.File
)

func init() {
	file, e := ini.Load("app.ini")
	if e != nil {
		fmt.Println("Fail to load conf/app.ini" + e.Error())
		return
	}
	iniFile = file
}

func GetSection(sectionName string) *ini.Section {
	section, e := iniFile.GetSection(sectionName)
	if e != nil {
		fmt.Println("未找到对应的配置信息:" + sectionName + e.Error())
		return nil
	}
	return section
}

func GetSectionMap(sectionName string) map[string]string {
	section, e := iniFile.GetSection(sectionName)
	if e != nil {
		fmt.Println("未找到对应的配置信息:" + sectionName + e.Error())
		return nil
	}
	section_map := make(map[string]string, 0)
	for _, e := range section.Keys() {
		section_map[e.Name()] = e.Value()
	}
	return section_map
}

func GetVal(sectionName string, key string) string {
	var temp_val string
	section := GetSection(sectionName)
	if nil != section {
		temp_val = section.Key(key).Value()
	}
	return temp_val
}
