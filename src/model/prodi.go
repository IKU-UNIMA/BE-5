package model

type Prodi struct {
	ID         int `gorm:"primaryKey"`
	IdFakultas int
	Nama       string `gorm:"type:varchar(255);unique"`
	Jenjang    string `gorm:"type:varchar(60)"`
}
