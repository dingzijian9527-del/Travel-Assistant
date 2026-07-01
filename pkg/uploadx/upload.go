package uploadx

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

// ValidateFile 校验上传文件大小和扩展名。
func ValidateFile(filename string, size int64, cfg config.UploadConfig) error {
	if cfg.MaxSizeMB > 0 && size > cfg.MaxSizeMB*1024*1024 {
		return fmt.Errorf("文件大小不能超过 %dMB", cfg.MaxSizeMB)
	}
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowed := range cfg.AllowedExtensions {
		if ext == strings.ToLower(allowed) {
			return nil
		}
	}
	return fmt.Errorf("不支持的文件类型: %s", ext)
}
