package alamat

import "github.com/google/wire"

var AlamatSet = wire.NewSet(
	NewAlamatRepository,
	NewAlamatService,
	NewAlamatHandler,
)
