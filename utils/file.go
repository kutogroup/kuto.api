package utils

import (
	"bytes"
	"io"
	"os"
	"time"
)

//FileCopy 复制文件
func FileCopy(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

//FileMove 移动文件
func FileMove(dstName, srcName string) (written int64, err error) {
	defer FileDelete(srcName)
	return FileCopy(dstName, srcName)
}

//FileExists 判断文件是否存在
func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

//FileDelete 删除文件
func FileDelete(file string) bool {
	return os.Remove(file) == nil
}

//FileGetSize 获取文件大小
func FileGetSize(file string) int64 {
	info, err := os.Stat(file)
	if err != nil {
		return -1
	}

	return info.Size()
}

//FileGetModTime 获取文件修改时间
func FileGetModTime(file string) *time.Time {
	info, err := os.Stat(file)
	if err != nil {
		return nil
	}

	t := info.ModTime()
	return &t
}

//FileWriteBytes 文件写入字节
func FileWriteBytes(file string, bytes []byte) (int, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(bytes)
}

//FileWriteString 文件写入字符串
func FileWriteString(file string, content string) (int, error) {
	return FileWriteBytes(file, []byte(content))
}

//FileReadBytes 读取文件字节
func FileReadBytes(file string) []byte {
	f, err := os.Open(file)
	if err != nil {
		return nil
	}
	defer f.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(f)
	return buf.Bytes()
}

//FileReadString 读取文件字符串
func FileReadString(file string) string {
	b := FileReadBytes(file)

	if b == nil || len(b) == 0 {
		return ""
	}

	return string(b)
}
