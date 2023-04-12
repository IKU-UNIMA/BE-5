package request

type Dokumen struct {
	IdJenisDokumen int    `json:"id_jenis_dokumen"`
	Nama           string `json:"nama"`
	Keterangan     string `json:"keterangan"`
}
