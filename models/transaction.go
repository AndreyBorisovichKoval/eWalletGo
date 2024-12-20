// C:\GoProject\src\eWalletGo_TestTask\models\transaction.go

package models

import "time"

// Transaction represents a wallet account transaction
type Transaction struct {
	ID        uint       `gorm:"primaryKey" json:"id"`                        // Unique transaction identifier...
	AccountID uint       `gorm:"not null" json:"account_id"`                  // Account ID (AccountID)...
	Amount    float64    `gorm:"not null" json:"amount"`                      // Transaction amount...
	Type      string     `gorm:"size:20;not null" json:"type"`                // Transaction type (recharge, withdrawal)...
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"` // Creation time...
	UpdatedAt *time.Time `json:"updated_at"`                                  // Last update time...
	DeletedAt *time.Time `json:"deleted_at"`                                  // Deletion time...
	IsDeleted bool       `gorm:"default:false" json:"is_deleted"`             // Deletion flag...
}

func (Transaction) TableName() string {
	return "transactions"
}
