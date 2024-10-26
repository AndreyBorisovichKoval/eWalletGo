// Файл: C:\GoProject\src\eWalletGo_TestTask\models\wallet.go

package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`                          // Уникальный идентификатор кошелька...
	PhoneID            uint       `gorm:"not null" json:"phone_id"`                      // ID телефона...
	UserID             *uint      `json:"user_id"`                                       // ID пользователя, если он идентифицирован...
	ClientType         string     `gorm:"size:20;not null" json:"client_type"`           // Тип клиента...
	Status             string     `gorm:"size:20;not null;default:active" json:"status"` // Статус кошелька (active, blocked, canceled)...
	WalletNumber       string     `gorm:"size:16;unique;not null" json:"wallet_number"`  // Полный номер кошелька (16-значный)...
	MaskedWalletNumber string     `gorm:"size:19;not null" json:"masked_wallet_number"`  // Маскированный номер кошелька (например, "1234 **** **** 5678")...
	CreatedAt          time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`   // Время создания...
	UpdatedAt          *time.Time `json:"updated_at"`                                    // Время последнего обновления...
	DeletedAt          *time.Time `json:"deleted_at"`                                    // Время удаления...
	IsDeleted          bool       `gorm:"default:false" json:"is_deleted"`               // Флаг удаления...

	// Связи
	Accounts []Account `gorm:"foreignKey:WalletID;constraint:OnDelete:CASCADE" json:"accounts"` // Связь с таблицей счетов...
}

// BeforeCreate хук для установки маскированного номера кошелька перед созданием записи
func (w *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	w.MaskedWalletNumber = w.generateMaskedNumber()
	return nil
}

// generateMaskedNumber формирует маскированный номер кошелька, например, "1234 **** **** 5678"
func (w *Wallet) generateMaskedNumber() string {
	if len(w.WalletNumber) == 16 {
		return fmt.Sprintf("%s **** **** %s", w.WalletNumber[:4], w.WalletNumber[12:])
	}
	return w.WalletNumber
}

func (Wallet) TableName() string {
	return "wallets"
}
