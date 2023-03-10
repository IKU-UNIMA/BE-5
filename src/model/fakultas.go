package model

type Fakultas struct {
	ID   int    `gorm:"primaryKey"`
	Nama string `gorm:"type:varchar(255);unique"`
}
