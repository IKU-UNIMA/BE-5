package model

type JenisDokumen struct {
	ID   int    `gorm:"primaryKey"`
	Nama string `gorm:"type:varchar(255)"`
}
