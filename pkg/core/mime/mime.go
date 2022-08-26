package mime

import (
	"database/sql/driver"
	"fmt"
)

// Type represents a mime type
type Type struct {
	MimeType string
	Ext      string
}

var (
	// Png is a png image
	Png = Type{MimeType: "image/png", Ext: "png"}
	// Jpg is a jpeg image
	//	Jpg = Type{MimeType: "image/jpeg", Ext: "jpeg"}
	// HTML is a html page
	HTML = Type{MimeType: "text/html", Ext: "html"}
)

var allTypes = []Type{Png, HTML}

// Scan scan value into mime.Type, implements sql.Scanner interface
func (t *Type) Scan(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot parse time: %v (%T) is not string", value, value)
	}
	for _, mt := range allTypes {
		if mt.MimeType == s {
			t.MimeType = mt.MimeType
			t.Ext = mt.Ext
			return nil
		}
	}
	return fmt.Errorf("no such mime type: %s", s)
}

// Value return MimeType value, implement driver.Valuer interface
func (t Type) Value() (driver.Value, error) {
	return t.MimeType, nil
}
