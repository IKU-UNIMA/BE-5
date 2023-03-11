package response

type Prodi struct {
	ID         int      `json:"id"`
	IdFakultas int      `json:"-"`
	Nama       string   `json:"nama"`
	Jenjang    string   `json:"jenjang"`
	Fakultas   Fakultas `gorm:"foreignKey:IdFakultas" json:"fakultas"`
}
