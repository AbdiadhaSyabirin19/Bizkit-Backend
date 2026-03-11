package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

type Permissions map[string][]string

func (p Permissions) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Permissions) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, p)
}

type Role struct {
	gorm.Model
	Name        string      `json:"name" gorm:"not null"`
	Permissions Permissions `json:"permissions" gorm:"type:json"`
}