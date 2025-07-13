package symbol

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SymbolFuturesPrices struct {
	GUID          uuid.UUID `gorm:"primaryKey"`
	Symbol        string
	UnifiedSymbol string
	Price         string
	MarkPrice     string
	FundingRate   string
	Exchange      string
	ChainId       string
	Base          string
	Quote         string
	Timestamp     uint64
}

type symbolFuturesPricesDB struct {
	gorm *gorm.DB
}

func NewSymbolFuturesPricesDB(db *gorm.DB) SymbolFuturesPricesDB {
	return &symbolFuturesPricesDB{
		gorm: db,
	}
}

type SymbolFuturesPricesDB interface {
	SaveSymbolFuturesPrices(*[]SymbolFuturesPrices) error
	UpdateSymbolFuturesPrices(*[]SymbolFuturesPrices) error
}

func (db *symbolFuturesPricesDB) SaveSymbolFuturesPrices(futuresPrices *[]SymbolFuturesPrices) error {
	result := db.gorm.CreateInBatches(&futuresPrices, len(*futuresPrices))
	return result.Error
}

func (db *symbolFuturesPricesDB) UpdateSymbolFuturesPrices(futuresPrices *[]SymbolFuturesPrices) error {
	result := db.gorm.Save(&futuresPrices)
	return result.Error
}
