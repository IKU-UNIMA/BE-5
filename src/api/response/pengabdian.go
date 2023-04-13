package response

type (
	Pengabdian struct {
		ID               int            `json:"id"`
		IdDosen          int            `json:"-"`
		Dosen            DosenReference `gorm:"foreignKey:IdDosen" json:"dosen"`
		TahunPelaksanaan string         `json:"tahun_pelaksanaan"`
		LamaKegiatan     uint           `json:"lama_kegiatan"`
	}

	DetailPengabdian struct {
		ID                      int                 `json:"id"`
		IdKategori              int                 `json:"id_kategori"`
		IdDosen                 int                 `json:"-"`
		Dosen                   Dosen               `gorm:"foreignKey:IdDosen" json:"dosen"`
		TahunKegiatan           int                 `json:"tahun_anggaran"`
		Afiliasi                string              `json:"afiliasi"`
		KelompokBidang          string              `json:"kelompok_bidang"`
		JenisSkim               string              `json:"jenis_skim"`
		NoSkPenugasan           string              `json:"no_sk_penugasan"`
		TglSkPenugasan          string              `json:"tgl_sk_penugasan"`
		LamaKegiatan            uint                `json:"lama_kegiatan"`
		Judul                   string              `json:"judul_kegiatan"`
		LokasiKegiatan          string              `json:"lokasi_kegiatan"`
		TahunPelaksanaanKe      uint                `json:"tahun_pelaksanaan_ke"`
		DanaDariDikti           float64             `json:"dana_dari_dikti"`
		DanaDariPerguruanTinggi float64             `json:"dana_dari_perguruan_tinggi"`
		DanaDariInstitusiLain   float64             `json:"dana_dari_institusi_lain"`
		InKind                  string              `json:"in_kind"`
		AnggotaDosen            []AnggotaPengabdian `gorm:"foreignKey:IdPengabdian" json:"anggota_dosen"`
		AnggotaMahasiswa        []AnggotaPengabdian `gorm:"foreignKey:IdPengabdian" json:"anggota_mahasiswa"`
		AnggotaEksternal        []AnggotaPengabdian `gorm:"foreignKey:IdPengabdian" json:"anggota_eksternal"`
		Dokumen                 []DokumenPengabdian `gorm:"foreignKey:IdPengabdian" json:"dokumen"`
	}

	KategoriPengabdian struct {
		ID                   int    `json:"id"`
		IdJenisKategoriPaten int    `json:"-"`
		Nama                 string `json:"nama"`
	}

	JenisKategoriPengabdian struct {
		ID            int             `json:"id"`
		Nama          string          `json:"nama"`
		KategoriPaten []KategoriPaten `gorm:"foreignKey:IdJenisKategoriPaten" json:"kategori_paten"`
	}

	AnggotaPengabdian struct {
		ID           int    `json:"id"`
		IdPengabdian int    `json:"-"`
		Nama         string `json:"nama"`
		Peran        string `json:"peran"`
		IsActive     bool   `json:"is_active"`
	}

	DokumenPengabdian struct {
		ID             string       `json:"id"`
		IdPengabdian   int          `json:"-"`
		IdJenisDokumen int          `json:"-"`
		Nama           string       `json:"nama"`
		NamaFile       string       `json:"nama_file"`
		JenisFile      string       `json:"jenis_file"`
		Keterangan     string       `json:"keterangan"`
		Url            string       `json:"url"`
		TanggalUpload  string       `json:"tanggal_upload"`
		JenisDokumen   JenisDokumen `gorm:"foreignKey:IdJenisDokumen" json:"jenis_dokumen"`
	}
)
