package web

import (
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	auth  middleware.Authentication
	cloud *cloudinary.Cloudinary
	l     logx.Logger
}

func NewFileHandler(cloud *cloudinary.Cloudinary, auth middleware.Authentication, l logx.Logger) *FileHandler {
	return &FileHandler{
		auth:  auth,
		cloud: cloud,
		l:     l,
	}
}

func (handler *FileHandler) RegisterRoutes(server *gin.RouterGroup) {
	fileGroup := server.Group("/file", handler.auth.CheckLogin(), handler.monitor())
	fileGroup.POST("/avatar", ginx.Wrap(handler.l, handler.UploadAvatar))
	fileGroup.POST("/cover", ginx.Wrap(handler.l, handler.UploadCover))
	fileGroup.POST("/image", ginx.Wrap(handler.l, handler.UploadImage))
}

func (handler *FileHandler) monitor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 监控用户异常上传行为
	}
}

func (handler *FileHandler) UploadAvatar(ctx *gin.Context) (ginx.Result, error) {
	// 获取图片流
	multi, err := ctx.FormFile("avatar")
	if err != nil {
		return ginx.InvalidParam(), err
	}

	// TODO 压缩图片

	reader, err := multi.Open()
	if err != nil {
		return ginx.InternalError(), err
	}

	resp, err := handler.cloud.Upload.Upload(ctx, reader, uploader.UploadParams{})
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(resp.URL), nil
}

func (handler *FileHandler) UploadCover(ctx *gin.Context) (ginx.Result, error) {
	// 获取图片流
	multi, err := ctx.FormFile("cover")
	if err != nil {
		return ginx.InvalidParam(), err
	}

	reader, err := multi.Open()
	if err != nil {
		return ginx.InternalError(), err
	}

	resp, err := handler.cloud.Upload.Upload(ctx, reader, uploader.UploadParams{})
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(resp.URL), nil
}

func (handler *FileHandler) UploadImage(ctx *gin.Context) (ginx.Result, error) {
	multi, err := ctx.FormFile("image")
	if err != nil {
		return ginx.InvalidParam(), err
	}

	reader, err := multi.Open()
	if err != nil {
		return ginx.InternalError(), err
	}

	resp, err := handler.cloud.Upload.Upload(ctx, reader, uploader.UploadParams{})
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(resp.URL), nil
}
