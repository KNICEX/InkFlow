package ioc

import (
	"github.com/KNICEX/InkFlow/pkg/gorsex"
)

func InitGorseCli() *gorsex.Client {
	return gorsex.NewClient("http://localhost:8088", "")
}
