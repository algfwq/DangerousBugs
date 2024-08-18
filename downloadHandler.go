package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// 获取客户端IP地址
	clientIP := getClientIP(r)

	// 从URL路径中获取文件名
	fileName := strings.TrimPrefix(r.URL.Path, "/download/")
	if fileName == "" {
		http.Error(w, "Missing 'file' parameter", http.StatusBadRequest)
		log.Printf("Missing 'file' parameter, clientIP: %s, FileName: %s", clientIP, fileName)
		return
	}

	// 构建文件路径
	filePath := filepath.Join("files", fileName)

	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		log.Printf("File not found, clientIP: %s, FileName: %s", clientIP, fileName)
		return
	}

	// 设置响应头，告诉浏览器这是一个文件下载
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")

	// 打开文件并将其内容写入响应
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		log.Printf("Failed to open file, clientIP: %s, FileName: %s", clientIP, fileName)
		return
	}
	//defer file.Close()
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("关闭文件失败: %v", err)
		}
	}()

	_, err = file.Seek(0, 0)
	if err != nil {
		http.Error(w, "Failed to seek file", http.StatusInternalServerError)
		log.Printf("Failed to seek file, clientIP: %s, FileName: %s", clientIP, fileName)
		return
	}

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Failed to write file to response", http.StatusInternalServerError)
		log.Printf("Failed to write file to response, clientIP: %s, FileName: %s", clientIP, fileName)
		return
	}

	// 记录下载日志
	log.Printf("客户端IP: %s, 下载文件: %s", clientIP, fileName)
}
