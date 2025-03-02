package uuidx

import (
	"encoding/base64"
	"github.com/google/uuid"
)

func NewShort() string {
	u, _ := uuid.New().MarshalBinary()
	return base64.StdEncoding.EncodeToString(u)[:8]
}

// NewShortN 生成指定长度的短uuid,  n不应该超过32
func NewShortN(n int) string {
	u, _ := uuid.New().MarshalBinary()
	return base64.StdEncoding.EncodeToString(u)[:n]
}
