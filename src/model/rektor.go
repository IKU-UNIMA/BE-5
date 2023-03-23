package model

type Rektor struct {
	ID   int    `gorm:"primaryKey"`
	Nama string `gorm:"type:varchar(255)"`
	Nip  string `gorm:"type:varchar(255);unique"`
}
