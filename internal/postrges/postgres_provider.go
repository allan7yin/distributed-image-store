package postrges

import "github.com/google/wire"

// ProviderSet for the postgres package
var ProviderSet = wire.NewSet(NewConnectionHandler)
