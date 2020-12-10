package filetool

import (
	"archive/tar"
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/issue9/logs"
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

//Find 深度遍历文件
//@baseDir 目标文件或文件夹所在目录
//@fileInfo 目标文件或文件夹
//@callback 遍历时处理 @param1 文件所在绝对路径 @param2 文件
func Find(baseDir string, fileInfo os.FileInfo, callback func(baseDir string, file *os.File)) {
	var find func(baseDir string, fileInfo os.FileInfo)
	find = func(baseDir string, fileInfo os.FileInfo) {
		fullName := filepath.Join(baseDir, fileInfo.Name())
		file, err := os.Open(fullName)
		if err != nil {
			panic(err)
		}
		if fileInfo.IsDir() {
			dirItems, err := file.Readdir(-1)
			if err != nil {
				panic(err)
			}
			for _, dirItem := range dirItems {
				if dirItem.IsDir() {
					find(fullName, dirItem)
				} else {
					file, err := os.Open(filepath.Join(fullName, dirItem.Name()))
					if err != nil {
						panic(err)
					}
					callback(fullName, file)
				}
			}
		} else {
			callback("", file)
		}
	}
	find(baseDir, fileInfo)
}

//Compress 打包文件到目标目录
//@baseDir 目标文件所在目录
//@fileInfo 目标文件
//@destDir 压缩到文件夹
//@destFileName 压缩文件名
func Compress(baseDir string, fileInfo os.FileInfo, destDir, destFileName string) (err error) {
	if destDir == "" {
		panic("目标目录不允许为空")
	}
	if destFileName == "" {
		panic("目标文件名不允许为空")
	}
	if _, err := os.Stat(baseDir); err != nil { //创建目标目录,err 不为nil时表示文件不存在
		logs.Infof("创建文件夹:%s", baseDir)
		if err := os.MkdirAll(baseDir, os.ModeDir); err != nil {
			return err
		}
	}
	if _, err := os.Stat(destDir); err != nil { //创建目标目录,err 不为nil时表示文件不存在
		logs.Infof("创建文件夹:%s", destDir)
		if err := os.MkdirAll(destDir, os.ModeDir); err != nil {
			return err
		}
	}
	var writer *tar.Writer
	var targetFile *os.File
	defer func() {
		if writer != nil {
			writer.Close()
		}
		if targetFile != nil {
			targetFile.Close()
		}
	}()
	once := sync.Once{}
	Find(baseDir, fileInfo, func(baseDir string, file *os.File) {
		once.Do(func() {
			targetFile, err = os.Create(filepath.Join(destDir, destFileName))
			if err != nil {
				panic("压缩文件创建失败")
			}
			writer = tar.NewWriter(targetFile)
		})

		fInfo, err := file.Stat()
		if err != nil {
			logs.Error(err)
			return
		}
		header, err := tar.FileInfoHeader(fInfo, "")
		if err != nil {
			logs.Error(err)
			return
		}
		header.Name = strings.TrimPrefix(file.Name(), baseDir)
		if err := writer.WriteHeader(header); err != nil {
			logs.Error(err)
			return
		}
		_, err = io.Copy(writer, file)
		if err != nil {
			logs.Error(err)
			return
		}
		file.Close()
		if err := os.Remove(file.Name()); err != nil {
			logs.Error(err)
		}
	})
	return nil
}

//DeCompress 解压文件到目标目录
//@baseDir 压缩文件所在目录
//@src 压缩文件
//@dest 目标目录
func DeCompress(baseDir string, src os.FileInfo, dest string) error {
	if baseDir == "" {
		panic("读取目录不允许为空")
	}
	if dest == "" {
		panic("目标目录不允许为空")
	}
	if _, err := os.Stat(baseDir); err != nil {
		if err := os.MkdirAll(baseDir, os.ModeDir); err != nil {
			return err
		}
	}
	if _, err := os.Stat(dest); err != nil {
		if err := os.MkdirAll(dest, os.ModeDir); err != nil {
			return err
		}
	}
	Find(baseDir, src, func(baseDir string, file *os.File) {
		reader := tar.NewReader(file)
		for {
			header, err := reader.Next()
			if err != nil {
				break
			}
			fullName := filepath.Join(dest, header.Name)
			if _, err := os.Stat(filepath.Dir(fullName)); err != nil {
				if err := os.MkdirAll(filepath.Dir(fullName), os.ModeDir); err != nil {
					panic(err)
				}
			}
			targetFile, err := os.Create(fullName)
			if err != nil {
				panic(err)
			}
			if _, err := io.Copy(targetFile, reader); err != nil {
				targetFile.Close()
				panic(err)
			}
			targetFile.Close()
		}
		file.Close()
		if err := os.Remove(file.Name()); err != nil {
			panic(err)
		}
	})
	return nil
}
