package middleware

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type LoggerBuilder struct {
	allowReqBody  bool
	allowRespBody bool
	maxLogUrlLen  int
	maxLogBodyLen int
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

// NewLoggerBuilder 创建一个 LoggerBuilder
// 默认不允许记录请求和响应的 Body
// 默认最大记录 URL 长度为 1024
// 默认最大记录 Body 长度为 8192
func NewLoggerBuilder(lf func(ctx context.Context, al *AccessLog)) *LoggerBuilder {
	return &LoggerBuilder{
		maxLogUrlLen:  1024,
		maxLogBodyLen: 8192,
		loggerFunc:    lf,
	}
}

func (b *LoggerBuilder) AllowReqBody() *LoggerBuilder {
	b.allowReqBody = true
	return b
}

func (b *LoggerBuilder) AllowRespBody() *LoggerBuilder {
	b.allowRespBody = true
	return b
}

func (b *LoggerBuilder) WithMaxLogBodyLen(maxLogBodyLen int) *LoggerBuilder {
	b.maxLogBodyLen = maxLogBodyLen
	return b
}

func (b *LoggerBuilder) WithMaxLogUrlLen(maxLogUrlLen int) *LoggerBuilder {
	b.maxLogUrlLen = maxLogUrlLen
	return b

}

func (b *LoggerBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > b.maxLogUrlLen {
			// ? 防范可能的故意长 URL 攻击
			url = url[:b.maxLogUrlLen]
		}
		al := &AccessLog{
			Method: ctx.Request.Method,
			Url:    url,
		}
		if b.allowReqBody && ctx.Request.Body != nil {
			body, _ := ctx.GetRawData()
			// Body 是一个 io.ReadCloser，读取后需要重新设置回去
			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))

			if len(body) > b.maxLogBodyLen {
				body = body[:b.maxLogBodyLen]
			}
			al.ReqBody = string(body)
		}

		if b.allowRespBody {
			w := &responseWriter{
				al:             al,
				ResponseWriter: ctx.Writer,
				maxLogBodyLen:  b.maxLogBodyLen,
			}
			ctx.Writer = w
		}
		defer func() {
			al.Duration = time.Since(start).String()
			al.Status = ctx.Writer.Status()
			b.loggerFunc(ctx.Request.Context(), al)
		}()
		ctx.Next()

	}
}

// responseWriter 重写了 Write 方法，用于记录响应 Body
type responseWriter struct {
	al            *AccessLog
	maxLogBodyLen int
	gin.ResponseWriter
}

func (w *responseWriter) Write(data []byte) (int, error) {
	if len(data) > w.maxLogBodyLen {
		w.al.RespBody = string(data[:w.maxLogBodyLen])
	} else {
		w.al.RespBody = string(data)
	}
	return w.ResponseWriter.Write(data)
}

func (w *responseWriter) WriteString(data string) (int, error) {
	if len(data) > w.maxLogBodyLen {
		w.al.RespBody = data[:w.maxLogBodyLen]
	} else {
		w.al.RespBody = data
	}
	return w.ResponseWriter.WriteString(data)
}

type AccessLog struct {
	Method   string
	Status   int
	Url      string
	ReqBody  string
	RespBody string
	Duration string
}
