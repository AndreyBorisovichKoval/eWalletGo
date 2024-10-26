// C:\GoProject\src\eWalletGo_TestTask\models\limit_settings.go

package models

import "time"

type LimitSettings struct {
	ID            uint       `gorm:"primaryKey" json:"id"`                        // Уникальный идентификатор настройки лимита...
	ClientType    string     `gorm:"size:20;unique;not null" json:"client_type"`  // Тип клиента...
	DefaultLimit  float64    `gorm:"not null" json:"default_limit"`               // Лимит по умолчанию...
	CustomLimit   *float64   `json:"custom_limit"`                                // Индивидуальный лимит...
	AgreementDate *time.Time `json:"agreement_date"`                              // Дата соглашения об индивидуальном лимите...
	CreatedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"` // Время создания...
	UpdatedAt     *time.Time `json:"updated_at"`                                  // Время последнего обновления...
	DeletedAt     *time.Time `json:"deleted_at"`                                  // Время удаления...
	IsDeleted     bool       `gorm:"default:false" json:"is_deleted"`             // Флаг удаления...
}

func (LimitSettings) TableName() string {
	return "limit_settings"
}

// WalletWithLimit включает данные о кошельке и лимите...
type WalletWithLimit struct {
	Balance   float64 `json:"balance"`    // Баланс кошелька
	AccountID uint    `json:"account_id"` // Идентификатор аккаунта
	MaxLimit  float64 `json:"max_limit"`  // Применяемый лимит
}
