package model

type Dosen struct {
	ID                int    `gorm:"primaryKey"`
	Nama              string `gorm:"type:varchar(255)"`
	Nidn              string `gorm:"type:varchar(255);unique"`
	Nip               string `gorm:"type:varchar(255);unique"`
	IdFakultas        int
	IdProdi           int
	StatusKepegawaian string   `gorm:"type:varchar(255)"`
	Fakultas          Fakultas `gorm:"foreignKey:IdFakultas;constraint:OnDelete:SET NULL"`
	Prodi             Prodi    `gorm:"foreignKey:IdProdi;constraint:OnDelete:SET NULL"`
}
