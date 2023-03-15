package model

type JenisPenelitian struct {
	ID   int    `gorm:"primaryKey"`
	Nama string `gorm:"type:varchar(255)"`
}
