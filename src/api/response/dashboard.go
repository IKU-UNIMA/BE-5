package response

type (
	Dashboard struct {
		Target     string                       `json:"target"`
		Total      int                          `json:"total"`
		TotalDosen int                          `json:"total_dosen"`
		Pencapaian string                       `json:"pencapaian"`
		Detail     []DashboardDetailPerFakultas `json:"detail"`
	}

	DashboardDetailPerFakultas struct {
		ID          int    `json:"id"`
		Fakultas    string `json:"fakultas"`
		JumlahDosen int    `json:"jumlah_dosen"`
		Jumlah      int    `json:"jumlah"`
		Persentase  string `json:"persentase"`
	}

	DashboardPerProdi struct {
		Fakultas   string                    `json:"fakultas"`
		Total      int                       `json:"total"`
		TotalDosen int                       `json:"total_dosen"`
		Pencapaian string                    `json:"pencapaian"`
		Detail     []DashboardDetailPerProdi `json:"detail"`
	}

	DashboardDetailPerProdi struct {
		Prodi       string `json:"prodi"`
		JumlahDosen int    `json:"jumlah_dosen"`
		Jumlah      int    `json:"jumlah"`
		Persentase  string `json:"persentase"`
	}

	DashboardTotal struct {
		Nama  string `json:"nama"`
		Total int    `json:"total"`
	}

	DashboardUmum struct {
		Fakultas  int `json:"fakultas"`
		Prodi     int `json:"prodi"`
		Dosen     int `json:"dosen"`
		Mahasiswa int `json:"mahasiswa"`
	}
)
