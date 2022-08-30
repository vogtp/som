package msg

import (
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	"github.com/vogtp/som/pkg/core/mime"
)

// FileMsgItem contains a file
type FileMsgItem struct {
	ID       uuid.UUID `json:"ID"  gorm:"primaryKey;type:uuid"`
	Name     string    `json:"Name"`
	Type     mime.Type
	Size     int `json:"Size"`
	Payload  []byte
}

// NewFileMsgItem creates a new FileMsgItem
func NewFileMsgItem(name string, mtype mime.Type, payload []byte) *FileMsgItem {
	return &FileMsgItem{
		ID:      uuid.New(),
		Name:    strcase.ToLowerCamel(name),
		Type:    mtype,
		Size:    len(payload),
		Payload: payload,
	}
}
