package wire

//// Provider sets for different components
//var DataStoreProviderSet = wire.NewSet(
//	postrges.ProviderSet,
//	image.ProviderSet,
//	s3.ProviderSet,
//)
//
//var ServiceProviderSet = wire.NewSet(
//	services.ProviderSet,
//)
//
//var HandlerProviderSet = wire.NewSet(
//	handlers.ProviderSet,
//)
//
//// Aggregate provider set
//var AppProviderSet = wire.NewSet(
//	DataStoreProviderSet,
//	ServiceProviderSet,
//	HandlerProviderSet,
//)
//
//// Injector functions
//// InitializeImageHandler initializes the ImageHandler.
//func InitializeImageHandler() (*handlers.ImageHandler, error) {
//	wire.Build(AppProviderSet)
//	return nil, nil
//}

// InitializeUserHandler initializes the UserHandler.
//func InitializeUserHandler() (*handlers.UserHandler, error) {
//	wire.Build(AppProviderSet)
//	return nil, nil
//}
