// C:\GoProject\src\eWalletGo_TestTask\models\account.go

package models

import (
	"time"
)

type Account struct {
	ID        uint       `gorm:"primaryKey" json:"id"`                        // Уникальный идентификатор счета...
	WalletID  uint       `gorm:"not null" json:"wallet_id"`                   // ID кошелька...
	UserID    *uint      `json:"user_id"`                                     // ID пользователя, если он идентифицирован...
	Balance   float64    `gorm:"not null;default:0" json:"balance"`           // Текущий баланс счета...
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"` // Время создания...
	UpdatedAt *time.Time `json:"updated_at"`                                  // Время последнего обновления...
	DeletedAt *time.Time `json:"deleted_at"`                                  // Время удаления...
	IsDeleted bool       `gorm:"default:false" json:"is_deleted"`             // Флаг удаления...

	// Связи
	Transactions []Transaction `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE" json:"transactions"` // Связь с таблицей транзакций...
}

func (Account) TableName() string {
	return "accounts"
}