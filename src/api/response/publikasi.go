package response

type (
	Publikasi struct {
		ID                int               `json:"id"`
		IdDosen           int               `json:"-"`
		Dosen             DosenReference    `gorm:"foreignKey:IdDosen" json:"dosen"`
		IdKategori        int               `json:"-"`
		IdJenisPenelitian int               `json:"-"`
		Judul             string            `json:"judul"`
		TanggalTerbit     string            `json:"tanggal_terbit"`
		Kategori          KategoriPublikasi `gorm:"foreignKey:IdKategori" json:"kategori_kegiatan"`
		JenisPenelitian   JenisPenelitian   `gorm:"foreignKey:IdJenisPenelitian" json:"jenis"`
	}

	DetailPublikasi struct {
		ID                  int                `json:"id"`
		IdKategori          int                `json:"id_kategori"`
		IdDosen             int                `json:"-"`
		Dosen               Dosen              `gorm:"foreignKey:IdDosen" json:"dosen"`
		IdJenisPenelitian   int                `json:"-"`
		IdKategoriCapaian   int                `json:"-"`
		Judul               string             `json:"judul"`
		JudulAsli           string             `json:"judul_asli"`
		JudulChapter        string             `json:"judul_chapter"`
		NamaJurnal          string             `json:"nama_jurnal"`
		NamaKoranMajalah    string             `json:"nama_koran_majalah"`
		NamaSeminar         string             `json:"nama_seminar"`
		TautanLamanJurnal   string             `json:"tautan_laman_jurnal"`
		TanggalTerbit       string             `json:"tanggal_terbit"`
		WaktuPelaksanaan    string             `json:"waktu_pelaksanaan"`
		Volume              string             `json:"volume"`
		Edisi               string             `json:"edisi"`
		Nomor               string             `json:"nomor"`
		Halaman             string             `json:"halaman"`
		JumlahHalaman       int                `json:"jumlah_halaman"`
		Penerbit            string             `json:"penerbit"`
		Penyelenggara       string             `json:"penyelenggara"`
		KotaPenyelenggaraan string             `json:"kota_penyelenggaraan"`
		IsSeminar           bool               `json:"is_seminar"`
		IsProsiding         bool               `json:"is_prosiding"`
		Bahasa              string             `json:"bahasa"`
		Doi                 string             `json:"doi"`
		Isbn                string             `json:"isbn"`
		Issn                string             `json:"issn"`
		EIssn               string             `json:"e_issn"`
		Tautan              string             `json:"tautan"`
		Keterangan          string             `json:"keterangan"`
		JenisPenelitian     JenisPenelitian    `gorm:"foreignKey:IdJenisPenelitian" json:"jenis_penelitian"`
		KategoriCapaian     KategoriCapaian    `gorm:"foreignKey:IdKategoriCapaian" json:"kategori_capaian"`
		PenulisDosen        []PenulisPublikasi `gorm:"foreignKey:IdPublikasi" json:"penulis_dosen"`
		PenulisMahasiswa    []PenulisPublikasi `gorm:"foreignKey:IdPublikasi" json:"penulis_mahasiswa"`
		PenulisLain         []PenulisPublikasi `gorm:"foreignKey:IdPublikasi" json:"penulis_lain"`
		Dokumen             []DokumenPublikasi `gorm:"foreignKey:IdPublikasi" json:"dokumen"`
	}

	KategoriPublikasi struct {
		ID                       int    `json:"id"`
		IdJenisKategoriPublikasi int    `json:"-"`
		Nama                     string `json:"nama"`
	}

	JenisKategoriPublikasi struct {
		ID                int                 `json:"id"`
		Nama              string              `json:"nama"`
		KategoriPublikasi []KategoriPublikasi `gorm:"foreignKey:IdJenisKategoriPublikasi" json:"kategori_publikasi"`
	}

	PenulisPublikasi struct {
		ID          int    `json:"id"`
		IdPublikasi int    `json:"-"`
		Nama        string `json:"nama"`
		Urutan      int    `json:"urutan"`
		Afiliasi    string `json:"afiliasi"`
		Peran       string `json:"peran"`
		IsAuthor    bool   `json:"is_author"`
	}

	DokumenPublikasi struct {
		ID             string       `json:"id"`
		IdPublikasi    int          `json:"-"`
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
