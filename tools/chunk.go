package tools

import (
	"CloudDisk/internal/model"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

// 读取块文件合并到 []byte
func MergeChunk(fileName string, file model.FileData) ([]byte, error) {

	fileSize, _ := strconv.Atoi(file.FileSize)
	chunkCount := fileSize / (5 * 1024 * 1024)

	var fd []byte
	for i := 0; i <= chunkCount; i++ {
		f, err := os.Open(file.FileAddr + "/" + strconv.Itoa(i+1))
		if err != nil {
			fmt.Println("read file fail", err.Error())
			return nil, err
		}
		defer f.Close()

		temp, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println("read to fd fail", err.Error)
			return nil, err
		}
		fd = append(fd, temp...)

	}
	return fd, nil
}
