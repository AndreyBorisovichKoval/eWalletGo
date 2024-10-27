// File: C:\GoProject\src\eWalletGo_TestTask\models\wallet.go

package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`                          // Unique wallet identifier...
	PhoneID            uint       `gorm:"not null" json:"phone_id"`                      // Phone ID...
	UserID             *uint      `json:"user_id"`                                       // User ID, if identified...
	ClientType         string     `gorm:"size:20;not null" json:"client_type"`           // Client type...
	Status             string     `gorm:"size:20;not null;default:active" json:"status"` // Wallet status (active, blocked, canceled)...
	WalletNumber       string     `gorm:"size:16;unique;not null" json:"wallet_number"`  // Full wallet number (16 digits)...
	MaskedWalletNumber string     `gorm:"size:19;not null" json:"masked_wallet_number"`  // Masked wallet number (e.g., "1234 **** **** 5678")...
	CreatedAt          time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`   // Creation time...
	UpdatedAt          *time.Time `json:"updated_at"`                                    // Last update time...
	DeletedAt          *time.Time `json:"deleted_at"`                                    // Deletion time...
	IsDeleted          bool       `gorm:"default:false" json:"is_deleted"`               // Deletion flag...

	// Relationships
	Accounts []Account `gorm:"foreignKey:WalletID;constraint:OnDelete:CASCADE" json:"accounts"` // Relationship with the accounts table...
}

// BeforeCreate hook to set the masked wallet number before creating a record
func (w *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	w.MaskedWalletNumber = w.generateMaskedNumber()
	return nil
}

// generateMaskedNumber generates the masked wallet number, e.g., "1234 **** **** 5678"
func (w *Wallet) generateMaskedNumber() string {
	if len(w.WalletNumber) == 16 {
		return fmt.Sprintf("%s **** **** %s", w.WalletNumber[:4], w.WalletNumber[12:])
	}
	return w.WalletNumber
}

func (Wallet) TableName() string {
	return "wallets"
}
