package service

type HandlerSymbolAdaptor interface {
	InitSpotSymbol() error
	InitFeatureSymbol() error
}
