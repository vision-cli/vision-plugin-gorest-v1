package tables

import (
	"gorm.io/datatypes"
)

type Users struct {
	ID   string `gorm:"primaryKey"`
	Data datatypes.JSON
}
