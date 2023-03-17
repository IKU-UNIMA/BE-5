package model

type Akun struct {
	ID       int    `gorm:"primaryKey"`
	Email    string `gorm:"type:varchar(255)"`
	Password string `gorm:"type:varchar(255)"`
	Role     string `gorm:"type:enum('rektor', 'dosen', 'admin')"`
}
