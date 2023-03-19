package model

type Akun struct {
	ID       int    `gorm:"primaryKey"`
	Email    string `gorm:"type:varchar(255);unique"`
	Password string `gorm:"type:varchar(255)"`
	Role     string `gorm:"type:enum('rektor', 'dosen', 'admin')"`
	Admin    Admin  `gorm:"foreignKey:ID;constraint:OnDelete:CASCADE"`
	Dosen    Dosen  `gorm:"foreignKey:ID;constraint:OnDelete:CASCADE"`
}
