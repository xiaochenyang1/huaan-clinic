package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 文件上传相关常量
const (
	MaxAvatarSize  = 2 * 1024 * 1024  // 2MB
	MaxImageSize   = 5 * 1024 * 1024  // 5MB
	UploadBasePath = "uploads"        // 上传文件基础路径
	AvatarPath     = "uploads/avatar" // 头像上传路径
	ImagePath      = "uploads/images" // 图片上传路径
)

// AllowedImageTypes 允许的图片类型
var AllowedImageTypes = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

// FileUploadResult 文件上传结果
type FileUploadResult struct {
	FileName string `json:"file_name"` // 文件名
	FilePath string `json:"file_path"` // 文件相对路径
	FileURL  string `json:"file_url"`  // 访问URL
	FileSize int64  `json:"file_size"` // 文件大小（字节）
}

// SaveUploadedFile 保存上传的文件
// fileHeader: 上传的文件
// savePath: 保存路径（相对路径，如 "uploads/avatar"）
// maxSize: 最大文件大小（字节）
// allowedExts: 允许的文件扩展名（如 []string{".jpg", ".png"}）
func SaveUploadedFile(fileHeader *multipart.FileHeader, savePath string, maxSize int64, allowedExts []string) (*FileUploadResult, error) {
	// 检查文件大小
	if fileHeader.Size > maxSize {
		return nil, fmt.Errorf("文件大小不能超过 %dMB", maxSize/1024/1024)
	}

	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		return nil, fmt.Errorf("文件必须有扩展名")
	}

	// 检查文件类型
	if !InArray(ext, allowedExts) {
		return nil, fmt.Errorf("不支持的文件类型: %s，仅支持: %s", ext, strings.Join(allowedExts, ", "))
	}

	// 确保目录存在
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}

	// 生成唯一文件名：日期 + UUID + 扩展名
	dateStr := time.Now().Format("20060102")
	fileName := fmt.Sprintf("%s_%s%s", dateStr, GenerateShortUUID(), ext)
	filePath := filepath.Join(savePath, fileName)

	// 打开上传的文件
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("打开上传文件失败: %v", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 生成访问URL（相对路径，前端需要拼接域名）
	fileURL := "/" + strings.ReplaceAll(filePath, "\\", "/")

	return &FileUploadResult{
		FileName: fileName,
		FilePath: filePath,
		FileURL:  fileURL,
		FileSize: fileHeader.Size,
	}, nil
}

// SaveAvatar 保存头像文件
func SaveAvatar(fileHeader *multipart.FileHeader) (*FileUploadResult, error) {
	return SaveUploadedFile(fileHeader, AvatarPath, MaxAvatarSize, AllowedImageTypes)
}

// SaveImage 保存图片文件
func SaveImage(fileHeader *multipart.FileHeader) (*FileUploadResult, error) {
	return SaveUploadedFile(fileHeader, ImagePath, MaxImageSize, AllowedImageTypes)
}

// DeleteFile 删除文件
func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	// 如果是URL格式，转换为文件路径
	if strings.HasPrefix(filePath, "/") {
		filePath = strings.TrimPrefix(filePath, "/")
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // 文件不存在，不报错
	}

	return os.Remove(filePath)
}

// GetFileSize 获取文件大小（字节）
func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// IsImageFile 检查是否为图片文件
func IsImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return InArray(ext, AllowedImageTypes)
}
