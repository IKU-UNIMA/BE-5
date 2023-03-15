package model

type KategoriCapaian struct {
	ID   int    `gorm:"primaryKey"`
	Nama string `gorm:"type:varchar(255)"`
}
