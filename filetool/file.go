package filetool

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"os/exec"
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
// @rerurn 读取到的文件数据包 @bs 读取到的数据 @n 有效数据长度  @err 表示读取异常
func ReadPagingBuffer(pagIndex, pagSize, off int, buffer io.Reader) (bs []byte, n int, err error) {
	bufReader := bufio.NewReader(buffer)
	_, err = bufReader.Discard(pagIndex*pagSize + off) //跳过指定字节数
	if err != nil {
		return
	}
	bs = make([]byte, pagSize)
	n, err = bufReader.Read(bs)
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
	fn := func(line string) {
		line = strings.TrimSpace(line)
		line = strings.Trim(line, "\r")
		select {
		case <-ctx.Done():
			return
		case lineChan <- line:
		}
	}
	go func() {
		defer close(lineChan)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fn(line)
				return
			}
			fn(line)
		}
	}()
	return lineChan
}

//GetDirFullPath 获取当前程序所在目录
func GetDirFullPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	return ret, nil
}

//GetDirFullPath2 获取当前程序所在目录
func GetDirFullPath2() (string, error) {
	return filepath.Abs("./")
}
