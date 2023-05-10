package model

type Target struct {
	Bagian string `gorm:"type:varchar(255)"`
	Target float32
	Tahun  int
}
