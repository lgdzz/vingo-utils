package vingo

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 创建目录
func Mkdir(dirPath string) string {
	// 将路径中的反斜杠替换为正斜杠，以支持 Windows 目录
	dirPath = filepath.ToSlash(dirPath)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// 目录不存在，创建目录
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			panic(fmt.Sprintf("创建目录失败：%v", err.Error()))
		} else {
			return dirPath
		}
	} else if err != nil {
		// 其他错误
		panic(fmt.Sprintf("判断目录是否存在时发生错误：%v", err.Error()))
	} else {
		// 目录存在
		return dirPath
	}
}

// 保存文件
// 将字节数据保存为文件
// dirPath 保存文件所在目录
// fileName 保存文件的名称
// data 保存文件的内容
func SaveFile(dirPath string, fileName string, data []byte) {
	targetFile := filepath.Join(dirPath, fileName)
	if err := os.WriteFile(targetFile, data, 0644); err != nil {
		panic(fmt.Sprintf("保存文件失败：%v", err.Error()))
	}
}

// 判断文件是否存在
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// 判断目录是否有读写权限
func HasDirReadWritePermission(dirPath string) bool {
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		LogError(fmt.Sprintf("目录 %v 不存在或无法访问: %v", dirPath, err.Error()))
		return false
	}

	mode := fileInfo.Mode()
	if !mode.IsDir() {
		LogError(fmt.Sprintf("%v 不是一个目录", dirPath))
		return false
	}

	perm := mode.Perm()
	if perm&(1<<(uint(7))) == 0 || perm&(1<<(uint(6))) == 0 {
		LogError(fmt.Sprintf("目录 %v 没有读和写权限", dirPath))
		return false
	}

	return true
}

func FileUpload(path string, request *http.Request) *FileInfo {
	var (
		requestFile multipart.File
		header      *multipart.FileHeader
		err         error
	)
	requestFile, header, err = request.FormFile("file")
	if err != nil {
		panic(err.Error())
	}
	defer requestFile.Close()

	// 获取文件大小
	fileSize := header.Size

	// 获取文件名称、类型、后缀
	fileName := header.Filename
	fileType := header.Header.Get("Content-Type")
	fileSuffix := filepath.Ext(fileName)

	// 获取当前日期，用于存储文件
	dateString := time.Now().Format(DateFormat)

	// 指定存储目录，如果不存在则创建
	dirPath := Mkdir(filepath.Join(path, dateString))

	// 创建文件
	filePath := filepath.Join(dirPath, fmt.Sprintf("%v%v", GetUUID(), fileSuffix))
	newFile, err := os.Create(filePath)
	if err != nil {
		panic(err.Error())
	}
	defer newFile.Close()

	// 将文件内容拷贝到新文件中
	if _, err = io.Copy(newFile, requestFile); err != nil {
		panic(err.Error())
	}

	// 返回结果
	return &FileInfo{
		Name:      fileName,
		Mimetype:  fileType,
		Extension: fileSuffix,
		Size:      fileSize,
		Realpath:  strings.Replace(filePath, "\\", "/", -1),
	}
}

// 复制文件
// src 文件位置
// dstDir 要复制到的位置
func FileCopy(src, dstDir string) string {
	// Open the source file for reading.
	srcFile, err := os.Open(src)
	if err != nil {
		panic(err.Error())
	}
	defer srcFile.Close()

	// 获取当前日期，用于存储文件
	dateString := time.Now().Format(DateFormat)

	// 指定存储目录，如果不存在则创建
	dstDir = filepath.Join(dstDir, dateString)
	if _, err = os.Stat(dstDir); os.IsNotExist(err) {
		if err = os.MkdirAll(dstDir, 0755); err != nil {
			panic(err.Error())
		}
	}

	// Generate a unique file name with the same extension as the source file.
	srcExt := filepath.Ext(src)
	dstFileName := uuid.New().String() + srcExt

	// Create the destination file with the generated file name.
	dstFile, err := os.Create(filepath.Join(dstDir, dstFileName))
	if err != nil {
		panic(err.Error())
	}
	defer dstFile.Close()

	// Copy the contents of the source file to the destination file.
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		panic(err.Error())
	}

	// Return the path to the copied file.
	return filepath.Join(dstDir, dstFileName)
}

func FileDelete(path string, showErr bool) {
	if FileExists(path) {
		// 文件存在，删除文件
		if err := os.Remove(path); err != nil {
			if showErr {
				panic(fmt.Sprintf("删除文件失败：%v", err.Error()))
			}
		}
	}
}
