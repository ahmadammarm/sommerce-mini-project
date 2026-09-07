package user

import "github.com/google/wire"

var UserSet = wire.NewSet(
	NewUserRepository,
	NewUserService,
	NewUserHandler,
)
