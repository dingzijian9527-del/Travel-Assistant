package upload

import (
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/uploadx"
)

// Check 校验网关上传请求中的文件元信息。
func Check(filename string, size int64, cfg config.UploadConfig) error {
	return uploadx.ValidateFile(filename, size, cfg)
}
