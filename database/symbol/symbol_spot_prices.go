package symbol

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SymbolSpotPrices struct {
	GUID          uuid.UUID `gorm:"primaryKey"`
	Symbol        string
	UnifiedSymbol string
	Price         string
	Exchange      string
	ChainId       string
	Base          string
	Quote         string
	Timestamp     uint64
}

type symbolSpotPricesDB struct {
	gorm *gorm.DB
}

func NewSymbolSpotPricesDB(db *gorm.DB) SymbolSpotPricesDB {
	return &symbolSpotPricesDB{
		gorm: db,
	}
}

type SymbolSpotPricesDB interface {
	SaveSymbolSpotPrices(*[]SymbolSpotPrices) error
	UpdateSymbolSpotPrices(*[]SymbolSpotPrices) error
}

func (db *symbolSpotPricesDB) SaveSymbolSpotPrices(symbolSpotPrices *[]SymbolSpotPrices) error {
	result := db.gorm.CreateInBatches(&symbolSpotPrices, len(*symbolSpotPrices))
	return result.Error
}

func (db *symbolSpotPricesDB) UpdateSymbolSpotPrices(symbolSpotPrices *[]SymbolSpotPrices) error {
	result := db.gorm.Save(&symbolSpotPrices)
	return result.Error
}
