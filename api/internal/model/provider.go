package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	return string(b), err
}

func (s *StringSlice) Scan(val interface{}) error {
	if val == nil {
		*s = StringSlice{}
		return nil
	}
	var bytes []byte
	switch v := val.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	}
	return json.Unmarshal(bytes, s)
}

type InterleavedType struct {
	Data interface{}
}

func (i InterleavedType) Value() (driver.Value, error) {
	if i.Data == nil {
		return nil, nil
	}
	return json.Marshal(i.Data)
}

func (i *InterleavedType) Scan(val interface{}) error {
	if val == nil {
		i.Data = nil
		return nil
	}
	var bytes []byte
	switch v := val.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	}
	return json.Unmarshal(bytes, &i.Data)
}

type Provider struct {
	ID         uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderID string      `gorm:"uniqueIndex;size:100;not null" json:"provider_id"`
	Name       string      `gorm:"size:255;not null" json:"name"`
	Npm        string      `gorm:"size:255;not null;default:''" json:"npm"`
	Env        StringSlice `gorm:"type:json;not null" json:"env"`
	DocURL     string      `gorm:"size:512;default:''" json:"doc_url"`
	ApiURL     string      `gorm:"size:512;default:''" json:"api_url"`
	SyncedAt   *time.Time  `gorm:"index:idx_providers_synced_at" json:"synced_at"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

func (Provider) TableName() string { return "providers" }
