package static

const (
	// HTTP Message
	Success            = "success"
	ToManyRequest      = "terjadi kesalahan coba beberapa saat lagi"
	Authorization      = "terjadi kesalahan akses ditolak"
	SomethingWrong     = "terjadi kesalahan pada sistem"
	InvalidAccessLogin = "email atau kata sandi salah"
	BadRequest         = "data payload tidak benar"
	Conflict           = "%v telah ditemukan"

	// General Message
	ValueAreadyExist = "%v telah digunakan"
	DataNotFound     = "%v tidak ditemukan"
	MinValue         = "%v harus lebih dari %v"
	ValueNotValid    = "format %v tidak benar"
	EmptyValue       = "%v tidak boleh kosong"
)
