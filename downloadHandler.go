package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func getContentType(fileName string) string {
	ext := filepath.Ext(fileName)
	switch ext {
	case ".txt":
		return "text/plain"
	case ".html", ".htm":
		return "text/html"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".json":
		return "application/json"
	default:
		return "application/octet-stream"
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	//数据库临时测试
	collection := clientInstance.Database("test").Collection("trainer")
	user := User{
		Username: "john_do",
		Email:    "john@example.com",
	}
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		panic(err.Error())
	}
	// end

	// 设置CORS头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 处理预检请求
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

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

	// 打开文件并将其内容写入响应
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		log.Printf("Failed to open file, clientIP: %s, FileName: %s", clientIP, fileName)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("关闭文件失败: %v", err)
		}
	}()

	// 设置响应头，根据文件扩展名设置适当的 Content-Type
	contentType := getContentType(fileName)
	w.Header().Set("Content-Type", contentType)

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Failed to write file to response", http.StatusInternalServerError)
		log.Printf("Failed to write file to response, clientIP: %s, FileName: %s", clientIP, fileName)
		return
	}

	//// 设置响应头，告诉浏览器这是一个文件下载
	//w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	//w.Header().Set("Content-Type", "application/octet-stream")
	//
	//// 打开文件并将其内容写入响应
	//file, err := os.Open(filePath)
	//if err != nil {
	//	http.Error(w, "Failed to open file", http.StatusInternalServerError)
	//	log.Printf("Failed to open file, clientIP: %s, FileName: %s", clientIP, fileName)
	//	return
	//}
	////defer file.Close()
	//defer func() {
	//	if err := file.Close(); err != nil {
	//		log.Printf("关闭文件失败: %v", err)
	//	}
	//}()
	//
	//_, err = file.Seek(0, 0)
	//if err != nil {
	//	http.Error(w, "Failed to seek file", http.StatusInternalServerError)
	//	log.Printf("Failed to seek file, clientIP: %s, FileName: %s", clientIP, fileName)
	//	return
	//}
	//
	//_, err = io.Copy(w, file)
	//if err != nil {
	//	http.Error(w, "Failed to write file to response", http.StatusInternalServerError)
	//	log.Printf("Failed to write file to response, clientIP: %s, FileName: %s", clientIP, fileName)
	//	return
	//}

	// 记录下载日志
	log.Printf("客户端IP: %s, 下载文件: %s", clientIP, fileName)
}
