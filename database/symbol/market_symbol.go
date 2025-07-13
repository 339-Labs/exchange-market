package symbol

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MarketSymbol struct {
	GUID          uuid.UUID `gorm:"primaryKey"`
	Symbol        string
	UnifiedSymbol string
	InstType      string
	Exchange      string
	ChainId       string
	Base          string
	Quote         string
	Timestamp     uint64
}

type marketSymbolDB struct {
	gorm *gorm.DB
}

func NewMarketSymbolDB(db *gorm.DB) MarketSymbolDB {
	return &marketSymbolDB{
		gorm: db,
	}
}

type MarketSymbolDB interface {
	SaveMarketSymbol(*[]MarketSymbol) error
	UpdateMarketSymbol(*[]MarketSymbol) error
}

func (db *marketSymbolDB) SaveMarketSymbol(symbolMappings *[]MarketSymbol) error {
	result := db.gorm.CreateInBatches(&symbolMappings, len(*symbolMappings))
	return result.Error
}

func (db *marketSymbolDB) UpdateMarketSymbol(symbolMappings *[]MarketSymbol) error {
	result := db.gorm.Save(&symbolMappings)
	return result.Error
}
