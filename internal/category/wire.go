package category

import "github.com/google/wire"

var CategorySet = wire.NewSet(
	NewCategoryRepository,
	NewCategoryService,
	NewCategoryHandler,
)
