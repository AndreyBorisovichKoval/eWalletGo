// C:\GoProject\src\eWalletGo_TestTask\models\transaction.go

package models

import "time"

type Transaction struct {
	ID        uint       `gorm:"primaryKey" json:"id"`                        // Уникальный идентификатор транзакции...
	AccountID uint       `gorm:"not null" json:"account_id"`                  // ID счета...
	Amount    float64    `gorm:"not null" json:"amount"`                      // Сумма транзакции...
	Type      string     `gorm:"size:20;not null" json:"type"`                // Тип транзакции...
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"` // Время создания...
	UpdatedAt *time.Time `json:"updated_at"`                                  // Время последнего обновления...
	DeletedAt *time.Time `json:"deleted_at"`                                  // Время удаления...
	IsDeleted bool       `gorm:"default:false" json:"is_deleted"`             // Флаг удаления...
}

func (Transaction) TableName() string {
	return "transactions"
}