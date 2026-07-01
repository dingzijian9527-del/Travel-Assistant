package rpcuser

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
	_ "github.com/go-sql-driver/mysql"
)

type mysqlUserRepo struct {
	db *sql.DB
}

func newMySQLUserRepo(cfg config.MySQLConfig) (*mysqlUserRepo, error) {
	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, err
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
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	repo := &mysqlUserRepo{db: db}
	if err := repo.ensureSchema(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return repo, nil
}

func (r *mysqlUserRepo) ensureSchema(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(64) NOT NULL COMMENT '用户主键，字符串类型',
    phone VARCHAR(20) NOT NULL COMMENT '手机号，前端登录注册主账号',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    nickname VARCHAR(64) NOT NULL DEFAULT '' COMMENT '昵称',
    avatar_url VARCHAR(512) DEFAULT NULL COMMENT '头像地址',
    home_city VARCHAR(64) DEFAULT NULL COMMENT '常住城市',
    current_city VARCHAR(64) DEFAULT NULL COMMENT '当前位置城市',
    member_level VARCHAR(32) NOT NULL DEFAULT 'normal' COMMENT '会员等级',
    account_status TINYINT NOT NULL DEFAULT 1 COMMENT '账号状态：1正常、2禁用、3注销',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME DEFAULT NULL COMMENT '软删除时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_users_phone (phone),
    KEY idx_users_account_status (account_status),
    KEY idx_users_created_at (created_at),
    KEY idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表'`)
	return err
}

func (r *mysqlUserRepo) Create(ctx context.Context, input *userModel) (*userModel, error) {
	now := time.Now()
	storedUser := cloneUser(input)
	storedUser.ID = uuid.New().String()
	storedUser.CreatedAt = now
	storedUser.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, `
INSERT INTO users(id, phone, password_hash, nickname, avatar_url, home_city, current_city, member_level, account_status, created_at, updated_at)
VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		storedUser.ID,
		storedUser.Phone,
		storedUser.PasswordHash,
		storedUser.Nickname,
		nullString(storedUser.AvatarURL),
		nullString(storedUser.HomeCity),
		nullString(storedUser.CurrentCity),
		storedUser.MemberLevel,
		storedUser.AccountStatus,
		storedUser.CreatedAt,
		storedUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return storedUser, nil
}

func (r *mysqlUserRepo) FindByPhone(ctx context.Context, phone string) (*userModel, bool) {
	return r.findOne(ctx, "phone = ?", phone)
}

func (r *mysqlUserRepo) FindByID(ctx context.Context, userID string) (*userModel, bool) {
	return r.findOne(ctx, "id = ?", userID)
}

func (r *mysqlUserRepo) Update(ctx context.Context, userID string, apply func(*userModel)) (*userModel, bool) {
	storedUser, ok := r.FindByID(ctx, userID)
	if !ok {
		return nil, false
	}
	apply(storedUser)
	storedUser.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `
UPDATE users
SET nickname = ?, avatar_url = ?, home_city = ?, current_city = ?, updated_at = ?
WHERE id = ? AND deleted_at IS NULL`,
		storedUser.Nickname,
		nullString(storedUser.AvatarURL),
		nullString(storedUser.HomeCity),
		nullString(storedUser.CurrentCity),
		storedUser.UpdatedAt,
		storedUser.ID,
	)
	if err != nil {
		return nil, false
	}
	return storedUser, true
}

func (r *mysqlUserRepo) findOne(ctx context.Context, where string, arg any) (*userModel, bool) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, phone, password_hash, nickname, avatar_url, home_city, current_city, member_level, account_status, created_at, updated_at, deleted_at
FROM users
WHERE `+where+` AND deleted_at IS NULL
LIMIT 1`, arg)
	var item userModel
	var avatarURL, homeCity, currentCity sql.NullString
	var deletedAt sql.NullTime
	err := row.Scan(
		&item.ID,
		&item.Phone,
		&item.PasswordHash,
		&item.Nickname,
		&avatarURL,
		&homeCity,
		&currentCity,
		&item.MemberLevel,
		&item.AccountStatus,
		&item.CreatedAt,
		&item.UpdatedAt,
		&deletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false
	}
	if err != nil {
		return nil, false
	}
	item.AvatarURL = avatarURL.String
	item.HomeCity = homeCity.String
	item.CurrentCity = currentCity.String
	if deletedAt.Valid {
		item.DeletedAt = &deletedAt.Time
	}
	return &item, true
}

func nullString(value string) sql.NullString {
	return sql.NullString{String: value, Valid: value != ""}
}
