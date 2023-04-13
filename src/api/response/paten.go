package response

type (
	Paten struct {
		ID                int             `json:"id"`
		IdDosen           int             `json:"-"`
		Dosen             DosenReference  `gorm:"foreignKey:IdDosen" json:"dosen"`
		IdKategori        int             `json:"-"`
		IdJenisPenelitian int             `json:"-"`
		Judul             string          `json:"judul"`
		Tanggal           string          `json:"tanggal"`
		Kategori          KategoriPaten   `gorm:"foreignKey:IdKategori" json:"kategori_kegiatan"`
		JenisPenelitian   JenisPenelitian `gorm:"foreignKey:IdJenisPenelitian" json:"jenis"`
	}

	DetailPaten struct {
		ID                int             `json:"id"`
		IdKategori        int             `json:"id_kategori"`
		IdDosen           int             `json:"-"`
		Dosen             Dosen           `gorm:"foreignKey:IdDosen" json:"dosen"`
		IdJenisPenelitian int             `json:"-"`
		IdKategoriCapaian int             `json:"-"`
		Judul             string          `json:"judul"`
		Tanggal           string          `json:"tanggal"`
		JumlahHalaman     int             `json:"jumlah_halaman"`
		Penyelenggara     string          `json:"penyelenggara"`
		Penerbit          string          `json:"penerbit"`
		Isbn              string          `json:"isbn"`
		TautanEksternal   string          `json:"tautan_eksternal"`
		Keterangan        string          `json:"keterangan"`
		JenisPenelitian   JenisPenelitian `gorm:"foreignKey:IdJenisPenelitian" json:"jenis_penelitian"`
		KategoriCapaian   KategoriCapaian `gorm:"foreignKey:IdKategoriCapaian" json:"kategori_capaian"`
		PenulisDosen      []PenulisPaten  `gorm:"foreignKey:IdPaten" json:"penulis_dosen"`
		PenulisMahasiswa  []PenulisPaten  `gorm:"foreignKey:IdPaten" json:"penulis_mahasiswa"`
		PenulisLain       []PenulisPaten  `gorm:"foreignKey:IdPaten" json:"penulis_lain"`
		Dokumen           []DokumenPaten  `gorm:"foreignKey:IdPaten" json:"dokumen"`
	}

	KategoriPaten struct {
		ID                   int    `json:"id"`
		IdJenisKategoriPaten int    `json:"-"`
		Nama                 string `json:"nama"`
	}

	JenisKategoriPaten struct {
		ID            int             `json:"id"`
		Nama          string          `json:"nama"`
		KategoriPaten []KategoriPaten `gorm:"foreignKey:IdJenisKategoriPaten" json:"kategori_paten"`
	}

	PenulisPaten struct {
		ID       int    `json:"id"`
		IdPaten  int    `json:"-"`
		Nama     string `json:"nama"`
		Urutan   int    `json:"urutan"`
		Afiliasi string `json:"afiliasi"`
		Peran    string `json:"peran"`
		IsAuthor bool   `json:"is_author"`
	}

	DokumenPaten struct {
		ID             string       `json:"id"`
		IdPaten        int          `json:"-"`
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
