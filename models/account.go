package models

import (
	"time"
)

type Account struct {
	UUID      string `gorm:"primaryKey"`
	PlayerID  uint
	Player    Player `gorm:"foreignKey:PlayerID"`
	EconomyID uint
	Economy   Economy `gorm:"foreignKey:EconomyID"`
	SettingID uint
	Setting   Setting `gorm:"foreignKey:SettingID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccountResponse struct {
	UUID      string    `json:"uuid"`
	Player    Player    `json:"player"`
	Economy   Economy   `json:"economy"`
	Setting   Setting   `json:"settings"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Player struct {
	ID      uint    `json:"id" gorm:"primaryKey"`
	OldName *string `json:"old_name"`
	NewName string  `json:"new_name"`
}

type Economy struct {
	ID     uint    `json:"id" gorm:"primaryKey"`
	Cash   float64 `json:"cash" gorm:"default:0"`
	Vault  float64 `json:"vault" gorm:"default:0"`
	Bank   float64 `json:"bank" gorm:"default:0"`
	Crypto float64 `json:"crypto" gorm:"default:0"`
	Total  float64 `json:"total" gorm:"default:0"`
}

type Setting struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Language string `json:"language"`
}
