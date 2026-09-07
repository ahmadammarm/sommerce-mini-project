package upload

import "github.com/google/wire"

// UploadSet provides the Dependency Injection setup for the Upload domain.
var UploadSet = wire.NewSet(
	NewUploadHandler,
)
