package response

type Prodi struct {
	ID         int               `json:"id"`
	IdFakultas int               `json:"-"`
	KodeProdi  int               `json:"kode_prodi"`
	Nama       string            `json:"nama"`
	Jenjang    string            `json:"jenjang"`
	Fakultas   FakultasReference `gorm:"foreignKey:IdFakultas" json:"fakultas"`
}

type FakultasReference struct {
	ID   int    `json:"id"`
	Nama string `json:"nama"`
}

func (FakultasReference) TableName() string {
	return "fakultas"
}
