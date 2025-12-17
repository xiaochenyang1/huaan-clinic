package handler

import (
	"github.com/gin-gonic/gin"

	"huaan-medical/pkg/errorcode"
	"huaan-medical/pkg/response"
	"huaan-medical/pkg/utils"
)

// UploadHandler 文件上传处理器
type UploadHandler struct{}

// NewUploadHandler 创建文件上传处理器实例
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// UploadAvatarResponse 上传头像响应
type UploadAvatarResponse struct {
	URL      string `json:"url"`       // 访问URL
	FileName string `json:"file_name"` // 文件名
	FileSize int64  `json:"file_size"` // 文件大小（字节）
}

// UploadAvatar 上传头像
// @Summary 上传头像
// @Description 上传医生头像图片
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param file formData file true "头像文件"
// @Success 200 {object} response.Response{data=UploadAvatarResponse}
// @Router /api/admin/upload/avatar [post]
func (h *UploadHandler) UploadAvatar(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage(c, errorcode.ErrInvalidParams, "请选择要上传的文件")
		return
	}

	// 保存文件
	result, err := utils.SaveAvatar(file)
	if err != nil {
		response.FailWithMessage(c, errorcode.ErrInternalServer, err.Error())
		return
	}

	// 返回结果
	resp := &UploadAvatarResponse{
		URL:      result.FileURL,
		FileName: result.FileName,
		FileSize: result.FileSize,
	}

	response.Success(c, resp)
}

// UploadImage 上传图片
// @Summary 上传图片
// @Description 上传普通图片
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param file formData file true "图片文件"
// @Success 200 {object} response.Response{data=UploadAvatarResponse}
// @Router /api/admin/upload/image [post]
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage(c, errorcode.ErrInvalidParams, "请选择要上传的文件")
		return
	}

	// 保存文件
	result, err := utils.SaveImage(file)
	if err != nil {
		response.FailWithMessage(c, errorcode.ErrInternalServer, err.Error())
		return
	}

	// 返回结果
	resp := &UploadAvatarResponse{
		URL:      result.FileURL,
		FileName: result.FileName,
		FileSize: result.FileSize,
	}

	response.Success(c, resp)
}
