package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/vogtp/som/pkg/core/cfg"
)

type MinMaxTime struct {
	t time.Time
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (mmt *MinMaxTime) Scan(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot parse time: %v (%T) is not string", value, value)
	}
	format := "2006-01-02 15:04:05.9999999"
	t, err := time.Parse(format, s[:len(format)])
	mmt.t = t
	return err
}

// Value return json value, implement driver.Valuer interface
func (mmt MinMaxTime) Value() (driver.Value, error) {
	return mmt.t, nil
}

func (mmt MinMaxTime) String() string {
	return mmt.t.Format(cfg.TimeFormatString)
}

func (mmt MinMaxTime) IsZero() bool {
	return mmt.t.IsZero()
}
