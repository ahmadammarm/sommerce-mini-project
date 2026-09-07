package produk

import "github.com/google/wire"

var ProdukSet = wire.NewSet(
	NewProdukRepository,
	NewProdukService,
	NewProdukHandler,
)
