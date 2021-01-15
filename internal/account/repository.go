package account

import (
	"github.com/Stardome-Team/Service-Template/pkg/database"
	"github.com/Stardome-Team/Service-Template/pkg/logset"
	"gorm.io/gorm"
)

type repository struct {
	db     *database.DB
	logger logset.Logger
}

// Account model
type Account struct {
	gorm.Model
	AccountID string `gorm:"column:account_id"`
	Username  string `gorm:"column:user_name"`
}

// Repository contains interfaces for authentication services
type Repository interface {
	CreateAccount(account Account) (int, error)
}

// NewRepository creates a new instance for authentication repository
func NewRepository(db *database.DB, l logset.Logger) Repository {
	return &repository{db: db, logger: l}
}

// CreateAccountWithProfile creates Account Profile with association
func (r *repository) CreateAccount(account Account) (int, error) {

	db := r.db.DB()

	result := db.Create(account)

	return int(result.RowsAffected), result.Error
}
