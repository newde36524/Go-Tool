package filetool

import (
	"bufio"
	"os"
)

// 分段获取指定文件的数据包
// @pagIndex 表示数据包的下标
// @pagSize 表示一次获取的数据包大小
// @filePath 表示文件路径
// @rerurn 读取到的文件数据包  error 表示读取异常
func ReadPagingFile(pagIndex int, pagSize int, filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bufReader := bufio.NewReader(file)
	_, err = bufReader.Discard(pagIndex * pagSize) //跳过指定字节数
	if err != nil {
		return nil, err
	}

	fileData, err := bufReader.Peek(pagSize) //拿到指定字节数的数据
	if err != nil {
		return nil, err
	}
	return fileData[:pagSize], nil
}
