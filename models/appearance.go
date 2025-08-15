package models

import (
	"database/sql"
	"time"
)

// Appearance 外觀模型，對應你的 PostgreSQL 表結構
type Appearance struct {
	ID             string         `db:"id" json:"id"`
	LogoURL        sql.NullString `db:"logo_url" json:"logo_url"`
	PrimaryColor   sql.NullString `db:"primary_color" json:"primary_color"`
	SecondaryColor sql.NullString `db:"secondary_color" json:"secondary_color"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at" json:"updated_at"`
	LogoKey        sql.NullString `db:"logo_key" json:"logo_key"`
	LogoSize       sql.NullInt64  `db:"logo_size" json:"logo_size"`
	LogoType       sql.NullString `db:"logo_type" json:"logo_type"`
}

// AppearanceRepository 外觀數據庫操作接口
type AppearanceRepository interface {
	GetByID(id string) (*Appearance, error)
	Update(appearance *Appearance) error
}

// AppearanceRepositoryImpl 外觀數據庫操作實現
type AppearanceRepositoryImpl struct {
	db *sql.DB
}

// NewAppearanceRepository 創建新的外觀倉庫實例
func NewAppearanceRepository(db *sql.DB) AppearanceRepository {
	return &AppearanceRepositoryImpl{db: db}
}

// GetByID 根據 ID 獲取外觀
func (r *AppearanceRepositoryImpl) GetByID(id string) (*Appearance, error) {
	query := `
		SELECT id, logo_url, primary_color, secondary_color, created_at, updated_at, logo_key, logo_size, logo_type
		FROM appearance WHERE id = $1
	`

	appearance := &Appearance{}
	err := r.db.QueryRow(query, id).Scan(
		&appearance.ID,
		&appearance.LogoURL,
		&appearance.PrimaryColor,
		&appearance.SecondaryColor,
		&appearance.CreatedAt,
		&appearance.UpdatedAt,
		&appearance.LogoKey,
		&appearance.LogoSize,
		&appearance.LogoType,
	)

	if err != nil {
		return nil, err
	}

	return appearance, nil
}

// Update 更新外觀信息
func (r *AppearanceRepositoryImpl) Update(appearance *Appearance) error {
	query := `
		UPDATE appearance 
		SET logo_url = $2, primary_color = $3, secondary_color = $4, updated_at = $5, 
		    logo_key = $6, logo_size = $7, logo_type = $8
		WHERE id = $1
	`

	appearance.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		appearance.ID,
		appearance.LogoURL,
		appearance.PrimaryColor,
		appearance.SecondaryColor,
		appearance.UpdatedAt,
		appearance.LogoKey,
		appearance.LogoSize,
		appearance.LogoType,
	)

	return err
}
