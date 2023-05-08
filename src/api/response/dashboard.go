package response

type (
	Dashboard struct {
		Target     string            `json:"target"`
		Total      int               `json:"total"`
		TotalDosen int               `json:"total_dosen"`
		Pencapaian string            `json:"pencapaian"`
		Detail     []DashboardDetail `json:"detail"`
	}

	DashboardDetail struct {
		ID          int    `json:"id"`
		Fakultas    string `json:"fakultas"`
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
