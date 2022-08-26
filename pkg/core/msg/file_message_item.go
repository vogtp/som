package msg

import (
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	"github.com/vogtp/som/pkg/core/mime"
)

// FileMsgItem contains a file
type FileMsgItem struct {
	ID       uint      `json:"-" gorm:"primaryKey"`
	ParentID uuid.UUID `json:"-"  gorm:"type:uuid"`
	Name     string    `json:"Name"`
	Type     mime.Type
	Size     int `json:"Size"`
	Payload  []byte
}

// NewFileMsgItem creates a new FileMsgItem
func NewFileMsgItem(name string, mtype mime.Type, payload []byte) *FileMsgItem {
	return &FileMsgItem{
		Name:    strcase.ToLowerCamel(name),
		Type:    mtype,
		Size:    len(payload),
		Payload: payload,
	}
}
