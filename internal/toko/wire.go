package toko

import "github.com/google/wire"

var TokoSet = wire.NewSet(
	NewTokoRepository,
	NewTokoService,
	NewTokoHandler,
)
