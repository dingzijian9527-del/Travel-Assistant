package uploadx

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

func UploadAvatarToQiniu(ctx context.Context, cfg config.UploadQiniuConfig, filename string, file multipart.File) (string, string, error) {
	if err := validateQiniuConfig(cfg); err != nil {
		return "", "", err
	}

	objectKey := buildObjectKey(filename)
	objectName := objectKey
	cred := credentials.NewCredentials(strings.TrimSpace(cfg.AccessKey), strings.TrimSpace(cfg.SecretKey))
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: cred,
		},
	})

	if _, err := file.Seek(0, 0); err != nil {
		return "", "", fmt.Errorf("重置上传文件读取位置失败: %w", err)
	}

	if err := uploadManager.UploadReader(ctx, file, &uploader.ObjectOptions{
		BucketName: strings.TrimSpace(cfg.Bucket),
		ObjectName: &objectName,
		FileName:   filepath.Base(filename),
	}, nil); err != nil {
		return "", "", fmt.Errorf("七牛云上传失败: %w", err)
	}

	return buildPublicURL(cfg.URL, objectKey), objectKey, nil
}

func validateQiniuConfig(cfg config.UploadQiniuConfig) error {
	if strings.TrimSpace(cfg.AccessKey) == "" {
		return fmt.Errorf("七牛云 AccessKey 未配置")
	}
	if strings.TrimSpace(cfg.SecretKey) == "" {
		return fmt.Errorf("七牛云 SecretKey 未配置")
	}
	if strings.TrimSpace(cfg.Bucket) == "" {
		return fmt.Errorf("七牛云存储空间未配置")
	}
	if strings.TrimSpace(cfg.URL) == "" {
		return fmt.Errorf("七牛云访问域名未配置")
	}
	return nil
}

func buildObjectKey(filename string) string {
	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(filename)))
	if ext == "" {
		ext = ".png"
	}
	return "avatars/" + uuid.New().String() + ext
}

func buildPublicURL(baseURL string, objectKey string) string {
	return strings.TrimRight(strings.TrimSpace(baseURL), "/") + "/" + strings.TrimLeft(strings.TrimSpace(objectKey), "/")
}
