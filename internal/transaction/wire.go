package transaction

import "github.com/google/wire"

var TransactionSet = wire.NewSet(
	NewTransactionRepository,
	NewTransactionService,
	NewTransactionHandler,
)
