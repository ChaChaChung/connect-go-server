package models

import (
	"database/sql"
	"time"
)

// Company 公司模型，對應你的 PostgreSQL 表結構
type Company struct {
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

// CompanyRepository 公司數據庫操作接口
type CompanyRepository interface {
	Create(company *Company) error
	GetByID(id string) (*Company, error)
	GetAll() ([]*Company, error)
	Update(company *Company) error
	Delete(id string) error
}

// CompanyRepositoryImpl 公司數據庫操作實現
type CompanyRepositoryImpl struct {
	db *sql.DB
}

// NewCompanyRepository 創建新的公司倉庫實例
func NewCompanyRepository(db *sql.DB) CompanyRepository {
	return &CompanyRepositoryImpl{db: db}
}

// Create 創建新公司
func (r *CompanyRepositoryImpl) Create(company *Company) error {
	query := `
		INSERT INTO appearance (id, logo_url, primary_color, secondary_color, created_at, updated_at, logo_key, logo_size, logo_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	now := time.Now()
	company.CreatedAt = now
	company.UpdatedAt = now

	_, err := r.db.Exec(query,
		company.ID,
		company.LogoURL,
		company.PrimaryColor,
		company.SecondaryColor,
		company.CreatedAt,
		company.UpdatedAt,
		company.LogoKey,
		company.LogoSize,
		company.LogoType,
	)

	return err
}

// GetByID 根據 ID 獲取公司
func (r *CompanyRepositoryImpl) GetByID(id string) (*Company, error) {
	query := `
		SELECT id, logo_url, primary_color, secondary_color, created_at, updated_at, logo_key, logo_size, logo_type
		FROM appearance WHERE id = $1
	`

	company := &Company{}
	err := r.db.QueryRow(query, id).Scan(
		&company.ID,
		&company.LogoURL,
		&company.PrimaryColor,
		&company.SecondaryColor,
		&company.CreatedAt,
		&company.UpdatedAt,
		&company.LogoKey,
		&company.LogoSize,
		&company.LogoType,
	)

	if err != nil {
		return nil, err
	}

	return company, nil
}

// GetAll 獲取所有公司
func (r *CompanyRepositoryImpl) GetAll() ([]*Company, error) {
	query := `
		SELECT id, logo_url, primary_color, secondary_color, created_at, updated_at, logo_key, logo_size, logo_type
		FROM appearance ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appearance []*Company
	for rows.Next() {
		company := &Company{}
		err := rows.Scan(
			&company.ID,
			&company.LogoURL,
			&company.PrimaryColor,
			&company.SecondaryColor,
			&company.CreatedAt,
			&company.UpdatedAt,
			&company.LogoKey,
			&company.LogoSize,
			&company.LogoType,
		)
		if err != nil {
			return nil, err
		}
		appearance = append(appearance, company)
	}

	return appearance, nil
}

// Update 更新公司信息
func (r *CompanyRepositoryImpl) Update(company *Company) error {
	query := `
		UPDATE appearance 
		SET logo_url = $2, primary_color = $3, secondary_color = $4, updated_at = $5, 
		    logo_key = $6, logo_size = $7, logo_type = $8
		WHERE id = $1
	`

	company.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		company.ID,
		company.LogoURL,
		company.PrimaryColor,
		company.SecondaryColor,
		company.UpdatedAt,
		company.LogoKey,
		company.LogoSize,
		company.LogoType,
	)

	return err
}

// Delete 刪除公司
func (r *CompanyRepositoryImpl) Delete(id string) error {
	query := `DELETE FROM appearance WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
