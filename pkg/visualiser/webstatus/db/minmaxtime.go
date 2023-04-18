package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core/cfg"
)

// MinMaxTime is a wrapper around time to support SQL max and min
type MinMaxTime struct {
	t time.Time
}

// Scan scan value the time, implements sql.Scanner interface
func (mmt *MinMaxTime) Scan(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot parse time: %v (%T) is not string", value, value)
	}
	if s == "0001-01-01 00:00:00+00:00" {
		// no time
		mmt.t = time.Time{}
		return nil
	}
	format := "2006-01-02 15:04:05.9999999 -0700 MST"
	// l := len(format)
	// if len(s) < l {
	// 	return fmt.Errorf("cannot parse time: %v", s)
	// }
	t, err := time.Parse(format, s)
	if err != nil {
		hcl.Error("Cannot parse time", "error", s)
	}
	mmt.t = t
	return err
}

// Value return the time, implement driver.Valuer interface
func (mmt MinMaxTime) Value() (driver.Value, error) {
	return mmt.t, nil
}

// String formats the time
func (mmt MinMaxTime) String() string {
	return mmt.t.Format(cfg.TimeFormatString)
}

// IsZero is wraps time.Time.IsZero
func (mmt MinMaxTime) IsZero() bool {
	return mmt.t.IsZero()
}

// Before is wraps time.Time.Before
func (mmt MinMaxTime) Before(t MinMaxTime) bool {
	return mmt.t.Before(t.t)
}

// After is wraps time.Time.After
func (mmt MinMaxTime) After(t MinMaxTime) bool {
	return mmt.t.After(t.t)
}

// Time returns the underlying time
func (mmt MinMaxTime) Time() time.Time {
	return mmt.t
}

// Since returns the underlying time
func (mmt MinMaxTime) Since(t MinMaxTime) time.Duration {
	return mmt.t.Sub(t.t).Truncate(time.Second)
}
