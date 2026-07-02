package mysqlx

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
)

// ErrMissingDSN 表示缺少 MySQL 连接串。
var ErrMissingDSN = errors.New("mysql dsn 未配置")

// New 创建并探活 MySQL 连接池。
func New(cfg config.MySQLConfig) (*sql.DB, error) {
	if strings.TrimSpace(cfg.DSN) == "" {
		return nil, ErrMissingDSN
	}
	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, err
	}
	ApplyPoolConfig(db, cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

// ApplyPoolConfig 统一设置 MySQL 连接池参数。
func ApplyPoolConfig(db *sql.DB, cfg config.MySQLConfig) {
	if db == nil {
		return
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.ConnMaxLifetimeSeconds > 0 {
		db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeSeconds) * time.Second)
	}
}
