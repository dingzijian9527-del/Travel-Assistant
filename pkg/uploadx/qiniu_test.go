package uploadx

import (
	"strings"
	"testing"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

func TestBuildPublicURL(t *testing.T) {
	url := buildPublicURL("http://t7aqo4inq.hn-bkt.clouddn.com/", "avatars/test.png")
	if url != "http://t7aqo4inq.hn-bkt.clouddn.com/avatars/test.png" {
		t.Fatalf("unexpected public url: %s", url)
	}
}

func TestBuildObjectKeyKeepsImageExtension(t *testing.T) {
	key := buildObjectKey("头像.JPG")
	if !strings.HasPrefix(key, "avatars/") {
		t.Fatalf("unexpected key prefix: %s", key)
	}
	if !strings.HasSuffix(key, ".jpg") {
		t.Fatalf("unexpected key suffix: %s", key)
	}
}

func TestValidateQiniuConfig(t *testing.T) {
	cfg := config.UploadQiniuConfig{
		AccessKey: "test-ak",
		SecretKey: "test-sk",
		Bucket:    "dingzijian",
		URL:       "http://example.com/",
	}
	err := validateQiniuConfig(cfg)
	if err != nil {
		t.Fatalf("validateQiniuConfig returned error: %v", err)
	}
}
