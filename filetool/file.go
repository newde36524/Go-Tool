package filetool

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ReadPagingFile 分段获取指定文件的数据包
// @pagIndex 表示数据包的下标
// @pagSize 表示一次获取的数据包大小
// @filePath 表示文件路径
// @rerurn 读取到的文件数据包 @bs 读取到的数据 @n 有效数据长度  @e 表示读取异常
func ReadPagingFile(pagIndex, pagSize, off int, filePath string) (bs []byte, n int, e error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, 0, err
	}
	bs, n, e = ReadPagingBuffer(pagIndex, pagSize, off, file)
	return
}

//ReadPagingBuffer 分段读取数据
// @pagIndex 表示数据包的下标
// @pagSize 表示一次获取的数据包大小
// @filePath 表示文件路径
// @rerurn 读取到的文件数据包 @bs 读取到的数据 @n 有效数据长度  @e 表示读取异常
func ReadPagingBuffer(pagIndex, pagSize, off int, buffer io.Reader) (bs []byte, n int, e error) {
	bs = make([]byte, pagSize)
	bufReader := bufio.NewReader(buffer)
	bufReader.Discard(off)
	_, e = bufReader.Discard(pagIndex * pagSize) //跳过指定字节数
	if e != nil {
		return
	}
	for i := 0; i < pagSize; i++ {
		b, e := bufReader.ReadByte()
		if e != nil {
			break
		}
		n = i + 1
		// bs = append(bs, b) //读取超出接线默认用0补足
		bs[i] = b
	}
	return
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
