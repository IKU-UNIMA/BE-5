package response

type Fakultas struct {
	ID   int    `json:"id"`
	Nama string `json:"nama"`
}

type DetailFakultas struct {
	ID    int              `json:"id"`
	Nama  string           `json:"nama"`
	Prodi []ProdiReference `gorm:"foreignKey:IdFakultas" json:"prodi"`
}

type ProdiReference struct {
	ID         int    `json:"id"`
	IdFakultas int    `json:"-"`
	KodeProdi  int    `json:"kode_prodi"`
	Nama       string `json:"nama"`
	Jenjang    string `json:"jenjang"`
}

func (ProdiReference) TableName() string {
	return "prodi"
}
