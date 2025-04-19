package symbol

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SymbolMapping struct {
	GUID          uuid.UUID `gorm:"primaryKey"`
	UnifiedSymbol string
	Symbol        string
	Exchange      string
	Chain         string
	Base          string
	Quote         string
	Timestamp     uint64
}

type symbolMappingDB struct {
	gorm *gorm.DB
}

func NewSymbolMappingDB(db *gorm.DB) SymbolMappingDB {
	return &symbolMappingDB{
		gorm: db,
	}
}

type SymbolMappingDB interface {
	SaveSymbolMapping(*[]SymbolMapping) error
	UpdateSymbolMapping(*[]SymbolMapping) error
}

func (db *symbolMappingDB) SaveSymbolMapping(symbolMappings *[]SymbolMapping) error {
	result := db.gorm.CreateInBatches(&symbolMappings, len(*symbolMappings))
	return result.Error
}

func (db *symbolMappingDB) UpdateSymbolMapping(symbolMappings *[]SymbolMapping) error {
	result := db.gorm.Save(&symbolMappings)
	return result.Error
}
