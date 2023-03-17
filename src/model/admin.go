package model

type Admin struct {
	ID     int    `gorm:"primaryKey"`
	Nama   string `gorm:"type:varchar(255)"`
	Nip    string `gorm:"type:varchar(255)"`
	Bagian string `gorm:"type:varchar(255)"`
}
