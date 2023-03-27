package response

type (
	Dosen struct {
		ID                int            `json:"id"`
		Nama              string         `json:"nama"`
		Nidn              string         `json:"nidn"`
		Nip               string         `json:"nip"`
		IdFakultas        int            `json:"-"`
		IdProdi           int            `json:"-"`
		StatusKepegawaian string         `json:"status_kepegawain"`
		Fakultas          Fakultas       `gorm:"foreignKey:IdFakultas" json:"fakultas"`
		Prodi             ProdiReference `gorm:"foreignKey:IdProdi" json:"prodi"`
	}

	DetailDosen struct {
		ID                int            `json:"id"`
		Nama              string         `json:"nama"`
		Nidn              string         `json:"nidn"`
		Nip               string         `json:"nip"`
		IdFakultas        int            `json:"-"`
		IdProdi           int            `json:"-"`
		Email             string         `json:"email"`
		StatusKepegawaian string         `json:"status_kepegawain"`
		Fakultas          Fakultas       `gorm:"foreignKey:IdFakultas" json:"fakultas"`
		Prodi             ProdiReference `gorm:"foreignKey:IdProdi" json:"prodi"`
	}
)
