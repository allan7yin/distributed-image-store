package services

import "github.com/google/wire"

// ProviderSet for ImageService
var ProviderSet = wire.NewSet(NewImageService)
