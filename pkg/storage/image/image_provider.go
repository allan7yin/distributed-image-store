package image

import (
	"github.com/google/wire"
)

// ProviderSet for the imagestore package
var ProviderSet = wire.NewSet(NewImageStore)
