package filetool

import (
	"bufio"
	"os"
	"path/filepath"
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
	var result []byte
	for i := 0; i < pagSize; i++ {
		b, _ := bufReader.ReadByte()
		result = append(result, b) //读取超出接线默认用0补足
	}
	return result, nil
}

//获取文件大小
func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}
