package v1

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/flanksource/duty/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

const (
	SQLServerType = "sqlserver"
	PostgresType  = "postgres"
	SqliteType    = "sqlite"
	text          = "TEXT"
	jsonType      = "json"
	jsonbType     = "JSONB"
	nvarcharType  = "NVARCHAR(MAX)"
)

type ResourceSelectors []ResourceSelector

type ComponentChecks []ComponentCheck

func (rs ResourceSelectors) Value() (driver.Value, error) {
	if len(rs) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(rs)
}

func (rs *ResourceSelectors) Scan(val interface{}) error {
	if val == nil {
		*rs = ResourceSelectors{}
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal ResourceSelectors value:", val))
	}
	return json.Unmarshal(ba, rs)
}

// GormDataType gorm common data type
func (rs ResourceSelectors) GormDataType() string {
	return "resourceSelectors"
}

// GormDBDataType gorm db data type
func (ResourceSelectors) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case SqliteType:
		return jsonType
	case PostgresType:
		return jsonbType
	case SQLServerType:
		return nvarcharType
	}
	return ""
}

func (rs ResourceSelectors) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := json.Marshal(rs)
	return gorm.Expr("?", string(data))
}

func (cs ComponentChecks) Value() (driver.Value, error) {
	if len(cs) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(cs)
}

func (cs *ComponentChecks) Scan(val interface{}) error {
	if val == nil {
		*cs = ComponentChecks{}
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal componentChecks value:", val))
	}
	return json.Unmarshal(ba, cs)
}

// GormDataType gorm common data type
func (cs ComponentChecks) GormDataType() string {
	return "componentChecks"
}

// GormDBDataType gorm db data type
func (ComponentChecks) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case SqliteType:
		return jsonType
	case PostgresType:
		return jsonbType
	case SQLServerType:
		return nvarcharType
	}
	return ""
}

func (cs ComponentChecks) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := json.Marshal(cs)
	return gorm.Expr("?", string(data))
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (s Summary) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (s *Summary) Scan(val interface{}) error {
	if val == nil {
		*s = Summary{}
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal properties value:", val))
	}
	err := json.Unmarshal(ba, s)
	return err
}

// GormDataType gorm common data type
func (Summary) GormDataType() string {
	return "summary"
}

func (Summary) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case SqliteType:
		return text
	case PostgresType:
		return jsonbType
	case SQLServerType:
		return nvarcharType
	}
	return ""
}

func (s Summary) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := json.Marshal(s)
	return gorm.Expr("?", data)
}

// LogSelector ...
type LogSelector struct {
	Name   string            `json:"name,omitempty" yaml:"name,omitempty"`
	Type   string            `json:"type,omitempty" yaml:"type,omitempty" template:"true"`
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty" template:"true"`
}

type LogSelectors []LogSelector

func (t LogSelectors) Value() (driver.Value, error) {
	return types.GenericStructValue(t, true)
}

func (t *LogSelectors) Scan(val any) error {
	return types.GenericStructScan(&t, val)
}

func (LogSelectors) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case types.SqliteType:
		return types.JSONType
	case types.PostgresType:
		return types.JSONBType
	case types.SQLServerType:
		return types.NVarcharType
	}

	return ""
}

func (rs LogSelectors) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := json.Marshal(rs)
	return gorm.Expr("?", string(data))
}
