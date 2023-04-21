package response

type (
	Dashboard struct {
		ID     int    `json:"id"`
		Nama   string `json:"nama"`
		Jumlah int    `json:"jumlah"`
	}

	DetailDashboard struct {
		IdProdi      int    `json:"id_prodi"`
		KodeProdi    int    `json:"kode_prodi"`
		NamaProdi    string `json:"nama_prodi"`
		JenjangProdi string `json:"jenjang_prodi"`
		IdFakultas   int    `json:"id_fakultas"`
		NamaFakultas string `json:"nama_fakultas"`
		Jumlah       int    `json:"jumlah"`
	}

	DashboardTotal struct {
		TotalPublikasi  int `json:"total_publikasi"`
		TotalPaten      int `json:"total_paten"`
		TotalPengabdian int `json:"total_pengabdian"`
	}

	DashboardUmum struct {
		Fakultas  int `json:"fakultas"`
		Prodi     int `json:"prodi"`
		Dosen     int `json:"dosen"`
		Mahasiswa int `json:"mahasiswa"`
	}
)
