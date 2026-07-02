package mysqlx

import (
	"errors"
	"testing"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

func TestNewRejectsEmptyDSN(t *testing.T) {
	_, err := New(config.MySQLConfig{})
	if !errors.Is(err, ErrMissingDSN) {
		t.Fatalf("空数据库连接串应返回 ErrMissingDSN，实际: %v", err)
	}
}
