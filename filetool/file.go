package filetool

import (
	"bufio"
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ReadPagingFile 分段获取指定文件的数据包
// @pagIndex 表示数据包的下标
// @pagSize 表示一次获取的数据包大小
// @filePath 表示文件路径
// @rerurn 读取到的文件数据包 @[]byte 读取到的数据 @int 有效数据长度  @error 表示读取异常
func ReadPagingFile(pagIndex int, pagSize int, off int, filePath string) ([]byte, int, error) {
	var resultData []byte   //返回数据
	var resultDataSize int  //有效数据长度
	var isRemenberSize bool //是否已记住
	file, err := os.Open(filePath)
	if err != nil {
		return nil, resultDataSize, err
	}
	defer file.Close()
	bufReader := bufio.NewReader(file)
	bufReader.Discard(off)
	_, err = bufReader.Discard(pagIndex * pagSize) //跳过指定字节数
	if err != nil {
		return nil, resultDataSize, err
	}
	for i := 0; i < pagSize; i++ {
		b, e := bufReader.ReadByte()
		if e != nil {
			if !isRemenberSize {
				isRemenberSize = true //只记录一次
				resultDataSize = i
			}
		}
		resultData = append(resultData, b) //读取超出接线默认用0补足
	}
	return resultData, resultDataSize, nil
}

// GetFileSize 获取文件大小
// @filePath 文件路径
func GetFileSize(filePath string) int64 {
	var result int64
	filepath.Walk(filePath, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

// GetFileSize2 获取文件大小
func GetFileSize2(filePath string) (n int64) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}

//GetFileModTime 检测文件更新时间
//@filePath 文件路径
func GetFileModTime(filePath string) int64 {
	f, err := os.Open(filePath)
	if err != nil {
		log.Println("open file error")
		return 0
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return 0
	}

	return fi.ModTime().Unix()
}

//ReadLines 读取文本所有行
//@ctx 上下文
//@filePath 文本文件路径
//@return 行文本信道
func ReadLines(ctx context.Context, filePath string) <-chan string {
	lineChan := make(chan string, 1)
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	go func() {
		defer close(lineChan)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			line = strings.TrimSpace(line)
			line = strings.Trim(line, "\r")
			select {
			case <-ctx.Done():
				return
			case lineChan <- line:
			}
		}
	}()
	return lineChan
}
